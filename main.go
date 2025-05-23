package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secmail/config"
	"secmail/controllers"

	"secmail/models"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	db, err := gorm.Open(postgres.Open(config.GlobalConfig.Database.DSN()))
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate models
	db.AutoMigrate(&models.TempEmail{}, &models.Message{}, &models.Attachment{})

	r := gin.Default()

	// Routes
	r.POST("/api/email", controllers.CreateTempEmail)
	r.GET("/api/email/:id", controllers.GetTempEmail)
	r.GET("/api/email/:id/messages", controllers.GetMessages)
	r.GET("/api/message/:id", controllers.GetMessage)
	r.GET("/api/message/:id/attachment/:attachmentId", controllers.GetAttachment)
	r.DELETE("/api/email/:id", controllers.DeleteTempEmail)

	r.Run(":8080")
}
