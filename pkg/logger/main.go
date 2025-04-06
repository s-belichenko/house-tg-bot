package logger

import (
	"time"
)

type LogLevel string

type Logger interface {
	Trace(message string, context LogContext)
	Debug(message string, context LogContext)
	Info(message string, context LogContext)
	Warn(message string, context LogContext)
	Error(message string, context LogContext)
	Fatal(message string, context LogContext)
}

const (
	TRACE LogLevel = "TRACE"
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogContext map[string]interface{}

type Record struct {
	Message   string     `json:"message"`
	Level     LogLevel   `json:"level"`
	Stream    string     `json:"stream_name"`
	Timestamp time.Time  `json:"timestamp"`
	Context   LogContext `json:"extra"`
}

func InitLog(logStreamName string) Logger {
	return newYandexLogger(logStreamName)
}
