package usecases

import (
	"context"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/interfaces"
)

func New(
	tagsService interfaces.TagsService,
	categoriesService interfaces.CategoriesService,
	mastersService interfaces.MastersService,
	toysService interfaces.ToysService,
	ssoService interfaces.SsoService,
) *UseCases {
	return &UseCases{
		tagsService:       tagsService,
		categoriesService: categoriesService,
		mastersService:    mastersService,
		toysService:       toysService,
		ssoService:        ssoService,
	}
}

type UseCases struct {
	tagsService       interfaces.TagsService
	categoriesService interfaces.CategoriesService
	mastersService    interfaces.MastersService
	toysService       interfaces.ToysService
	ssoService        interfaces.SsoService
}

func (useCases *UseCases) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	return useCases.tagsService.GetTagByID(ctx, id)
}

func (useCases *UseCases) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	return useCases.tagsService.GetAllTags(ctx)
}

func (useCases *UseCases) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	return useCases.categoriesService.GetCategoryByID(ctx, id)
}

func (useCases *UseCases) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	return useCases.categoriesService.GetAllCategories(ctx)
}

func (useCases *UseCases) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	return useCases.toysService.GetToyByID(ctx, id)
}

func (useCases *UseCases) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	return useCases.toysService.GetAllToys(ctx)
}

func (useCases *UseCases) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	return useCases.toysService.GetMasterToys(ctx, masterID)
}

func (useCases *UseCases) GetUserToys(ctx context.Context, userID uint64) ([]entities.Toy, error) {
	master, err := useCases.mastersService.GetMasterByUserID(ctx, userID)
	if err != nil {
		return nil, nil
	}

	return useCases.GetMasterToys(ctx, master.ID)
}

func (useCases *UseCases) AddToy(ctx context.Context, rawToyData entities.RawAddToyDTO) (uint64, error) {
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

func (useCases *UseCases) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByID(ctx, id)
}

func (useCases *UseCases) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	return useCases.mastersService.GetMasterByUserID(ctx, userID)
}

func (useCases *UseCases) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	return useCases.mastersService.GetAllMasters(ctx)
}

func (useCases *UseCases) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	if _, err := useCases.ssoService.GetUserByID(ctx, masterData.UserID); err != nil {
		return 0, err
	}

	return useCases.mastersService.RegisterMaster(ctx, masterData)
}

func (useCases *UseCases) CreateTags(ctx context.Context, tagsData []entities.CreateTagDTO) ([]uint32, error) {
	allTags, err := useCases.tagsService.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	allTagsSet := make(map[string]uint32)
	for _, tag := range allTags {
		allTagsSet[tag.Name] = tag.ID
	}

	var existingTagIDs []uint32
	var tagsToCreate []entities.CreateTagDTO
	for _, tag := range tagsData {
		if _, ok := allTagsSet[tag.Name]; ok {
			existingTagIDs = append(existingTagIDs, allTagsSet[tag.Name])

			continue
		}

		tagsToCreate = append(tagsToCreate, tag)
	}

	createdTagIDs, err := useCases.tagsService.CreateTags(ctx, tagsToCreate)
	if err != nil {
		return nil, err
	}

	createdTagIDs = append(createdTagIDs, existingTagIDs...)
	return createdTagIDs, nil
}

func (useCases *UseCases) DeleteToy(ctx context.Context, id uint64) error {
	if _, err := useCases.toysService.GetToyByID(ctx, id); err != nil {
		return err
	}

	return useCases.toysService.DeleteToy(ctx, id)
}

func (useCases *UseCases) UpdateToy(ctx context.Context, rawToyData entities.RawUpdateToyDTO) error {
	toy, err := useCases.toysService.GetToyByID(ctx, rawToyData.ID)
	if err != nil {
		return err
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

func (useCases *UseCases) UpdateMaster(ctx context.Context, masterData entities.UpdateMasterDTO) error {
	if _, err := useCases.mastersService.GetMasterByID(ctx, masterData.ID); err != nil {
		return err
	}

	return useCases.mastersService.UpdateMaster(ctx, masterData)
}
