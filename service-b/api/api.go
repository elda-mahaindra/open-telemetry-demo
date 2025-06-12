package api

import (
	"service-b/api/pb"
	"service-b/service"
)

type Api struct {
	pb.UnimplementedBServiceServer

	serviceName string

	service *service.Service
}

func NewApi(serviceName string, service *service.Service) *Api {
	return &Api{
		serviceName: serviceName,
		service:     service,
	}
}
