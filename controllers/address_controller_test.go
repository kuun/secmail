package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"secmail/config"
	"secmail/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.EmailAddress{}, &models.Message{}, &models.Attachment{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	return r
}

func TestCreateTempEmail(t *testing.T) {
	// Initialize test config
	config.GlobalConfig = config.Config{
		EmailDomain: "test.com",
	}

	tests := []struct {
		name         string
		setupDB      func(*gorm.DB)
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "Success",
			setupDB:      func(db *gorm.DB) {},
			expectedCode: http.StatusCreated,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response CreateEmailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Validate email format
				pattern := regexp.MustCompile(`^[a-z0-9]{10}@test\.com$`)
				assert.True(t, pattern.MatchString(response.Address))

				// Validate expiration time
				expectedExpiry := time.Now().Add(emailLifespan)
				assert.WithinDuration(t, expectedExpiry, response.ExpiresAt, 2*time.Second)
			},
		},
		{
			name: "Database Error",
			setupDB: func(db *gorm.DB) {
				// Drop the table to simulate DB error
				db.Migrator().DropTable(&models.EmailAddress{})
			},
			expectedCode: http.StatusInternalServerError,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], "Failed to create email")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			tt.setupDB(db)
			router := setupTestRouter(db)
			router.POST("/email", CreateTempEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/email", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestGetTempEmail(t *testing.T) {
	config.GlobalConfig = config.Config{
		EmailDomain: "test.com",
	}

	tests := []struct {
		name         string
		emailAddress string
		setupDB      func(*gorm.DB)
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "Success",
			emailAddress: "abcd123456@test.com",
			setupDB: func(db *gorm.DB) {
				db.Create(&models.EmailAddress{
					Address:   "abcd123456@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				})
			},
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response models.EmailAddress
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "abcd123456@test.com", response.Address)
			},
		},
		{
			name:         "Invalid Format",
			emailAddress: "invalid@email",
			setupDB:      func(db *gorm.DB) {},
			expectedCode: http.StatusBadRequest,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid email format")
			},
		},
		{
			name:         "Expired Email",
			emailAddress: "abcd123456@test.com",
			setupDB: func(db *gorm.DB) {
				db.Create(&models.EmailAddress{
					Address:   "abcd123456@test.com",
					ExpiresAt: time.Now().Add(-time.Hour), // expired
				})
			},
			expectedCode: http.StatusGone,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Email address has expired")
			},
		},
		{
			name:         "Not Found",
			emailAddress: "notfound12@test.com",
			setupDB:      func(db *gorm.DB) {},
			expectedCode: http.StatusNotFound,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Email address not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			tt.setupDB(db)
			router := setupTestRouter(db)
			router.GET("/email/:id", GetTempEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/email/"+tt.emailAddress, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestDeleteTempEmail(t *testing.T) {
	config.GlobalConfig = config.Config{
		EmailDomain: "test.com",
	}

	tests := []struct {
		name         string
		emailAddress string
		setupDB      func(*gorm.DB)
		expectedCode int
		validate     func(*testing.T, *gorm.DB)
	}{
		{
			name:         "Success",
			emailAddress: "abcd123456@test.com",
			setupDB: func(db *gorm.DB) {
				email := models.EmailAddress{
					Address:   "abcd123456@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				db.Create(&email)

				// Create associated message and attachment
				msg := models.Message{
					EmailID: email.ID,
					Subject: "Test",
				}
				db.Create(&msg)
			},
			expectedCode: http.StatusNoContent,
			validate: func(t *testing.T, db *gorm.DB) {
				// Verify email was deleted
				var count int64
				db.Model(&models.EmailAddress{}).Where("address = ?", "abcd123456@test.com").Count(&count)
				assert.Equal(t, int64(0), count)

				// Verify messages were deleted
				db.Model(&models.Message{}).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name:         "Invalid Format",
			emailAddress: "invalid@email",
			setupDB:      func(db *gorm.DB) {},
			expectedCode: http.StatusBadRequest,
			validate:     func(t *testing.T, db *gorm.DB) {},
		},
		{
			name:         "Not Found",
			emailAddress: "notfound12@test.com",
			setupDB:      func(db *gorm.DB) {},
			expectedCode: http.StatusNotFound,
			validate:     func(t *testing.T, db *gorm.DB) {},
		},
		{
			name:         "Database Error",
			emailAddress: "abcd123456@test.com",
			setupDB: func(db *gorm.DB) {
				email := models.EmailAddress{
					Address:   "abcd123456@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				db.Create(&email)
				// Drop tables to simulate DB error
				db.Migrator().DropTable(&models.Message{})
			},
			expectedCode: http.StatusInternalServerError,
			validate:     func(t *testing.T, db *gorm.DB) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			tt.setupDB(db)
			router := setupTestRouter(db)
			router.DELETE("/email/:id", DeleteTempEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/email/"+tt.emailAddress, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.validate != nil {
				tt.validate(t, db)
			}
		})
	}
}
