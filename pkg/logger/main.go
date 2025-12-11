package logger

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

type SystemLogger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Output(calldepth int, s string) error
}
