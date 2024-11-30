package services

import (
	"log/slog"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/logging"
)

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
	tagsRepository interfaces.TagsRepository
	logger         *slog.Logger
}

func (service *CommonToysService) GetToyByID(id uint64) (*entities.Toy, error) {
	toy, err := service.toysRepository.GetToyByID(id)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get toy by id",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.ToyNotFoundError{}
	}

	if err = processToyTags(toy, service.tagsRepository, service.logger); err != nil {
		return nil, err
	}

	return toy, nil
}

func (service *CommonToysService) GetAllToys() ([]*entities.Toy, error) {
	toys, err := service.toysRepository.GetAllToys()
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get all toys",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, err
	}

	for _, toy := range toys {
		if err = processToyTags(toy, service.tagsRepository, service.logger); err != nil {
			return nil, err
		}
	}

	return toys, nil
}

func (service *CommonToysService) AddToy(toyData entities.AddToyDTO) (uint64, error) {
	return service.toysRepository.AddToy(toyData)
}

func NewCommonToysService(
	toysRepository interfaces.ToysRepository,
	tagsRepository interfaces.TagsRepository,
	logger *slog.Logger,
) *CommonToysService {
	return &CommonToysService{
		toysRepository: toysRepository,
		tagsRepository: tagsRepository,
		logger:         logger,
	}
}
