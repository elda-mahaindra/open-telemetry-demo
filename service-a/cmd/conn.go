package main

import (
	"fmt"

	"service-a/adapter/service_b_adapter"
	"service-a/util/config"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createServiceBAdapter(config config.ServiceB, logger *logrus.Logger, tracer trace.Tracer) (*service_b_adapter.Adapter, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithPropagators(propagation.TraceContext{}),
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s grpc server: %w", config.Name, err)
	}

	grpcAdapter := service_b_adapter.NewAdapter(config.Name, logger, tracer, conn)

	return grpcAdapter, nil
}
