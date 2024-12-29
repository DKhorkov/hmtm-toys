package entities

import "time"

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterMasterDTO struct {
	UserID uint64 `json:"user_id"`
	Info   string `json:"info"`
}
