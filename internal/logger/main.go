package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogLevel string

type YandexLoggerInterface interface {
	Trace(message string, context LogContext)
	Debug(message string, context LogContext)
	Info(message string, context LogContext)
	Warning(message string, context LogContext)
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

type Logger struct {
	stream string
	logger *log.Logger
}

func InitLog(logStreamName string) *Logger {
	return newLogger(logStreamName)
}

func newLogger(streamName string) *Logger {
	logger := log.New(os.Stdout, "", 0) // Отключаем все флаги
	return &Logger{
		stream: streamName,
		logger: logger,
	}
}

func (l *Logger) createEntry(level LogLevel, message string, context LogContext) *Record {
	return &Record{
		Message:   message,
		Level:     level,
		Stream:    l.stream,
		Timestamp: time.Now(),
		Context:   context,
	}
}

func (l *Logger) write(entry *Record) {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Println("Ошибка при маршалинге JSON:", err)
		return
	}

	l.logger.Output(2, string(jsonBytes))
}

func (l *Logger) Trace(message string, context LogContext) {
	entry := l.createEntry(TRACE, message, context)
	l.write(entry)
}
func (l *Logger) Debug(message string, context LogContext) {
	entry := l.createEntry(DEBUG, message, context)
	l.write(entry)
}

func (l *Logger) Info(message string, context LogContext) {
	entry := l.createEntry(INFO, message, context)
	l.write(entry)
}

func (l *Logger) Warn(message string, context LogContext) {
	entry := l.createEntry(WARN, message, context)
	l.write(entry)
}

func (l *Logger) Error(message string, context LogContext) {
	entry := l.createEntry(ERROR, message, context)
	l.write(entry)
}

func (l *Logger) Fatal(message string, context LogContext) {
	entry := l.createEntry(FATAL, message, context)
	l.write(entry)
}
