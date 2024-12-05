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
	toysRepository.EXPECT().GetToyByID(uint64(1)).Return(&entities.Toy{ID: 1}, nil).MaxTimes(1)
	toysRepository.EXPECT().GetToyByID(uint64(2)).Return(
		nil,
		&customerrors.ToyNotFoundError{},
	).MaxTimes(1)

	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetToyTags(uint64(1)).Return([]*entities.Tag{}, nil).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	toysService := services.NewCommonToysService(toysRepository, tagsRepository, logger)

	for _, tc := range testCases {
		toy, err := toysService.GetToyByID(tc.toyID)
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
		expectedTags := []*entities.Tag{
			{ID: 1},
		}

		expectedToys := []*entities.Toy{
			{
				ID:   1,
				Tags: expectedTags,
			},
		}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys().DoAndReturn(
			func() ([]*entities.Toy, error) {
				return expectedToys, nil
			},
		).MaxTimes(1)

		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetToyTags(uint64(1)).DoAndReturn(
			func(_ uint64) ([]*entities.Tag, error) {
				return expectedTags, nil
			},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, tagsRepository, logger)

		toys, err := toysService.GetAllToys()
		require.NoError(t, err)
		assert.Len(t, toys, len(expectedToys))
		assert.Equal(t, expectedToys, toys)
	})

	t.Run("all toys without existing toys", func(t *testing.T) {
		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys().Return([]*entities.Toy{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, nil, logger)

		toys, err := toysService.GetAllToys()
		require.NoError(t, err)
		assert.Empty(t, toys)
	})

	t.Run("all toys fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetAllToys().Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, nil, logger)

		toys, err := toysService.GetAllToys()
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

		expectedTags := []*entities.Tag{
			{ID: 1},
		}

		expectedToys := []*entities.Toy{
			{
				ID:       toyID,
				MasterID: masterID,
				Tags:     expectedTags,
			},
		}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetMasterToys(masterID).DoAndReturn(
			func(_ uint64) ([]*entities.Toy, error) {
				return expectedToys, nil
			},
		).MaxTimes(1)

		tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
		tagsRepository.EXPECT().GetToyTags(toyID).DoAndReturn(
			func(_ uint64) ([]*entities.Tag, error) {
				return expectedTags, nil
			},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, tagsRepository, logger)

		toys, err := toysService.GetMasterToys(masterID)
		require.NoError(t, err)
		assert.Len(t, toys, len(expectedToys))
		assert.Equal(t, expectedToys, toys)
	})

	t.Run("master toys with non-existing masterID", func(t *testing.T) {
		const masterID uint64 = 1

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().GetMasterToys(masterID).Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, nil, logger)

		toys, err := toysService.GetMasterToys(masterID)
		require.Error(t, err)
		assert.Nil(t, toys)
	})
}

func TestCommonToysServiceAddToy(t *testing.T) {
	t.Run("add toy success", func(t *testing.T) {
		const expectedToyID = uint64(1)

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().AddToy(gomock.Any()).Return(expectedToyID, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, nil, logger)

		toyID, err := toysService.AddToy(entities.AddToyDTO{})
		require.NoError(t, err)
		assert.Equal(t, expectedToyID, toyID)
	})

	t.Run("add toy fail", func(t *testing.T) {
		const expectedToyID = uint64(0)
		var expectedError = &customerrors.ToyAlreadyExistsError{}

		mockController := gomock.NewController(t)
		toysRepository := mockrepositories.NewMockToysRepository(mockController)
		toysRepository.EXPECT().AddToy(gomock.Any()).Return(expectedToyID, expectedError).MaxTimes(1)
		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		toysService := services.NewCommonToysService(toysRepository, nil, logger)

		toyID, err := toysService.AddToy(entities.AddToyDTO{})
		require.Error(t, err)
		require.IsType(t, expectedError, err)
		assert.Equal(t, expectedToyID, toyID)
	})
}
