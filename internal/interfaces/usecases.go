package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

//go:generate mockgen -source=usecases.go -destination=../../mocks/usecases/usecases.go -package=mockusecases
type UseCases interface {
	// Tags cases:
	TagsService

	// Categories cases:
	CategoriesService

	// Masters cases:
	GetMasters(
		ctx context.Context,
		pagination *entities.Pagination,
		filters *entities.MastersFilters,
	) ([]entities.Master, error)
	CountMasters(ctx context.Context, filters *entities.MastersFilters) (uint64, error)
	GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error)
	GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error)
	RegisterMaster(
		ctx context.Context,
		rawMasterData entities.RegisterMasterDTO,
	) (masterID uint64, err error)
	UpdateMaster(ctx context.Context, masterData entities.UpdateMasterDTO) error

	// Toys cases:
	AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (toyID uint64, err error)
	GetToys(ctx context.Context, pagination *entities.Pagination, filters *entities.ToysFilters) ([]entities.Toy, error)
	CountToys(ctx context.Context, filters *entities.ToysFilters) (uint64, error)
	GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error)
	GetMasterToys(
		ctx context.Context,
		masterID uint64,
		pagination *entities.Pagination,
		filters *entities.ToysFilters,
	) ([]entities.Toy, error)
	CountMasterToys(ctx context.Context, masterID uint64, filters *entities.ToysFilters) (uint64, error)
	GetUserToys(
		ctx context.Context,
		userID uint64,
		pagination *entities.Pagination,
		filters *entities.ToysFilters,
	) ([]entities.Toy, error)
	CountUserToys(ctx context.Context, userID uint64, filters *entities.ToysFilters) (uint64, error)
	DeleteToy(ctx context.Context, id uint64) error
	UpdateToy(ctx context.Context, rawToyData entities.RawUpdateToyDTO) error
}
