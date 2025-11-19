# MiGaleria Backend

A RESTful API backend for MiGaleria - a multimedia gallery application for sharing and managing photos and short videos. Built with Go, PostgreSQL, and JWT authentication.

[![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?logo=go)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-336791?logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

> This project is the backend for [MiGaleria](https://github.com/Martins-Iroka/MiGaleria), which is currently undergoing refactoring.

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
- [Running the Application](#-running-the-application)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Database Schema](#-database-schema)
- [Authentication](#-authentication)
- [Development](#-development)

## ✨ Features

### Authentication & Authorization
- ✅ **User Registration** with email verification via Twilio OTP
- ✅ **JWT Authentication** with access and refresh tokens
- ✅ **Token Refresh** mechanism for secure session management
- ✅ **Secure Password Hashing** using bcrypt
- ✅ **Automatic Token Cleanup** - Daily background job removes expired refresh tokens

### Photo Management
- ✅ **Photo Posts** - Create and view photo galleries
- ✅ **Comments** - Add comments to photo posts
- ✅ **Pagination** - Efficient data retrieval with configurable limits

### Video Management
- ✅ **Video Posts** with thumbnail support
- ✅ **Multiple Download Formats** - Store multiple video quality/size options
- ✅ **Video Comments** - Engage with video content
- ✅ **Duration Tracking** - Store and retrieve video duration metadata

### Security & Performance
- ✅ **CORS Configuration** - Secure cross-origin resource sharing
- ✅ **Request Timeout Protection** - 60-second request timeout
- ✅ **Database Connection Pooling** - Optimized database performance
- ✅ **Index Optimization** - Strategic database indexes for query performance
- ✅ **SQL Injection Protection** - Parameterized queries throughout

## 🛠 Tech Stack

### Core
- **[Go 1.25.3](https://golang.org/)** - High-performance backend language
- **[Chi Router](https://github.com/go-chi/chi)** - Lightweight, composable HTTP router
- **[PostgreSQL](https://www.postgresql.org/)** - Robust relational database

### Authentication & Security
- **[golang-jwt/jwt](https://github.com/golang-jwt/jwt)** - JWT token generation and validation
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** - Password hashing
- **[Twilio Verify](https://www.twilio.com/verify)** - OTP verification service

### Database & Migrations
- **[lib/pq](https://github.com/lib/pq)** - PostgreSQL driver
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - Database migration management

### Development Tools
- **[Zap](https://github.com/uber-go/zap)** - High-performance structured logging
- **[Swagger](https://swagger.io/)** - API documentation generation
- **[Air](https://github.com/cosmtrek/air)** - Live reload for Go apps
- **[Docker](https://www.docker.com/)** - PostgreSQL containerization

### Validation & Middleware
- **[go-playground/validator](https://github.com/go-playground/validator)** - Request validation
- **[Chi CORS](https://github.com/go-chi/cors)** - CORS middleware

## 🏗 Architecture

```
┌─────────────┐
│   Client    │
│  (Frontend) │
└──────┬──────┘
       │ HTTP/REST
       ▼
┌─────────────────────────────────┐
│      Chi Router (Middleware)     │
│  ┌────────────────────────────┐ │
│  │ CORS │ Auth │ Timeout │Log │ │
│  └────────────────────────────┘ │
└────────────┬────────────────────┘
             │
    ┌────────┴────────┐
    ▼                 ▼
┌─────────┐      ┌──────────┐
│  Auth   │      │   API    │
│ Handlers│      │ Handlers │
└────┬────┘      └─────┬────┘
     │                 │
     └────────┬────────┘
              ▼
       ┌─────────────┐
       │   Storage   │
       │   Layer     │
       └──────┬──────┘
              ▼
       ┌─────────────┐
       │ PostgreSQL  │
       │  Database   │
       └─────────────┘
```

### Layered Architecture

1. **Router Layer** (`cmd/api/`)
   - HTTP routing and middleware
   - Request/response handling
   - CORS and authentication middleware

2. **Handler Layer** (`cmd/api/authApi.go`, etc.)
   - Business logic coordination
   - Request validation
   - Response formatting

3. **Storage Layer** (`internal/user/`, `internal/video/`, `internal/photo/`)
   - Database operations
   - Data access logic
   - Query optimization

4. **Authentication Layer** (`internal/auth/`)
   - JWT token management
   - Password hashing
   - OTP verification

## 🚀 Getting Started

### Prerequisites

- **Go 1.25.3+** - [Download](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **golang-migrate** - [Installation Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- **Make** (optional but recommended)
- **Twilio Account** - [Sign up](https://www.twilio.com/try-twilio) for OTP verification

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Martins-Iroka/MyGallery-Backend.git
   cd MyGallery-Backend
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Install development tools** (optional)
   ```bash
   # Install Air for live reloading
   go install github.com/cosmtrek/air@latest

   # Install staticcheck for linting
   go install honnef.co/go/tools/cmd/staticcheck@latest

   # Install Swagger CLI
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### Configuration

1. **Create `.envrc` file** in the project root:
   ```bash
   # Server Configuration
   export ADDR=":8080"
   export EXTERNAL_URL="localhost:8080"
   export ENV="development"

   # Database Configuration
   export DB_ADDR="postgres://admin:adminpassword@localhost/mygallery?sslmode=disable"
   export DB_MAX_OPEN_CONNS=30
   export DB_MAX_IDLE_CONNS=30
   export DB_MAX_IDLE_TIME="15m"

   # Twilio Configuration (for OTP verification)
   export TWILIO_ACCOUNT_SID="your_account_sid"
   export TWILIO_AUTH_TOKEN="your_auth_token"
   export TWILIO_SID="your_service_sid"

   # JWT Configuration
   export AUTH_TOKEN_SECRET="your_secret_key_here"

   # CORS Configuration
   export CORS_ALLOWED_ORIGIN="http://localhost:3000"
   ```

2. **Load environment variables**
   ```bash
   source .envrc
   # Or use direnv: https://direnv.net/
   ```

### Database Setup

1. **Start PostgreSQL with Docker**
   ```bash
   docker-compose up -d
   ```

2. **Verify database is running**
   ```bash
   docker ps
   # Look for 'mygallery-db' container
   ```

3. **Run database migrations**
   ```bash
   make migrate-up
   ```

4. **Verify migration status**
   ```bash
   docker exec -it mygallery-db psql -U admin -d mygallery -c "SELECT * FROM schema_migrations ORDER BY version;"
   ```

### Create New Migrations (Optional)

```bash
# Create a new migration
make migrations your_migration_name

# This creates two files:
# - cmd/migrate/migrations/000XXX_your_migration_name.up.sql
# - cmd/migrate/migrations/000XXX_your_migration_name.down.sql
```

## 🏃 Running the Application

### Development Mode

**Option 1: Using Air (with live reload)**
```bash
air
```

**Option 2: Using Go directly**
```bash
go run cmd/api/*.go
```

**Option 3: Using Make**
```bash
make build
./bin/main
```

The server will start at `http://localhost:8080`

### Production Build

```bash
# Build binary
go build -o bin/main cmd/api/*.go

# Run binary
./bin/main
```

### Health Check

```bash
curl http://localhost:8080/v1/health
```

Expected response:
```json
{
  "status": "ok",
  "environment": "development",
  "version": "1.4.0"
}
```

## 📚 API Documentation

### Swagger UI

Once the server is running, access interactive API documentation at:

```
http://localhost:8080/v1/swagger/index.html
```

### Generate/Update Swagger Docs

```bash
make gen-docs
```

### API Endpoints Overview

#### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/v1/authentication/register` | Register new user | No |
| POST | `/v1/authentication/verify` | Verify email with OTP | No |
| POST | `/v1/authentication/login` | Login user | No |
| POST | `/v1/authentication/refresh` | Refresh access token | No |
| POST | `/v1/authentication/logout` | Logout and revoke tokens | Yes |

#### Photos
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v1/photos` | Get paginated photos | Yes |
| POST | `/v1/photos/{postID}/create-comment` | Add comment to photo | Yes |
| GET | `/v1/photos/{postID}/comments` | Get photo comments | Yes |

#### Videos
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/v1/videos` | Get paginated videos | Yes |
| POST | `/v1/videos/{postID}/create-comment` | Add comment to video | Yes |
| GET | `/v1/videos/{postID}/comments` | Get video comments | Yes |

### Authentication Flow

1. **Register**: `POST /v1/authentication/register`
   ```json
   {
     "username": "johndoe",
     "email": "john@example.com",
     "password": "securePassword123"
   }
   ```

2. **Verify**: OTP sent via Twilio SMS
   ```json
   {
     "email": "john@example.com",
     "code": "123456"
   }
   ```

3. **Login**: Receive access token (15mins) and refresh token (7 days)
   ```json
   {
     "email": "john@example.com",
     "password": "securePassword123"
   }
   ```

4. **Use Access Token**: Add to Authorization header
   ```
   Authorization: Bearer <access_token>
   ```

5. **Refresh Token**: When access token expires
   ```json
   {
     "refresh_token": "<refresh_token>"
   }
   ```

## 📁 Project Structure

```
MyGallery-Backend/
├── cmd/
│   ├── api/                    # API handlers and routing
│   │   ├── main.go            # Application entry point
│   │   ├── api.go             # Router setup and middleware
│   │   ├── authApi.go         # Authentication endpoints
│   │   └── util/              # API utilities
│   │       ├── errors.go      # Error handling
│   │       └── json.go        # JSON helpers
│   └── migrate/
│       └── migrations/        # Database migration files
│           ├── 000001_create_users.up.sql
│           ├── 000002_users_verification.up.sql
│           ├── 000003_create_photo_post.up.sql
│           ├── 000004_create_video_posts.up.sql
│           ├── 000005_create_video_download_files.up.sql
│           ├── 000006_create_photos_comment.up.sql
│           ├── 000007_create_videos_comment.up.sql
│           ├── 000008_create_refresh_tokens.up.sql
│           ├── 000009_add_indexes.up.sql
│           └── 000010_add_image_column_in_video_post.up.sql
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── storage.go             # Storage layer aggregation
│   ├── auth/                  # Authentication logic
│   │   ├── jwt.go            # JWT token management
│   │   ├── password.go       # Password hashing
│   │   └── twilio.go         # OTP verification
│   ├── db/
│   │   └── db.go             # Database connection
│   ├── env/
│   │   └── env.go            # Environment variable helpers
│   ├── user/
│   │   └── userstore.go      # User data access
│   ├── photo/
│   │   └── photostore.go     # Photo data access
│   └── video/
│       └── videostore.go     # Video data access
├── docs/                      # Swagger documentation (auto-generated)
├── bin/                       # Compiled binaries
├── docker-compose.yml         # PostgreSQL Docker setup
├── Makefile                   # Build and migration commands
├── go.mod                     # Go module dependencies
└── README.md                  # This file
```

## 🗄 Database Schema

### Users
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Refresh Tokens
```sql
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
```

### Photo Posts
```sql
CREATE TABLE photo_posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    caption TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Video Posts
```sql
CREATE TABLE video_posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_url TEXT NOT NULL,
    video_image TEXT,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Comments
```sql
CREATE TABLE photos_comment (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES photo_posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE video_comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES video_posts(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔐 Authentication

### JWT Token Strategy

**Access Token**
- **Expiration**: 15 minutes
- **Purpose**: API authentication
- **Storage**: Client-side (memory/state, not localStorage)
- **Claims**: `user_id`, `email`, `exp`, `iss`

**Refresh Token**
- **Expiration**: 7 days
- **Purpose**: Generate new access tokens
- **Storage**: Database (hashed) + Client-side (httpOnly cookie recommended)
- **Revocation**: Supports logout and automatic cleanup

### Token Refresh Flow

1. Access token expires (15 minutes)
2. Client calls `/v1/authentication/refresh` with refresh token
3. Server validates refresh token (not expired, not revoked, user exists)
4. Server generates new access token
5. Client uses new access token

### Security Features

- ✅ **Passwords**: Bcrypt hashing with salt
- ✅ **Tokens**: HMAC-SHA256 signed JWT
- ✅ **Refresh Tokens**: SHA-256 hashed before database storage
- ✅ **Token Cleanup**: Daily cron job removes expired tokens
- ✅ **SQL Injection**: Parameterized queries
- ✅ **CORS**: Configurable allowed origins

### Token Expiration Handling

When refresh token expires (after 7 days):
- User must log in again with credentials
- All existing tokens for that session are invalid
- Frontend should redirect to login page

## 🔧 Development

### Code Linting

```bash
make staticcheck
```

### Database Management

```bash
# Migrate up
make migrate-up

# Migrate down (rollback 1 migration)
make migrate-down

# Migrate down N steps
make migrate-down 2

# Create new migration
make migrations add_user_profile
```

### Background Jobs

The application runs a daily cleanup job that removes expired refresh tokens:

```go
// Runs daily at 24-hour intervals
app.startCleanupJob()
```

This ensures the `refresh_tokens` table doesn't grow indefinitely.


### Manual API Testing

Use the Swagger UI or tools like:
- **Postman**: Import from Swagger JSON
- **cURL**: Command-line testing
- **HTTPie**: User-friendly HTTP client

Example cURL:
```bash
# Health check
curl http://localhost:8080/v1/health

# Register user
curl -X POST http://localhost:8080/v1/authentication/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"test123"}'

# Login
curl -X POST http://localhost:8080/v1/authentication/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](http://www.apache.org/licenses/LICENSE-2.0.html) file for details.

## 👥 Authors

- **Ikechukwu Iroka** - [GitHub](https://github.com/Martins-Iroka)
