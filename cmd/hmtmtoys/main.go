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

	defer dbConnector.CloseConnection()

	mastersRepository := repositories.NewCommonMastersRepository(dbConnector)
	tagsRepository := repositories.NewCommonTagsRepository(dbConnector)
	categoriesRepository := repositories.NewCommonCategoriesRepository(dbConnector)
	toysRepository := repositories.NewCommonToysRepository(dbConnector)
	tagsService := services.NewCommonTagsService(
		tagsRepository,
		logger,
	)

	categoriesService := services.NewCommonCategoriesService(
		categoriesRepository,
		logger,
	)

	mastersService := services.NewCommonMastersService(
		mastersRepository,
		logger,
	)

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
