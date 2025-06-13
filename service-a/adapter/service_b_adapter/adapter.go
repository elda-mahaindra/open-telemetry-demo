package service_b_adapter

import (
	"service-a/adapter/service_b_adapter/pb"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a wrapper around the grpc client
type Adapter struct {
	serviceName string

	logger *logrus.Logger
	tracer trace.Tracer

	serviceBClient pb.BServiceClient
}

// NewAdapter creates a new grpc adapter
func NewAdapter(
	serviceName string,
	logger *logrus.Logger,
	tracer trace.Tracer,
	cc *grpc.ClientConn,
) *Adapter {
	serviceBClient := pb.NewBServiceClient(cc)

	return &Adapter{
		serviceName: serviceName,

		logger: logger,
		tracer: tracer,

		serviceBClient: serviceBClient,
	}
}
