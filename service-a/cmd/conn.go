package main

import (
	"fmt"

	"service-a/adapter/service_b_adapter"
	"service-a/util/config"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createServiceBAdapter(config config.ServiceB) (*service_b_adapter.Adapter, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s grpc server: %w", config.Name, err)
	}

	grpcAdapter := service_b_adapter.NewAdapter(config.Name, conn)

	return grpcAdapter, nil
}
