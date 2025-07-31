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

func TestMastersServer_GetMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMasterIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetMasterOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMasterIn{
				ID: masterID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := mastersServer.GetMaster(ctx, tc.in)
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

func TestMastersServer_GetMasterByUser(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMasterByUserIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetMasterOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMasterByUserIn{
				UserID: userID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := mastersServer.GetMasterByUser(ctx, tc.in)
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

func TestMastersServer_GetMasters(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMastersIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
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
				Filters: &toys.MastersFilters{
					Search:              pointers.New("test"),
					CreatedAtOrderByAsc: pointers.New[bool](true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.MastersFilters{
							Search:              pointers.New("test"),
							CreatedAtOrderByAsc: pointers.New[bool](true),
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
				Filters: &toys.MastersFilters{
					Search:              pointers.New("test"),
					CreatedAtOrderByAsc: pointers.New[bool](true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasters(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.MastersFilters{
							Search:              pointers.New("test"),
							CreatedAtOrderByAsc: pointers.New[bool](true),
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := mastersServer.GetMasters(ctx, tc.in)
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

func TestMastersServer_CountMasters(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.CountMastersIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.CountOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.CountMastersIn{
				Filters: &toys.MastersFilters{
					Search:              pointers.New("test"),
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountMasters(
						gomock.Any(),
						&entities.MastersFilters{
							Search:              pointers.New("test"),
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(uint64(1), nil).
					Times(1)
			},
			expected: &toys.CountOut{
				Count: 1,
			},
		},
		{
			name: "error",
			in: &toys.CountMastersIn{
				Filters: &toys.MastersFilters{
					Search:              pointers.New("test"),
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountMasters(
						gomock.Any(),
						&entities.MastersFilters{
							Search:              pointers.New("test"),
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(uint64(0), errors.New("some error")).
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := mastersServer.CountMasters(ctx, tc.in)
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

func TestMastersServer_RegisterMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.RegisterMasterIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
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
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := mastersServer.RegisterMaster(ctx, tc.in)
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

func TestMastersServer_UpdateMaster(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.UpdateMasterIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.UpdateMasterIn{
				ID:   masterID,
				Info: pointers.New[string]("test"),
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
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
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	mastersServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			_, err := mastersServer.UpdateMaster(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
