package services_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
)

func TestCommonCategoriesServiceGetCategoryByID(t *testing.T) {
	testCases := []struct {
		categoryID    uint32
		resultLength  int
		errorExpected bool
		err           error
	}{
		{
			categoryID:    1,
			resultLength:  1,
			errorExpected: false,
		},
		{
			categoryID:    2,
			errorExpected: true,
			err:           &customerrors.CategoryNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
	categoriesRepository.EXPECT().GetCategoryByID(uint32(1)).Return(&entities.Category{}, nil).MaxTimes(1)
	categoriesRepository.EXPECT().GetCategoryByID(uint32(2)).Return(
		nil,
		&customerrors.CategoryNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	categoriesService := services.NewCommonCategoriesService(categoriesRepository, logger)

	for _, tc := range testCases {
		category, err := categoriesService.GetCategoryByID(tc.categoryID)
		if tc.errorExpected {
			require.Error(t, err)
			require.IsType(t, tc.err, err)
			assert.Nil(t, category)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestCommonCategoriesServiceGetAllCategories(t *testing.T) {
	t.Run("all categories with existing categories", func(t *testing.T) {
		expectedCategories := []*entities.Category{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories().DoAndReturn(
			func() ([]*entities.Category, error) {
				return expectedCategories, nil
			},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCommonCategoriesService(categoriesRepository, logger)

		categories, err := categoriesService.GetAllCategories()
		require.NoError(t, err)
		assert.Len(t, categories, len(expectedCategories))
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("all categories without existing categories", func(t *testing.T) {
		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories().Return([]*entities.Category{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCommonCategoriesService(categoriesRepository, logger)

		categories, err := categoriesService.GetAllCategories()
		require.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("all categories fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories().Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCommonCategoriesService(categoriesRepository, logger)

		categories, err := categoriesService.GetAllCategories()
		require.Error(t, err)
		assert.Nil(t, categories)
	})
}
