package services

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	customerrors "github.com/DKhorkov/hmtm-toys/internal/errors"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

func NewTagsService(
	tagsRepository interfaces.TagsRepository,
	logger logging.Logger,
) *TagsService {
	return &TagsService{
		tagsRepository: tagsRepository,
		logger:         logger,
	}
}

type TagsService struct {
	tagsRepository interfaces.TagsRepository
	logger         logging.Logger
}

func (service *TagsService) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
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

func (service *TagsService) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return service.tagsRepository.GetAllTags(ctx)
}

func (service *TagsService) CreateTags(ctx context.Context, tagsData []entities.CreateTagDTO) ([]uint32, error) {
	return service.tagsRepository.CreateTags(ctx, tagsData)
}
