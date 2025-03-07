package entities

import "time"

type Toy struct {
	ID          uint64       `json:"id"`
	MasterID    uint64       `json:"master_id"`
	CategoryID  uint32       `json:"category_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       float32      `json:"price"`
	Quantity    uint32       `json:"quantity"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Tags        []Tag        `json:"tags,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	ID        uint64    `json:"id"`
	ToyID     uint64    `json:"toy_id"`
	Link      string    `json:"link"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddToyDTO struct {
	MasterID    uint64   `json:"master_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type RawAddToyDTO struct {
	UserID      uint64   `json:"user_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

type UpdateToyDTO struct {
	ID                    uint64   `json:"id"`
	CategoryID            uint32   `json:"category_id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Price                 float32  `json:"price"`
	Quantity              uint32   `json:"quantity"`
	TagIDsToAdd           []uint32 `json:"tag_ids_to_add,omitempty"`
	TagIDsToDelete        []uint32 `json:"tag_ids_to_delete,omitempty"`
	AttachmentsToAdd      []string `json:"attachments_to_add,omitempty"`
	AttachmentIDsToDelete []uint64 `json:"attachment_ids_to_delete,omitempty"`
}

type RawUpdateToyDTO struct {
	ID          uint64   `json:"id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagIDs      []uint32 `json:"tag_ids,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}
