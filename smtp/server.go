package smtp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"

	"secmail/config"
	"secmail/models"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
	"github.com/kuun/slog"
	"gorm.io/gorm"
)

type __logger struct{}

var log = slog.GetLogger(__logger{})

type Backend struct {
	db *gorm.DB
}

type Session struct {
	backend *Backend
	from    string
	to      string
}

func NewBackend(db *gorm.DB) *Backend {
	return &Backend{db: db}
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{backend: b}, nil
}

func (s *Session) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	// read the entire message into a buffer
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	// parse the email using enmime
	env, err := enmime.ReadEnvelope(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	// find the recipient email address in the database
	var addr models.EmailAddress
	if err := s.backend.db.Where("address = ? AND expires_at > ?", s.to, time.Now()).
		First(&addr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warnf("Email address %s not found or expired", s.to)
		} else {
			log.Errorf("Database error: %+v", err)
		}
		return nil
	}

	// create a new message
	msg := models.Message{
		ID:          uuid.New(),
		EmailID:     addr.ID,
		From:        s.from,
		Subject:     env.GetHeader("Subject"),
		Content:     env.Text,
		HTMLContent: env.HTML,
	}

	// save attachments if they exist
	for _, a := range env.Attachments {
		att := models.Attachment{
			ID:          uuid.New(),
			MessageID:   msg.ID,
			FileName:    a.FileName,
			ContentType: a.ContentType,
			Data:        a.Content,
		}
		msg.Attachments = append(msg.Attachments, att)
	}

	// save the message to the database
	return s.backend.db.Create(&msg).Error
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func createSMTPServer(be *Backend) *smtp.Server {
	s := smtp.NewServer(be)
	s.Domain = config.GlobalConfig.EmailDomain
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 1
	s.AllowInsecureAuth = true

	if config.GlobalConfig.SMTP.TLS.Enable {
		cert, err := tls.LoadX509KeyPair(
			config.GlobalConfig.SMTP.TLS.CertFile,
			config.GlobalConfig.SMTP.TLS.KeyFile,
		)
		if err != nil {
			log.Errorf("Failed to load TLS certificate: %v", err)
			return nil
		}
		s.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	}

	return s
}

func StartSMTPServer(db *gorm.DB) error {
	be := NewBackend(db)
	server := createSMTPServer(be)
	server.Addr = fmt.Sprintf("%s:%d",
		config.GlobalConfig.SMTP.Host,
		config.GlobalConfig.SMTP.Port)

	log.Infof("Starting SMTP server at %s (STARTTLS: %v)",
		server.Addr,
		config.GlobalConfig.SMTP.TLS.Enable)
	return server.ListenAndServe()
}
