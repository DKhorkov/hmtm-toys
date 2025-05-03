package masters

import (
	"context"
	"errors"
	"testing"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	mockusecases "github.com/DKhorkov/hmtm-toys/mocks/usecases"
	mocklogger "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ctx    = context.Background()
	master = &entities.Master{
		ID:        masterID,
		UserID:    masterID,
		Info:      pointers.New[string]("test"),
		CreatedAt: now,
		UpdatedAt: now,
	}
)

const (
	masterID uint64 = 1
	userID   uint64 = 1
)

func TestTagsServer_GetMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMasterIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetMasterOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMasterIn{
				ID: masterID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(master, nil).
					Times(1)
			},
			expected: mappedMaster,
		},
		{
			name: "Master not found",
			in: &toys.GetMasterIn{
				ID: masterID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, &customerrors.MasterNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.NotFound,
		},
		{
			name: "internal error",
			in: &toys.GetMasterIn{
				ID: masterID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, errors.New("some error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.Internal,
		},
	}

	ctrl := gomock.NewController(t)
	usecases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	tagsServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := tagsServer.GetMaster(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTagsServer_GetMasterByUser(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMasterByUserIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetMasterOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMasterByUserIn{
				UserID: userID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(master, nil).
					Times(1)
			},
			expected: mappedMaster,
		},
		{
			name: "Master not found",
			in: &toys.GetMasterByUserIn{
				UserID: userID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, &customerrors.MasterNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.NotFound,
		},
		{
			name: "internal error",
			in: &toys.GetMasterByUserIn{
				UserID: userID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, errors.New("some error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.Internal,
		},
	}

	ctrl := gomock.NewController(t)
	usecases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	tagsServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := tagsServer.GetMasterByUser(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTagsServer_GetMasters(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMastersIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetMastersOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMastersIn{
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
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
							*master,
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetMastersOut{
				Masters: []*toys.GetMasterOut{
					mappedMaster,
				},
			},
		},
		{
			name: "error",
			in: &toys.GetMastersIn{
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
					).
					Return(nil, errors.New("some error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.Internal,
		},
	}

	ctrl := gomock.NewController(t)
	usecases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	tagsServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := tagsServer.GetMasters(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTagsServer_RegisterMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.RegisterMasterIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.RegisterMasterOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.RegisterMasterIn{
				UserID: userID,
				Info:   pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						entities.RegisterMasterDTO{
							UserID: userID,
							Info:   pointers.New[string]("test"),
						},
					).
					Return(masterID, nil).
					Times(1)
			},
			expected: &toys.RegisterMasterOut{
				MasterID: masterID,
			},
		},
		{
			name: "Master already exists",
			in: &toys.RegisterMasterIn{
				UserID: userID,
				Info:   pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						entities.RegisterMasterDTO{
							UserID: userID,
							Info:   pointers.New[string]("test"),
						},
					).
					Return(uint64(0), &customerrors.MasterAlreadyExistsError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.AlreadyExists,
		},
		{
			name: "internal error",
			in: &toys.RegisterMasterIn{
				UserID: userID,
				Info:   pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						entities.RegisterMasterDTO{
							UserID: userID,
							Info:   pointers.New[string]("test"),
						},
					).
					Return(uint64(0), errors.New("test error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.Internal,
		},
	}

	ctrl := gomock.NewController(t)
	usecases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	tagsServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := tagsServer.RegisterMaster(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestTagsServer_UpdateMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.UpdateMasterIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.UpdateMasterIn{
				ID:   masterID,
				Info: pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   masterID,
							Info: pointers.New[string]("test"),
						},
					).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Master not found",
			in: &toys.UpdateMasterIn{
				ID:   masterID,
				Info: pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   masterID,
							Info: pointers.New[string]("test"),
						},
					).
					Return(&customerrors.MasterNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.NotFound,
		},
		{
			name: "internal error",
			in: &toys.UpdateMasterIn{
				ID:   masterID,
				Info: pointers.New[string]("test"),
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   masterID,
							Info: pointers.New[string]("test"),
						},
					).
					Return(errors.New("test error")).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.Internal,
		},
	}

	ctrl := gomock.NewController(t)
	usecases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	tagsServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			_, err := tagsServer.UpdateMaster(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
