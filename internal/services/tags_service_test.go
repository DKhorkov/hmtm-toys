package services_test

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/services"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
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
	tagsRepository.EXPECT().GetTagByID(uint32(1)).Return(&entities.Tag{}, nil).MaxTimes(1)
	tagsRepository.EXPECT().GetTagByID(uint32(2)).Return(
		nil,
		&customerrors.TagNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	tagsService := services.NewCommonTagsService(tagsRepository, logger)

	for _, tc := range testCases {
		tag, err := tagsService.GetTagByID(tc.tagID)
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
		expectedTags := []*entities.Tag{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags().DoAndReturn(
			func() ([]*entities.Tag, error) {
				return expectedTags, nil
			},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)

		tags, err := tagsService.GetAllTags()
		require.NoError(t, err)
		assert.Len(t, tags, len(expectedTags))
		assert.Equal(t, expectedTags, tags)
	})

	t.Run("all tags without existing tags", func(t *testing.T) {
		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags().Return([]*entities.Tag{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)

		tags, err := tagsService.GetAllTags()
		require.NoError(t, err)
		assert.Empty(t, tags)
	})

	t.Run("all tags fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetAllTags().Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		tagsService := services.NewCommonTagsService(tagsRepository, logger)

		tags, err := tagsService.GetAllTags()
		require.Error(t, err)
		assert.Nil(t, tags)
	})
}
