// copied from https://github.com/hibiken/asynq/blob/master/internal/log/log.go

package log

import (
	"conf-reload/internal/app"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"sync"
)

// Base supports logging at various log levels.
type Base interface {
	// Debug logs a message at Debug level.
	Debug(args ...interface{})

	// Info logs a message at Info level.
	Info(args ...interface{})

	// Warn logs a message at Warning level.
	Warn(args ...interface{})

	// Error logs a message at Error level.
	Error(args ...interface{})

	// Fatal logs a message at Fatal level
	// and process will exit with status set to 1.
	Fatal(args ...interface{})
}

type baseLogger struct {
	*stdlog.Logger
}

func (l *baseLogger) Debug(args ...interface{}) {
	l.prefixPrint("DEBUG: ", args...)
}

func (l *baseLogger) Info(args ...interface{}) {
	l.prefixPrint("INFO: ", args...)
}

func (l *baseLogger) Warn(args ...interface{}) {
	l.prefixPrint("WARN: ", args...)
}

func (l *baseLogger) Error(args ...interface{}) {
	l.prefixPrint("ERROR: ", args...)
}

func (l *baseLogger) Fatal(args ...interface{}) {
	l.prefixPrint("FATAL: ", args...)
	os.Exit(1)
}

func (l *baseLogger) prefixPrint(prefix string, args ...interface{}) {
	args = append([]interface{}{prefix}, args...)
	l.Print(args...)
}

func newBase(out io.Writer) *baseLogger {
	prefix := fmt.Sprintf("%s@%s: pid=%d ", app.PackageName, app.Version, os.Getpid())
	return &baseLogger{
		stdlog.New(out, prefix, stdlog.Ldate|stdlog.Ltime|stdlog.Lmicroseconds|stdlog.LUTC),
	}
}

type Logger struct {
	base  Base
	mu    sync.Mutex
	level Level
}

func NewLogger(base Base) *Logger {
	if base == nil {
		base = newBase(os.Stderr)
	}
	return &Logger{base: base, level: DebugLevel}
}

type Level int32

const (
	DebugLevel Level = iota

	InfoLevel

	WarnLevel

	ErrorLevel

	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warning"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l *Logger) canLogAt(v Level) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return v >= l.level
}

func (l *Logger) Debug(args ...interface{}) {
	if !l.canLogAt(DebugLevel) {
		return
	}
	l.base.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	if !l.canLogAt(InfoLevel) {
		return
	}
	l.base.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	if !l.canLogAt(WarnLevel) {
		return
	}
	l.base.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	if !l.canLogAt(ErrorLevel) {
		return
	}
	l.base.Error(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	if !l.canLogAt(FatalLevel) {
		return
	}
	l.base.Fatal(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// SetLevel sets the logger level.
// It panics if v is less than DebugLevel or greater than FatalLevel.
func (l *Logger) SetLevel(v Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v < DebugLevel || v > FatalLevel {
		panic("log: invalid log level")
	}
	l.level = v
}
