package models

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"index;not null" json:"token"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userId"`
	Type      string    `gorm:"not null" json:"type"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	Blacklisted bool    `gorm:"default:false" json:"blacklisted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}