package interfaces

import "github.com/DKhorkov/hmtm-toys/pkg/entities"

type ToysRepository interface {
	AddToy(toyData entities.AddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]*entities.Toy, error)
	GetToyByID(id uint64) (*entities.Toy, error)
	GetMasterToys(masterID uint64) ([]*entities.Toy, error)
}

type MastersRepository interface {
	GetAllMasters() ([]*entities.Master, error)
	GetMasterByID(id uint64) (*entities.Master, error)
	GetMasterByUserID(userID uint64) (*entities.Master, error)
	RegisterMaster(masterData entities.RegisterMasterDTO) (masterID uint64, err error)
}

type CategoriesRepository interface {
	GetAllCategories() ([]*entities.Category, error)
	GetCategoryByID(id uint32) (*entities.Category, error)
}

type TagsRepository interface {
	GetAllTags() ([]*entities.Tag, error)
	GetTagByID(id uint32) (*entities.Tag, error)
	GetToyTags(toyID uint64) ([]*entities.Tag, error)
}
