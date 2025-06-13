package service

import (
	"service-a/adapter/service_b_adapter"
)

type Service struct {
	serviceBAdapter *service_b_adapter.Adapter
}

func NewService(serviceBAdapter *service_b_adapter.Adapter) *Service {
	return &Service{
		serviceBAdapter: serviceBAdapter,
	}
}
