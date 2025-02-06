package main

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-toys/internal/app"
	"github.com/DKhorkov/hmtm-toys/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-toys/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-toys/internal/repositories"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	"github.com/DKhorkov/hmtm-toys/internal/usecases"
)

func main() {
	settings := config.New()
	logger := logging.GetInstance(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	dbConnector, err := db.New(
		db.BuildDsn(settings.Database),
		settings.Database.Driver,
		logger,
		db.WithMaxOpenConnections(settings.Database.Pool.MaxOpenConnections),
		db.WithMaxIdleConnections(settings.Database.Pool.MaxIdleConnections),
		db.WithMaxConnectionLifetime(settings.Database.Pool.MaxConnectionLifetime),
		db.WithMaxConnectionIdleTime(settings.Database.Pool.MaxConnectionIdleTime),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err = dbConnector.Close(); err != nil {
			logging.LogError(logger, "Failed to close db connections pool", err)
		}
	}()

	traceProvider, err := tracing.New(settings.Tracing.Server)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = traceProvider.Shutdown(context.Background()); err != nil {
			logging.LogError(logger, "Error shutting down tracer", err)
		}
	}()

	tagsRepository := repositories.NewCommonTagsRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Tags,
	)

	tagsService := services.NewCommonTagsService(
		tagsRepository,
		logger,
	)

	categoriesRepository := repositories.NewCommonCategoriesRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Categories,
	)

	categoriesService := services.NewCommonCategoriesService(
		categoriesRepository,
		logger,
	)

	mastersRepository := repositories.NewCommonMastersRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Masters,
	)

	mastersService := services.NewCommonMastersService(
		mastersRepository,
		logger,
	)

	toysRepository := repositories.NewCommonToysRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Toys,
	)

	toysService := services.NewCommonToysService(
		toysRepository,
		logger,
	)

	useCases := usecases.NewCommonUseCases(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
	)

	controller := grpccontroller.New(
		settings.HTTP.Host,
		settings.HTTP.Port,
		useCases,
		logger,
		traceProvider,
		settings.Tracing.Spans.Root,
	)

	application := app.New(controller)
	application.Run()
}
