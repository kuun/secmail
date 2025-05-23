package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TempEmailID uint
	From        string
	Subject     string
	Content     string
	HTMLContent string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Attachments []Attachment
}

type Attachment struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	TempEmailID uint
	MessageID   uint
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
