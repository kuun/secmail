package jobs

import (
	"time"

	"secmail/models"

	"github.com/kuun/slog"
	"gorm.io/gorm"
)

type __logger struct{}

var log = slog.GetLogger(__logger{})

func CleanupExpiredEmails(db *gorm.DB) {
	result := db.Where("expires_at < ?", time.Now()).Unscoped().Delete(&models.EmailAddress{})
	if result.Error != nil {
		log.Warnf("Failed to cleanup expired emails: %v", result.Error)
		return
	}
	log.Warnf("Cleaned up %d expired emails", result.RowsAffected)
}

func StartCleanupJob(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			CleanupExpiredEmails(db)
		}
	}()
}
