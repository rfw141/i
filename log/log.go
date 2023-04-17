package log

import (
	"github.com/petermattis/goid"
	"strings"
)

func Log(level Level, option *Option, params ...interface{}) { log(level, option, "", params...) }
func Logf(level Level, option *Option, format string, params ...interface{}) {
	log(level, option, format, params...)
}
func Debug(params ...interface{})                 { log(DEBUG, nil, "", params...) }
func Debugf(format string, params ...interface{}) { log(DEBUG, nil, format, params...) }
func Info(params ...interface{})                  { log(INFO, nil, "", params...) }
func Infof(format string, params ...interface{})  { log(INFO, nil, format, params...) }
func Print(params ...interface{})                 { log(INFO, nil, "", params...) }
func Printf(format string, params ...interface{}) { log(INFO, nil, format, params...) }
func Warn(params ...interface{})                  { log(WARN, nil, "", params...) }
func Warnf(format string, params ...interface{})  { log(WARN, nil, format, params...) }
func Error(params ...interface{})                 { log(ERROR, nil, "", params...) }
func Errorf(format string, params ...interface{}) { log(ERROR, nil, format, params...) }
func Panic(params ...interface{})                 { log(PANIC, nil, "", params...) }
func Panicf(format string, params ...interface{}) { log(PANIC, nil, format, params...) }
func Fatal(params ...interface{})                 { log(FATAL, nil, "", params...) }
func Fatalf(format string, params ...interface{}) { log(FATAL, nil, format, params...) }

func GetGoId() int64 { return goid.Get() }
func NewTraceId(traceIds ...string) string {
	if len(traceIds) > 0 {
		return strings.Join(traceIds, ".") + "." + strRandom(16)
	}
	return strRandom(16)
}
func DelTraceId(gid int64) {
	if !allowTrace {
		return
	}
	if gid == 0 {
		gid = GetGoId()
	}
	states.Delete(gid)
}
