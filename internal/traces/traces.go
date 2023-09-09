package traces

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type traceContextKey string

const (
	TraceIDKey traceContextKey = "trace.id"
)

func WithTraceID(ctx context.Context) context.Context {
	id, err := uuid.NewRandom()
	if err != nil {
		logrus.Warn(errors.Wrap(err, "failed to generate trace id"))
		return ctx
	}
	return context.WithValue(ctx, TraceIDKey, id.String())
}

func Logger(ctx context.Context) *logrus.Entry {
	return logrus.WithField("trace.id", TraceID(ctx))
}

// Creates a new instance of the logger which only logs WARNING events.
// We use it as logger for the GRPC events which are hardcoded as INFO level and pollute the logs.
func WarnLogger(ctx context.Context) *logrus.Entry {
	warnLevel, _ := logrus.ParseLevel("warn")
	warnLogger := logrus.New();
	warnLogger.SetLevel(warnLevel)
	return warnLogger.WithField("trace.id", TraceID(ctx))
}

func TraceID(ctx context.Context) string {
	if id, ok := ctx.Value(TraceIDKey).(string); ok {
		return id
	}
	return "<no-trace-id>"
}
