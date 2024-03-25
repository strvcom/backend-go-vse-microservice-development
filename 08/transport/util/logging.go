package util

import (
	"context"
	"log/slog"
	"sync"

	netx "go.strv.io/net"
	logx "go.strv.io/net/logger"
)

const (
	requestIDFieldName = "request_id"
)

var (
	// once ensures that no one will override server log level.
	once = &sync.Once{}

	// serverLogLevel is by default info log level.
	serverLogLevel = slog.LevelInfo
)

// SetServerLogLevel sets log level that will be used for all new server loggers.
// Only first call of this function is valid. Other calls are ignored to prevent unwanted behavior.
func SetServerLogLevel(l slog.Level) {
	once.Do(func() {
		slog.SetLogLoggerLevel(serverLogLevel)
	})
}

// ServerLogger is a logger that is used by go.strv.io/net/http.Server to log errors and debug messages.
type ServerLogger struct {
	*slog.Logger
}

func NewServerLogger(caller string) ServerLogger {
	return ServerLogger{slog.Default()}
}

// With is a wrapper around zap.With using zap.Any as a field type.
func (l ServerLogger) With(fields ...logx.Field) logx.ServerLogger {
	f := make([]any, 0, len(fields))
	for _, field := range fields {
		f = append(f, slog.Any(field.Key, field.Value))
	}
	return ServerLogger{Logger: l.Logger.With(f...)}
}

func (l ServerLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l ServerLogger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l ServerLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l ServerLogger) Error(msg string, err error) {
	l.Logger.Error(msg, slog.Any("error", err))
}

// WithCtx returns logger with fields extracted from context.
func WithCtx(ctx context.Context, l *slog.Logger) *slog.Logger {
	return l.With(
		slog.String(requestIDFieldName, netx.RequestIDFromCtx(ctx)),
	)
}
