package interfaces

import "github.com/DKhorkov/hmtm-toys/pkg/entities"

type ToysService interface {
	AddToy(toyData entities.AddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]*entities.Toy, error)
	GetToyByID(id uint64) (*entities.Toy, error)
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
