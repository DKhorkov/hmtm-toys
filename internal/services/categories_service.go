package services

import (
	"log/slog"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/logging"
)

type CommonCategoriesService struct {
	categoriesRepository interfaces.CategoriesRepository
	logger               *slog.Logger
}

func (service *CommonCategoriesService) GetCategoryByID(id uint32) (*entities.Category, error) {
	category, err := service.categoriesRepository.GetCategoryByID(id)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get category by id",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.CategoryNotFoundError{}
	}

	return category, nil
}

func (service *CommonCategoriesService) GetAllCategories() ([]*entities.Category, error) {
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
