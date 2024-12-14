package services

import (
	"context"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/entities"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/libs/logging"
)

type CommonMastersService struct {
	mastersRepository interfaces.MastersRepository
	logger            *slog.Logger
}

func (service *CommonMastersService) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByID(id)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get master by id", err)
		return nil, &customerrors.MasterNotFoundError{BaseErr: err}
	}

	return master, nil
}

func (service *CommonMastersService) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByUserID(userID)
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get master by userID", err)
		return nil, &customerrors.MasterNotFoundError{BaseErr: err}
	}

	return master, nil
}

func (service *CommonMastersService) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	masters, err := service.mastersRepository.GetAllMasters()
	if err != nil {
		logging.LogErrorContext(ctx, service.logger, "Error occurred while trying to get all masters", err)
		return nil, err
	}

	return masters, nil
}

func (service *CommonMastersService) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	master, _ := service.mastersRepository.GetMasterByUserID(masterData.UserID)
	if master != nil {
		return 0, &customerrors.MasterAlreadyExistsError{}
	}

	return service.mastersRepository.RegisterMaster(masterData)
}

func NewCommonMastersService(
	mastersRepository interfaces.MastersRepository,
	logger *slog.Logger,
) *CommonMastersService {
	return &CommonMastersService{
		mastersRepository: mastersRepository,
		logger:            logger,
	}
}
