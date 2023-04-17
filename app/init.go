package app

import "github.com/rfw141/i/log"

func Init(project, service string, isCmd bool) {
	if isCmd {
		log.SetName("")
		log.EnableStdout()
	} else {
		log.SetName(service)
		log.EnableStdout()
		log.EnableTrace()
		log.SetStackLevel(log.ERROR)
		log.SetLevel(log.DEBUG)
		log.SetSourceLevel(log.ERROR)
	}
	LoadEnv(project)
}
