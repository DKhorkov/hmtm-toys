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
		tagID         uint32
		resultLength  int
		errorExpected bool
		err           error
	}{
		{
			tagID:         1,
			resultLength:  1,
			errorExpected: false,
		},
		{
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
		tag, err := tagsService.GetTagByID(ctx, tc.tagID)
		if tc.errorExpected {
			require.Error(t, err)
			require.IsType(t, tc.err, err)
			assert.Nil(t, tag)
		} else {
			require.NoError(t, err)
		}
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
