package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
	loggermock "github.com/DKhorkov/libs/logging/mocks"
)

var (
	email         = "testUser@mail.ru"
	userID uint64 = 1
)

func TestSsoService_GetUserByID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint64
		expected      *entities.User
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "successfully got User by id",
			userID:   userID,
			expected: &entities.User{ID: userID},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, _ *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(&entities.User{ID: userID}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:   "failed to get User by id",
			userID: userID,
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("test error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	ssoService := services.NewSsoService(ssoRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			user, err := ssoService.GetUserByID(ctx, tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_GetUserByEmail(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		expected      *entities.User
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "successfully got User by email",
			email:    email,
			expected: &entities.User{ID: userID, Email: email},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, _ *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(&entities.User{ID: userID, Email: email}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "failed to get User by email",
			email: email,
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Return(nil, errors.New("test error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	ssoService := services.NewSsoService(ssoRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			user, err := ssoService.GetUserByEmail(ctx, tc.email)
			if tc.errorExpected {
				require.Error(t, err)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSsoService_GetAllUsers(t *testing.T) {
	testCases := []struct {
		name          string
		expected      []entities.User
		setupMocks    func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "all Users with existing Users",
			expected: []entities.User{{ID: userID}},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, _ *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(
						[]entities.User{
							{ID: userID},
						},
						nil,
					).
					Times(1)
			},
		},
		{
			name:     "all Users without existing Users",
			expected: []entities.User{},
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, _ *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return([]entities.User{}, nil).
					Times(1)
			},
		},
		{
			name: "all Users error",
			setupMocks: func(ssoRepository *mockrepositories.MockSsoRepository, logger *loggermock.MockLogger) {
				ssoRepository.
					EXPECT().
					GetAllUsers(gomock.Any()).
					Return(nil, errors.New("test error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	ssoRepository := mockrepositories.NewMockSsoRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	ssoService := services.NewSsoService(ssoRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(ssoRepository, logger)
			}

			users, err := ssoService.GetAllUsers(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, users, len(tc.expected))
			assert.Equal(t, tc.expected, users)
		})
	}
}
