package api

import (
	"service-b/api/pb"
	"service-b/service"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Api struct {
	pb.UnimplementedBServiceServer

	logger *logrus.Logger
	tracer trace.Tracer

	service *service.Service
}

func NewApi(
	logger *logrus.Logger,
	tracer trace.Tracer,
	service *service.Service,
) *Api {
	return &Api{
		logger: logger,
		tracer: tracer,

		service: service,
	}
}
