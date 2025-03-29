package tags

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
	tagID uint32 = 1
)

func TestTagsServer_GetTag(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetTagIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetTagOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetTagIn{
				ID: tagID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(
						&entities.Tag{
							ID:   tagID,
							Name: "test",
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetTagOut{
				ID:        tagID,
				Name:      "test",
				CreatedAt: timestamppb.New(time.Time{}),
				UpdatedAt: timestamppb.New(time.Time{}),
			},
		},
		{
			name: "Tag not found",
			in: &toys.GetTagIn{
				ID: tagID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(nil, &customerrors.TagNotFoundError{}).
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
			in: &toys.GetTagIn{
				ID: tagID,
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
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

			actual, err := tagsServer.GetTag(ctx, tc.in)
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

func TestTagsServer_GetTags(t *testing.T) {
	testCases := []struct {
		name          string
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetTagsOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(
						[]entities.Tag{
							{
								ID:   tagID,
								Name: "test",
							},
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetTagsOut{
				Tags: []*toys.GetTagOut{
					{
						ID:        tagID,
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
					GetAllTags(gomock.Any()).
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

			actual, err := tagsServer.GetTags(ctx, &emptypb.Empty{})
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

func TestTagsServer_CreateTags(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.CreateTagsIn
		setupMocks    func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.CreateTagsOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.CreateTagsIn{
				Tags: []*toys.CreateTagIn{
					{
						Name: "test",
					},
				},
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					CreateTags(
						gomock.Any(),
						[]entities.CreateTagDTO{
							{
								Name: "test",
							},
						},
					).
					Return([]uint32{tagID}, nil).
					Times(1)
			},
			expected: &toys.CreateTagsOut{
				Tags: []*toys.CreateTagOut{
					{
						ID: tagID,
					},
				},
			},
		},
		{
			name: "error",
			in: &toys.CreateTagsIn{
				Tags: []*toys.CreateTagIn{
					{
						Name: "test",
					},
				},
			},
			setupMocks: func(usecases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				usecases.
					EXPECT().
					CreateTags(
						gomock.Any(),
						[]entities.CreateTagDTO{
							{
								Name: "test",
							},
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

			actual, err := tagsServer.CreateTags(ctx, tc.in)
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
