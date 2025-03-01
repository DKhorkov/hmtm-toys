package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DKhorkov/hmtm-toys/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	loggermock "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestToysServiceGetToyByID(t *testing.T) {
	testCases := []struct {
		name          string
		toyID         uint64
		expected      *entities.Toy
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "successfully got Toy by id",
			toyID:    1,
			expected: &entities.Toy{ID: 1},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(1)).
					Return(&entities.Toy{ID: 1}, nil).
					MaxTimes(1)
			},
			errorExpected: false,
		},
		{
			name:  "failed to get Toy by id",
			toyID: 2,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetToyByID(gomock.Any(), uint64(2)).
					Return(nil, &customerrors.ToyNotFoundError{}).
					MaxTimes(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(1)
			},
			errorExpected: true,
			err:           &customerrors.ToyNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	toysService := services.NewToysService(toysRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toy, err := toysService.GetToyByID(ctx, tc.toyID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, toy)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysServiceGetAllToys(t *testing.T) {
	testCases := []struct {
		name          string
		expected      []entities.Toy
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "all Toys with existing Toys",
			expected: []entities.Toy{{ID: 1}},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return(
						[]entities.Toy{
							{ID: 1},
						},
						nil,
					).
					MaxTimes(1)
			},
		},
		{
			name:     "all Toys without existing Toys",
			expected: []entities.Toy{},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return([]entities.Toy{}, nil).
					MaxTimes(1)
			},
		},
		{
			name: "all Toys error",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return(nil, errors.New("test error")).
					MaxTimes(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	toysService := services.NewToysService(toysRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toys, err := toysService.GetAllToys(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, toys, len(tc.expected))
			assert.Equal(t, tc.expected, toys)
		})
	}
}

func TestToysServiceGetMasterToys(t *testing.T) {
	testCases := []struct {
		name          string
		masterID      uint64
		expected      []entities.Toy
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name: "get Master Toys with existing masterID",
			expected: []entities.Toy{
				{
					ID:       1,
					MasterID: 1,
				},
			},
			masterID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(1)).
					Return(
						[]entities.Toy{
							{
								ID:       1,
								MasterID: 1,
							},
						},
						nil,
					).
					MaxTimes(1)
			},
		},
		{
			name:     "get Master Toys error",
			masterID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(1)).
					Return(nil, errors.New("test error")).
					MaxTimes(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	toysService := services.NewToysService(toysRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			toys, err := toysService.GetMasterToys(ctx, tc.masterID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, toys, len(tc.expected))
			assert.Equal(t, tc.expected, toys)
		})
	}
}

func TestToysServiceAddToy(t *testing.T) {
	testCases := []struct {
		name          string
		toy           entities.AddToyDTO
		expected      uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "add Toy success",
			toy:      entities.AddToyDTO{MasterID: 1, Description: "test", Name: "test", CategoryID: 1},
			expected: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(1)).
					Return([]entities.Toy{}, nil).
					MaxTimes(1)

				toysRepository.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.AddToyDTO{
							MasterID:    1,
							Description: "test",
							Name:        "test",
							CategoryID:  1},
					).
					Return(uint64(1), nil).
					MaxTimes(1)
			},
			errorExpected: false,
		},
		{
			name: "add Toy fail - already exists",
			toy:  entities.AddToyDTO{MasterID: 1, Description: "test", Name: "test", CategoryID: 1},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetMasterToys(gomock.Any(), uint64(1)).
					Return(
						[]entities.Toy{
							{
								ID:          1,
								MasterID:    1,
								Name:        "test",
								Description: "test",
								CategoryID:  1,
							},
						}, nil,
					).
					MaxTimes(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					MaxTimes(1)
			},
			errorExpected: true,
			err:           &customerrors.ToyAlreadyExistsError{},
		},
	}

	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	toysService := services.NewToysService(toysRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(toysRepository, logger)
			}

			masterID, err := toysService.AddToy(ctx, tc.toy)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Zero(t, masterID)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// func TestToysServiceAddToy(t *testing.T) {
//	t.Run("add toy success", func(t *testing.T) {
//		const expectedToyID = uint64(1)
//
//		mockController := gomock.NewController(t)
//		toysRepository := mockrepositories.NewMockToysRepository(mockController)
//		toysRepository.EXPECT().AddToy(gomock.Any(), gomock.Any()).Return(expectedToyID, nil).MaxTimes(1)
//		toysRepository.EXPECT().GetMasterToys(gomock.Any(), gomock.Any()).Return([]entities.Toy{}, nil).MaxTimes(1)
//
//		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
//		toysService := NewToysService(toysRepository, logger)
//		ctx := context.Background()
//
//		toyID, err := toysService.AddToy(ctx, entities.AddToyDTO{})
//		require.NoError(t, err)
//		assert.Equal(t, expectedToyID, toyID)
//	})
//
//	t.Run("add toy fail already exists", func(t *testing.T) {
//		var expectedError = &customerrors.ToyAlreadyExistsError{}
//		const (
//			expectedToyID                 = uint64(0)
//			expectedMasterID       uint64 = 1
//			expectedToyName               = "test Toy"
//			expectedToyDescription        = "test Toy description"
//			expectedToyCategory    uint32 = 1
//		)
//
//		mockController := gomock.NewController(t)
//		toysRepository := mockrepositories.NewMockToysRepository(mockController)
//		toysRepository.EXPECT().GetMasterToys(gomock.Any(), gomock.Any()).Return(
//			[]entities.Toy{
//				{
//					ID:          expectedToyID,
//					MasterID:    expectedMasterID,
//					Name:        expectedToyName,
//					Description: expectedToyDescription,
//					CategoryID:  expectedToyCategory,
//				},
//			},
//			nil).MaxTimes(1)
//
//		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
//		toysService := NewToysService(toysRepository, logger)
//		ctx := context.Background()
//
//		toyID, err := toysService.AddToy(
//			ctx,
//			entities.AddToyDTO{
//				MasterID:    expectedMasterID,
//				Name:        expectedToyName,
//				Description: expectedToyDescription,
//				CategoryID:  expectedToyCategory,
//			},
//		)
//
//		require.Error(t, err)
//		require.IsType(t, expectedError, err)
//		assert.Equal(t, expectedToyID, toyID)
//	})
//}
