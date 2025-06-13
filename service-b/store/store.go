package store

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Store struct {
	logger *logrus.Logger
	tracer trace.Tracer
}

func NewStore(
	logger *logrus.Logger,
	tracer trace.Tracer,
) *Store {
	return &Store{
		logger: logger,
		tracer: tracer,
	}
}
