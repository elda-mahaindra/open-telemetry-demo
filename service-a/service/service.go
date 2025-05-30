package service

import (
	"service-a/store"
)

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}
