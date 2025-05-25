package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	EmailID      uint
	EmailAddress string
	Action       string // e.g. "create", "delete", "access"
	IP           string
	UserAgent    string
	CreatedAt    time.Time
}
