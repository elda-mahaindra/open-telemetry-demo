package api

import (
	"service-a/middleware"
	"service-a/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Api struct {
	logger *logrus.Logger
	tracer trace.Tracer

	service *service.Service
}

func NewApi(
	logger *logrus.Logger,
	tracer trace.Tracer,
	service *service.Service,
) *Api {
	return &Api{
		logger: logger,
		tracer: tracer,

		service: service,
	}
}

func (api *Api) SetupRoutes(app *fiber.App) *fiber.App {
	// Error handler middleware
	app.Use(middleware.ErrorHandler())

	// Ping Routes
	ping := app.Group("/ping")
	ping.Get("/", api.Ping)

	return app
}
