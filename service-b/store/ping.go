package store

import (
	"context"
	"fmt"
	"strings"
	"time"

	"service-b/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type PingArgs struct {
	PingMessage string `json:"message"`
}

type PingData struct {
	PongMessage string `json:"message"`
}

func (store *Store) Ping(ctx context.Context, args *PingArgs) (*PingData, error) {
	const op = "store.Store.Ping"

	// Start span
	ctx, span := store.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("store.operation", "ping"),
		attribute.String("store.input.args", fmt.Sprintf("%+v", args)),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, store.logger)

	logger.WithFields(logrus.Fields{
		"[op]": op,
		"args": args,
	}).Info()

	if strings.Contains(args.PingMessage, "error") {
		return nil, fmt.Errorf("error in store.ping")
	}

	message := fmt.Sprintf("pong %s", args.PingMessage)

	// Simulate a database operation
	span.AddEvent("database_query_start")
	time.Sleep(500 * time.Millisecond)
	span.AddEvent("database_query_end")

	data := &PingData{
		PongMessage: message,
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("store.output.data", fmt.Sprintf("%+v", data)),
	)
	span.SetStatus(codes.Ok, "store operation completed successfully")

	return data, nil
}
