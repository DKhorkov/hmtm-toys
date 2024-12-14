package usecases

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/libs/security"
)

type CommonUseCases struct {
	tagsService       interfaces.TagsService
	categoriesService interfaces.CategoriesService
	mastersService    interfaces.MastersService
	toysService       interfaces.ToysService
	jwtConfig         security.JWTConfig
}

// parseAccessToken parses JWT and gets UserID from it for further purposes.
func (useCases *CommonUseCases) parseAccessToken(accessToken string) (uint64, error) {
	accessTokenPayload, err := security.ParseJWT(accessToken, useCases.jwtConfig.SecretKey)
	if err != nil {
		return 0, &security.InvalidJWTError{}
	}

	userID := uint64(accessTokenPayload.(float64))
	return userID, nil
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

func (useCases *CommonUseCases) AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (uint64, error) {
	userID, err := useCases.parseAccessToken(rawToyData.AccessToken)
	if err != nil {
		return 0, err
	}

	master, err := useCases.mastersService.GetMasterByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	if _, err = useCases.GetCategoryByID(ctx, rawToyData.CategoryID); err != nil {
		return 0, err
	}

	for _, tagID := range rawToyData.TagsIDs {
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
		TagsIDs:     rawToyData.TagsIDs,
	}

	return useCases.toysService.AddToy(ctx, toyData)
}

func (useCases *CommonUseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByID(ctx, id)
}

func (useCases *CommonUseCases) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return useCases.mastersService.GetAllMasters(ctx)
}

func (useCases *CommonUseCases) RegisterMaster(
	ctx context.Context,
	rawMasterData entities.RawRegisterMasterDTO,
) (uint64, error) {
	userID, err := useCases.parseAccessToken(rawMasterData.AccessToken)
	if err != nil {
		return 0, err
	}

	masterData := entities.RegisterMasterDTO{
		UserID: userID,
		Info:   rawMasterData.Info,
	}

	return useCases.mastersService.RegisterMaster(ctx, masterData)
}

func NewCommonUseCases(
	tagsService interfaces.TagsService,
	categoriesService interfaces.CategoriesService,
	mastersService interfaces.MastersService,
	toysService interfaces.ToysService,
	jwtConfig security.JWTConfig,
) *CommonUseCases {
	return &CommonUseCases{
		tagsService:       tagsService,
		categoriesService: categoriesService,
		mastersService:    mastersService,
		toysService:       toysService,
		jwtConfig:         jwtConfig,
	}
}
