package usecases

import (
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
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

func (useCases *CommonUseCases) GetTagByID(id uint32) (*entities.Tag, error) {
	return useCases.tagsService.GetTagByID(id)
}

func (useCases *CommonUseCases) GetAllTags() ([]*entities.Tag, error) {
	return useCases.tagsService.GetAllTags()
}

func (useCases *CommonUseCases) GetCategoryByID(id uint32) (*entities.Category, error) {
	return useCases.categoriesService.GetCategoryByID(id)
}

func (useCases *CommonUseCases) GetAllCategories() ([]*entities.Category, error) {
	return useCases.categoriesService.GetAllCategories()
}

func (useCases *CommonUseCases) GetToyByID(id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(id)
}

func (useCases *CommonUseCases) GetAllToys() ([]*entities.Toy, error) {
	return useCases.toysService.GetAllToys()
}

func (useCases *CommonUseCases) GetMasterToys(masterID uint64) ([]*entities.Toy, error) {
	return useCases.toysService.GetMasterToys(masterID)
}

func (useCases *CommonUseCases) AddToy(rawToyData entities.RawAddToyDTO) (uint64, error) {
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

	toyData := entities.AddToyDTO{
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

func (useCases *CommonUseCases) GetMasterByID(id uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByID(id)
}

func (useCases *CommonUseCases) GetAllMasters() ([]*entities.Master, error) {
	return useCases.mastersService.GetAllMasters()
}

func (useCases *CommonUseCases) RegisterMaster(rawMasterData entities.RawRegisterMasterDTO) (uint64, error) {
	userID, err := useCases.parseAccessToken(rawMasterData.AccessToken)
	if err != nil {
		return 0, err
	}

	masterData := entities.RegisterMasterDTO{
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
