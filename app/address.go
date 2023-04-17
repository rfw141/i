package app

import (
	"strings"

	"github.com/rfw141/i/config"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/slicex"
)

type Addr struct {
	Name string   `json:"name"`
	Addr []string `json:"addr"`
}

func NewAddr(name, addr string) *Addr {
	var a Addr
	a.Name = name
	if addr == "" {
		addr = ":0"
	}
	a.Addr = []string{addr}
	return &a
}

func (a *Addr) GetAddr() string {
	if len(a.Addr) == 0 {
		return ""
	}
	return a.Addr[0]
}

func (a *Addr) SetAddr(addr string) {
	a.Addr = []string{addr}
}

func (a *Addr) GetAddress() string {
	address := a.GetAddr()
	if i := strings.Index(address, "://"); i > 0 {
		return address[i+3:]
	}
	return address
}

func (a *Addr) GetNetwork() string {
	address := a.GetAddr()
	if i := strings.Index(address, "://"); i > 0 {
		return address[:i]
	}
	return "tcp"
}

func (a *Addr) BuildPaths(path string) []string {
	var ls []string
	for _, v := range a.Addr {
		var addr = v
		if !strings.HasPrefix(v, "://") {
			addr = "http://" + addr
		}
		ls = append(ls, addr+path)
	}
	return ls
}

type Addrs struct {
	*config.Config
	Services map[string]*Addr `json:"services"`
}

func NewAddrs(c *config.Config) *Addrs {
	var a Addrs
	a.Services = make(map[string]*Addr)
	a.Config = config.NewConfig(c, &a)
	a.Load()
	return &a
}

func (a *Addrs) Add(addr *Addr) error {
	log.Debugf("(addrs) add addr: %v", addr)
	if _, ok := a.Services[addr.Name]; !ok {
		a.Services[addr.Name] = addr
	} else {
		a.Services[addr.Name].Addr = slicex.StrUniq(append(a.Services[addr.Name].Addr, addr.Addr...))
	}
	return a.Store()
}

func (a *Addrs) Remove(addr *Addr) error {
	if _, ok := a.Services[addr.Name]; !ok {
		return nil
	}
	log.Debugf("(addrs) remove addr: %v", addr)
	for _, v := range addr.Addr {
		a.Services[addr.Name].Addr = slicex.StrRemove(a.Services[addr.Name].Addr, v)
	}
	if len(a.Services[addr.Name].Addr) == 0 {
		delete(a.Services, addr.Name)
	}
	return a.Store()
}
