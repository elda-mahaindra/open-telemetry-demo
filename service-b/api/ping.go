package api

import (
	"context"
	"log"
	"time"

	"service-b/api/pb"
	"service-b/service"
	"service-b/util/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (api *Api) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	log.Println("Api.Ping")

	// Start a new span for the API request
	tracer := tracing.GetTracer("service-a-api")
	ctx, span := tracer.Start(ctx, "api.ping")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("api.endpoint", "/ping"),
		attribute.String("api.method", "GET"),
		attribute.String("ping.message", request.GetPingMessage()),
	)

	// Simulate a validation operation
	time.Sleep(250 * time.Millisecond)

	param := &service.PingParam{
		PingMessage: request.GetPingMessage(),
	}

	result, err := api.service.Ping(ctx, param)
	if err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("response.pong_message", result.PongMessage),
	)
	span.SetStatus(codes.Ok, "request completed successfully")

	return &pb.PingResponse{
		PongMessage: result.PongMessage,
	}, nil
}
