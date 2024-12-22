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

func NewCommonTagsService(
	tagsRepository interfaces.TagsRepository,
	logger *slog.Logger,
) *CommonTagsService {
	return &CommonTagsService{
		tagsRepository: tagsRepository,
		logger:         logger,
	}
}

type CommonTagsService struct {
	tagsRepository interfaces.TagsRepository
	logger         *slog.Logger
}

func (service *CommonTagsService) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	tag, err := service.tagsRepository.GetTagByID(ctx, id)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			service.logger,
			fmt.Sprintf("Error occurred while trying to get Tag with ID=%d", id),
			err,
		)

		return nil, &customerrors.TagNotFoundError{}
	}

	return tag, nil
}

func (service *CommonTagsService) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return service.tagsRepository.GetAllTags(ctx)
}
