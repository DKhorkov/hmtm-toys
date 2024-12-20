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

func TestCommonMastersServiceGetMasterByID(t *testing.T) {
	testCases := []struct {
		masterID      uint64
		resultLength  int
		errorExpected bool
		err           error
	}{
		{
			masterID:      1,
			resultLength:  1,
			errorExpected: false,
		},
		{
			masterID:      2,
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	mastersRepository.EXPECT().GetMasterByID(uint64(1)).Return(&entities.Master{ID: 1}, nil).MaxTimes(1)
	mastersRepository.EXPECT().GetMasterByID(uint64(2)).DoAndReturn(
		func(_ uint64) (*entities.Master, error) {
			return nil, &customerrors.MasterNotFoundError{}
		},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	mastersService := services.NewCommonMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		master, err := mastersService.GetMasterByID(ctx, tc.masterID)
		if tc.errorExpected {
			require.Error(t, err)
			require.IsType(t, tc.err, err)
			assert.Nil(t, master)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestCommonMastersServiceGetMasterByUserID(t *testing.T) {
	testCases := []struct {
		userID        uint64
		resultLength  int
		errorExpected bool
		err           error
	}{
		{
			userID:        1,
			resultLength:  1,
			errorExpected: false,
		},
		{
			userID:        2,
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	mastersRepository.EXPECT().GetMasterByUserID(uint64(1)).Return(&entities.Master{ID: 1}, nil).MaxTimes(1)
	mastersRepository.EXPECT().GetMasterByUserID(uint64(2)).DoAndReturn(
		func(_ uint64) (*entities.Master, error) {
			return nil, &customerrors.MasterNotFoundError{}
		},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	mastersService := services.NewCommonMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		master, err := mastersService.GetMasterByUserID(ctx, tc.userID)
		if tc.errorExpected {
			require.Error(t, err)
			require.IsType(t, tc.err, err)
			assert.Nil(t, master)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestCommonMastersServiceGetAllMasters(t *testing.T) {
	t.Run("all masters with existing masters", func(t *testing.T) {
		expectedMasters := []entities.Master{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().GetAllMasters().DoAndReturn(
			func() ([]entities.Master, error) {
				return expectedMasters, nil
			},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masters, err := mastersService.GetAllMasters(ctx)
		require.NoError(t, err)
		assert.Len(t, masters, len(expectedMasters))
		assert.Equal(t, expectedMasters, masters)
	})

	t.Run("all masters without existing masters", func(t *testing.T) {
		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().GetAllMasters().Return([]entities.Master{}, nil).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masters, err := mastersService.GetAllMasters(ctx)
		require.NoError(t, err)
		assert.Empty(t, masters)
	})

	t.Run("all masters fail", func(t *testing.T) {
		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().GetAllMasters().Return(nil, errors.New("test error")).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masters, err := mastersService.GetAllMasters(ctx)
		require.Error(t, err)
		assert.Nil(t, masters)
	})
}

func TestCommonMastersServiceRegisterMaster(t *testing.T) {
	t.Run("register master success", func(t *testing.T) {
		const expectedMasterID = uint64(1)

		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().RegisterMaster(gomock.Any()).Return(expectedMasterID, nil).MaxTimes(1)
		mastersRepository.EXPECT().GetMasterByUserID(gomock.Any()).Return(
			nil,
			&customerrors.MasterNotFoundError{},
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masterID, err := mastersService.RegisterMaster(ctx, entities.RegisterMasterDTO{})
		require.NoError(t, err)
		assert.Equal(t, expectedMasterID, masterID)
	})

	t.Run("register master fail", func(t *testing.T) {
		const expectedMasterID = uint64(0)
		var expectedError = &customerrors.MasterAlreadyExistsError{}

		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().GetMasterByUserID(gomock.Any()).Return(&entities.Master{}, nil).MaxTimes(1)
		mastersRepository.EXPECT().RegisterMaster(gomock.Any()).Return(expectedMasterID, expectedError).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masterID, err := mastersService.RegisterMaster(ctx, entities.RegisterMasterDTO{})
		require.Error(t, err)
		require.IsType(t, expectedError, err)
		assert.Equal(t, expectedMasterID, masterID)
	})
}
