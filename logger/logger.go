package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
)

var contextKey struct{} // context key

type Logger interface {
	Logger() *logrus.Entry
	SetLogger(logger *logrus.Entry)
}

// New creates a new loggger
func New(level string, fields ...logrus.Fields) *logrus.Entry {
	logger := logrus.New()
	entry := logrus.NewEntry(logger)
	for _, f := range fields {
		entry = entry.WithFields(f)
	}

	askedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		askedLevel = logrus.TraceLevel
		entry.WithField("level", level).Warn("bad log level")
	}
	logger.SetLevel(askedLevel)

	return entry
}

// WithLogger adds the logger as context value
func WithLogger(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, contextKey, entry)
}

// FromContext returns the logger set with WithLogger or nil
func FromContext(ctx context.Context, noopLogger bool) *logrus.Entry {
	entry, ok := ctx.Value(contextKey).(*logrus.Entry)
	if !ok {
		if noopLogger {
			l := logrus.New()
			l.SetOutput(io.Discard)
			return l.WithField("no-op-logger", true)
		}
		return nil
	}
	return entry
}
