package controllers

import (
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"secmail/config"
	"secmail/models"

	"github.com/gin-gonic/gin"
	"github.com/kuun/slog"
	"gorm.io/gorm"
)

type _logger struct {
}

var log = slog.GetLogger(_logger{})

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
	Address   string    `json:"address"`
	ExpiresAt time.Time `json:"expiresAt"`
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

	// Get client info
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Generate random email address
	var email models.EmailAddress
	var exists bool

	// Keep generating until we find an unused address
	for {
		randomPart := generateRandomString(emailLength)
		emailAddress := randomPart + "@" + config.GlobalConfig.EmailDomain

		exists = db.Where("address = ?", emailAddress).First(&email).Error == nil
		if !exists {
			email = models.EmailAddress{
				Address:      emailAddress,
				ExpiresAt:   time.Now().Add(emailLifespan),
				CreatorIP:   clientIP,
				CreatorAgent: userAgent,
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
		Address:   email.Address,
		ExpiresAt: email.ExpiresAt,
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

	var email models.EmailAddress
	if err := db.Where("address = ?", emailAddress).First(&email).Error; err != nil {
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

	// email to JSON CreateEmailResponse
	emailResponse := CreateEmailResponse{
		Address:   email.Address,
		ExpiresAt: email.ExpiresAt,
	}
	c.JSON(http.StatusOK, emailResponse)
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
	var email models.EmailAddress
	if err := tx.Where("address = ?", emailAddress).First(&email).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Set client info for audit log
	email.CreatorIP = c.ClientIP()
	email.CreatorAgent = c.GetHeader("User-Agent")

	// Delete associated data (attachments will be deleted by cascade)
	if err := tx.Where("email_id = ?", email.ID).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete messages"})
		return
	}

	// Delete the email
	if err := tx.Delete(&email).Error; err != nil {
		tx.Rollback()
		log.Errorf("Failed to delete email address: %s, error: %s", email.Address, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Errorf("Failed to commit transaction for email address: %s, error: %s", email.Address, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.Status(http.StatusNoContent)
}
