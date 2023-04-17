package log

import (
	"io"
	"os"
	"strconv"
	"sync"
)

var (
	pid    string
	name   = "default"
	states sync.Map

	defaultLevel = INFO
	stackLevel   = PANIC
	sourceLevel  = ERROR

	writer io.Writer

	defaultSourceSkip = 5

	allowTrace  = false
	allowStdout = false
)

func init() {
	pid = strconv.Itoa(os.Getpid())
}

type state struct {
	traceId string
}

func getState(gids ...int64) *state {
	var gid int64
	if len(gids) == 0 {
		gid = GetGoId()
	} else {
		gid = gids[0]
	}
	v, ok := states.Load(gid)
	if !ok {
		s := state{
			traceId: NewTraceId(),
		}
		states.Store(gid, &s)
		return &s
	}
	return v.(*state)
}

func SetWriter(w io.Writer)      { writer = w }
func EnableStdout()              { allowStdout = true }
func SetName(n string)           { name = n }
func EnableTrace()               { allowTrace = true }
func SetLevel(level Level)       { defaultLevel = level }
func SetStackLevel(level Level)  { stackLevel = level }
func SetSourceLevel(level Level) { sourceLevel = level }
