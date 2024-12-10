package services

import (
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/entities"

	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
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
