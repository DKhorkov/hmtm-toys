package services

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/DKhorkov/hmtm-toys/internal/interfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
)

func TestCommonToysServiceGetToyByID(t *testing.T) {
	testCases := []struct {
		toyID         uint64
		resultLength  int
		errorExpected bool
	}{
		{
			toyID:         1,
			resultLength:  1,
			errorExpected: false,
		},
		{
			toyID:         2,
			errorExpected: true,
		},
	}

	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	toysRepository.EXPECT().GetToyByID(gomock.Any(), uint64(1)).Return(&entities.Toy{ID: 1}, nil).MaxTimes(1)
	toysRepository.EXPECT().GetToyByID(gomock.Any(), uint64(2)).Return(
		nil,
		&customerrors.ToyNotFoundError{},
	).MaxTimes(1)

	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetToyTags(gomock.Any(), uint64(1)).Return([]entities.Tag{}, nil).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	toysService := NewCommonToysService(toysRepository, tagsRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		toy, err := toysService.GetToyByID(ctx, tc.toyID)
		if tc.errorExpected {
			require.Error(t, err)
			assert.Nil(t, toy)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestCommonToysServiceGetAllToys(t *testing.T) {
	t.Run("all toys with existing toys", func(t *testing.T) {
		expectedTags := []entities.Tag{
			{ID: 1},
		}

		expectedToys := []entities.Toy{
			{
				ID:   1,
				Tags: expectedTags,
			},
		}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys(gomock.Any()).Return(expectedToys, nil).MaxTimes(1)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetToyTags(gomock.Any(), uint64(1)).Return(expectedTags, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, tagsRepository, logger)
		ctx := context.Background()

		toys, err := toysService.GetAllToys(ctx)
		require.NoError(t, err)
		assert.Len(t, toys, len(expectedToys))
		assert.Equal(t, expectedToys, toys)
	})

	t.Run("all toys without existing toys", func(t *testing.T) {
		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys(gomock.Any()).Return([]entities.Toy{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, nil, logger)
		ctx := context.Background()

		toys, err := toysService.GetAllToys(ctx)
		require.NoError(t, err)
		assert.Empty(t, toys)
	})

	t.Run("all toys fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys(gomock.Any()).Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, nil, logger)
		ctx := context.Background()

		toys, err := toysService.GetAllToys(ctx)
		require.Error(t, err)
		assert.Nil(t, toys)
	})
}

func TestCommonToysServiceGetMasterToys(t *testing.T) {
	t.Run("master toys with existing masterID", func(t *testing.T) {
		const (
			masterID uint64 = 1
			toyID    uint64 = 1
		)

		expectedTags := []entities.Tag{
			{ID: 1},
		}

		expectedToys := []entities.Toy{
			{
				ID:       toyID,
				MasterID: masterID,
				Tags:     expectedTags,
			},
		}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetMasterToys(gomock.Any(), masterID).Return(expectedToys, nil).MaxTimes(1)
		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetToyTags(gomock.Any(), toyID).Return(expectedTags, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, tagsRepository, logger)
		ctx := context.Background()

		toys, err := toysService.GetMasterToys(ctx, masterID)
		require.NoError(t, err)
		assert.Len(t, toys, len(expectedToys))
		assert.Equal(t, expectedToys, toys)
	})

	t.Run("master toys with non-existing masterID", func(t *testing.T) {
		const masterID uint64 = 1

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetMasterToys(gomock.Any(), masterID).Return(
			nil,
			errors.New("test error"),
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, nil, logger)
		ctx := context.Background()

		toys, err := toysService.GetMasterToys(ctx, masterID)
		require.Error(t, err)
		assert.Nil(t, toys)
	})
}

func TestCommonToysServiceAddToy(t *testing.T) {
	t.Run("add toy success", func(t *testing.T) {
		const expectedToyID = uint64(1)

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().AddToy(gomock.Any(), gomock.Any()).Return(expectedToyID, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, nil, logger)
		ctx := context.Background()

		toyID, err := toysService.AddToy(ctx, entities.AddToyDTO{})
		require.NoError(t, err)
		assert.Equal(t, expectedToyID, toyID)
	})

	t.Run("add toy fail", func(t *testing.T) {
		const expectedToyID = uint64(0)
		var expectedError = &customerrors.ToyAlreadyExistsError{}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().AddToy(gomock.Any(), gomock.Any()).Return(expectedToyID, expectedError).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := NewCommonToysService(toysRepository, nil, logger)
		ctx := context.Background()

		toyID, err := toysService.AddToy(ctx, entities.AddToyDTO{})
		require.Error(t, err)
		require.IsType(t, expectedError, err)
		assert.Equal(t, expectedToyID, toyID)
	})
}

func TestCommonToysServiceProcessToyTags(t *testing.T) {
	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetToyTags(gomock.Any(), uint64(1)).Return([]entities.Tag{}, nil).MaxTimes(1)
	tagsRepository.EXPECT().GetToyTags(gomock.Any(), uint64(2)).Return(
		nil,
		errors.New("test error"),
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	toysService := NewCommonToysService(nil, tagsRepository, logger)
	ctx := context.Background()

	testCases := []struct {
		toy           *entities.Toy
		repository    interfaces.TagsRepository
		logger        *slog.Logger
		errorExpected bool
	}{
		{
			toy:           &entities.Toy{ID: 1},
			repository:    tagsRepository,
			logger:        logger,
			errorExpected: false,
		},
		{
			toy:           &entities.Toy{ID: 2},
			repository:    tagsRepository,
			logger:        logger,
			errorExpected: true,
		},
	}

	for _, tc := range testCases {
		err := toysService.processToyTags(ctx, tc.toy)
		if tc.errorExpected {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
