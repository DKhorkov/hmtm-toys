package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	loggermock "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestCategoriesServiceGetCategoryByID(t *testing.T) {
	testCases := []struct {
		name          string
		categoryID    uint32
		expected      *entities.Category
		setupMocks    func(categoriesRepository *mockrepositories.MockCategoriesRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:       "successfully got Category by id",
			categoryID: 1,
			expected:   &entities.Category{ID: 1},
			setupMocks: func(categoriesRepository *mockrepositories.MockCategoriesRepository, _ *loggermock.MockLogger) {
				categoriesRepository.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(1)).
					Return(&entities.Category{ID: 1}, nil).
					MaxTimes(1)
			},
			errorExpected: false,
		},
		{
			name:       "failed to get Category by id",
			categoryID: 2,
			setupMocks: func(categoriesRepository *mockrepositories.MockCategoriesRepository, logger *loggermock.MockLogger) {
				categoriesRepository.
					EXPECT().
					GetCategoryByID(gomock.Any(), uint32(2)).
					Return(nil, &customerrors.CategoryNotFoundError{}).
					MaxTimes(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(1)
			},
			errorExpected: true,
			err:           &customerrors.CategoryNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	categoriesService := services.NewCategoriesService(categoriesRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(categoriesRepository, logger)
			}

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
	testCases := []struct {
		name          string
		expected      []entities.Category
		setupMocks    func(categoriesRepository *mockrepositories.MockCategoriesRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "all Categories with existing Categories",
			expected: []entities.Category{{ID: 1}},
			setupMocks: func(categoriesRepository *mockrepositories.MockCategoriesRepository, _ *loggermock.MockLogger) {
				categoriesRepository.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(
						[]entities.Category{
							{ID: 1},
						},
						nil,
					).
					MaxTimes(1)
			},
		},
		{
			name:     "all Categories without existing Categories",
			expected: []entities.Category{},
			setupMocks: func(categoriesRepository *mockrepositories.MockCategoriesRepository, _ *loggermock.MockLogger) {
				categoriesRepository.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return([]entities.Category{}, nil).
					MaxTimes(1)
			},
		},
		{
			name: "all Categories error",
			setupMocks: func(categoriesRepository *mockrepositories.MockCategoriesRepository, _ *loggermock.MockLogger) {
				categoriesRepository.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(nil, errors.New("test error")).
					MaxTimes(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	categoriesRepository := mockrepositories.NewMockCategoriesRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	categoriesService := services.NewCategoriesService(categoriesRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(categoriesRepository, logger)
			}

			categories, err := categoriesService.GetAllCategories(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, categories, len(tc.expected))
			assert.Equal(t, tc.expected, categories)
		})
	}
}
