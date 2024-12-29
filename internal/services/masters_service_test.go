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
		name          string
		masterID      uint64
		errorExpected bool
		err           error
	}{
		{
			name:          "successfully got Master by id",
			masterID:      1,
			errorExpected: false,
		},
		{
			name:          "failed to get Master by id",
			masterID:      2,
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	mastersRepository.EXPECT().GetMasterByID(gomock.Any(), uint64(1)).Return(&entities.Master{ID: 1}, nil).MaxTimes(1)
	mastersRepository.EXPECT().GetMasterByID(gomock.Any(), uint64(2)).Return(
		nil,
		&customerrors.MasterNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	mastersService := services.NewCommonMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			master, err := mastersService.GetMasterByID(ctx, tc.masterID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, master)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommonMastersServiceGetMasterByUserID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint64
		errorExpected bool
		err           error
	}{
		{
			name:          "successfully got Master by userID",
			userID:        1,
			errorExpected: false,
		},
		{
			name:          "failed to get Master by userID",
			userID:        2,
			errorExpected: true,
			err:           &customerrors.MasterNotFoundError{},
		},
	}

	mockController := gomock.NewController(t)
	mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
	mastersRepository.EXPECT().GetMasterByUserID(gomock.Any(), uint64(1)).Return(&entities.Master{ID: 1}, nil).MaxTimes(1)
	mastersRepository.EXPECT().GetMasterByUserID(gomock.Any(), uint64(2)).Return(
		nil,
		&customerrors.MasterNotFoundError{},
	).MaxTimes(1)

	logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
	mastersService := services.NewCommonMastersService(mastersRepository, logger)
	ctx := context.Background()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			master, err := mastersService.GetMasterByUserID(ctx, tc.userID)
			if tc.errorExpected {
				require.Error(t, err)
				require.IsType(t, tc.err, err)
				assert.Nil(t, master)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCommonMastersServiceGetAllMasters(t *testing.T) {
	t.Run("all masters with existing masters", func(t *testing.T) {
		expectedMasters := []entities.Master{
			{ID: 1},
		}

		mockController := gomock.NewController(t)
		mastersRepository := mockrepositories.NewMockMastersRepository(mockController)
		mastersRepository.EXPECT().GetAllMasters(gomock.Any()).Return(expectedMasters, nil).MaxTimes(1)

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
		mastersRepository.EXPECT().GetAllMasters(gomock.Any()).Return([]entities.Master{}, nil).MaxTimes(1)

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
		mastersRepository.EXPECT().GetAllMasters(gomock.Any()).Return(nil, errors.New("test error")).MaxTimes(1)

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
		mastersRepository.EXPECT().RegisterMaster(gomock.Any(), gomock.Any()).Return(expectedMasterID, nil).MaxTimes(1)
		mastersRepository.EXPECT().GetMasterByUserID(gomock.Any(), gomock.Any()).Return(
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
		mastersRepository.EXPECT().GetMasterByUserID(gomock.Any(), gomock.Any()).Return(
			&entities.Master{},
			nil,
		).MaxTimes(1)

		mastersRepository.EXPECT().RegisterMaster(gomock.Any(), gomock.Any()).Return(
			expectedMasterID,
			expectedError,
		).MaxTimes(1)

		logger := slog.New(slog.NewJSONHandler(bytes.NewBuffer(make([]byte, 1000)), nil))
		mastersService := services.NewCommonMastersService(mastersRepository, logger)
		ctx := context.Background()

		masterID, err := mastersService.RegisterMaster(ctx, entities.RegisterMasterDTO{})
		require.Error(t, err)
		require.IsType(t, expectedError, err)
		assert.Equal(t, expectedMasterID, masterID)
	})
}
