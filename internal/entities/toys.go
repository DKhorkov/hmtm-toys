package entities

import "time"

type Toy struct {
	ID          uint64    `json:"id"`
	MasterID    uint64    `json:"masterID"`
	CategoryID  uint32    `json:"categoryID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Quantity    uint32    `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `json:"tags"`
}

type AddToyDTO struct {
	MasterID    uint64   `json:"masterID"`
	CategoryID  uint32   `json:"categoryID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagsIDs     []uint32 `json:"tags"`
}

type RawAddToyDTO struct {
	AccessToken string   `json:"accessToken"`
	CategoryID  uint32   `json:"categoryID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float32  `json:"price"`
	Quantity    uint32   `json:"quantity"`
	TagsIDs     []uint32 `json:"tags"`
}
