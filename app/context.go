package app

import (
	"bytes"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/context"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
)

type Ctx interface {
	context.IContext
	GetInstant() interface{}
}

type RpcCtx struct {
	reqId string
}

func (ctx *RpcCtx) GetReqId() string {
	return ctx.reqId
}

func (ctx *RpcCtx) GetInstant() interface{} {
	return ctx
}

func NewRpcCtx() *RpcCtx {
	return &RpcCtx{
		reqId: log.NewTraceId(),
	}
}

type CmdCtx struct {
	reqId string
	req   string
}

func (ctx *CmdCtx) GetReqId() string {
	return ctx.reqId
}

func (ctx *CmdCtx) GetInstant() interface{} {
	return ctx
}

func (ctx *CmdCtx) MustParse(req proto.Message) proto.Message {
	if ctx.req == "" {
		return req
	}
	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	errorx.Must(unmarshaler.Unmarshal(bytes.NewReader([]byte(ctx.req)), req))
	return req
}

func NewCmdCtx() *CmdCtx {
	return &CmdCtx{
		reqId: log.NewTraceId(),
	}
}
