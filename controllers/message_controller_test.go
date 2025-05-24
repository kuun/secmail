package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"secmail/config"
	"secmail/models"
)

func TestGetMessages(t *testing.T) {
	config.GlobalConfig = config.Config{
		EmailDomain: "test.com",
	}

	tests := []struct {
		name         string
		emailAddress string
		query        string
		setupDB      func(*gorm.DB)
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:         "Success With Messages",
			emailAddress: "abcd123456@test.com",
			query:        "?page=1&size=2",
			setupDB: func(db *gorm.DB) {
				email := models.EmailAddress{
					Address:   "abcd123456@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				db.Create(&email)

				// Create test messages
				for i := 0; i < 3; i++ {
					msg := models.Message{
						ID:      uuid.New(),
						EmailID: email.ID,
						From:    fmt.Sprintf("sender%d@example.com", i),
						Subject: fmt.Sprintf("Test Subject %d", i),
					}
					db.Create(&msg)
				}
			},
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Messages []MessageResponse `json:"messages"`
					Total    int64             `json:"total"`
					Page     int               `json:"page"`
					Size     int               `json:"size"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(3), response.Total)
				assert.Equal(t, 2, len(response.Messages))
				assert.Equal(t, 1, response.Page)
				assert.Equal(t, 2, response.Size)
			},
		},
		{
			name:         "Empty Messages",
			emailAddress: "empty12345@test.com",
			query:        "",
			setupDB: func(db *gorm.DB) {
				email := models.EmailAddress{
					Address:   "empty12345@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				db.Create(&email)
			},
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response struct {
					Messages []MessageResponse `json:"messages"`
					Total    int64             `json:"total"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), response.Total)
				assert.Empty(t, response.Messages)
			},
		},
		{
			name:         "Invalid Email Format",
			emailAddress: "invalid@email",
			query:        "",
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
			emailAddress: "expired123@test.com",
			query:        "",
			setupDB: func(db *gorm.DB) {
				email := models.EmailAddress{
					Address:   "expired123@test.com",
					ExpiresAt: time.Now().Add(-time.Hour),
				}
				db.Create(&email)
			},
			expectedCode: http.StatusGone,
			validate: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Email address has expired")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			tt.setupDB(db)
			router := setupTestRouter(db)
			router.GET("/email/:id/messages", GetMessages)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/email/"+tt.emailAddress+"/messages"+tt.query, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.validate != nil {
				tt.validate(t, w)
			}
		})
	}
}

func TestGetAttachment(t *testing.T) {
	tests := []struct {
		name         string
		attachmentID string
		setupDB      func(*gorm.DB) *models.Attachment
		expectedCode int
		validate     func(*testing.T, *httptest.ResponseRecorder, *models.Attachment)
	}{
		{
			name:         "Success",
			attachmentID: "123e4567-e89b-12d3-a456-426614174000",
			setupDB: func(db *gorm.DB) *models.Attachment {
				email := models.EmailAddress{
					Address:   "test123456@test.com",
					ExpiresAt: time.Now().Add(time.Hour),
				}
				db.Create(&email)

				msg := models.Message{
					ID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
					EmailID: email.ID,
					From:    "sender@example.com",
				}
				db.Create(&msg)

				att := models.Attachment{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					MessageID:   msg.ID,
					FileName:    "test.txt",
					ContentType: "text/plain",
					Data:        []byte("test content"),
				}
				db.Create(&att)
				return &att
			},
			expectedCode: http.StatusOK,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, att *models.Attachment) {
				assert.Equal(t, att.ContentType, w.Header().Get("Content-Type"))
				assert.Equal(t, "attachment; filename="+att.FileName, w.Header().Get("Content-Disposition"))
				assert.Equal(t, string(att.Data), w.Body.String())
			},
		},
		{
			name:         "Invalid UUID",
			attachmentID: "invalid-uuid",
			setupDB:      func(db *gorm.DB) *models.Attachment { return nil },
			expectedCode: http.StatusBadRequest,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, _ *models.Attachment) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Invalid attachment ID")
			},
		},
		{
			name:         "Expired Email",
			attachmentID: "123e4567-e89b-12d3-a456-426614174000",
			setupDB: func(db *gorm.DB) *models.Attachment {
				email := models.EmailAddress{
					Address:   "expired@test.com",
					ExpiresAt: time.Now().Add(-time.Hour), // expired
				}
				db.Create(&email)

				msg := models.Message{
					ID:      uuid.New(),
					EmailID: email.ID,
				}
				db.Create(&msg)

				att := models.Attachment{
					ID:        uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					MessageID: msg.ID,
				}
				db.Create(&att)
				return &att
			},
			expectedCode: http.StatusGone,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, _ *models.Attachment) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Attachment not found or email expired")
			},
		},
		{
			name:         "Not Found",
			attachmentID: "123e4567-e89b-12d3-a456-426614174000",
			setupDB:      func(db *gorm.DB) *models.Attachment { return nil },
			expectedCode: http.StatusGone,
			validate: func(t *testing.T, w *httptest.ResponseRecorder, _ *models.Attachment) {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], "Attachment not found or email expired")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			att := tt.setupDB(db)
			router := setupTestRouter(db)
			router.GET("/attachment/:attachmentId", GetAttachment)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/attachment/"+tt.attachmentID, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			tt.validate(t, w, att)
		})
	}
}
