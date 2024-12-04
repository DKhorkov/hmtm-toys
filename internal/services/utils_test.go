package services

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	mockrepositories "github.com/DKhorkov/hmtm-toys/mocks/repositories"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
)

func TestServiceUtilsProcessToyTags(t *testing.T) {
	mockController := gomock.NewController(t)
	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetToyTags(uint64(1)).Return([]*entities.Tag{}, nil).MaxTimes(1)
	tagsRepository.EXPECT().GetToyTags(uint64(2)).Return(nil, errors.New("test error")).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))

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
		err := processToyTags(tc.toy, tc.repository, tc.logger)
		if tc.errorExpected {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestServiceUtilsProcessMasterToys(t *testing.T) {
	mockController := gomock.NewController(t)
	toysRepository := mockrepositories.NewMockToysRepository(mockController)
	toysRepository.EXPECT().GetMasterToys(uint64(1)).DoAndReturn(
		func(_ any) ([]*entities.Toy, error) {
			return []*entities.Toy{
				{ID: 1},
			}, nil
		},
	).MaxTimes(1)

	toysRepository.EXPECT().GetMasterToys(uint64(2)).Return(nil, errors.New("test error")).MaxTimes(1)
	toysRepository.EXPECT().GetMasterToys(uint64(3)).DoAndReturn(
		func(_ any) ([]*entities.Toy, error) {
			return []*entities.Toy{
				{ID: 2},
			}, nil
		},
	).MaxTimes(1)

	tagsRepository := mockrepositories.NewMockTagsRepository(mockController)
	tagsRepository.EXPECT().GetToyTags(uint64(1)).Return([]*entities.Tag{}, nil).MaxTimes(1)
	tagsRepository.EXPECT().GetToyTags(uint64(2)).Return(nil, errors.New("test error")).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))

	testCases := []struct {
		master         *entities.Master
		toysRepository interfaces.ToysRepository
		tagsRepository interfaces.TagsRepository
		logger         *slog.Logger
		errorExpected  bool
	}{
		{
			master:         &entities.Master{ID: 1},
			toysRepository: toysRepository,
			tagsRepository: tagsRepository,
			logger:         logger,
			errorExpected:  false,
		},
		{
			master:         &entities.Master{ID: 2},
			toysRepository: toysRepository,
			tagsRepository: tagsRepository,
			logger:         logger,
			errorExpected:  true,
		},
		{
			master:         &entities.Master{ID: 3},
			toysRepository: toysRepository,
			tagsRepository: tagsRepository,
			logger:         logger,
			errorExpected:  true,
		},
	}

	for _, tc := range testCases {
		err := processMasterToys(tc.master, tc.toysRepository, tc.tagsRepository, logger)
		if tc.errorExpected {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
