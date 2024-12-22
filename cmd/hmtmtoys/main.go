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
	defer func() {
		if err = tagsRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Tags repository", err)
		}
	}()

	tagsService := services.NewCommonTagsService(
		tagsRepository,
		logger,
	)

	categoriesRepository := repositories.NewCommonCategoriesRepository(dbConnector, logger)
	defer func() {
		if err = categoriesRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Categories repository", err)
		}
	}()

	categoriesService := services.NewCommonCategoriesService(
		categoriesRepository,
		logger,
	)

	mastersRepository := repositories.NewCommonMastersRepository(dbConnector, logger)
	defer func() {
		if err = mastersRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Masters repository", err)
		}
	}()

	mastersService := services.NewCommonMastersService(
		mastersRepository,
		logger,
	)

	toysRepository := repositories.NewCommonToysRepository(dbConnector, logger)
	defer func() {
		if err = toysRepository.Close(); err != nil {
			logging.LogError(logger, "Failed to close Toys repository", err)
		}
	}()

	toysService := services.NewCommonToysService(
		toysRepository,
		tagsRepository,
		logger,
	)

	useCases := usecases.NewCommonUseCases(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		settings.Security.JWT,
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
