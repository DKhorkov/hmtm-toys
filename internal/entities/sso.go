package entities

import (
	"time"
)

type User struct {
	ID                uint64    `json:"id"`
	DisplayName       string    `json:"displayName"`
	Email             string    `json:"email"`
	EmailConfirmed    bool      `json:"emailConfirmed"`
	Password          string    `json:"password"`
	Phone             *string   `json:"phone,omitempty"`
	PhoneConfirmed    bool      `json:"phoneConfirmed"`
	Telegram          *string   `json:"telegram,omitempty"`
	TelegramConfirmed bool      `json:"telegramConfirmed"`
	Avatar            *string   `json:"avatar,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
