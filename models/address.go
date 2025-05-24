package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailAddress struct {
	gorm.Model
	Address   string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
	Messages  []Message `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}
