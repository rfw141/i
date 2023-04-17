package app

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/rfw141/i/context"
	"github.com/rfw141/i/routine"

	"github.com/rfw141/i/x/jsonx"

	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
)

const (
	applicationJson      = "application/json"
	applicationXProtobuf = "application/x-protobuf"
)

const (
	HeaderErrCode = "X-ErrCode"
	HeaderErrMsg  = "X-ErrMsg"
)

type Handler struct {
	Path        string
	Handle      interface{}
	reqType     reflect.Type
	rspType     reflect.Type
	handleValue reflect.Value
}

type Svr struct {
	*Cfg
	*Addrs
	handlers map[string]Handler
	rpcs  map[string]*Rpc
}

func NewSvr(cfg *Cfg, addrs *Addrs, rpcs map[string]*Rpc) *Svr {
	var s Svr
	s.Cfg = cfg
	s.Addrs = addrs
	s.handlers = make(map[string]Handler)
	s.rpcs = rpcs
	return &s
}

func (s *Svr) RegisterHandlers(handlers []Handler) {
	registerDefHandlers(s.Name)
	for _, handler := range handlers {
		nHandler := handler
		nHandler.handleValue = reflect.ValueOf(nHandler.Handle)
		nHandler.reqType = nHandler.handleValue.Type().In(1).Elem()
		registerHandler(nHandler.Path, func(w http.ResponseWriter, r *http.Request) {
			doHandler(NewRpcCtx(), w, r, nHandler)
		})
		s.handlers[handler.Path] = nHandler
	}
}

func (s *Svr) Run(ctx Ctx) error {
	done := make(chan struct{})
	addr := NewAddr(s.Name, s.Addr)
	defer func() {
		if s.Addrs != nil {
			if err := s.Addrs.Remove(addr); err != nil {
				log.Errorf("err:%v", err)
			}
		}

		time.Sleep(100 * time.Millisecond)
		close(done)
	}()

	log.Debug("(svr) register signal handler")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	routine.Go(ctx, func(context context.IContext) error {
		switch <-sigs {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Debug("(svr) get signal SIGTERM, prepare exit")
			done <- struct{}{}
			return nil
		}
		return nil
	})

	log.Debug("(svr) run server")
	routine.Go(ctx, func(context context.IContext) error {
		var err error
		var ln net.Listener
		err = errorx.Call(func() error {
			if ln, err = net.Listen(addr.GetNetwork(), addr.GetAddress()); err != nil {
				return err
			}
			addr.SetAddr(ln.Addr().String())
			return nil
		}, func() error {
			log.Infof("(svr) %s run, addr: %v", s.Name, addr.GetAddr())
			if s.Addrs != nil {
				if err := s.Addrs.Add(addr); err != nil {
					return err
				}
			}
			return nil
		}, func() error {
			return http.Serve(ln, nil)
		})
		if err != nil {
			done <- struct{}{}
			return err
		}
		return nil
	})

	<-done

	return nil
}

func doHandler(ctx Ctx, w http.ResponseWriter, r *http.Request, handler Handler) {
	var handleBody, rspBody, reqBody []byte
	var contentType = r.Header.Get("content-type")
	var err error

	defer func() {
		// log.Infof("%v rsp:%v", strings.ReplaceAll(string(c.Request.Header.Header()), "\r\n", " "), string(jsonBody))
		w.Header().Set("Content-Type", contentType)
		if err != nil {
			log.Logf(log.ERROR, &log.Option{NoTrace: true}, "(svr) rpc err:%v", err)
			if e, ok := err.(errorx.IErr); ok {
				w.Header().Set(HeaderErrCode, fmt.Sprintf("%d", e.GetCode()))
				w.Header().Set(HeaderErrMsg, e.GetMsg())
			} else {
				w.Header().Set(HeaderErrCode, "-1")
				w.Header().Set(HeaderErrMsg, err.Error())
			}
			return
		}
		_, err = w.Write(rspBody)
		if err != nil {
			log.Errorf("(svr) err:%v", err)
		}
	}()

	req := reflect.New(handler.reqType)
	if reqBody, err = io.ReadAll(r.Body); err != nil {
		return
	}
	switch {
	case strings.Contains(contentType, applicationJson):
		unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
		if err = unmarshaler.Unmarshal(bytes.NewReader(reqBody), req.Interface().(proto.Message)); err != nil {
			return
		}
	case strings.Contains(contentType, applicationXProtobuf):
		if err = proto.Unmarshal(reqBody, req.Interface().(proto.Message)); err != nil {
			return
		}
	default:
		err = fmt.Errorf("(svr) invalid content-type:%s", contentType)
		return
	}

	res := handler.handleValue.Call([]reflect.Value{reflect.ValueOf(ctx), req})
	if !res[1].IsNil() {
		err = res[1].Interface().(error)
		return
	}
	var rsp proto.Message
	if !res[0].IsNil() {
		rsp = res[0].Interface().(proto.Message)
		if handleBody, err = proto.Marshal(rsp); err != nil {
			return
		}
	}

	switch {
	case strings.Contains(contentType, applicationJson):
		if rspBody, err = jsonx.Marshal(rsp); err != nil {
			return
		}
	case strings.Contains(contentType, applicationXProtobuf):
		rspBody = handleBody
	}
}
