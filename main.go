package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secmail/config"
	"secmail/controllers"
	"secmail/jobs"
	"secmail/models"
	"secmail/smtp"
)

// Middleware to inject DB into context
func injectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func main() {
	if err := config.InitConfig(); err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	db, err := gorm.Open(postgres.Open(config.GlobalConfig.Database.DSN()))
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate models
	db.AutoMigrate(
		&models.EmailAddress{},
		&models.Message{},
		&models.Attachment{},
		&models.AuditLog{},
	)

	r := gin.Default()

	// Add DB middleware to all routes
	r.Use(injectDB(db))

	// Routes
	r.POST("/api/email", controllers.CreateTempEmail)
	r.GET("/api/email/:id", controllers.GetTempEmail)
	r.GET("/api/email/:id/messages", controllers.GetMessages)
	r.GET("/api/message/:id", controllers.GetMessage)
	r.DELETE("/api/message/:id", controllers.DeleteMessage)
	r.GET("/api/message/:id/attachment/:attachmentId", controllers.GetAttachment)
	r.DELETE("/api/email/:id", controllers.DeleteTempEmail)

	// Start SMTP server in a goroutine
	go func() {
		if err := smtp.StartSMTPServer(db); err != nil {
			log.Printf("SMTP server error: %v", err)
		}
	}()

	// Start cleanup job
	jobs.StartCleanupJob(db)

	r.Run(":8080")
}
