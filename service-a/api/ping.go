package api

import (
	"log"
	"time"

	"service-a/service"
	"service-a/util/tracing"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (api *Api) Ping(c *fiber.Ctx) error {
	log.Println("Api.Ping")

	// Start a new span for the API request
	tracer := tracing.GetTracer("service-a-api")
	ctx, span := tracer.Start(c.Context(), "api.ping")
	defer span.End()

	message := c.Query("message")

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("api.endpoint", "/ping"),
		attribute.String("api.method", "GET"),
		attribute.String("ping.message", message),
	)

	// Simulate a validation operation
	time.Sleep(250 * time.Millisecond)

	param := &service.PingParam{
		PingMessage: message,
	}

	result, err := api.service.Ping(ctx, param)
	if err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("response.pong_message", result.PongMessage),
	)
	span.SetStatus(codes.Ok, "request completed successfully")

	return c.JSON(result)
}
