package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

func NewMastersService(
	mastersRepository interfaces.MastersRepository,
	logger logging.Logger,
) *MastersService {
	return &MastersService{
		mastersRepository: mastersRepository,
		logger:            logger,
	}
}

type MastersService struct {
	mastersRepository interfaces.MastersRepository
	logger            logging.Logger
}

func (service *MastersService) GetMasterByID(
	ctx context.Context,
	id uint64,
) (*entities.Master, error) {
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

func (service *MastersService) GetMasterByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.Master, error) {
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

func (service *MastersService) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return service.mastersRepository.GetAllMasters(ctx)
}

func (service *MastersService) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	master, _ := service.mastersRepository.GetMasterByUserID(ctx, masterData.UserID)
	if master != nil {
		return 0, &customerrors.MasterAlreadyExistsError{}
	}

	return service.mastersRepository.RegisterMaster(ctx, masterData)
}

func (service *MastersService) UpdateMaster(
	ctx context.Context,
	masterData entities.UpdateMasterDTO,
) error {
	return service.mastersRepository.UpdateMaster(ctx, masterData)
}
