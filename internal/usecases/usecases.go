package usecases

import (
	"context"
	"strings"

	"github.com/DKhorkov/libs/validation"

	"github.com/DKhorkov/hmtm-toys/internal/config"
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

const (
	priceCeil     = 1_000_000
	priceFloor    = 1
	quantityCeil  = 1_000
	quantityFloor = 1
)

type UseCases struct {
	tagsService       interfaces.TagsService
	categoriesService interfaces.CategoriesService
	mastersService    interfaces.MastersService
	toysService       interfaces.ToysService
	ssoService        interfaces.SsoService
	validationConfig  config.ValidationConfig
}

func New(
	tagsService interfaces.TagsService,
	categoriesService interfaces.CategoriesService,
	mastersService interfaces.MastersService,
	toysService interfaces.ToysService,
	ssoService interfaces.SsoService,
	validationConfig config.ValidationConfig,
) *UseCases {
	return &UseCases{
		tagsService:       tagsService,
		categoriesService: categoriesService,
		mastersService:    mastersService,
		toysService:       toysService,
		ssoService:        ssoService,
		validationConfig:  validationConfig,
	}
}

func (useCases *UseCases) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	return useCases.tagsService.GetTagByID(ctx, id)
}

func (useCases *UseCases) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return useCases.tagsService.GetAllTags(ctx)
}

func (useCases *UseCases) GetCategoryByID(
	ctx context.Context,
	id uint32,
) (*entities.Category, error) {
	return useCases.categoriesService.GetCategoryByID(ctx, id)
}

func (useCases *UseCases) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return useCases.categoriesService.GetAllCategories(ctx)
}

func (useCases *UseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *UseCases) GetToys(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	return useCases.toysService.GetToys(ctx, pagination, filters)
}

func (useCases *UseCases) CountToys(ctx context.Context, filters *entities.ToysFilters) (uint64, error) {
	return useCases.toysService.CountToys(ctx, filters)
}

func (useCases *UseCases) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID, pagination, filters)
}

func (useCases *UseCases) CountMasterToys(
	ctx context.Context,
	masterID uint64,
	filters *entities.ToysFilters,
) (uint64, error) {
	return useCases.toysService.CountMasterToys(ctx, masterID, filters)
}

func (useCases *UseCases) GetUserToys(
	ctx context.Context,
	userID uint64,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	master, err := useCases.GetMasterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return useCases.GetMasterToys(ctx, master.ID, pagination, filters)
}

func (useCases *UseCases) CountUserToys(
	ctx context.Context,
	userID uint64,
	filters *entities.ToysFilters,
) (uint64, error) {
	master, err := useCases.GetMasterByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return useCases.toysService.CountMasterToys(ctx, master.ID, filters)
}

func (useCases *UseCases) AddToy(
	ctx context.Context,
	rawToyData entities.RawAddToyDTO,
) (uint64, error) {
	if !validation.ValidateValueByRules(
		rawToyData.Name,
		useCases.validationConfig.Toy.Name,
	) || validation.ContainsForbiddenWords(
		rawToyData.Name,
	) {
		return 0, &validation.Error{Message: "invalid toy name"}
	}

	if !validation.ValidateValueByRules(
		rawToyData.Description,
		useCases.validationConfig.Toy.Description,
	) || validation.ContainsForbiddenWords(
		rawToyData.Description,
	) {
		return 0, &validation.Error{Message: "invalid toy description"}
	}

	if rawToyData.Price > priceCeil || rawToyData.Price < priceFloor {
		return 0, &validation.Error{Message: "invalid toy price"}
	}

	if rawToyData.Quantity > quantityCeil || rawToyData.Quantity < quantityFloor {
		return 0, &validation.Error{Message: "invalid toy quantity"}
	}

	master, err := useCases.GetMasterByUserID(ctx, rawToyData.UserID)
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

func (useCases *UseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByID(ctx, id)
}

func (useCases *UseCases) GetMasterByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByUserID(ctx, userID)
}

func (useCases *UseCases) GetMasters(ctx context.Context, pagination *entities.Pagination) ([]entities.Master, error) {
	return useCases.mastersService.GetMasters(ctx, pagination)
}

func (useCases *UseCases) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	if masterData.Info != nil &&
		(!validation.ValidateValueByRules(
			*masterData.Info,
			useCases.validationConfig.Master.Info,
		) || validation.ContainsForbiddenWords(
			*masterData.Info,
		)) {
		return 0, &validation.Error{Message: "invalid master info"}
	}

	if _, err := useCases.ssoService.GetUserByID(ctx, masterData.UserID); err != nil {
		return 0, err
	}

	return useCases.mastersService.RegisterMaster(ctx, masterData)
}

func (useCases *UseCases) CreateTags(
	ctx context.Context,
	tagsData []entities.CreateTagDTO,
) ([]uint32, error) {
	for _, tag := range tagsData {
		if !validation.ValidateValueByRules(
			tag.Name,
			useCases.validationConfig.Tag.Name,
		) || validation.ContainsForbiddenWords(
			tag.Name,
		) {
			return nil, &validation.Error{Message: "invalid tag name: " + tag.Name}
		}
	}

	existingTags, err := useCases.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	existingTagsSet := make(map[string]uint32)
	for _, tag := range existingTags {
		existingTagsSet[tag.Name] = tag.ID
	}

	uniqueTags := make(map[string]struct{}, len(tagsData))

	for _, tag := range tagsData {
		lowerCaseTagName := strings.ToLower(tag.Name)
		if _, ok := uniqueTags[lowerCaseTagName]; !ok {
			uniqueTags[lowerCaseTagName] = struct{}{}
		}
	}

	var tagsToCreate []entities.CreateTagDTO

	var existingTagIDs []uint32

	for tag := range uniqueTags {
		if _, ok := existingTagsSet[tag]; ok {
			existingTagIDs = append(existingTagIDs, existingTagsSet[tag])

			continue
		}

		tagsToCreate = append(tagsToCreate, entities.CreateTagDTO{Name: tag})
	}

	createdTagIDs, err := useCases.tagsService.CreateTags(ctx, tagsToCreate)
	if err != nil {
		return nil, err
	}

	createdTagIDs = append(createdTagIDs, existingTagIDs...)

	return createdTagIDs, nil
}

func (useCases *UseCases) DeleteToy(ctx context.Context, id uint64) error {
	if _, err := useCases.GetToyByID(ctx, id); err != nil {
		return err
	}

	return useCases.toysService.DeleteToy(ctx, id)
}

func (useCases *UseCases) UpdateToy(
	ctx context.Context,
	rawToyData entities.RawUpdateToyDTO,
) error {
	if rawToyData.Name != nil &&
		(!validation.ValidateValueByRules(
			*rawToyData.Name,
			useCases.validationConfig.Toy.Name,
		) || validation.ContainsForbiddenWords(
			*rawToyData.Name,
		)) {
		return &validation.Error{Message: "invalid toy name"}
	}

	if rawToyData.Description != nil &&
		(!validation.ValidateValueByRules(
			*rawToyData.Description,
			useCases.validationConfig.Toy.Description,
		) || validation.ContainsForbiddenWords(
			*rawToyData.Description,
		)) {
		return &validation.Error{Message: "invalid toy description"}
	}

	if rawToyData.Price != nil && (*rawToyData.Price > priceCeil || *rawToyData.Price < priceFloor) {
		return &validation.Error{Message: "invalid toy price"}
	}

	if rawToyData.Quantity != nil && (*rawToyData.Quantity > quantityCeil || *rawToyData.Quantity < quantityFloor) {
		return &validation.Error{Message: "invalid toy quantity"}
	}

	toy, err := useCases.GetToyByID(ctx, rawToyData.ID)
	if err != nil {
		return err
	}

	if rawToyData.CategoryID != nil {
		if _, err = useCases.GetCategoryByID(ctx, *rawToyData.CategoryID); err != nil {
			return err
		}
	}

	for _, tagID := range rawToyData.TagIDs {
		if _, err = useCases.GetTagByID(ctx, tagID); err != nil {
			return err
		}
	}

	// Old Toy Tags IDs set:
	oldTagIDsSet := make(map[uint32]struct{}, len(toy.Tags))
	for _, tag := range toy.Tags {
		oldTagIDsSet[tag.ID] = struct{}{}
	}

	// New Toy Tags IDs set:
	newTagIDsSet := make(map[uint32]struct{}, len(rawToyData.TagIDs))
	for _, tagID := range rawToyData.TagIDs {
		newTagIDsSet[tagID] = struct{}{}
	}

	// Add new Tag if it is not already exists:
	tagIDsToAdd := make([]uint32, 0)

	for _, tagID := range rawToyData.TagIDs {
		if _, ok := oldTagIDsSet[tagID]; !ok {
			tagIDsToAdd = append(tagIDsToAdd, tagID)
		}
	}

	// Delete old Tag if it is not used by Toy now:
	tagIDsToDelete := make([]uint32, 0)

	for _, tag := range toy.Tags {
		if _, ok := newTagIDsSet[tag.ID]; !ok {
			tagIDsToDelete = append(tagIDsToDelete, tag.ID)
		}
	}

	// Old Toy Attachments set:
	oldAttachmentsSet := make(map[string]struct{}, len(toy.Attachments))
	for _, attachment := range toy.Attachments {
		oldAttachmentsSet[attachment.Link] = struct{}{}
	}

	// New Toy Attachments set:
	newAttachmentsSet := make(map[string]struct{}, len(rawToyData.Attachments))
	for _, attachment := range rawToyData.Attachments {
		newAttachmentsSet[attachment] = struct{}{}
	}

	// Add new Attachments if it is not already exists:
	attachmentsToAdd := make([]string, 0)

	for _, attachment := range rawToyData.Attachments {
		if _, ok := oldAttachmentsSet[attachment]; !ok {
			attachmentsToAdd = append(attachmentsToAdd, attachment)
		}
	}

	// Delete old Attachments if it is not used by Toy now:
	attachmentsToDelete := make([]uint64, 0)

	for _, attachment := range toy.Attachments {
		if _, ok := newAttachmentsSet[attachment.Link]; !ok {
			attachmentsToDelete = append(attachmentsToDelete, attachment.ID)
		}
	}

	toyData := entities.UpdateToyDTO{
		ID:                    rawToyData.ID,
		CategoryID:            rawToyData.CategoryID,
		Name:                  rawToyData.Name,
		Description:           rawToyData.Description,
		Price:                 rawToyData.Price,
		Quantity:              rawToyData.Quantity,
		TagIDsToAdd:           tagIDsToAdd,
		TagIDsToDelete:        tagIDsToDelete,
		AttachmentsToAdd:      attachmentsToAdd,
		AttachmentIDsToDelete: attachmentsToDelete,
	}

	return useCases.toysService.UpdateToy(ctx, toyData)
}

func (useCases *UseCases) UpdateMaster(
	ctx context.Context,
	masterData entities.UpdateMasterDTO,
) error {
	if masterData.Info != nil &&
		(!validation.ValidateValueByRules(
			*masterData.Info,
			useCases.validationConfig.Master.Info,
		) || validation.ContainsForbiddenWords(
			*masterData.Info,
		)) {
		return &validation.Error{Message: "invalid master info"}
	}

	if _, err := useCases.GetMasterByID(ctx, masterData.ID); err != nil {
		return err
	}

	return useCases.mastersService.UpdateMaster(ctx, masterData)
}
