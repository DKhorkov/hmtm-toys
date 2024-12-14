package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/entities"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
)

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
	tagsRepository interfaces.TagsRepository
	logger         *slog.Logger
}

func (service *CommonToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	toy, err := service.toysRepository.GetToyByID(id)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get toy by id", err)
		return nil, &customerrors.ToyNotFoundError{BaseErr: err}
	}

	if err = processToyTags(toy, service.tagsRepository, service.logger); err != nil {
		return nil, err
	}

	return toy, nil
}

func (service *CommonToysService) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetAllToys()
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all toys", err)
		return nil, err
	}

	for _, toy := range toys {
		if err = processToyTags(&toy, service.tagsRepository, service.logger); err != nil {
			return nil, err
		}
	}

	return toys, nil
}

func (service *CommonToysService) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	toys, err := service.toysRepository.GetMasterToys(masterID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get all toys for master with ID=%d", masterID),
			err,
		)

		return nil, err
	}

	for _, toy := range toys {
		if err = processToyTags(&toy, service.tagsRepository, service.logger); err != nil {
			return nil, err
		}
	}

	return toys, nil
}

func (service *CommonToysService) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
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
