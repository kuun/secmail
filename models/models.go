package models

import (
	"time"

	"gorm.io/gorm"
)

type TempEmail struct {
	gorm.Model
	EmailAddress string `gorm:"uniqueIndex"`
	ExpiresAt    time.Time
}
