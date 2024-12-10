package interfaces

import (
	entities2 "github.com/DKhorkov/hmtm-toys/internal/entities"
)

type UseCases interface {
	TagsService
	CategoriesService
	GetAllMasters() ([]entities2.Master, error)
	GetMasterByID(id uint64) (*entities2.Master, error)
	RegisterMaster(rawMasterData entities2.RawRegisterMasterDTO) (masterID uint64, err error)
	AddToy(rawToyData entities2.RawAddToyDTO) (toyID uint64, err error)
	GetAllToys() ([]entities2.Toy, error)
	GetToyByID(id uint64) (*entities2.Toy, error)
	GetMasterToys(masterID uint64) ([]entities2.Toy, error)
}
