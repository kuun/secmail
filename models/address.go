package models

import (
	"time"

	"gorm.io/gorm"
)

type EmailAddress struct {
	gorm.Model
	Address      string `gorm:"uniqueIndex"`
	ExpiresAt    time.Time
	Messages     []Message `gorm:"foreignKey:EmailID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	CreatorIP    string    `gorm:"-"`
	CreatorAgent string    `gorm:"-"`
}

func (e *EmailAddress) AfterCreate(tx *gorm.DB) error {
	return tx.Create(&AuditLog{
		EmailID:      e.ID,
		EmailAddress: e.Address,
		Action:       "create",
		IP:           e.CreatorIP,
		UserAgent:    e.CreatorAgent,
		CreatedAt:    time.Now(),
	}).Error
}

func (e *EmailAddress) BeforeDelete(tx *gorm.DB) error {
	return tx.Create(&AuditLog{
		EmailID:      e.ID,
		EmailAddress: e.Address,
		Action:       "delete",
		IP:           e.CreatorIP,
		UserAgent:    e.CreatorAgent,
		CreatedAt:    time.Now(),
	}).Error
}
