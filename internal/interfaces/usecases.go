package interfaces

import "github.com/DKhorkov/hmtm-toys/pkg/entities"

type UseCases interface {
	TagsService
	CategoriesService
	GetAllMasters() ([]*entities.Master, error)
	GetMasterByID(id uint64) (*entities.Master, error)
	RegisterMaster(rawMasterData entities.RawRegisterMasterDTO) (masterID uint64, err error)
	AddToy(rawToyData entities.RawAddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]*entities.Toy, error)
	GetToyByID(id uint64) (*entities.Toy, error)
	GetMasterToys(masterID uint64) ([]*entities.Toy, error)
}
