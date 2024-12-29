package entities

import "time"

type Toy struct {
	ID          uint64    `json:"id"`
	MasterID    uint64    `json:"master_id"`
	CategoryID  uint32    `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []Tag     `json:"tags"`
}

type AddToyDTO struct {
	MasterID    uint64   `json:"master_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagsIDs     []uint32 `json:"tag_ids"`
}

type RawAddToyDTO struct {
	UserID      uint64   `json:"user_id"`
	CategoryID  uint32   `json:"category_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagsIDs     []uint32 `json:"tag_ids"`
}
