package app

import (
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/filex"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	gEnv *env
)

type env struct {
	prefix string
	v      *viper.Viper
}

const (
	EnvModule    = "MODULE"
	EnvWorkPath  = "WORK_PATH"
	EnvProtoc    = "PROTOC"
	EnvProtoPath = "PROTO_PATH"
	EnvTplPath   = "TPL_PATH"
)

func findEnv(base string) []string {
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	var ls []string
	var envFile string
	for i := 0; i < 5; i++ {
		envFile = filepath.ToSlash(filepath.Join(base, ".env"))
		if filex.Exist(envFile) {
			ls = append(ls, envFile)
		}
		envFile = filepath.ToSlash(filepath.Join(base, "env"))
		if filex.Exist(envFile) {
			ls = append(ls, envFile)
		}
		base = filepath.Join(base, "..")
	}
	return ls
}

func LoadEnv(prefix string) {
	gEnv = &env{
		prefix: prefix,
		v:      viper.New(),
	}
	gEnv.v.SetEnvPrefix(prefix)
	envFiles := findEnv(".")
	for _, envFile := range envFiles {
		gEnv.v.SetConfigType("env")
		gEnv.v.SetConfigName(filepath.Base(envFile))
		gEnv.v.AddConfigPath(filepath.Dir(envFile))
		break
	}
	err := gEnv.v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	for k, v := range gEnv.v.AllSettings() {
		log.Debugf("(env) %s=%s", k, v)
	}
}

func Env(key string) string {
	return gEnv.v.GetString(strings.ToUpper(gEnv.prefix) + "_" + key)
}

func Envs() map[string]string {
	all := gEnv.v.AllSettings()
	m := map[string]string{}
	for k, v := range all {
		m[strings.ToUpper(k)] = v.(string)
	}
	return m
}
