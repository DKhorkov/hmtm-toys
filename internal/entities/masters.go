package entities

import "time"

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userId"`
	Info      *string   `json:"info,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterMasterDTO struct {
	UserID uint64  `json:"userId"`
	Info   *string `json:"info,omitempty"`
}

type UpdateMasterDTO struct {
	ID   uint64  `json:"id"`
	Info *string `json:"info,omitempty"`
}

type MastersFilters struct {
	Search              *string `json:"search,omitempty"`
	CreatedAtOrderByAsc *bool   `json:"createdAtOrderByAsc,omitempty"`
}
