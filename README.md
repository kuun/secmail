# SecMail - Secure Temporary Email Service

A secure and easy-to-use temporary email service built with Go and Vue.js.

## Features

- One-click temporary email address generation
- Real-time inbox with auto-refresh
- Support for HTML emails and attachments
- Email address expiration (1 hour by default)
- Mobile-responsive design
- Audit logging for security
- Browser-based persistence
- Copy to clipboard functionality

## Tech Stack

### Backend
- Go
- Gin Web Framework
- GORM
- PostgreSQL
- SMTP Server

### Frontend
- Vue 3
- Pinia for state management
- Vue Router
- TailwindCSS
- Heroicons

## Setup

### Prerequisites
- Go 1.19+
- Node.js 16+
- PostgreSQL 13+

### Backend Setup
```bash
# Clone repository
git clone https://github.com/yourusername/secmail.git
cd secmail

# Install Go dependencies
go mod tidy

# Configure database
cp etc/secmail.example.yaml etc/secmail.yaml
# Edit secmail.yaml with your database credentials

# Start the server
go run main.go
```

### Frontend Setup
```bash
# Navigate to webapp directory
cd webapp

# Install dependencies
npm install

# Start development server
npm run dev
```

## Configuration

Edit `config/config.yaml` to configure:
- Email domain
- SMTP server settings
- Database connection
- Server port

## Security Features

- Email address expiration
- Audit logging of creation and deletion
- IP tracking
- User agent logging
- No registration required
- Data auto-cleanup

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
