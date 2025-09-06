# Social Media Scheduler - Go Full Stack Application

A comprehensive social media scheduling platform built with Go, featuring a complete web interface for managing posts, analytics, and social media accounts across multiple platforms.

## 🚀 Features

### Core Features
- **Multi-Platform Support**: Schedule posts to Facebook, Twitter, Instagram, LinkedIn, and TikTok
- **Smart Scheduling**: AI-powered optimal posting time recommendations
- **Analytics Dashboard**: Track engagement, reach, and performance metrics
- **Team Collaboration**: Role-based access control for team management
- **Media Library**: Organize and manage images, videos, and other media assets
- **Subscription Management**: Free, Basic, Professional, and Agency tiers

### Technical Features
- **Full Stack Go Application**: Backend API and web frontend served from a single binary
- **JWT Authentication**: Secure user authentication and session management
- **SQLite Database**: Lightweight database for development (easily switchable to PostgreSQL)
- **Responsive Design**: Bootstrap-based UI that works on all devices
- **RESTful API**: Complete REST API for all functionality
- **Auto-Migration**: Database schema automatically created and maintained

## 🛠 Technology Stack

- **Backend**: Go (Gin framework, GORM ORM)
- **Frontend**: HTML templates with Bootstrap, JavaScript
- **Database**: SQLite (development) / PostgreSQL (production)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Styling**: Bootstrap 5 with Font Awesome icons
- **Charts**: Chart.js for analytics visualization

## 📦 Installation

### Prerequisites
- Go 1.21 or later
- Git

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/immunesh/automated-social-media-scheduler.git
   cd automated-social-media-scheduler
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Build the application**
   ```bash
   go build -o scheduler
   ```

4. **Run the application**
   ```bash
   ./scheduler
   ```

5. **Open your browser**
   Navigate to `http://localhost:8080`

## 🔧 Configuration

The application can be configured using environment variables:

```bash
export DATABASE_PATH="./scheduler.db"          # Database file path
export JWT_SECRET="your-secret-key"            # JWT signing secret
export PORT="8080"                             # Server port
export GIN_MODE="release"                      # Gin mode (debug/release)
```

## 🐳 Docker Support

### Build Docker Image
```bash
docker build -t social-scheduler .
```

### Run with Docker
```bash
docker run -p 8080:8080 social-scheduler
```

### Docker Compose (Development)
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=debug
      - DATABASE_PATH=/data/scheduler.db
    volumes:
      - ./data:/data
```

## 📖 API Documentation

### Authentication Endpoints
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user

### Post Management
- `GET /api/posts` - Get all posts for authenticated user
- `POST /api/posts` - Create a new post
- `PUT /api/posts/:id` - Update an existing post
- `DELETE /api/posts/:id` - Delete a post

### Social Accounts
- `GET /api/social-accounts` - Get connected social media accounts
- `POST /api/social-accounts` - Connect a new social media account
- `DELETE /api/social-accounts/:id` - Disconnect a social media account

### Analytics
- `GET /api/analytics` - Get analytics data for user's posts

### User Profile
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update user profile

## 🏗 Project Structure

```
├── main.go                          # Application entry point
├── cmd/
│   └── server/                      # Server setup and routing
├── internal/
│   ├── auth/                        # Authentication utilities
│   ├── database/                    # Database connection and setup
│   ├── handlers/                    # HTTP request handlers
│   ├── middleware/                  # HTTP middleware
│   ├── models/                      # Database models
│   └── services/                    # Business logic services
├── web/
│   ├── templates/                   # HTML templates
│   └── static/                      # CSS, JavaScript, images
├── config/                          # Configuration management
├── migrations/                      # Database migrations
├── Dockerfile                       # Docker configuration
└── README.md                        # This file
```

## 👤 Usage

### Getting Started

1. **Register Account**: Create a new user account at `/register`
2. **Connect Social Media**: Go to Settings and connect your social media accounts
3. **Create Posts**: Use the Posts page to create and schedule content
4. **View Analytics**: Monitor your post performance in the Analytics dashboard

### Demo Mode

The application includes a demo mode for testing purposes:
- Social media connections are simulated
- Analytics data is generated for demonstration
- All core features are functional without actual API integrations

### Production Setup

For production deployment:

1. **Use PostgreSQL**: Switch from SQLite to PostgreSQL
2. **Set Environment Variables**: Configure proper JWT secrets and database credentials
3. **Enable HTTPS**: Use a reverse proxy like Nginx for SSL termination
4. **Social Media APIs**: Implement actual API integrations for each platform
5. **Monitoring**: Add application monitoring and logging

## 🔐 Security Features

- **Password Hashing**: Bcrypt for secure password storage
- **JWT Authentication**: Stateless authentication with configurable expiration
- **CSRF Protection**: Built-in protection against cross-site request forgery
- **Input Validation**: Comprehensive input validation and sanitization
- **SQL Injection Prevention**: GORM ORM provides built-in protection

## 📈 Subscription Plans

### Free Tier
- 1 user account
- 3 social media accounts
- 10 scheduled posts per month
- Basic analytics

### Basic Tier ($19/month)
- 2 user accounts
- 10 social media accounts
- 100 scheduled posts per month
- Standard analytics
- Basic AI suggestions

### Professional Tier ($49/month)
- 5 user accounts
- 25 social media accounts
- 500 scheduled posts per month
- Advanced analytics
- Full AI content suite

### Agency Tier ($99/month)
- 10 user accounts
- Unlimited social media accounts
- Unlimited scheduled posts
- White-label reports
- Priority support

## 🚧 Future Enhancements

- [ ] Real social media API integrations
- [ ] Advanced AI content generation
- [ ] Competitor analysis tools
- [ ] Automated A/B testing
- [ ] Content calendar planning
- [ ] White-label solution
- [ ] Mobile applications
- [ ] Advanced team collaboration features

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support, please create an issue in the GitHub repository or contact the development team.

## 🔧 Development

### Running in Development Mode
```bash
export GIN_MODE=debug
go run main.go
```

### Building for Production
```bash
CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o scheduler
```

### Running Tests
```bash
go test ./...
```

---

Built with ❤️ using Go and modern web technologies.