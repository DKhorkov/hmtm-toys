package toys

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
	ctx = context.Background()
	toy = &entities.Toy{
		ID:          toyID,
		MasterID:    masterID,
		CategoryID:  categoryID,
		Name:        "test toy",
		Description: "test description",
		Quantity:    1,
		Price:       110,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tags: []entities.Tag{
			{
				ID:   tagID,
				Name: "test tag",
			},
		},
		Attachments: []entities.Attachment{
			{
				ID:        attachmentID,
				Link:      "https://example.com/attachment",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
)

const (
	toyID        uint64 = 1
	masterID     uint64 = 1
	categoryID   uint32 = 1
	tagID        uint32 = 1
	attachmentID uint64 = 1
	userID       uint64 = 1
)

func TestToysServer_GetToy(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetToyIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetToyOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(toy, nil).
					Times(1)
			},
			expected: mappedToy,
		},
		{
			name: "Toy not found",
			in: &toys.GetToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, &customerrors.ToyNotFoundError{}).
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
			in: &toys.GetToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.GetToy(ctx, tc.in)
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

func TestToysServer_GetToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetToysOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetToysIn{
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(
						[]entities.Toy{
							*toy,
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetToysOut{
				Toys: []*toys.GetToyOut{
					mappedToy,
				},
			},
		},
		{
			name: "error",
			in: &toys.GetToysIn{
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetToys(
						gomock.Any(),
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.GetToys(ctx, tc.in)
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

func TestToysServer_CountToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.CountToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.CountOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.CountToysIn{
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountToys(
						gomock.Any(),
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
			in: &toys.CountToysIn{
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountToys(
						gomock.Any(),
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.CountToys(ctx, tc.in)
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

func TestToysServer_CountMasterToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.CountMasterToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.CountOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.CountMasterToysIn{
				MasterID: masterID,
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountMasterToys(
						gomock.Any(),
						masterID,
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
			in: &toys.CountMasterToysIn{
				MasterID: masterID,
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountMasterToys(
						gomock.Any(),
						masterID,
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.CountMasterToys(ctx, tc.in)
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

func TestToysServer_CountUserToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.CountUserToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.CountOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.CountUserToysIn{
				UserID: userID,
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountUserToys(
						gomock.Any(),
						userID,
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
			in: &toys.CountUserToysIn{
				UserID: userID,
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					CountUserToys(
						gomock.Any(),
						userID,
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.CountUserToys(ctx, tc.in)
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

func TestToysServer_GetMasterToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetMasterToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetToysOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetMasterToysIn{
				MasterID: masterID,
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						masterID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(
						[]entities.Toy{
							*toy,
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetToysOut{
				Toys: []*toys.GetToyOut{
					mappedToy,
				},
			},
		},
		{
			name: "error",
			in: &toys.GetMasterToysIn{
				MasterID: masterID,
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetMasterToys(
						gomock.Any(),
						masterID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.GetMasterToys(ctx, tc.in)
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

func TestToysServer_GetUserToys(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.GetUserToysIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.GetToysOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.GetUserToysIn{
				UserID: userID,
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:              pointers.New("toy2"),
					PriceCeil:           pointers.New[float32](1000),
					PriceFloor:          pointers.New[float32](10),
					QuantityFloor:       pointers.New[uint32](1),
					CategoryIDs:         []uint32{1},
					TagIDs:              []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						userID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:        pointers.New("toy2"),
							PriceCeil:     pointers.New[float32](1000),
							PriceFloor:    pointers.New[float32](10),
							QuantityFloor: pointers.New[uint32](1),
							CategoryIDs:   []uint32{1}, TagIDs: []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
						},
					).
					Return(
						[]entities.Toy{
							*toy,
						},
						nil,
					).
					Times(1)
			},
			expected: &toys.GetToysOut{
				Toys: []*toys.GetToyOut{
					mappedToy,
				},
			},
		},
		{
			name: "error",
			in: &toys.GetUserToysIn{
				UserID: userID,
				Pagination: &toys.Pagination{
					Limit:  pointers.New[uint64](1),
					Offset: pointers.New[uint64](1),
				},
				Filters: &toys.ToysFilters{
					Search:        pointers.New("toy2"),
					PriceCeil:     pointers.New[float32](1000),
					PriceFloor:    pointers.New[float32](10),
					QuantityFloor: pointers.New[uint32](1),
					CategoryIDs:   []uint32{1}, TagIDs: []uint32{1},
					CreatedAtOrderByAsc: pointers.New(true),
				},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					GetUserToys(
						gomock.Any(),
						userID,
						&entities.Pagination{
							Limit:  pointers.New[uint64](1),
							Offset: pointers.New[uint64](1),
						},
						&entities.ToysFilters{
							Search:              pointers.New("toy2"),
							PriceCeil:           pointers.New[float32](1000),
							PriceFloor:          pointers.New[float32](10),
							QuantityFloor:       pointers.New[uint32](1),
							CategoryIDs:         []uint32{1},
							TagIDs:              []uint32{1},
							CreatedAtOrderByAsc: pointers.New(true),
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.GetUserToys(ctx, tc.in)
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

func TestToysServer_AddToy(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.AddToyIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		expected      *toys.AddToyOut
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(toyID, nil).
					Times(1)
			},
			expected: &toys.AddToyOut{
				ToyID: toyID,
			},
		},
		{
			name: "Toy already exists",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(uint64(0), &customerrors.ToyAlreadyExistsError{}).
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
			name: "Tag not found",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(uint64(0), &customerrors.TagNotFoundError{}).
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
			name: "Category not found",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(uint64(0), &customerrors.CategoryNotFoundError{}).
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
			name: "Master not found",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(uint64(0), &customerrors.MasterNotFoundError{}).
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
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
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
		{
			name: "validation error",
			in: &toys.AddToyIn{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test toy",
				Description: "test description",
				Quantity:    1,
				Price:       110,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.RawAddToyDTO{
							UserID:      userID,
							CategoryID:  categoryID,
							Name:        "test toy",
							Description: "test description",
							Quantity:    1,
							Price:       110,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(uint64(0), validationError).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.FailedPrecondition,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			actual, err := toysServer.AddToy(ctx, tc.in)
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

func TestToysServer_DeleteToy(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.DeleteToyIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.DeleteToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					DeleteToy(gomock.Any(), toyID).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Toy not found",
			in: &toys.DeleteToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					DeleteToy(gomock.Any(), toyID).
					Return(&customerrors.ToyNotFoundError{}).
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
			in: &toys.DeleteToyIn{
				ID: toyID,
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					DeleteToy(gomock.Any(), toyID).
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
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			_, err := toysServer.DeleteToy(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestToysServer_UpdateToy(t *testing.T) {
	testCases := []struct {
		name          string
		in            *toys.UpdateToyIn
		setupMocks    func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger)
		errorExpected bool
		errorCode     codes.Code
	}{
		{
			name: "success",
			in: &toys.UpdateToyIn{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("test toy"),
				Description: pointers.New[string]("test description"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110),
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, _ *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.RawUpdateToyDTO{
							ID:          toyID,
							CategoryID:  pointers.New[uint32](categoryID),
							Name:        pointers.New[string]("test toy"),
							Description: pointers.New[string]("test description"),
							Quantity:    pointers.New[uint32](1),
							Price:       pointers.New[float32](110),
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Toy not found",
			in: &toys.UpdateToyIn{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("test toy"),
				Description: pointers.New[string]("test description"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110),
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.RawUpdateToyDTO{
							ID:          toyID,
							CategoryID:  pointers.New[uint32](categoryID),
							Name:        pointers.New[string]("test toy"),
							Description: pointers.New[string]("test description"),
							Quantity:    pointers.New[uint32](1),
							Price:       pointers.New[float32](110),
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(&customerrors.ToyNotFoundError{}).
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
			in: &toys.UpdateToyIn{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("test toy"),
				Description: pointers.New[string]("test description"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110),
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.RawUpdateToyDTO{
							ID:          toyID,
							CategoryID:  pointers.New[uint32](categoryID),
							Name:        pointers.New[string]("test toy"),
							Description: pointers.New[string]("test description"),
							Quantity:    pointers.New[uint32](1),
							Price:       pointers.New[float32](110),
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
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
		{
			name: "validation error",
			in: &toys.UpdateToyIn{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("test toy"),
				Description: pointers.New[string]("test description"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110),
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test attachment"},
			},
			setupMocks: func(useCases *mockusecases.MockUseCases, logger *mocklogger.MockLogger) {
				useCases.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.RawUpdateToyDTO{
							ID:          toyID,
							CategoryID:  pointers.New[uint32](categoryID),
							Name:        pointers.New[string]("test toy"),
							Description: pointers.New[string]("test description"),
							Quantity:    pointers.New[uint32](1),
							Price:       pointers.New[float32](110),
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test attachment"},
						},
					).
					Return(validationError).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			errorCode:     codes.FailedPrecondition,
		},
	}

	ctrl := gomock.NewController(t)
	useCases := mockusecases.NewMockUseCases(ctrl)
	logger := mocklogger.NewMockLogger(ctrl)
	toysServer := &ServerAPI{
		logger:   logger,
		useCases: useCases,
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(useCases, logger)
			}

			_, err := toysServer.UpdateToy(ctx, tc.in)
			if tc.errorExpected {
				require.Error(t, err)
				require.Equal(t, tc.errorCode, status.Code(err))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
