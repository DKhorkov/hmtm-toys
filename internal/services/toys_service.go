package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

type ToysService struct {
	toysRepository interfaces.ToysRepository
	logger         logging.Logger
}

func NewToysService(
	toysRepository interfaces.ToysRepository,
	logger logging.Logger,
) *ToysService {
	return &ToysService{
		toysRepository: toysRepository,
		logger:         logger,
	}
}

func (service *ToysService) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
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

func (service *ToysService) GetToys(ctx context.Context, pagination *entities.Pagination) ([]entities.Toy, error) {
	return service.toysRepository.GetToys(ctx, pagination)
}

func (service *ToysService) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
) ([]entities.Toy, error) {
	return service.toysRepository.GetMasterToys(ctx, masterID, pagination)
}

func (service *ToysService) AddToy(
	ctx context.Context,
	toyData entities.AddToyDTO,
) (uint64, error) {
	if service.checkToyExistence(ctx, toyData) {
		return 0, &customerrors.ToyAlreadyExistsError{}
	}

	return service.toysRepository.AddToy(ctx, toyData)
}

func (service *ToysService) DeleteToy(ctx context.Context, id uint64) error {
	return service.toysRepository.DeleteToy(ctx, id)
}

func (service *ToysService) UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error {
	return service.toysRepository.UpdateToy(ctx, toyData)
}

func (service *ToysService) checkToyExistence(
	ctx context.Context,
	toyData entities.AddToyDTO,
) bool {
	toys, err := service.toysRepository.GetMasterToys(ctx, toyData.MasterID, nil)
	if err == nil {
		for _, toy := range toys {
			if toy.Name == toyData.Name && toy.CategoryID == toyData.CategoryID &&
				toy.Description == toyData.Description {
				return true
			}
		}
	}

	return false
}
