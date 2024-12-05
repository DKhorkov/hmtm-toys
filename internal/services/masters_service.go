package services

import (
	"log/slog"

	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/logging"
)

type CommonMastersService struct {
	mastersRepository interfaces.MastersRepository
	logger            *slog.Logger
}

func (service *CommonMastersService) GetMasterByID(id uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByID(id)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get master by id",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.MasterNotFoundError{}
	}

	return master, nil
}

func (service *CommonMastersService) GetMasterByUserID(userID uint64) (*entities.Master, error) {
	master, err := service.mastersRepository.GetMasterByUserID(userID)
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get master by userID",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, &customerrors.MasterNotFoundError{}
	}

	return master, nil
}

func (service *CommonMastersService) GetAllMasters() ([]*entities.Master, error) {
	masters, err := service.mastersRepository.GetAllMasters()
	if err != nil {
		service.logger.Error(
			"Error occurred while trying to get all masters",
			"Traceback",
			logging.GetLogTraceback(),
			"Error",
			err,
		)

		return nil, err
	}

	return masters, nil
}

func (service *CommonMastersService) RegisterMaster(masterData entities.RegisterMasterDTO) (uint64, error) {
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
