package log

import "fmt"

type Logger interface {
	Log(Level, ...interface{})
}

type Level uint32

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

type nilLogger struct{}

func (nilLogger) Log(Level, ...interface{}) {
	// do nothing
}

var defaultLogger Logger = nilLogger{}

func SetLogger(l Logger) {
	if l == nil {
		defaultLogger = nilLogger{}
	}
	defaultLogger = l
}

func Log(l Level, args ...interface{}) {
	defaultLogger.Log(l, args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Log(DebugLevel, args...)
}

func Debugf(format string, args ...interface{}) {
	defaultLogger.Log(DebugLevel, fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	defaultLogger.Log(InfoLevel, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Log(InfoLevel, fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	defaultLogger.Log(ErrorLevel, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Log(ErrorLevel, fmt.Sprintf(format, args...))
}

func Fatal(args ...interface{}) {
	defaultLogger.Log(FatalLevel, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLogger.Log(FatalLevel, fmt.Sprintf(format, args...))
}
