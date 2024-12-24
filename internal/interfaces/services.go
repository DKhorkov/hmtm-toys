package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

type ToysService interface {
	ToysRepository
}

type TagsService interface {
	GetAllTags(ctx context.Context) ([]entities.Tag, error)
	GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error)
}

type MastersService interface {
	MastersRepository
}

type CategoriesService interface {
	CategoriesRepository
}
