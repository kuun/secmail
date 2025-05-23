package controllers

import (
	"net/http"
	"secmail/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"receivedAt"`
}

type AttachmentResponse struct {
	ID          uuid.UUID `json:"id"`
	FileName    string    `json:"fileName"`
	ContentType string    `json:"contentType"`
}

type MessageDetailResponse struct {
	ID          uuid.UUID            `json:"id"`
	From        string               `json:"from"`
	Subject     string               `json:"subject"`
	Content     string               `json:"content"`
	HTMLContent string               `json:"htmlContent"`
	CreatedAt   time.Time            `json:"receivedAt"`
	Attachments []AttachmentResponse `json:"attachments"`
}

func GetMessages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	emailAddress := c.Param("id")

	// Validate email format
	if !isValidEmailFormat(emailAddress) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Check if email exists and not expired
	var email models.TempEmail
	if err := db.Where("email_address = ?", emailAddress).First(&email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email address not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if time.Now().After(email.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Email address has expired"})
		return
	}

	// Get messages
	var messages []models.Message
	var total int64

	db.Model(&models.Message{}).Where("temp_email_id = ?", email.ID).Count(&total)

	if err := db.Where("temp_email_id = ?", email.ID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	// Convert to response format
	response := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = MessageResponse{
			ID:        msg.ID,
			From:      msg.From,
			Subject:   msg.Subject,
			CreatedAt: msg.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": response,
		"total":    total,
		"page":     page,
		"size":     pageSize,
	})
}

func GetMessage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	var message models.Message
	if err := db.Preload("Attachments").First(&message, messageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if associated email has expired
	var email models.TempEmail
	if err := db.First(&email, message.TempEmailID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if time.Now().After(email.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Email has expired"})
		return
	}

	// Convert attachments to response format
	attachments := make([]AttachmentResponse, len(message.Attachments))
	for i, att := range message.Attachments {
		attachments[i] = AttachmentResponse{
			ID:          att.ID,
			FileName:    att.FileName,
			ContentType: att.ContentType,
		}
	}

	c.JSON(http.StatusOK, MessageDetailResponse{
		ID:          message.ID,
		From:        message.From,
		Subject:     message.Subject,
		Content:     message.Content,
		HTMLContent: message.HTMLContent,
		CreatedAt:   message.CreatedAt,
		Attachments: attachments,
	})
}

func GetAttachment(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Parse attachment ID
	attachmentID, err := uuid.Parse(c.Param("attachmentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	// Get attachment with associated message
	var attachment models.Attachment
	if err := db.Preload("Message").First(&attachment, attachmentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if associated email has expired
	var email models.TempEmail
	if err := db.First(&email, attachment.TempEmailID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if time.Now().After(email.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Email has expired"})
		return
	}

	// Set response headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+attachment.FileName)
	c.Header("Content-Type", attachment.ContentType)
	c.Data(http.StatusOK, attachment.ContentType, attachment.Data)
}
