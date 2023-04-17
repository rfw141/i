package log

import "fmt"

type Level int8

const (
	DEBUG Level = iota - 1
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

func (l Level) String() string {
	switch l {
	case DEBUG:
		return "dbg"
	case INFO:
		return "inf"
	case WARN:
		return "war"
	case ERROR:
		return "err"
	case PANIC:
		return "pan"
	case FATAL:
		return "fat"
	default:
		return fmt.Sprintf("lvl(%d)", l)
	}
}

func (l Level) ShortString() string {
	switch l {
	case DEBUG:
		return "D"
	case INFO:
		return "I"
	case WARN:
		return "W"
	case ERROR:
		return "E"
	case PANIC:
		return "P"
	case FATAL:
		return "F"
	default:
		return fmt.Sprintf("%d", l)
	}
}

func (l Level) ColorShortString() string {
	switch l {
	case DEBUG:
		return "\u001B[106m \u001B[0m"
	case INFO:
		return "\u001B[102m \u001B[0m"
	case WARN:
		return "\u001B[103m \u001B[0m"
	case ERROR:
		return "\u001B[101m \u001B[0m"
	case PANIC:
		return "\u001B[105m \u001B[0m"
	case FATAL:
		return "\u001B[107m \u001B[0m"
	default:
		return "\u001B[104m \u001B[0m"
	}
}

func (l Level) ColorString(val string) string {
	switch l {
	case DEBUG:
		return fmt.Sprintf("\u001B[36m%s\u001B[0m", val)
	case INFO:
		return fmt.Sprintf("\u001B[32m%s\u001B[0m", val)
	case WARN:
		return fmt.Sprintf("\u001B[33m%s\u001B[0m", val)
	case ERROR:
		return fmt.Sprintf("\u001B[31m%s\u001B[0m", val)
	case PANIC:
		return fmt.Sprintf("\u001B[35m%s\u001B[0m", val)
	case FATAL:
		return fmt.Sprintf("\u001B[37m%s\u001B[0m", val)
	default:
		return fmt.Sprintf("\u001B[34m%s\u001B[0m", val)
	}
}
