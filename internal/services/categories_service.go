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

type CommonCategoriesService struct {
	categoriesRepository interfaces.CategoriesRepository
	logger               *slog.Logger
}

func (service *CommonCategoriesService) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	category, err := service.categoriesRepository.GetCategoryByID(id)
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

func (service *CommonCategoriesService) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return service.categoriesRepository.GetAllCategories()
}

func NewCommonCategoriesService(
	categoriesRepository interfaces.CategoriesRepository,
	logger *slog.Logger,
) *CommonCategoriesService {
	return &CommonCategoriesService{
		categoriesRepository: categoriesRepository,
		logger:               logger,
	}
}
