package api

import (
	"context"
	"fmt"
	"time"

	"service-b/api/pb"
	"service-b/service"
	"service-b/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (api *Api) Ping(ctx context.Context, request *pb.PingRequest) (*pb.PingResponse, error) {
	const op = "api.Api.Ping"

	// Start span
	ctx, span := api.tracer.Start(ctx, op)
	defer span.End()

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, api.logger)

	logger.WithFields(logrus.Fields{
		"[op]":    op,
		"request": request,
	}).Info()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("api.operation", "ping"),
		attribute.String("api.input.request", fmt.Sprintf("%+v", request)),
	)

	// Simulate a validation operation
	time.Sleep(250 * time.Millisecond)

	params := &service.PingParams{
		PingMessage: request.GetPingMessage(),
	}

	result, err := api.service.Ping(ctx, params)
	if err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	response := &pb.PingResponse{
		PongMessage: result.PongMessage,
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("api.output.response", fmt.Sprintf("%+v", response)),
	)
	span.SetStatus(codes.Ok, "request completed successfully")

	return response, nil
}
