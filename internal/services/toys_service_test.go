package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-toys/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	loggermock "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestToysService_GetToyByID(t *testing.T) {
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
					Times(1)
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
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
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

func TestToysService_GetAllToys(t *testing.T) {
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
					Times(1)
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
					Times(1)
			},
		},
		{
			name: "all Toys error",
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					GetAllToys(gomock.Any()).
					Return(nil, errors.New("test error")).
					Times(1)
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

func TestToysService_GetMasterToys(t *testing.T) {
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
					Times(1)
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
					Times(1)
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

func TestToysService_AddToy(t *testing.T) {
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
					Times(1)

				toysRepository.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.AddToyDTO{
							MasterID:    1,
							Description: "test",
							Name:        "test",
							CategoryID:  1,
						},
					).
					Return(uint64(1), nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name: "add Toy fail - already exists",
			toy:  entities.AddToyDTO{MasterID: 1, Description: "test", Name: "test", CategoryID: 1},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
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
					Times(1)
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

func TestToysService_DeleteToy(t *testing.T) {
	testCases := []struct {
		name          string
		toyID         uint64
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:  "delete Toy success",
			toyID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "delete Toy fail - not found",
			toyID: 1,
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					DeleteToy(gomock.Any(), uint64(1)).
					Return(errors.New("test error")).
					Times(1)
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

			err := toysService.DeleteToy(ctx, tc.toyID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysService_UpdateToy(t *testing.T) {
	testCases := []struct {
		name          string
		toy           entities.UpdateToyDTO
		setupMocks    func(toysRepository *mockrepositories.MockToysRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name: "add Toy success",
			toy: entities.UpdateToyDTO{
				ID:                    1,
				Description:           pointers.New[string]("test"),
				Name:                  pointers.New[string]("test"),
				CategoryID:            pointers.New[uint32](1),
				Price:                 pointers.New[float32](10),
				Quantity:              pointers.New[uint32](1),
				TagIDsToAdd:           []uint32{1, 2},
				TagIDsToDelete:        []uint32{3, 4},
				AttachmentsToAdd:      []string{"test"},
				AttachmentIDsToDelete: []uint64{1},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.UpdateToyDTO{
							ID:                    1,
							Description:           pointers.New[string]("test"),
							Name:                  pointers.New[string]("test"),
							CategoryID:            pointers.New[uint32](1),
							Price:                 pointers.New[float32](10),
							Quantity:              pointers.New[uint32](1),
							TagIDsToAdd:           []uint32{1, 2},
							TagIDsToDelete:        []uint32{3, 4},
							AttachmentsToAdd:      []string{"test"},
							AttachmentIDsToDelete: []uint64{1},
						},
					).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "add Toy fail - already exists",
			toy: entities.UpdateToyDTO{
				ID:                    1,
				Description:           pointers.New[string]("test"),
				Name:                  pointers.New[string]("test"),
				CategoryID:            pointers.New[uint32](1),
				Price:                 pointers.New[float32](10),
				Quantity:              pointers.New[uint32](1),
				TagIDsToAdd:           []uint32{1, 2},
				TagIDsToDelete:        []uint32{3, 4},
				AttachmentsToAdd:      []string{"test"},
				AttachmentIDsToDelete: []uint64{1},
			},
			setupMocks: func(toysRepository *mockrepositories.MockToysRepository, _ *loggermock.MockLogger) {
				toysRepository.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.UpdateToyDTO{
							ID:                    1,
							Description:           pointers.New[string]("test"),
							Name:                  pointers.New[string]("test"),
							CategoryID:            pointers.New[uint32](1),
							Price:                 pointers.New[float32](10),
							Quantity:              pointers.New[uint32](1),
							TagIDsToAdd:           []uint32{1, 2},
							TagIDsToDelete:        []uint32{3, 4},
							AttachmentsToAdd:      []string{"test"},
							AttachmentIDsToDelete: []uint64{1},
						},
					).
					Return(errors.New("test")).
					Times(1)
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

			err := toysService.UpdateToy(ctx, tc.toy)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
