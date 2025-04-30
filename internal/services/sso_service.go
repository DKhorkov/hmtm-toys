package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

type SsoService struct {
	ssoRepository interfaces.SsoRepository
	logger        logging.Logger
}

func NewSsoService(ssoRepository interfaces.SsoRepository, logger logging.Logger) *SsoService {
	return &SsoService{
		ssoRepository: ssoRepository,
		logger:        logger,
	}
}

func (service *SsoService) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	user, err := service.ssoRepository.GetUserByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get User with ID=%d", id),
			err,
		)
	}

	return user, err
}

func (service *SsoService) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	user, err := service.ssoRepository.GetUserByEmail(ctx, email)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			"Error occurred while trying to get User with Email="+email,
			err,
		)
	}

	return user, err
}
