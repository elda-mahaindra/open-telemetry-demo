package service

import (
	"context"
	"fmt"
	"time"

	"service-b/store"
	"service-b/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type PingParams struct {
	PingMessage string `json:"ping_message"`
}

type PingResult struct {
	PongMessage string `json:"pong_message"`
}

func (service *Service) Ping(ctx context.Context, params *PingParams) (*PingResult, error) {
	const op = "service.Service.Ping"

	// Start span
	ctx, span := service.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("service.operation", "ping"),
		attribute.String("service.input.params", fmt.Sprintf("%+v", params)),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, service.logger)

	logger.WithFields(logrus.Fields{
		"[op]":   op,
		"params": params,
	}).Info()

	// Simulate a service operation
	time.Sleep(500 * time.Millisecond)

	arg := &store.PingArgs{
		PingMessage: params.PingMessage,
	}

	data, err := service.store.Ping(ctx, arg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":   op,
			"params": params,
			"error":  err,
		}).Error()

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
		attribute.String("service.output.result", fmt.Sprintf("%+v", result)),
	)
	span.SetStatus(codes.Ok, "service operation completed successfully")

	return result, nil
}
