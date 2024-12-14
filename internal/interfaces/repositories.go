package interfaces

import (
	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/toys_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,TagsRepository -package=mockrepositories
type ToysRepository interface {
	AddToy(toyData entities.AddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]entities.Toy, error)
	GetToyByID(id uint64) (*entities.Toy, error)
	GetMasterToys(masterID uint64) ([]entities.Toy, error)
}

//go:generate mockgen -source=repositories.go  -destination=../../mocks/repositories/masters_repository.go -exclude_interfaces=TagsRepository,CategoriesRepository,ToysRepository -package=mockrepositories
type MastersRepository interface {
	GetAllMasters() ([]entities.Master, error)
	GetMasterByID(id uint64) (*entities.Master, error)
	GetMasterByUserID(userID uint64) (*entities.Master, error)
	RegisterMaster(masterData entities.RegisterMasterDTO) (masterID uint64, err error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/categories_repository.go -exclude_interfaces=MastersRepository,TagsRepository,ToysRepository -package=mockrepositories
type CategoriesRepository interface {
	GetAllCategories() ([]entities.Category, error)
	GetCategoryByID(id uint32) (*entities.Category, error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/tags_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,ToysRepository -package=mockrepositories
type TagsRepository interface {
	GetAllTags() ([]entities.Tag, error)
	GetTagByID(id uint32) (*entities.Tag, error)
	GetToyTags(toyID uint64) ([]entities.Tag, error)
}
