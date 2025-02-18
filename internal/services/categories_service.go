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

func NewCategoriesService(
	categoriesRepository interfaces.CategoriesRepository,
	logger *slog.Logger,
) *CategoriesService {
	return &CategoriesService{
		categoriesRepository: categoriesRepository,
		logger:               logger,
	}
}

type CategoriesService struct {
	categoriesRepository interfaces.CategoriesRepository
	logger               *slog.Logger
}

func (service *CategoriesService) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	category, err := service.categoriesRepository.GetCategoryByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Category with ID=%d", id),
			err,
		)

		return nil, &customerrors.CategoryNotFoundError{}
	}

	return category, nil
}

func (service *CategoriesService) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return service.categoriesRepository.GetAllCategories(ctx)
}
