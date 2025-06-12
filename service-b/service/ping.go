package service

import (
	"context"
	"log"
	"time"

	"service-b/store"
	"service-b/util/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type PingParam struct {
	PingMessage string `json:"ping_message"`
}

type PingResult struct {
	PongMessage string `json:"pong_message"`
}

func (s *Service) Ping(ctx context.Context, param *PingParam) (*PingResult, error) {
	log.Println("Service.Ping")

	// Start a new span for the service operation
	tracer := tracing.GetTracer("service-a-service")
	ctx, span := tracer.Start(ctx, "service.ping")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("service.operation", "ping"),
		attribute.String("service.input.message", param.PingMessage),
	)

	// Simulate a service operation
	time.Sleep(500 * time.Millisecond)

	arg := &store.PingArg{
		PingMessage: param.PingMessage,
	}

	data, err := s.store.Ping(ctx, arg)
	if err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	result := &PingResult{
		PongMessage: data.PongMessage,
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("service.output.message", result.PongMessage),
	)
	span.SetStatus(codes.Ok, "service operation completed successfully")

	return result, nil
}
