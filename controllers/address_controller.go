package controllers

import (
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"secmail/config"
	"secmail/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	emailLength   = 10
	emailLifespan = 1 * time.Hour
	charset       = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// Update email pattern to use dynamic domain
func getEmailPattern() string {
	escapedDomain := regexp.QuoteMeta(config.GlobalConfig.EmailDomain)
	return `^[a-z0-9]{10}@` + escapedDomain + `$`
}

type CreateEmailResponse struct {
	EmailAddress string    `json:"emailAddress"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func isValidEmailFormat(email string) bool {
	pattern := getEmailPattern()
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func CreateTempEmail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Generate random email address
	var email models.TempEmail
	var exists bool

	// Keep generating until we find an unused address
	for {
		randomPart := generateRandomString(emailLength)
		emailAddress := randomPart + "@" + config.GlobalConfig.EmailDomain

		exists = db.Where("email_address = ?", emailAddress).First(&email).Error == nil
		if !exists {
			email = models.TempEmail{
				EmailAddress: emailAddress,
				ExpiresAt:    time.Now().Add(emailLifespan),
			}
			break
		}
	}

	// Save to database
	if err := db.Create(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	// Return response
	c.JSON(http.StatusCreated, CreateEmailResponse{
		EmailAddress: email.EmailAddress,
		ExpiresAt:    email.ExpiresAt,
	})
}

func GetTempEmail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	emailAddress := c.Param("id")

	// Add email format validation
	if !isValidEmailFormat(emailAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var email models.TempEmail
	if err := db.Where("email_address = ?", emailAddress).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if email has expired
	if time.Now().After(email.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Email address has expired"})
		return
	}

	c.JSON(http.StatusOK, email)
}

func DeleteTempEmail(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	emailAddress := c.Param("id")

	// Validate email format
	if !isValidEmailFormat(emailAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Begin transaction
	tx := db.Begin()

	// Find email
	var email models.TempEmail
	if err := tx.Where("email_address = ?", emailAddress).First(&email).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Delete associated data (attachments will be deleted by cascade)
	if err := tx.Where("temp_email_id = ?", email.ID).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete messages"})
		return
	}

	// Delete the email
	if err := tx.Delete(&email).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}
