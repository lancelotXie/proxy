package logger

import (
	"context"
	"io"
)

var defaultlogger = NewLogger()

// Set : 设置日志输出级别
func Set(level LogLevel) { defaultlogger.Set(level) }

// Info :
func Info(ctx context.Context, m ...interface{}) { defaultlogger.Info(ctx, m...) }

// Warning :
func Warning(ctx context.Context, m ...interface{}) { defaultlogger.Warning(ctx, m...) }

// Error :
func Error(ctx context.Context, m ...interface{}) { defaultlogger.Error(ctx, m...) }

// Logger :
type Logger interface {
	Set(LogLevel)
	Info(ctx context.Context, m ...interface{})
	Warning(ctx context.Context, m ...interface{})
	Error(ctx context.Context, m ...interface{})
	SetOutput(w io.Writer)
}

// NewLogger :
func NewLogger() Logger {
	return &logger{val: 0}
}

// LogLevel :
type LogLevel int

// INFO :
const (
	_LogLevel = iota
	INFO
	WARN
	ERROR
	NOLOG
)
