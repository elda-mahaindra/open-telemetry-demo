package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"service-b/api"
	"service-b/service"
	"service-b/store"
	"service-b/util/config"
	"service-b/util/tracing"
)

func start() {
	// Load environment variables from .env file
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("failed to load config: %v", err)

		os.Exit(1)
	}

	log.Printf("config: %v", config)
	log.Printf("Starting %s service ...", config.App.Name)

	// Initialize OpenTelemetry
	cleanup, err := tracing.InitTracer(config.OpenTelemetry.TracerName, config.OpenTelemetry.CollectorURL)
	if err != nil {
		log.Printf("failed to initialize tracer: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			log.Printf("failed to cleanup tracer: %v", err)
		}
	}()

	// Init store layer
	store := store.NewStore()

	// Init service layer
	service := service.NewService(store)

	// Init api layer
	restApi := api.NewApi(config.App.Name, service)

	// Run rest server
	runGrpcServer(config.App.Port, restApi)

	// wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch

	log.Printf("end of program...")
}
