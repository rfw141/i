package app

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Cli struct {
	*Cfg
	*Addrs
}

func NewCli(cfg *Cfg, addrs *Addrs) *Cli {
	var c Cli
	c.Cfg = cfg
	c.Addrs = addrs
	return &c
}

func (c *Cli) call(ctx Ctx, url string, req, rsp proto.Message) error {
	body, err := proto.Marshal(req)
	if err != nil {
		return errorx.New(err.Error())
	}
	client := http.Client{Timeout: time.Duration(1) * time.Second}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Add("Content-Type", applicationXProtobuf)
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return errors.New(err.Error())
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}(httpRsp.Body)
	if httpRsp.StatusCode != http.StatusOK {
		return errorx.New("status not ok, %d", httpRsp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return errorx.New(err.Error())
	}
	if err := proto.Unmarshal(bodyBytes, rsp); err != nil {
		return errorx.New(err.Error())
	}
	return nil
}

func (c *Cli) Call(ctx Ctx, service, path string, req, rsp proto.Message) error {
	if addrs, ok := c.Addrs.Services[service]; ok {
		size := len(addrs.Addr)
		for i, urlPath := range addrs.BuildPaths(path) {
			err := c.call(ctx, urlPath, req, rsp)
			if err != nil {
				if i == size-1 {
					return err
				}
				log.Errorf("err:%v", err)
				continue
			}
			return nil
		}
	}
	return errorx.New("service %s not found", service)
}
