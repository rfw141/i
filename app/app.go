package app

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/config"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"reflect"
)

var defApp *App

type RpcDoc struct {
	Path string
	Func interface{}
}

type Rpc struct {
	Path string
	Req  reflect.Type
	Rsp  reflect.Type
}

func (r *Rpc) Example() {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: true,
		Indent:       "",
		AnyResolver:  nil,
	}
	reqV := newDefaultInstance(r.Req)
	rspV := newDefaultInstance(r.Rsp)
	req, err := marshaler.MarshalToString(reqV.Interface().(proto.Message))
	errorx.Must(err)
	rsp, err := marshaler.MarshalToString(rspV.Interface().(proto.Message))
	errorx.Must(err)
	log.Infof("path: %s", r.Path)
	log.Infof("req: %s", req)
	log.Infof("rsp: %s", rsp)
}

type App struct {
	addrs *Addrs
	svr   *Svr
	cli   *Cli
	cmd   *Cmd
}

func NewApp(cfg Cfg, rpcs []RpcDoc) *App {
	var a App
	a.addrs = NewAddrs(&config.Config{
		Filename: cfg.AddrsFile,
	})

	m := make(map[string]*Rpc)
	for _, v := range rpcs {
		rpc, err := parseRpc(v)
		errorx.Must(err)
		m[rpc.Path] = rpc
	}

	a.cli = NewCli(&cfg, a.addrs)
	a.svr = NewSvr(&cfg, a.addrs, m)
	a.cmd = NewCmd(&cfg, m)

	return &a
}

func (a *App) Run(ctx Ctx) error {
	return a.svr.Run(ctx)
}

func SetDefApp(app *App) {
	defApp = app
}

func SvrRun(ctx Ctx, cfg Cfg, handlers []Handler, rpcs []RpcDoc) error {
	defApp = NewApp(cfg, rpcs)
	defApp.svr.RegisterHandlers(handlers)
	return defApp.Run(ctx)
}

func CliCall(ctx Ctx, service, path string, req, rsp proto.Message) error {
	return defApp.cli.Call(ctx, service, path, req, rsp)
}

func CmdExecute(ctx Ctx, cfg Cfg, executors []Executor, rpcs []RpcDoc) error {
	defApp = NewApp(cfg, rpcs)
	defApp.cmd.RegisterExecutors(executors)
	return defApp.cmd.Execute(ctx)
}

func parseRpc(doc RpcDoc) (*Rpc, error) {
	t := reflect.TypeOf(doc.Func)
	if t.Kind() != reflect.Func {
		return nil, errorx.New("rpc must be func")
	}
	if t.NumIn() != 2 {
		return nil, errorx.New("rpc func must have 2 args")
	}
	ctx := t.In(0)
	if !ctx.Implements(reflect.TypeOf((*Ctx)(nil)).Elem()) {
		return nil, errorx.New("rpc func first arg must be Ctx")
	}
	req := t.In(1)
	if !req.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return nil, errorx.New("rpc func second arg must be proto.Message")
	}
	if t.NumOut() != 2 {
		return nil, errorx.New("rpc func must have 2 results")
	}
	rsp := t.Out(0)
	if !rsp.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return nil, errorx.New("rpc func first result must be proto.Message")
	}
	err := t.Out(1)
	if !err.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, errorx.New("rpc func second result must be error")
	}
	return &Rpc{
		Path: doc.Path,
		Req:  req,
		Rsp:  rsp,
	}, nil
}

func newDefaultInstance(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		return newDefaultInstance(t).Addr()
	case reflect.Slice:
		eleType := t.Elem()
		s := reflect.MakeSlice(t, 0, 0)
		s = reflect.Append(s, newDefaultInstance(eleType))
		return s
	case reflect.Map:
		m := reflect.MakeMap(t)
		keyType := t.Key()
		valType := t.Elem()
		m.SetMapIndex(newDefaultInstance(keyType), newDefaultInstance(valType))
		return m
	case reflect.Struct:
		v := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			v.Field(i).Set(newDefaultInstance(t.Field(i).Type))
		}
		return v
	default:
		return reflect.Zero(t)
	}
}
