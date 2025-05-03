package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	loggermock "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestMastersService_GetMasterByID(t *testing.T) {
	testCases := []struct {
		name          string
		masterID      uint64
		expected      *entities.Master
		setupMocks    func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "successfully got Master by id",
			masterID: 1,
			expected: &entities.Master{ID: 1},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(1)).
					Return(&entities.Master{ID: 1}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:     "failed to get Master by id",
			masterID: 2,
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByID(gomock.Any(), uint64(2)).
					Return(nil, &customerrors.MasterNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	mastersService := services.NewMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mastersRepository, logger)
			}

			master, err := mastersService.GetMasterByID(ctx, tc.masterID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, master)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMastersService_GetMasterByUserID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint64
		expected      *entities.Master
		setupMocks    func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "successfully got Master by userID",
			userID:   1,
			expected: &entities.Master{ID: 1},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(&entities.Master{ID: 1}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:   "failed to get Master by userID",
			userID: 2,
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(2)).
					Return(nil, &customerrors.MasterNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	mastersService := services.NewMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mastersRepository, logger)
			}

			master, err := mastersService.GetMasterByUserID(ctx, tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, master)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMastersService_GetMasters(t *testing.T) {
	testCases := []struct {
		name          string
		pagination    *entities.Pagination
		expected      []entities.Master
		setupMocks    func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name: "all Masters with existing Masters",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			expected: []entities.Master{{ID: 1}},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(
						[]entities.Master{
							{ID: 1},
						},
						nil,
					).
					Times(1)
			},
		},
		{
			name: "all Masters without existing Masters",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			expected: []entities.Master{},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return([]entities.Master{}, nil).
					Times(1)
			},
		},
		{
			name: "all Masters error",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("test error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	mastersService := services.NewMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mastersRepository, logger)
			}

			masters, err := mastersService.GetMasters(ctx, tc.pagination)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, masters, len(tc.expected))
			assert.Equal(t, tc.expected, masters)
		})
	}
}

func TestMastersService_RegisterMaster(t *testing.T) {
	testCases := []struct {
		name          string
		master        entities.RegisterMasterDTO
		expected      uint64
		setupMocks    func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "register Master success",
			master:   entities.RegisterMasterDTO{Info: pointers.New[string]("test"), UserID: 1},
			expected: 1,
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(nil, &customerrors.MasterNotFoundError{}).
					Times(1)

				mastersRepository.
					EXPECT().
					RegisterMaster(gomock.Any(), entities.RegisterMasterDTO{Info: pointers.New[string]("test"), UserID: 1}).
					Return(uint64(1), nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:   "register Master fail - already exists",
			master: entities.RegisterMasterDTO{Info: pointers.New[string]("test"), UserID: 1},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					GetMasterByUserID(gomock.Any(), uint64(1)).
					Return(&entities.Master{}, nil).
					Times(1)
			},
			errorExpected: true,
			err:           &customerrors.MasterAlreadyExistsError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	mastersService := services.NewMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mastersRepository, logger)
			}

			masterID, err := mastersService.RegisterMaster(ctx, tc.master)
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

func TestMastersService_UpdateMaster(t *testing.T) {
	testCases := []struct {
		name          string
		master        entities.UpdateMasterDTO
		setupMocks    func(mastersRepository *mockrepositories.MockMastersRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name: "update Master success",
			master: entities.UpdateMasterDTO{
				ID:   1,
				Info: pointers.New[string]("test"),
			},
			setupMocks: func(mastersRepository *mockrepositories.MockMastersRepository, _ *loggermock.MockLogger) {
				mastersRepository.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   1,
							Info: pointers.New[string]("test"),
						}).
					Return(nil).
					Times(1)
			},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	mastersService := services.NewMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(mastersRepository, logger)
			}

			err := mastersService.UpdateMaster(ctx, tc.master)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
