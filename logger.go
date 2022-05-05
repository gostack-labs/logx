package logx

import "strings"

const (
	LEVEL_DEBUG LogLevel = iota - 1
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

type LogLevel int8

func (l LogLevel) Int() int {
	return int(l)
}

func (l LogLevel) String() string {
	switch l {
	case LEVEL_DEBUG:
		return "DEBUG"
	case LEVEL_INFO:
		return "INFO"
	case LEVEL_WARN:
		return "WARN"
	case LEVEL_ERROR:
		return "ERROR"
	default:
		return ""
	}
}

func ParseLogLevel(s string) LogLevel {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LEVEL_DEBUG
	case "INFO":
		return LEVEL_INFO
	case "WARN":
		return LEVEL_WARN
	case "ERROR":
		return LEVEL_ERROR
	case "OFF":
		return LEVEL_OFF
	}
	return LEVEL_INFO
}

type Logger interface {
	SetLevel(level string)

	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}
