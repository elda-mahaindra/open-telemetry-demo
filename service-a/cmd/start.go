package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"service-a/api"
	"service-a/service"
	"service-a/util/config"
	"service-a/util/tracing"

	"github.com/sirupsen/logrus"
)

func start() {
	const op = "main.start"

	// --- Init logger ---
	var logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	// --- Load config ---
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "LoadConfig",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}

	// --- Init otel tracer ---
	cleanup, err := tracing.InitTracer(config.OtelTracer.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "InitTracer",
			"err":   err.Error(),
		}).Error()
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.WithFields(logrus.Fields{
				"[op]":  op,
				"scope": "CleanupTracer",
				"err":   err.Error(),
			}).Error()
		}
	}()

	tracer := tracing.GetTracer(config.OtelTracer.Name)

	logger.WithFields(logrus.Fields{
		"[op]":   op,
		"config": fmt.Sprintf("%+v", config),
	}).Infof("Starting '%s' service ...", config.App.Name)

	// --- Init service-b adapter ---
	serviceBAdapter, err := createServiceBAdapter(config.ServiceB, logger, tracer)
	if err != nil {
		log.Printf("failed to create service-b adapter: %v", err)
		os.Exit(1)
	}

	// --- Init service layer ---
	service := service.NewService(logger, tracer, serviceBAdapter)

	// --- Init api layer ---
	restApi := api.NewApi(logger, tracer, service)

	// --- Run servers ---
	runRestServer(config.App.Port, restApi)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// --- Block until signal is received ---
	<-ch

	log.Printf("end of program...")
}
