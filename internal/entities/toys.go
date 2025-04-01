package entities

import "time"

type Toy struct {
	ID          uint64       `json:"id"`
	MasterID    uint64       `json:"masterId"`
	CategoryID  uint32       `json:"categoryId"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       float32      `json:"price"`
	Quantity    uint32       `json:"quantity"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Tags        []Tag        `json:"tags,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	ID        uint64    `json:"id"`
	ToyID     uint64    `json:"toyId"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AddToyDTO struct {
	MasterID    uint64   `json:"masterId"`
	CategoryID  uint32   `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type RawAddToyDTO struct {
	UserID      uint64   `json:"userId"`
	CategoryID  uint32   `json:"categoryId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type UpdateToyDTO struct {
	ID                    uint64   `json:"id"`
	CategoryID            *uint32  `json:"categoryId,omitempty"`
	Name                  *string  `json:"name,omitempty"`
	Description           *string  `json:"description,omitempty"`
	Price                 *float32 `json:"price,omitempty"`
	Quantity              *uint32  `json:"quantity,omitempty"`
	TagIDsToAdd           []uint32 `json:"tagIdsToAdd,omitempty"`
	TagIDsToDelete        []uint32 `json:"tagIdsToDelete,omitempty"`
	AttachmentsToAdd      []string `json:"attachmentsToAdd,omitempty"`
	AttachmentIDsToDelete []uint64 `json:"attachmentIdsToDelete,omitempty"`
}

type RawUpdateToyDTO struct {
	ID          uint64   `json:"id"`
	CategoryID  *uint32  `json:"categoryId,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float32 `json:"price,omitempty"`
	Quantity    *uint32  `json:"quantity,omitempty"`
	TagIDs      []uint32 `json:"tagIds,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}
