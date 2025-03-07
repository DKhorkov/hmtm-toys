package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

type UseCases interface {
	TagsService
	CategoriesService
	GetAllMasters(ctx context.Context) ([]entities.Master, error)
	GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error)
	GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error)
	RegisterMaster(ctx context.Context, rawMasterData entities.RegisterMasterDTO) (masterID uint64, err error)
	AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (toyID uint64, err error)
	GetAllToys(ctx context.Context) ([]entities.Toy, error)
	GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error)
	GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error)
	GetUserToys(ctx context.Context, userID uint64) ([]entities.Toy, error)
	DeleteToy(ctx context.Context, id uint64) error
	UpdateToy(ctx context.Context, rawToyData entities.RawUpdateToyDTO) error
}
