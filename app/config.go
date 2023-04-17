package app

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"github.com/rfw141/i/x/filex"
	"github.com/rfw141/i/x/jsonx"
	"github.com/spf13/viper"
	"path/filepath"
)

var ErrConfigNotFound = errorx.New("config not found")

type Cfg struct {
	Name      string
	Addr      string
	Path      string
	AddrsFile string

	v *viper.Viper
}

func NewCfg(cfg Cfg) (*Cfg, error) {
	if cfg.Name == "" {
		panic("cfg.Name is empty")
	}
	if cfg.Path == "" {
		cfg.Path = filepath.ToSlash(filepath.Join(Env(EnvWorkPath), "config"))
	}
	if err := filex.MustExist(cfg.Path); err != nil {
		return nil, err
	}
	cfg.v = viper.New()
	cfg.v.AddConfigPath(cfg.Path)
	cfg.v.SetConfigType("toml")
	cfgFile := filepath.ToSlash(filepath.Join(cfg.Path, ".config.toml"))
	if filex.Exist(cfgFile) {
		log.Debugf("(cfg) read config from %s", cfgFile)
		cfg.v.SetConfigName(".config.toml")
	} else {
		cfgFile = filepath.ToSlash(filepath.Join(cfg.Path, "config.toml"))
		if err := filex.MustExist(cfgFile); err != nil {
			return nil, err
		}
		log.Debugf("(cfg) read config from %s", cfgFile)
		cfg.v.SetConfigName("config.toml")
	}
	if err := cfg.v.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg.v.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("(cfg) config file changed: %s", e.String())
	})
	cfg.v.WatchConfig()
	if cfg.AddrsFile == "" {
		cfg.AddrsFile = cfg.GetStr("addrs_file", "")
	}
	if cfg.Addr == "" {
		cfg.Addr = cfg.GetStr("service."+cfg.Name+".addr", "")
	}
	if cfg.Addr == "" {
		cfg.Addr = "127.0.0.1:0"
	}
	return &cfg, nil
}

func (c *Cfg) Get(key string, def interface{}) interface{} {
	if !c.v.IsSet(key) {
		log.Warnf("%s not found, use default val: %v", key, def)
		return def
	}
	return c.v.GetString(key)
}

func (c *Cfg) GetStr(key string, def string) string {
	if !c.v.IsSet(key) {
		log.Warnf("%s not found, use default val: %v", key, def)
		return def
	}
	return c.v.GetString(key)
}

func (c *Cfg) GetJsonData(key string, out interface{}) error {
	if !c.v.IsSet(key) {
		log.Warnf("%s not found", key)
		return ErrConfigNotFound
	}

	valByte, err := jsonx.Marshal(c.v.Get(key))
	if err != nil {
		return err
	}

	err = jsonx.Unmarshal(valByte, out)
	if err != nil {
		return err
	}

	return nil
}
