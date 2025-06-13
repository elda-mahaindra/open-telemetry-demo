package service_b_adapter

import (
	"context"
	"fmt"
	"log"

	"service-a/adapter/service_b_adapter/pb"
)

func (client *Adapter) Ping(ctx context.Context, message string) (*pb.PingResponse, error) {
	log.Println("Adapter.Ping")

	request := &pb.PingRequest{
		PingMessage: message,
	}

	log.Printf("Sending request to service-b: %v\n", request)

	response, err := client.serviceBClient.Ping(ctx, request)
	if err != nil {
		log.Printf("Error sending request to service-b: %v\n", err)

		return nil, fmt.Errorf("error sending request: %w", err)
	}

	log.Printf("Received response from service-b: %v\n", response)

	return response, nil
}
