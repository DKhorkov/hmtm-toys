package entities

import "time"

type Master struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"userID"`
	Info      string    `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Toys      []*Toy    `json:"toys"`
}

type RegisterMasterDTO struct {
	UserID uint64 `json:"userID"`
	Info   string `json:"info"`
}

type RawRegisterMasterDTO struct {
	AccessToken string `json:"accessToken"`
	Info        string `json:"info"`
}
