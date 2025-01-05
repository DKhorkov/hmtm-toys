package main

import (
	"github.com/DKhorkov/hmtm-toys/internal/app"
	"github.com/DKhorkov/hmtm-toys/internal/config"
	grpccontroller "github.com/DKhorkov/hmtm-toys/internal/controllers/grpc"
	"github.com/DKhorkov/hmtm-toys/internal/repositories"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	"github.com/DKhorkov/hmtm-toys/internal/usecases"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
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

	tagsRepository := repositories.NewCommonTagsRepository(dbConnector, logger)
	tagsService := services.NewCommonTagsService(
		tagsRepository,
		logger,
	)

	categoriesRepository := repositories.NewCommonCategoriesRepository(dbConnector, logger)
	categoriesService := services.NewCommonCategoriesService(
		categoriesRepository,
		logger,
	)

	mastersRepository := repositories.NewCommonMastersRepository(dbConnector, logger)
	mastersService := services.NewCommonMastersService(
		mastersRepository,
		logger,
	)

	toysRepository := repositories.NewCommonToysRepository(dbConnector, logger)
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
	)

	application := app.New(controller)
	application.Run()
}
