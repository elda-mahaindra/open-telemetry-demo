package service_b_adapter

import (
	"context"
	"fmt"

	"service-a/adapter/service_b_adapter/pb"
	"service-a/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
)

func (client *Adapter) Ping(ctx context.Context, message string) (*pb.PingResponse, error) {
	const op = "service_b_adapter.Adapter.Ping"

	// Start span
	ctx, span := client.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("service_b_adapter.operation", "ping"),
		attribute.String("service_b_adapter.input.message", message),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, client.logger)

	request := &pb.PingRequest{
		PingMessage: message,
	}

	logger.WithFields(logrus.Fields{
		"[op]":    op,
		"request": request,
		"type":    fmt.Sprintf("%T", request),
	}).Info()

	// Call service B
	response, err := client.serviceBClient.Ping(ctx, request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":    op,
			"request": request,
			"type":    fmt.Sprintf("%T", request),
			"error":   err.Error(),
		}).Error()

		return nil, fmt.Errorf("error sending request: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"[op]":     op,
		"response": response,
		"type":     fmt.Sprintf("%T", response),
	}).Info()

	return response, nil
}
