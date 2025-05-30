package api

import (
	"service-a/middleware"
	"service-a/service"

	"github.com/gofiber/fiber/v2"
)

type Api struct {
	serviceName string

	service *service.Service
}

func NewApi(serviceName string, service *service.Service) *Api {
	return &Api{
		serviceName: serviceName,
		service:     service,
	}
}

func (api *Api) DefineEndpoints(app *fiber.App) *fiber.App {
	// Error handler middleware
	app.Use(middleware.ErrorHandler())

	// Ping Routes
	ping := app.Group("/ping")
	ping.Get("/", api.Ping)

	return app
}
