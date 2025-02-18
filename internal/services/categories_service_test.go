package services_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestCategoriesServiceGetCategoryByID(t *testing.T) {
	testCases := []struct {
		name          string
		categoryID    uint32
		errorExpected bool
		err           error
	}{
		{
			name:          "successfully got Category by id",
			categoryID:    1,
			errorExpected: false,
		},
		{
			name:          "failed to get Category by id",
			categoryID:    2,
			errorExpected: true,
			err:           &customerrors.CategoryNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
	categoriesRepository.EXPECT().GetCategoryByID(gomock.Any(), uint32(1)).Return(&entities.Category{}, nil).MaxTimes(1)
	categoriesRepository.EXPECT().GetCategoryByID(gomock.Any(), uint32(2)).Return(
		nil,
		&customerrors.CategoryNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	categoriesService := services.NewCategoriesService(categoriesRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			category, err := categoriesService.GetCategoryByID(ctx, tc.categoryID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, category)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCategoriesServiceGetAllCategories(t *testing.T) {
	t.Run("all categories with existing categories", func(t *testing.T) {
		expectedCategories := []entities.Category{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories(gomock.Any()).Return(expectedCategories, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCategoriesService(categoriesRepository, logger)
		ctx := context.Background()

		categories, err := categoriesService.GetAllCategories(ctx)
		require.NoError(t, err)
		assert.Len(t, categories, len(expectedCategories))
		assert.Equal(t, expectedCategories, categories)
	})

	t.Run("all categories without existing categories", func(t *testing.T) {
		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories(gomock.Any()).Return([]entities.Category{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCategoriesService(categoriesRepository, logger)
		ctx := context.Background()

		categories, err := categoriesService.GetAllCategories(ctx)
		require.NoError(t, err)
		assert.Empty(t, categories)
	})

	t.Run("all categories fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
		categoriesRepository.EXPECT().GetAllCategories(gomock.Any()).Return(
			nil,
			errors.New("test error"),
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		categoriesService := services.NewCategoriesService(categoriesRepository, logger)
		ctx := context.Background()

		categories, err := categoriesService.GetAllCategories(ctx)
		require.Error(t, err)
		assert.Nil(t, categories)
	})
}
