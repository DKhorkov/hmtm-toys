package categories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DKhorkov/hmtm-toys/api/protobuf/generated/go/toys"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	mockusecases "github.com/DKhorkov/hmtm-toys/mocks/usecases"
	mocklogger "github.com/DKhorkov/libs/logging/mocks"
)

var (
	ctx = context.Background()
)

const (
	categoryID uint32 = 1
)

func TestCategoriesServer_GetCategory(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetCategoryIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetCategoryOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetCategoryIn{
				ID: categoryID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetCategoryOut{
				ID:        categoryID,
				Name:      "test",
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			},
		},
		{
			name: "Category not found",
			in: &toys.GetCategoryIn{
				ID: categoryID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(nil, &customerrors.CategoryNotFoundError{}).
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
			in: &toys.GetCategoryIn{
				ID: categoryID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
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
	categoriesServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := categoriesServer.GetCategory(ctx, tc.in)
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

func TestCategoriesServer_GetCategories(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetCategoriesOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetAllCategories(gomock.Any()).
					Return(
						[]entities.Category{
							{
								ID:   categoryID,
								Name: "test",
							},
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetCategoriesOut{
				Categories: []*toys.GetCategoryOut{
					{
						ID:        categoryID,
						Name:      "test",
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
		},
		{
			name: "error",
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetAllCategories(gomock.Any()).
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
	categoriesServer := &ServerAPI{
		logger:   logger,
		useCases: usecases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(usecases, logger)
			}

			actual, err := categoriesServer.GetCategories(ctx, &emptypb.Empty{})
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
