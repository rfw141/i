package app

import (
	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/x/errorx"
	"github.com/rfw141/i/x/jsonx"
	"reflect"
)

type TestCaseInput struct {
	Ctx      Ctx
	Rpc      interface{}
	Body     string
	Req      proto.Message
	BeforeFn func()
	AfterFn  func()
}

type TestCase struct {
	Description string
	Input       TestCaseInput
	Error       bool
	Result      interface{}
}

func TestGetReq(req proto.Message, str string, init ...func()) proto.Message {
	for _, v := range init {
		v()
	}
	jsonx.MustUnmarshal([]byte(str), req)
	return req
}

func TestCallByJson(rpcFn interface{}, ctx Ctx, jsonStr string, init ...func()) (
	proto.Message, error) {
	return TestCall(rpcFn, ctx, TestGetReq(
		reflect.New(reflect.ValueOf(rpcFn).Type().In(1).Elem()).Interface().(proto.Message),
		jsonStr,
	), init...)
}

func TestCall(rpcFn interface{}, ctx Ctx, req proto.Message, init ...func()) (
	proto.Message, error) {
	for _, v := range init {
		v()
	}
	res := reflect.ValueOf(rpcFn).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(req),
	})

	retErr := res[1].Interface()
	if retErr == nil {
		return res[0].Interface().(proto.Message), nil
	}

	err := retErr.(error)
	if err != nil {
		return nil, errorx.New(err.Error())
	}

	return res[0].Interface().(proto.Message), nil
}
