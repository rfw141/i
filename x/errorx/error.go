package errorx

import (
	"fmt"
	"github.com/rfw141/i/log"
)

type IErr interface {
	GetCode() int32
	GetMsg() string
}

type E struct {
	Code int32
	Msg  string
}

func (e *E) GetCode() int32 {
	return e.Code
}
func (e *E) GetMsg() string {
	return e.Msg
}
func (e *E) Error() string {
	return fmt.Sprintf("code:%d, msg:%s", e.Code, e.Msg)
}

func GetMsg(errmsg ...interface{}) string {
	var msg string
	msgLen := len(errmsg)
	if msgLen == 0 {
		msg = ""
	} else if msgLen == 1 {
		msg = errmsg[0].(string)
	} else {
		var args []interface{}
		for _, v := range errmsg[1:] {
			args = append(args, v)
		}
		msg = fmt.Sprintf(errmsg[0].(string), args...)
	}
	return msg
}

func Error(code int32, msg ...interface{}) error {
	err := &E{Code: code, Msg: GetMsg(msg...)}
	log.Logf(log.ERROR, &log.Option{AddSourceSkip: 1}, "err:%v", err)
	return err
}

func New(msg ...interface{}) error {
	err := &E{Code: -1, Msg: GetMsg(msg...)}
	log.Logf(log.ERROR, &log.Option{AddSourceSkip: 1}, "err:%v", err)
	return err
}

func IsCode(err error, code int32) bool {
	if e, ok := err.(IErr); ok {
		return e.GetCode() == code
	}
	return false
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Exit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
