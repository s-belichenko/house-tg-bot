package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type YandexLogger struct {
	stream string
	logger *log.Logger
}

func newYandexLogger(streamName string) Logger {
	logger := log.New(os.Stdout, "", 0) // Отключаем все флаги
	return &YandexLogger{
		stream: streamName,
		logger: logger,
	}
}

func (l *YandexLogger) createEntry(level LogLevel, message string, context LogContext) *Record {
	return &Record{
		Message:   message,
		Level:     level,
		Stream:    l.stream,
		Timestamp: time.Now(),
		Context:   context,
	}
}

func (l *YandexLogger) write(entry *Record) {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		log.Println("Ошибка при маршалинге JSON:", err)
		return
	}

	_ = l.logger.Output(2, string(jsonBytes))
}

func (l *YandexLogger) Error(message string, context LogContext) {
	entry := l.createEntry(ERROR, message, context)
	l.write(entry)
}

func (l *YandexLogger) Fatal(message string, context LogContext) {
	entry := l.createEntry(FATAL, message, context)
	l.write(entry)
}

func (l *YandexLogger) Info(message string, context LogContext) {
	entry := l.createEntry(INFO, message, context)
	l.write(entry)
}

func (l *YandexLogger) Warn(message string, context LogContext) {
	entry := l.createEntry(WARN, message, context)
	l.write(entry)
}

func (l *YandexLogger) Trace(message string, context LogContext) {
	entry := l.createEntry(TRACE, message, context)
	l.write(entry)
}

func (l *YandexLogger) Debug(message string, context LogContext) {
	entry := l.createEntry(DEBUG, message, context)
	l.write(entry)
}
