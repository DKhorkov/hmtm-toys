package main

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-toys/internal/app"
	ssogrpcclient "github.com/DKhorkov/hmtm-toys/internal/clients/sso/grpc"
	"github.com/DKhorkov/hmtm-toys/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-toys/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-toys/internal/repositories"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	"github.com/DKhorkov/hmtm-toys/internal/usecases"
)

func main() {
	settings := config.New()
	logger := logging.New(
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

	ssoClient, err := ssogrpcclient.New(
		settings.Clients.SSO.Host,
		settings.Clients.SSO.Port,
		settings.Clients.SSO.RetriesCount,
		settings.Clients.SSO.RetryTimeout,
		logger,
		traceProvider,
		settings.Tracing.Spans.Clients.SSO,
	)
	if err != nil {
		panic(err)
	}

	ssoRepository := repositories.NewSsoRepository(ssoClient)
	ssoService := services.NewSsoService(ssoRepository, logger)

	tagsRepository := repositories.NewTagsRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Tags,
	)

	tagsService := services.NewTagsService(
		tagsRepository,
		logger,
	)

	categoriesRepository := repositories.NewCategoriesRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Categories,
	)

	categoriesService := services.NewCategoriesService(
		categoriesRepository,
		logger,
	)

	mastersRepository := repositories.NewMastersRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Masters,
	)

	mastersService := services.NewMastersService(
		mastersRepository,
		logger,
	)

	toysRepository := repositories.NewToysRepository(
		dbConnector,
		logger,
		traceProvider,
		settings.Tracing.Spans.Repositories.Toys,
	)

	toysService := services.NewToysService(
		toysRepository,
		logger,
	)

	useCases := usecases.New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
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
