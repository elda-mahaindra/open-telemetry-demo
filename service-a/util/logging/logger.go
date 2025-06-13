package logging

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// LogWithTrace returns a logrus.Entry enriched with log_id
func LogWithTrace(ctx context.Context, logger *logrus.Logger) *logrus.Entry {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return logger.WithField("[log_id]", "unknown")
	}

	return logger.WithFields(logrus.Fields{
		"[log_id]": sc.TraceID().String(),
	})
}
