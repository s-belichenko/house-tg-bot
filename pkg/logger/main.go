package logger

import (
	"log"
	"os"
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

func InitLog(logStreamName string) *YandexLogger {
	logger := log.New(os.Stdout, "", 0) // Отключаем все флаги.

	return newYandexLogger(logStreamName, logger)
}
