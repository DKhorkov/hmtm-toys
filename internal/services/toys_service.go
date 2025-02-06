package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

func NewCommonToysService(
	toysRepository interfaces.ToysRepository,
	logger *slog.Logger,
) *CommonToysService {
	return &CommonToysService{
		toysRepository: toysRepository,
		logger:         logger,
	}
}

type CommonToysService struct {
	toysRepository interfaces.ToysRepository
	logger         *slog.Logger
}

func (service *CommonToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	toy, err := service.toysRepository.GetToyByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Toy with ID=%d", id),
			err,
		)

		return nil, &customerrors.ToyNotFoundError{}
	}

	return toy, nil
}

func (service *CommonToysService) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	return service.toysRepository.GetAllToys(ctx)
}

func (service *CommonToysService) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	return service.toysRepository.GetMasterToys(ctx, masterID)
}

func (service *CommonToysService) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	if service.checkToyExistence(ctx, toyData) {
		return 0, &customerrors.ToyAlreadyExistsError{}
	}

	return service.toysRepository.AddToy(ctx, toyData)
}

func (service *CommonToysService) checkToyExistence(ctx context.Context, toyData entities.AddToyDTO) bool {
	toys, err := service.toysRepository.GetMasterToys(ctx, toyData.MasterID)
	if err == nil {
		for _, toy := range toys {
			if toy.Name == toyData.Name && toy.CategoryID == toyData.CategoryID && toy.Description == toyData.Description {
				return true
			}
		}
	}

	return false
}
