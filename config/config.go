package config

import (
	"crypto/md5"
	"github.com/rfw141/i/context"
	"github.com/rfw141/i/routine"
	"github.com/rfw141/i/x/filex"
	"github.com/rfw141/i/x/jsonx"
	"sync"
	"time"

	"github.com/rfw141/i/log"
)

type IConfig interface {
	Load()
	Watch()
	Refresh()
}

type Config struct {
	Filename string `json:"-"`
	iConfig  IConfig
	content  []byte
	version  [16]byte
	mu       sync.RWMutex
}

func NewConfig(c *Config, iConfig IConfig) *Config {
	if c.Filename == "" {
		panic("(config) filename is empty")
	}

	log.Debugf("(config) config file: %v", c.Filename)
	if !filex.Exist(c.Filename) {
		log.Debugf("(config) create config file: %v", c.Filename)
		if err := filex.Write(c.Filename, "{}"); err != nil {
			panic(err)
		}
	}

	c.iConfig = iConfig
	routine.Go(nil, func(context context.IContext) error {
		c.iConfig.Watch()
		return nil
	})
	return c
}

func (c *Config) Refresh() {
	err := jsonx.Unmarshal(c.content, c.iConfig)
	if err != nil {
		log.Errorf("(config) unmarshal fail: %v", err)
		return
	}
	log.Infof("(config) refresh config: %s", string(c.content))
}

func (c *Config) Load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	content, err := filex.Read(c.Filename)
	if err != nil {
		return
	}
	version := md5.Sum(content)
	if version == c.version {
		return
	}
	c.version = version
	c.content = content
	c.iConfig.Refresh()
}

func (c *Config) Watch() {
	log.Debug("(config) default watcher start")
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for range ticker.C {
		c.iConfig.Load()
	}
}

func (c *Config) Store() error {
	var err error
	if c.content, err = jsonx.Marshal(c.iConfig); err != nil {
		return err
	}
	if err = filex.Write(c.Filename, string(c.content)); err != nil {
		return err
	}
	c.version = md5.Sum(c.content)
	return nil
}
