package app

import (
	"net/http"

	"github.com/rfw141/i/log"
)

func registerHandler(path string, handle func(w http.ResponseWriter, r *http.Request)) {
	log.Debugf("(handle) register path: %s", path)
	http.HandleFunc(path, handle)
}

func registerDefHandlers(name string) {
	registerHandler("/"+name+"/ping", ping)
}

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Errorf("err:%v", err)
	}
}
