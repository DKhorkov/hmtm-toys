package repositories

import (
	"context"

	"github.com/DKhorkov/hmtm-sso/api/protobuf/generated/go/sso"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

type SsoRepository struct {
	client interfaces.SsoClient
}

func NewSsoRepository(client interfaces.SsoClient) *SsoRepository {
	return &SsoRepository{client: client}
}

func (repo *SsoRepository) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	response, err := repo.client.GetUser(
		ctx,
		&sso.GetUserIn{
			ID: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return repo.processUserResponse(response), nil
}

func (repo *SsoRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	response, err := repo.client.GetUserByEmail(
		ctx,
		&sso.GetUserByEmailIn{
			Email: email,
		},
	)
	if err != nil {
		return nil, err
	}

	return repo.processUserResponse(response), nil
}

func (repo *SsoRepository) processUserResponse(userResponse *sso.GetUserOut) *entities.User {
	return &entities.User{
		ID:                userResponse.GetID(),
		DisplayName:       userResponse.GetDisplayName(),
		Email:             userResponse.GetEmail(),
		EmailConfirmed:    userResponse.GetEmailConfirmed(),
		Phone:             userResponse.Phone,
		PhoneConfirmed:    userResponse.GetPhoneConfirmed(),
		Telegram:          userResponse.Telegram,
		TelegramConfirmed: userResponse.GetTelegramConfirmed(),
		Avatar:            userResponse.Avatar,
		CreatedAt:         userResponse.GetCreatedAt().AsTime(),
		UpdatedAt:         userResponse.GetUpdatedAt().AsTime(),
	}
}
