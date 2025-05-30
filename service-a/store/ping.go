package store

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"service-a/util/tracing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type PingArg struct {
	PingMessage string `json:"message"`
}

type PingData struct {
	PongMessage string `json:"message"`
}

func (s *Store) Ping(ctx context.Context, arg *PingArg) (*PingData, error) {
	log.Println("Store.Ping")

	// Start a new span for the store operation
	tracer := tracing.GetTracer("service-a-store")
	_, span := tracer.Start(ctx, "store.ping")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("store.operation", "ping"),
		attribute.String("store.input.message", arg.PingMessage),
		attribute.String("store.operation_type", "database_query"),
	)

	if strings.Contains(arg.PingMessage, "error") {
		return nil, fmt.Errorf("error in store.ping")
	}

	message := fmt.Sprintf("pong %s", arg.PingMessage)

	// Simulate a database operation
	span.AddEvent("database_query_start")
	time.Sleep(1 * time.Second)
	span.AddEvent("database_query_end")

	data := &PingData{
		PongMessage: message,
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("store.output.message", data.PongMessage),
		attribute.Int("store.query_duration_ms", 1000),
	)
	span.SetStatus(codes.Ok, "store operation completed successfully")

	return data, nil
}
