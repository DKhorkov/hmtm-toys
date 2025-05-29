package interfaces

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/toys_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,TagsRepository,SsoRepository -package=mockrepositories
type ToysRepository interface {
	AddToy(ctx context.Context, toyData entities.AddToyDTO) (toyID uint64, err error)
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
	DeleteToy(ctx context.Context, id uint64) error
	UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error
}

//go:generate mockgen -source=repositories.go  -destination=../../mocks/repositories/masters_repository.go -exclude_interfaces=TagsRepository,CategoriesRepository,ToysRepository,SsoRepository -package=mockrepositories
type MastersRepository interface {
	GetMasters(ctx context.Context, pagination *entities.Pagination) ([]entities.Master, error)
	GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error)
	GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error)
	RegisterMaster(
		ctx context.Context,
		masterData entities.RegisterMasterDTO,
	) (masterID uint64, err error)
	UpdateMaster(ctx context.Context, masterData entities.UpdateMasterDTO) error
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/categories_repository.go -exclude_interfaces=MastersRepository,TagsRepository,ToysRepository,SsoRepository -package=mockrepositories
type CategoriesRepository interface {
	GetAllCategories(ctx context.Context) ([]entities.Category, error)
	GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/tags_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,ToysRepository,SsoRepository -package=mockrepositories
type TagsRepository interface {
	CreateTags(ctx context.Context, tagsData []entities.CreateTagDTO) ([]uint32, error)
	GetAllTags(ctx context.Context) ([]entities.Tag, error)
	GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error)
}

//go:generate mockgen -source=repositories.go -destination=../../mocks/repositories/sso_repository.go -exclude_interfaces=MastersRepository,CategoriesRepository,ToysRepository,TagsRepository -package=mockrepositories
type SsoRepository interface {
	GetUserByID(ctx context.Context, id uint64) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
}
