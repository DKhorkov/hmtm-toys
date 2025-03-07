package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	loggermock "github.com/DKhorkov/libs/logging/mocks"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestTagsServiceGetTagByID(t *testing.T) {
	testCases := []struct {
		name          string
		tagID         uint32
		expected      *entities.Tag
		setupMocks    func(tagsRepository *mockrepositories.MockTagsRepository, logger *loggermock.MockLogger)
		errorExpected bool
		err           error
	}{
		{
			name:     "successfully got Tag by id",
			tagID:    1,
			expected: &entities.Tag{ID: 1},
			setupMocks: func(tagsRepository *mockrepositories.MockTagsRepository, _ *loggermock.MockLogger) {
				tagsRepository.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(1)).
					Return(&entities.Tag{ID: 1}, nil).
					Times(1)
			},
			errorExpected: false,
		},
		{
			name:  "failed to get Tag by id",
			tagID: 2,
			setupMocks: func(tagsRepository *mockrepositories.MockTagsRepository, logger *loggermock.MockLogger) {
				tagsRepository.
					EXPECT().
					GetTagByID(gomock.Any(), uint32(2)).
					Return(nil, &customerrors.TagNotFoundError{}).
					Times(1)

				logger.
					EXPECT().
					ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
			errorExpected: true,
			err:           &customerrors.TagNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	tagsService := services.NewTagsService(tagsRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(tagsRepository, logger)
			}

			tag, err := tagsService.GetTagByID(ctx, tc.tagID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, tag)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTagsServiceGetAllTags(t *testing.T) {
	testCases := []struct {
		name          string
		expected      []entities.Tag
		setupMocks    func(tagsRepository *mockrepositories.MockTagsRepository, logger *loggermock.MockLogger)
		errorExpected bool
	}{
		{
			name:     "all Tags with existing Tags",
			expected: []entities.Tag{{ID: 1}},
			setupMocks: func(tagsRepository *mockrepositories.MockTagsRepository, _ *loggermock.MockLogger) {
				tagsRepository.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(
						[]entities.Tag{
							{ID: 1},
						},
						nil,
					).
					Times(1)
			},
		},
		{
			name:     "all Tags without existing Tags",
			expected: []entities.Tag{},
			setupMocks: func(tagsRepository *mockrepositories.MockTagsRepository, _ *loggermock.MockLogger) {
				tagsRepository.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return([]entities.Tag{}, nil).
					Times(1)
			},
		},
		{
			name: "all Tags error",
			setupMocks: func(tagsRepository *mockrepositories.MockTagsRepository, _ *loggermock.MockLogger) {
				tagsRepository.
					EXPECT().
					GetAllTags(gomock.Any()).
					Return(nil, errors.New("test error")).
					Times(1)
			},
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	tagsService := services.NewTagsService(tagsRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(tagsRepository, logger)
			}

			tags, err := tagsService.GetAllTags(ctx)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, tags, len(tc.expected))
			assert.Equal(t, tc.expected, tags)
		})
	}
}

func TestTagsServiceCreateTags(t *testing.T) {
	testCases := []struct {
		name           string
		tagsData       []entities.CreateTagDTO
		expectedResult []uint32
		errorExpected  bool
		setupMocks     func(tagsRepo *mockrepositories.MockTagsRepository)
	}{
		{
			name: "successfully created Tags",
			tagsData: []entities.CreateTagDTO{
				{
					Name: "test",
				},
				{
					Name: "test2",
				},
			},
			expectedResult: []uint32{1, 2},
			errorExpected:  false,
			setupMocks: func(tagsRepo *mockrepositories.MockTagsRepository) {
				tagsRepo.
					EXPECT().
					CreateTags(
						gomock.Any(),
						[]entities.CreateTagDTO{
							{
								Name: "test",
							},
							{
								Name: "test2",
							},
						},
					).
					Return(
						[]uint32{1, 2},
						nil,
					).
					Times(1)
			},
		},
	}

	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	logger := loggermock.NewMockLogger(mockController)
	tagsService := services.NewTagsService(tagsRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupMocks != nil {
				tc.setupMocks(tagsRepository)
			}

			tagIDs, err := tagsService.CreateTags(ctx, tc.tagsData)
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResult, tagIDs)
		})
	}
}
