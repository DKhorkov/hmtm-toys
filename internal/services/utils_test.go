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
