package usecases

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

func NewCommonUseCases(
	tagsService interfaces.TagsService,
	categoriesService interfaces.CategoriesService,
	mastersService interfaces.MastersService,
	toysService interfaces.ToysService,
) *CommonUseCases {
	return &CommonUseCases{
		tagsService:       tagsService,
		categoriesService: categoriesService,
		mastersService:    mastersService,
		toysService:       toysService,
	}
}

type CommonUseCases struct {
	tagsService       interfaces.TagsService
	categoriesService interfaces.CategoriesService
	mastersService    interfaces.MastersService
	toysService       interfaces.ToysService
}

func (useCases *CommonUseCases) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	return useCases.tagsService.GetTagByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return useCases.tagsService.GetAllTags(ctx)
}

func (useCases *CommonUseCases) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	return useCases.categoriesService.GetCategoryByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return useCases.categoriesService.GetAllCategories(ctx)
}

func (useCases *CommonUseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	return useCases.toysService.GetAllToys(ctx)
}

func (useCases *CommonUseCases) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID)
}

func (useCases *CommonUseCases) GetUserToys(ctx context.Context, userID uint64) ([]entities.Toy, error) {
	master, err := useCases.mastersService.GetMasterByUserID(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return useCases.GetMasterToys(ctx, master.ID)
}

func (useCases *CommonUseCases) AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (uint64, error) {
	master, err := useCases.mastersService.GetMasterByUserID(ctx, rawToyData.UserID)
	if err != nil {
		return 0, err
	}

	if _, err = useCases.GetCategoryByID(ctx, rawToyData.CategoryID); err != nil {
		return 0, err
	}

	for _, tagID := range rawToyData.TagIDs {
		if _, err = useCases.GetTagByID(ctx, tagID); err != nil {
			return 0, err
		}
	}

	toyData := entities.AddToyDTO{
		MasterID:    master.ID,
		Name:        rawToyData.Name,
		Description: rawToyData.Description,
		Price:       rawToyData.Price,
		Quantity:    rawToyData.Quantity,
		CategoryID:  rawToyData.CategoryID,
		TagIDs:      rawToyData.TagIDs,
		Attachments: rawToyData.Attachments,
	}

	return useCases.toysService.AddToy(ctx, toyData)
}

func (useCases *CommonUseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByID(ctx, id)
}

func (useCases *CommonUseCases) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByUserID(ctx, userID)
}

func (useCases *CommonUseCases) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return useCases.mastersService.GetAllMasters(ctx)
}

func (useCases *CommonUseCases) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	return useCases.mastersService.RegisterMaster(ctx, masterData)
}
