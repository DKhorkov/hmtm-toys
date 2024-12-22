package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
)

func NewCommonMastersService(
	mastersRepository interfaces.MastersRepository,
	logger *slog.Logger,
) *CommonMastersService {
	return &CommonMastersService{
		mastersRepository: mastersRepository,
		logger:            logger,
	}
}

type CommonMastersService struct {
	mastersRepository interfaces.MastersRepository
	logger            *slog.Logger
}

func (service *CommonMastersService) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Master with ID=%d", id),
			err,
		)

		return nil, &customerrors.MasterNotFoundError{}
	}

	return master, nil
}

func (service *CommonMastersService) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByUserID(ctx, userID)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Master by userID=%d", userID),
			err,
		)

		return nil, &customerrors.MasterNotFoundError{}
	}

	return master, nil
}

func (service *CommonMastersService) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return service.mastersRepository.GetAllMasters(ctx)
}

func (service *CommonMastersService) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	master, _ := service.mastersRepository.GetMasterByUserID(ctx, masterData.UserID)
	if master != nil {
		return 0, &customerrors.MasterAlreadyExistsError{}
	}

	return service.mastersRepository.RegisterMaster(ctx, masterData)
}
