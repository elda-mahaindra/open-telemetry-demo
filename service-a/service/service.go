package service

import (
	"service-a/adapter/service_b_adapter"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	logger *logrus.Logger
	tracer trace.Tracer

	serviceBAdapter *service_b_adapter.Adapter
}

func NewService(
	logger *logrus.Logger,
	tracer trace.Tracer,
	serviceBAdapter *service_b_adapter.Adapter,
) *Service {
	return &Service{
		logger: logger,
		tracer: tracer,

		serviceBAdapter: serviceBAdapter,
	}
}
