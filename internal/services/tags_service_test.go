package services_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestCommonTagsServiceGetTagByID(t *testing.T) {
	testCases := []struct {
		name          string
		tagID         uint32
		errorExpected bool
		err           error
	}{
		{
			name:          "successfully got Tag by id",
			tagID:         1,
			errorExpected: false,
		},
		{
			name:          "failed to get Tag by id",
			tagID:         2,
			errorExpected: true,
			err:           &customerrors.TagNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetTagByID(gomock.Any(), uint32(1)).Return(&entities.Tag{}, nil).MaxTimes(1)
	tagsRepository.EXPECT().GetTagByID(gomock.Any(), uint32(2)).Return(
		nil,
		&customerrors.TagNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	tagsService := services.NewCommonTagsService(tagsRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

func TestCommonTagsServiceGetAllTags(t *testing.T) {
	t.Run("all tags with existing tags", func(t *testing.T) {
		expectedTags := []entities.Tag{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags(gomock.Any()).Return(expectedTags, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)
		ctx := context.Background()

		tags, err := tagsService.GetAllTags(ctx)
		require.NoError(t, err)
		assert.Len(t, tags, len(expectedTags))
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("all tags without existing tags", func(t *testing.T) {
		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags(gomock.Any()).Return([]entities.Tag{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)
		ctx := context.Background()

		tags, err := tagsService.GetAllTags(ctx)
		require.NoError(t, err)
		assert.Empty(t, tags)
	})

	t.Run("all tags fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags(gomock.Any()).Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)
		ctx := context.Background()

		tags, err := tagsService.GetAllTags(ctx)
		require.Error(t, err)
		assert.Nil(t, tags)
	})
}

func TestCommonTagsServiceCreateTags(t *testing.T) {
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
					MaxTimes(1)
			},
		},
	}

	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	tagsService := services.NewCommonTagsService(tagsRepository, logger)
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
