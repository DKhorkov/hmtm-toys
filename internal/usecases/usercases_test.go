package usecases

import (
	"context"
	"errors"
	"github.com/DKhorkov/hmtm-toys/internal/config"
	"github.com/DKhorkov/libs/validation"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/pointers"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	mockservices "github.com/DKhorkov/hmtm-toys/mocks/services"
)

const (
	tagID      uint32 = 1
	categoryID uint32 = 1
	toyID      uint64 = 1
	masterID   uint64 = 1
	userID     uint64 = 1
)

var (
	ctx              = context.Background()
	cfg              = config.New()
	validationConfig = cfg.Validation
)

func TestUseCases_GetTagByID(t *testing.T) {
	testCases := []struct {
		name       string
		tagID      uint32
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      *entities.Tag
		errorExpected bool
	}{
		{
			name:  "success",
			tagID: tagID,
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(
						&entities.Tag{
							ID:   tagID,
							Name: "test",
						}, nil,
					).
					Times(1)
			},
			expected: &entities.Tag{
				ID:   tagID,
				Name: "test",
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetTagByID(ctx, tc.tagID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllTags(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Tag
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				tagsService.
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
			expected: []entities.Tag{
				{
					ID:   tagID,
					Name: "test",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetAllTags(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetCategoryByID(t *testing.T) {
	testCases := []struct {
		name       string
		categoryID uint32
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      *entities.Category
		errorExpected bool
	}{
		{
			name:       "success",
			categoryID: categoryID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						}, nil,
					).
					Times(1)
			},
			expected: &entities.Category{
				ID:   categoryID,
				Name: "test",
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetCategoryByID(ctx, tc.categoryID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetAllCategories(t *testing.T) {
	testCases := []struct {
		name       string
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Category
		errorExpected bool
	}{
		{
			name: "success",
			setupMocks: func(
				_ *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				categoriesService.
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
			expected: []entities.Category{
				{
					ID:   categoryID,
					Name: "test",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetAllCategories(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetToyByID(t *testing.T) {
	testCases := []struct {
		name       string
		toyID      uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      *entities.Toy
		errorExpected bool
	}{
		{
			name:  "success",
			toyID: toyID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(
						&entities.Toy{
							ID:   toyID,
							Name: "test",
						}, nil,
					).
					Times(1)
			},
			expected: &entities.Toy{
				ID:   toyID,
				Name: "test",
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetToyByID(ctx, tc.toyID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetToys(t *testing.T) {
	testCases := []struct {
		name       string
		pagination *entities.Pagination
		filters    *entities.ToysFilters
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
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
							{
								ID:   toyID,
								Name: "test",
							},
						},
						nil,
					).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:   toyID,
					Name: "test",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetToys(ctx, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_CountToys(t *testing.T) {
	testCases := []struct {
		name       string
		filters    *entities.ToysFilters
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
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
			expected: 1,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.CountToys(ctx, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_CountMasterToys(t *testing.T) {
	testCases := []struct {
		name       string
		masterID   uint64
		filters    *entities.ToysFilters
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name:     "success",
			masterID: masterID,
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
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
			expected: 1,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.CountMasterToys(ctx, tc.masterID, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_CountUserToys(t *testing.T) {
	testCases := []struct {
		name       string
		userID     uint64
		filters    *entities.ToysFilters
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name:   "success",
			userID: userID,
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{ID: masterID, UserID: userID},
						nil,
					).
					Times(1)

				toysService.
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
			expected: 1,
		},
		{
			name:   "master with provided userID not found",
			userID: userID,
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.CountUserToys(ctx, tc.userID, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMasterToys(t *testing.T) {
	testCases := []struct {
		name       string
		pagination *entities.Pagination
		filters    *entities.ToysFilters
		masterID   uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			masterID: masterID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
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
							{
								ID:   toyID,
								Name: "test",
							},
						},
						nil,
					).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:   toyID,
					Name: "test",
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetMasterToys(ctx, tc.masterID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetUserToys(t *testing.T) {
	testCases := []struct {
		name       string
		pagination *entities.Pagination
		filters    *entities.ToysFilters
		userID     uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Toy
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			filters: &entities.ToysFilters{
				Search:              pointers.New("toy2"),
				PriceCeil:           pointers.New[float32](1000),
				PriceFloor:          pointers.New[float32](10),
				QuantityFloor:       pointers.New[uint32](1),
				CategoryIDs:         []uint32{1},
				TagIDs:              []uint32{1},
				CreatedAtOrderByAsc: pointers.New(true),
			},
			userID: userID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{
							ID: masterID,
						},
						nil,
					).
					Times(1)

				toysService.
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
							{
								ID:   toyID,
								Name: "test",
							},
						},
						nil,
					).
					Times(1)
			},
			expected: []entities.Toy{
				{
					ID:   toyID,
					Name: "test",
				},
			},
		},
		{
			name: "master not found",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			userID: userID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, errors.New("master not found")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetUserToys(ctx, tc.userID, tc.pagination, tc.filters)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMasterByID(t *testing.T) {
	testCases := []struct {
		name       string
		masterID   uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name:     "success",
			masterID: masterID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(
						&entities.Master{
							ID:   masterID,
							Info: pointers.New[string]("test"),
						}, nil,
					).
					Times(1)
			},
			expected: &entities.Master{
				ID:   masterID,
				Info: pointers.New[string]("test"),
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetMasterByID(ctx, tc.masterID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMasters(t *testing.T) {
	testCases := []struct {
		name       string
		pagination *entities.Pagination
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []entities.Master
		errorExpected bool
	}{
		{
			name: "success",
			pagination: &entities.Pagination{
				Limit:  pointers.New[uint64](1),
				Offset: pointers.New[uint64](1),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
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
							{
								ID:   masterID,
								Info: pointers.New[string]("test"),
							},
						},
						nil,
					).
					Times(1)
			},
			expected: []entities.Master{
				{
					ID:   masterID,
					Info: pointers.New[string]("test"),
				},
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetMasters(ctx, tc.pagination)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_GetMasterByUserID(t *testing.T) {
	testCases := []struct {
		name       string
		userID     uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      *entities.Master
		errorExpected bool
	}{
		{
			name:   "success",
			userID: userID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{
							ID:   toyID,
							Info: pointers.New[string]("test"),
						},
						nil,
					).
					Times(1)
			},
			expected: &entities.Master{
				ID:   toyID,
				Info: pointers.New[string]("test"),
			},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.GetMasterByUserID(ctx, tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_AddToy(t *testing.T) {
	testCases := []struct {
		name       string
		toy        entities.RawAddToyDTO
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      uint64
		errorExpected bool
	}{
		{
			name: "success",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       110.5,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{
							ID:     masterID,
							UserID: userID,
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						}, nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(
						&entities.Tag{
							ID:   tagID,
							Name: "test",
						}, nil,
					).
					Times(1)

				toysService.
					EXPECT().
					AddToy(
						gomock.Any(),
						entities.AddToyDTO{
							MasterID:    masterID,
							CategoryID:  categoryID,
							Name:        "Игрушка",
							Description: "Тестовая игрушка",
							Quantity:    1,
							Price:       110.5,
							TagIDs:      []uint32{tagID},
							Attachments: []string{"test"},
						},
					).
					Return(toyID, nil).
					Times(1)
			},
			expected: toyID,
		},
		{
			name: "Master not found",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       110.5,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "Category not found",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       110.5,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{
							ID:     masterID,
							UserID: userID,
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "Tag not found",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       110.5,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByUserID(gomock.Any(), userID).
					Return(
						&entities.Master{
							ID:     masterID,
							UserID: userID,
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						}, nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "invalid quantity",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1_000_000,
				Price:       110.5,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "invalid price",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Игрушка",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       1_000_000_000,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "invalid name",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "Мразь",
				Description: "Тестовая игрушка",
				Quantity:    1,
				Price:       1_000_000_000,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			expected:      0,
			errorExpected: true,
		},
		{
			name: "invalid description",
			toy: entities.RawAddToyDTO{
				UserID:      userID,
				CategoryID:  categoryID,
				Name:        "test",
				Description: "Конченная мразь",
				Quantity:    1,
				Price:       1_000_000_000,
				TagIDs:      []uint32{tagID},
				Attachments: []string{"test"},
			},
			expected:      0,
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.AddToy(ctx, tc.toy)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_RegisterMaster(t *testing.T) {
	testCases := []struct {
		name       string
		master     entities.RegisterMasterDTO
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      uint64
		errorExpected bool
		expectedError error
	}{
		{
			name: "success",
			master: entities.RegisterMasterDTO{
				UserID: userID,
				Info:   pointers.New[string]("Мастер о себе"),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				ssoService *mockservices.MockSsoService,
			) {
				ssoService.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(
						&entities.User{
							ID: userID,
						},
						nil,
					).
					Times(1)

				mastersService.
					EXPECT().
					RegisterMaster(
						gomock.Any(),
						entities.RegisterMasterDTO{
							UserID: userID,
							Info:   pointers.New[string]("Мастер о себе"),
						},
					).
					Return(masterID, nil).
					Times(1)
			},
			expected: masterID,
		},
		{
			name: "User not found",
			master: entities.RegisterMasterDTO{
				UserID: userID,
				Info:   pointers.New[string]("Мастер"),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				ssoService *mockservices.MockSsoService,
			) {
				ssoService.
					EXPECT().
					GetUserByID(gomock.Any(), userID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "Invalid master info",
			master: entities.RegisterMasterDTO{
				Info: pointers.New[string]("Сука"),
			},
			errorExpected: true,
			expectedError: &validation.Error{},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.RegisterMaster(ctx, tc.master)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_UpdateMaster(t *testing.T) {
	testCases := []struct {
		name       string
		master     entities.UpdateMasterDTO
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		errorExpected bool
		expectedError error
	}{
		{
			name: "success",
			master: entities.UpdateMasterDTO{
				ID:   masterID,
				Info: pointers.New[string]("Мастер о себе"),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(
						&entities.Master{
							ID:     masterID,
							UserID: userID,
							Info:   pointers.New[string]("Какая-то инфа"),
						},
						nil,
					).
					Times(1)

				mastersService.
					EXPECT().
					UpdateMaster(
						gomock.Any(),
						entities.UpdateMasterDTO{
							ID:   masterID,
							Info: pointers.New[string]("Мастер о себе"),
						},
					).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Master not found",
			master: entities.UpdateMasterDTO{
				ID:   masterID,
				Info: pointers.New[string]("Мастер о себе"),
			},
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				mastersService.
					EXPECT().
					GetMasterByID(gomock.Any(), masterID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "Invalid master info",
			master: entities.UpdateMasterDTO{
				ID:   masterID,
				Info: pointers.New[string]("invalid master info that would not work"),
			},
			errorExpected: true,
			expectedError: &validation.Error{},
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			err := useCases.UpdateMaster(ctx, tc.master)
			if tc.errorExpected {
				require.Error(t, err)
				if tc.expectedError != nil {
					require.IsType(t, tc.expectedError, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_DeleteToy(t *testing.T) {
	testCases := []struct {
		name       string
		toyID      uint64
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		errorExpected bool
	}{
		{
			name:  "success",
			toyID: toyID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(
						&entities.Toy{
							ID:   toyID,
							Name: "test",
						},
						nil,
					).
					Times(1)

				toysService.
					EXPECT().
					DeleteToy(gomock.Any(), toyID).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "Toy not found",
			toyID: toyID,
			setupMocks: func(
				_ *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, errors.New("test")).
					Times(1)

			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			err := useCases.DeleteToy(ctx, tc.toyID)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUseCases_CreateTags(t *testing.T) {
	testCases := []struct {
		name       string
		tags       []entities.CreateTagDTO
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		expected      []uint32
		errorExpected bool
	}{
		{
			name: "success",
			tags: []entities.CreateTagDTO{
				{
					Name: "тестовыйТег",
				},
				{
					Name: "новыйТег",
				},
				{
					Name: "новыйТег",
				},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				tagsService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(
						[]entities.Tag{
							{
								ID:   tagID,
								Name: "тестовыйтег",
							},
						},
						nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					CreateTags(
						gomock.Any(),
						[]entities.CreateTagDTO{
							{
								Name: "новыйтег",
							},
						},
					).
					Return([]uint32{2, 3}, nil).
					Times(1)
			},
			expected: []uint32{2, 3, tagID},
		},
		{
			name: "get all Tags error",
			tags: []entities.CreateTagDTO{
				{
					Name: "тестовыйТег",
				},
				{
					Name: "новыйТег",
				},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				tagsService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "create Tags error",
			tags: []entities.CreateTagDTO{
				{
					Name: "тестовыйТег",
				},
				{
					Name: "новыйТег",
				},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				_ *mockservices.MockCategoriesService,
				_ *mockservices.MockMastersService,
				_ *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				tagsService.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(
						[]entities.Tag{
							{
								ID:   tagID,
								Name: "тестовыйтег",
							},
						},
						nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					CreateTags(
						gomock.Any(),
						[]entities.CreateTagDTO{
							{
								Name: "новыйтег",
							},
						},
					).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "validation error",
			tags: []entities.CreateTagDTO{
				{
					Name: "тестовыйТег",
				},
				{
					Name: "Сука",
				},
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			actual, err := useCases.CreateTags(ctx, tc.tags)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Len(t, actual, len(tc.expected))
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestUseCases_UpdateToy(t *testing.T) {
	testCases := []struct {
		name       string
		toy        entities.RawUpdateToyDTO
		setupMocks func(
			tagsService *mockservices.MockTagsService,
			categoriesService *mockservices.MockCategoriesService,
			mastersService *mockservices.MockMastersService,
			toysService *mockservices.MockToysService,
			ssoService *mockservices.MockSsoService,
		)
		errorExpected bool
	}{
		{
			name: "success",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(
						&entities.Toy{
							ID:          toyID,
							CategoryID:  categoryID,
							Name:        "Какая-то игрушка",
							Description: "Какое-то описание",
							Quantity:    1,
							Price:       110.5,
							Tags: []entities.Tag{
								{
									ID:   tagID,
									Name: "test",
								},
								{
									ID:   3,
									Name: "tagToDelete",
								},
							},
							Attachments: []entities.Attachment{
								{
									ID:   1,
									Link: "oldAttachment",
								},
								{
									ID:   2,
									Link: "attachmentToDelete",
								},
							},
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						}, nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(
						&entities.Tag{
							ID:   tagID,
							Name: "test",
						}, nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(2)).
					Return(
						&entities.Tag{
							ID:   2,
							Name: "tagToAdd",
						}, nil,
					).
					Times(1)

				toysService.
					EXPECT().
					UpdateToy(
						gomock.Any(),
						entities.UpdateToyDTO{
							ID:                    toyID,
							CategoryID:            pointers.New[uint32](categoryID),
							Name:                  pointers.New[string]("Игрушка"),
							Description:           pointers.New[string]("Тестовая игрушка"),
							Quantity:              pointers.New[uint32](1),
							Price:                 pointers.New[float32](110.5),
							TagIDsToAdd:           []uint32{2},
							TagIDsToDelete:        []uint32{3},
							AttachmentsToAdd:      []string{"newAttachment"},
							AttachmentIDsToDelete: []uint64{2},
						},
					).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "Tag not found",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(
						&entities.Toy{
							ID:          toyID,
							CategoryID:  categoryID,
							Name:        "Какая-то игрушка",
							Description: "Какое-то описание",
							Quantity:    1,
							Price:       110.5,
							Tags: []entities.Tag{
								{
									ID:   tagID,
									Name: "test",
								},
								{
									ID:   3,
									Name: "tagToDelete",
								},
							},
							Attachments: []entities.Attachment{
								{
									ID:   1,
									Link: "oldAttachment",
								},
								{
									ID:   2,
									Link: "attachmentToDelete",
								},
							},
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(
						&entities.Category{
							ID:   categoryID,
							Name: "test",
						}, nil,
					).
					Times(1)

				tagsService.
					EXPECT().
					GetTagByID(gomock.Any(), tagID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "Category not found",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(
						&entities.Toy{
							ID:          toyID,
							CategoryID:  categoryID,
							Name:        "Какая-то игрушка",
							Description: "Какое-то описание",
							Quantity:    1,
							Price:       110.5,
							Tags: []entities.Tag{
								{
									ID:   tagID,
									Name: "test",
								},
								{
									ID:   3,
									Name: "tagToDelete",
								},
							},
							Attachments: []entities.Attachment{
								{
									ID:   1,
									Link: "oldAttachment",
								},
								{
									ID:   2,
									Link: "attachmentToDelete",
								},
							},
						},
						nil,
					).
					Times(1)

				categoriesService.
					EXPECT().
					GetCategoryByID(gomock.Any(), categoryID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "Toy not found",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			setupMocks: func(
				tagsService *mockservices.MockTagsService,
				categoriesService *mockservices.MockCategoriesService,
				mastersService *mockservices.MockMastersService,
				toysService *mockservices.MockToysService,
				_ *mockservices.MockSsoService,
			) {
				toysService.
					EXPECT().
					GetToyByID(gomock.Any(), toyID).
					Return(nil, errors.New("test")).
					Times(1)
			},
			errorExpected: true,
		},
		{
			name: "invalid quantity",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1_000_000),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			errorExpected: true,
		},
		{
			name: "invalid price",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](1_000_000_000),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			errorExpected: true,
		},
		{
			name: "invalid name",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Мразь"),
				Description: pointers.New[string]("Тестовая игрушка"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			errorExpected: true,
		},
		{
			name: "invalid description",
			toy: entities.RawUpdateToyDTO{
				ID:          toyID,
				CategoryID:  pointers.New[uint32](categoryID),
				Name:        pointers.New[string]("Игрушка"),
				Description: pointers.New[string]("Сука"),
				Quantity:    pointers.New[uint32](1),
				Price:       pointers.New[float32](110.5),
				TagIDs:      []uint32{tagID, 2},
				Attachments: []string{"oldAttachment", "newAttachment"},
			},
			errorExpected: true,
		},
	}

	ctrl := gomock.NewController(t)
	tagsService := mockservices.NewMockTagsService(ctrl)
	categoriesService := mockservices.NewMockCategoriesService(ctrl)
	mastersService := mockservices.NewMockMastersService(ctrl)
	toysService := mockservices.NewMockToysService(ctrl)
	ssoService := mockservices.NewMockSsoService(ctrl)
	useCases := New(
		tagsService,
		categoriesService,
		mastersService,
		toysService,
		ssoService,
		validationConfig,
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(
					tagsService,
					categoriesService,
					mastersService,
					toysService,
					ssoService,
				)
			}

			err := useCases.UpdateToy(ctx, tc.toy)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
