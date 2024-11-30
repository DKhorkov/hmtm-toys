package services

import (
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/logging"
)

func processToyTags(
	toy *entities.Toy,
	tagsRepository interfaces.TagsRepository,
	logger *slog.Logger,
) error {
	toyTags, err := tagsRepository.GetToyTags(toy.ID)
	if err != nil {
		logger.Error(
			fmt.Sprintf("Error occurred while trying to get tags for toy with ID=%d", toy.ID),
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return err
	}

	toy.Tags = toyTags
	return nil
}

func processMasterToys(
	master *entities.Master,
	toysRepository interfaces.ToysRepository,
	tagsRepository interfaces.TagsRepository,
	logger *slog.Logger,
) error {
	masterToys, err := toysRepository.GetMasterToys(master.ID)
	if err != nil {
		logger.Error(
			fmt.Sprintf("Error occurred while trying to get toys for master with ID=%d", master.ID),
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return err
	}

	for _, toy := range masterToys {
		if err = processToyTags(toy, tagsRepository, logger); err != nil {
			return err
		}
	}

	master.Toys = masterToys
	return nil
}
