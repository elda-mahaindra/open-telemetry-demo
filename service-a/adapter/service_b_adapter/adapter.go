package service_b_adapter

import (
	"service-a/adapter/service_b_adapter/pb"

	"google.golang.org/grpc"
)

// Adapter is a wrapper around the grpc client
type Adapter struct {
	serviceName string

	serviceBClient pb.BServiceClient
}

// NewAdapter creates a new grpc adapter
func NewAdapter(
	svcName string,
	cc *grpc.ClientConn,
) *Adapter {
	serviceBClient := pb.NewBServiceClient(cc)

	return &Adapter{
		serviceName:    svcName,
		serviceBClient: serviceBClient,
	}
}
