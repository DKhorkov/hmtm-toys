package interfaces

import "github.com/DKhorkov/hmtm-toys/pkg/entities"

type ToysService interface {
	ToysRepository
}

type TagsService interface {
	GetAllTags() ([]*entities.Tag, error)
	GetTagByID(id uint32) (*entities.Tag, error)
}

type MastersService interface {
	MastersRepository
}

type CategoriesService interface {
	CategoriesRepository
}
