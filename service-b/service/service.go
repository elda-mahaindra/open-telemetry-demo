package service

import (
	"service-b/store"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	logger *logrus.Logger
	tracer trace.Tracer

	store *store.Store
}

func NewService(
	logger *logrus.Logger,
	tracer trace.Tracer,
	store *store.Store,
) *Service {
	return &Service{
		logger: logger,
		tracer: tracer,

		store: store,
	}
}
