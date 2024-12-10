package usecases

import (
	entities2 "github.com/DKhorkov/hmtm-toys/internal/entities"
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

func (useCases *CommonUseCases) GetTagByID(id uint32) (*entities2.Tag, error) {
	return useCases.tagsService.GetTagByID(id)
}

func (useCases *CommonUseCases) GetAllTags() ([]entities2.Tag, error) {
	return useCases.tagsService.GetAllTags()
}

func (useCases *CommonUseCases) GetCategoryByID(id uint32) (*entities2.Category, error) {
	return useCases.categoriesService.GetCategoryByID(id)
}

func (useCases *CommonUseCases) GetAllCategories() ([]entities2.Category, error) {
	return useCases.categoriesService.GetAllCategories()
}

func (useCases *CommonUseCases) GetToyByID(id uint64) (*entities2.Toy, error) {
	return useCases.toysService.GetToyByID(id)
}

func (useCases *CommonUseCases) GetAllToys() ([]entities2.Toy, error) {
	return useCases.toysService.GetAllToys()
}

func (useCases *CommonUseCases) GetMasterToys(masterID uint64) ([]entities2.Toy, error) {
	return useCases.toysService.GetMasterToys(masterID)
}

func (useCases *CommonUseCases) AddToy(rawToyData entities2.RawAddToyDTO) (uint64, error) {
	userID, err := useCases.parseAccessToken(rawToyData.AccessToken)
	if err != nil {
		return 0, err
	}

	master, err := useCases.mastersService.GetMasterByUserID(userID)
	if err != nil {
		return 0, err
	}

	if _, err = useCases.GetCategoryByID(rawToyData.CategoryID); err != nil {
		return 0, err
	}

	for _, tagID := range rawToyData.TagsIDs {
		if _, err = useCases.GetTagByID(tagID); err != nil {
			return 0, err
		}
	}

	toyData := entities2.AddToyDTO{
		MasterID:    master.ID,
		Name:        rawToyData.Name,
		Description: rawToyData.Description,
		Price:       rawToyData.Price,
		Quantity:    rawToyData.Quantity,
		CategoryID:  rawToyData.CategoryID,
		TagsIDs:     rawToyData.TagsIDs,
	}

	return useCases.toysService.AddToy(toyData)
}

func (useCases *CommonUseCases) GetMasterByID(id uint64) (*entities2.Master, error) {
	return useCases.mastersService.GetMasterByID(id)
}

func (useCases *CommonUseCases) GetAllMasters() ([]entities2.Master, error) {
	return useCases.mastersService.GetAllMasters()
}

func (useCases *CommonUseCases) RegisterMaster(rawMasterData entities2.RawRegisterMasterDTO) (uint64, error) {
	userID, err := useCases.parseAccessToken(rawMasterData.AccessToken)
	if err != nil {
		return 0, err
	}

	masterData := entities2.RegisterMasterDTO{
		UserID: userID,
		Info:   rawMasterData.Info,
	}

	return useCases.mastersService.RegisterMaster(masterData)
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
