package interfaces

import (
	entities2 "github.com/DKhorkov/hmtm-toys/internal/entities"
)

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/toys_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,TagsRepository -package=mockrepositories
type ToysRepository interface {
	AddToy(toyData entities2.AddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]entities2.Toy, error)
	GetToyByID(id uint64) (*entities2.Toy, error)
	GetMasterToys(masterID uint64) ([]entities2.Toy, error)
}

//go:generate mockgen -source=repositories.go  -destination=../../mocks/repositories/masters_repository.go -exclude_interfaces=TagsRepository,CategoriesRepository,ToysRepository -package=mockrepositories
type MastersRepository interface {
	GetAllMasters() ([]entities2.Master, error)
	GetMasterByID(id uint64) (*entities2.Master, error)
	GetMasterByUserID(userID uint64) (*entities2.Master, error)
	RegisterMaster(masterData entities2.RegisterMasterDTO) (masterID uint64, err error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/categories_repository.go -exclude_interfaces=MastersRepository,TagsRepository,ToysRepository -package=mockrepositories
type CategoriesRepository interface {
	GetAllCategories() ([]entities2.Category, error)
	GetCategoryByID(id uint32) (*entities2.Category, error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/tags_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,ToysRepository -package=mockrepositories
type TagsRepository interface {
	GetAllTags() ([]entities2.Tag, error)
	GetTagByID(id uint32) (*entities2.Tag, error)
	GetToyTags(toyID uint64) ([]entities2.Tag, error)
}
