package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	EmailID     uint
	From        string
	Subject     string
	Content     string
	HTMLContent string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Attachments []Attachment   `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

type Attachment struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	MessageID   uuid.UUID
	FileName    string
	ContentType string
	Data        []byte
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
