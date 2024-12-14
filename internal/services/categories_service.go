package services

import (
	"context"
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
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get category by id", err)
		return nil, &customerrors.CategoryNotFoundError{BaseErr: err}
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
