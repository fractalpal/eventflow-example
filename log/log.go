package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Warn(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type ctxLoggerTag struct{}

var ctxLoggerKey = &ctxLoggerTag{}

type ctxLogger struct {
	logger *logrus.Entry
	fields logrus.Fields
}

func NewLogrus(lvl logrus.Level) *logrus.Logger {
	l := logrus.New()
	l.Level = lvl
	return l
}

func FromContext(ctx context.Context) *logrus.Entry {
	return FromContextWithFields(ctx, logrus.Fields{})
}

func FromContextWithFields(ctx context.Context, fields logrus.Fields) *logrus.Entry {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		// return default one
		lvl := logrus.DebugLevel
		return NewLogrus(lvl).WithField("logger", []string{"default", lvl.String()})
	}
	for k, v := range l.fields {
		fields[k] = v
	}
	return l.logger.WithFields(fields)
}

func AddFields(ctx context.Context, fields logrus.Fields) {
	l, ok := ctx.Value(ctxLoggerKey).(*ctxLogger)
	if !ok || l == nil {
		return
	}
	for k, v := range fields {
		l.fields[k] = v
	}
}

func ContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	l := &ctxLogger{
		logger: logger,
		fields: logrus.Fields{},
	}
	return context.WithValue(ctx, ctxLoggerKey, l)
}

// We wants to attached fields only if error occurred.
// This way we can ensure (in defer) that given scope will exit with proper log fields
func AddFieldsForErr(ctx context.Context, fields logrus.Fields, err error) {
	if err != nil {
		AddFields(ctx, fields)
	}
}
