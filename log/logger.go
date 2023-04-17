package log

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

func linkf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
	//return fmt.Sprintf("\u001B]8;;%s\u0007\u001B[35m%s\u001B]8;;\u0007\u001B[0m", fmt.Sprintf(format, args...), fmt.Sprintf(format, args...))
}

func log(level Level, option *Option, format string, params ...interface{}) {
	if option == nil {
		option = &Option{}
	}

	write(level, option, format, params...)
	stack(level, option, 5)

	if level == FATAL {
		os.Exit(1)
	} else if level == PANIC {
		panic("")
	}
}

func write(level Level, option *Option, format string, args ...interface{}) {
	if level < defaultLevel {
		return
	}

	var msg string
	if format == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(format, args...)
	}

	msg = string(formatter(level, msg, option))
	if allowStdout {
		//fmt.Printf("%s", msg)
		_, err := os.Stdout.Write([]byte(msg))
		if err != nil {
			return
		}
		err = os.Stdout.Sync()
		if err != nil {
			return
		}
	}

	if writer != nil {
		_, err := writer.Write([]byte(msg))
		if err != nil {
			return
		}
	}
}

func stack(level Level, option *Option, skip int) {
	if level < stackLevel || option.NoSource {
		return
	}
	option.NoSource = true
	var cnt, max = 0, 10
	for ; ; skip++ {
		if cnt > max {
			break
		}
		cnt += 1
		callerName := ""
		pc, callerFile, callerLine, ok := runtime.Caller(skip)
		if ok {
			callerName = runtime.FuncForPC(pc).Name()
		}
		if callerName == "runtime.goexit" || callerName == "" {
			break
		}
		if strings.HasPrefix(callerName, "reflect.Value.") {
			continue
		}
		filePath, fileFunc := getPackageName(callerName)
		filePath = path.Join(filePath, path.Base(callerFile))
		if fileFunc == "" {
			write(level, option, "#STACK "+linkf("%s:%d", filePath, callerLine))
		} else {
			write(level, option, "#STACK "+linkf("%s:%d:%s", filePath, callerLine, fileFunc))
		}
	}
}

func getPackageName(f string) (string, string) {
	slashIndex := strings.LastIndex(f, "/")
	if slashIndex > 0 {
		idx := strings.Index(f[slashIndex:], ".") + slashIndex
		return f[:idx], f[idx+1:]
	}
	return f, ""
}

func source(skip int) string {
	var callerName, callerFile string
	var callerLine int
	var ok bool
	var pc uintptr
	pc, callerFile, callerLine, ok = runtime.Caller(skip)
	callerName = ""
	if ok {
		callerName = runtime.FuncForPC(pc).Name()
	}
	filePath, fileFunc := getPackageName(callerName)
	if fileFunc == "" {
		return linkf("%s:%d", path.Join(filePath, path.Base(callerFile)), callerLine)
	}
	return linkf("%s:%d:%s", path.Join(filePath, path.Base(callerFile)), callerLine, fileFunc)
}

func formatter(level Level, msg string, option *Option) []byte {
	var ss string
	if !option.NoSource {
		if level >= sourceLevel || option.SourceString != "" {
			if option.SourceString != "" {
				ss = option.SourceString
			} else {
				ss = source(option.AddSourceSkip + defaultSourceSkip)
			}
		}
	}

	var ls = level.ColorShortString()

	now := time.Now()

	var traceId string
	var logName string
	if allowTrace {
		gid := GetGoId()
		s := getState(gid)
		traceId = s.traceId
		logName = fmt.Sprintf("%s.%s.%d", name, pid, gid)
	} else {
		logName = fmt.Sprintf("%s.%s", name, pid)
	}
	if name == "" {
		logName = ""
	}

	items := []string{
		ls, now.Format("06-01-02T15:04:05.0000"),
	}
	if traceId != "" {
		items = append(items, traceId)
	}
	if logName != "" {
		items = append(items, logName)
	}
	items = append(items, level.ColorString(msg))
	if ss != "" {
		items = append(items, ss)
	}
	items = append(items, "\n")

	var b bytes.Buffer
	b.WriteString(strings.Join(items, " "))

	return b.Bytes()
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	strRandSrc = rand.NewSource(time.Now().UnixNano())
)

func strRandom(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, strRandSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = strRandSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return sb.String()
}
