# IAM Service - Identity & Access Management System

[![en](https://img.shields.io/badge/lang-en-blue.svg)](README.md) [![id](https://img.shields.io/badge/lang-id-red.svg)](README.id.md)

A comprehensive Identity & Access Management (IAM) system built with Go and modern cloud-native technologies. This service is application-agnostic and can be integrated with any system requiring authentication and authorization.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Running Locally](#running-locally)
- [API Documentation](#api-documentation)
- [Database Migrations](#database-migrations)
- [Development](#development)
- [Environment Variables](#environment-variables)

## Overview

IAM Service is a multi-tenant Identity & Access Management system that provides secure authentication, role-based access control, and comprehensive user management capabilities. It is designed to be consumed by any application requiring identity management.

## Features

### Authentication & Authorization
- 🔐 User registration with email verification (OTP)
- 🔑 Secure login with JWT (access & refresh tokens)
- 📱 6-digit PIN setup with security validations
- 🔄 Password reset with OTP verification
- 🚪 Logout with token revocation
- 👥 Multi-tenant support
- 🎭 Role-based access control (RBAC)
- 🛡️ Platform admin management

### User Management
- 📝 Profile completion workflow
- 👤 Special account registration (admin/approver)
- 🔒 Account security features (login attempts tracking, account locking)
- 📊 User activation tracking

### Security Features
- 🔐 Password strength validation
- 🔄 Password history tracking (last 5 passwords)
- 🔢 PIN validation (prevents weak patterns)
- 🚫 Rate limiting for OTP requests
- 🔒 Bcrypt password hashing
- 🎫 JWT-based authentication
- 🔑 HashiCorp Vault integration for secrets management

## Tech Stack

### Core
- **Language**: Go 1.24
- **Web Framework**: Fiber v2
- **Database**: PostgreSQL 18
- **ORM**: GORM

### Infrastructure
- **Cache**: Redis
- **Secrets Management**: HashiCorp Vault
- **Object Storage**: MinIO
- **Containerization**: Docker & Docker Compose

### Libraries & Tools
- **JWT**: golang-jwt/jwt v5
- **Validation**: go-playground/validator v10
- **Configuration**: Viper
- **Logging**: Uber Zap
- **Live Reload**: Air (development)

## Project Structure

```
iam-service/
├── cmd/                          # Application entry points
│   ├── http/                     # HTTP server entry point
│   ├── grpc/                     # gRPC server (future)
│   ├── graphql/                  # GraphQL server (future)
│   ├── websocket/                # WebSocket server (future)
│   ├── worker/                   # Background workers
│   └── migrate/                  # Database migration runner
├── config/                       # Configuration management
├── delivery/                     # Delivery layer (controllers/handlers)
│   ├── http/                     # HTTP controllers & routing
│   └── grpc/                     # gRPC handlers
├── entity/                       # Domain entities/models
├── iam/                          # Identity & Access Management modules
│   ├── auth/                     # Authentication use cases
│   ├── user/                     # User management use cases
│   ├── role/                     # Role management use cases
│   └── health/                   # Health check use cases
├── impl/                         # Infrastructure implementations
│   ├── postgres/                 # PostgreSQL repositories
│   ├── redis/                    # Redis cache implementation
│   ├── hashivault/               # Vault integration
│   ├── minio/                    # MinIO storage implementation
│   └── mailer/                   # Email service implementation
├── pkg/                          # Shared packages
│   ├── jwt/                      # JWT utilities
│   ├── errors/                   # Custom error types
│   ├── logger/                   # Logging utilities
│   ├── validator/                # Validation helpers
│   └── utils/                    # General utilities
├── infrastructure/               # Infrastructure setup
├── migration/                    # Database migration files
├── deployment/                   # Deployment configurations
│   ├── docker/                   # Docker & Docker Compose
│   └── k8s/                      # Kubernetes manifests
├── doc/                          # Documentation
│   ├── postman/                  # Postman API collection
│   └── openapi/                  # OpenAPI/Swagger specs
└── script/                       # Utility scripts
```

### Architecture Patterns

This project follows **Clean Architecture** principles:

- **Entity Layer**: Core business models (`entity/`)
- **Use Case Layer**: Business logic (`iam/*/iam/`)
- **Repository Layer**: Data access abstractions (`iam/*/contract/`)
- **Delivery Layer**: HTTP handlers (`delivery/http/`)
- **Infrastructure Layer**: External services (`impl/`, `infrastructure/`)

Each module (auth, user, role) follows this structure:
```
iam/module/
├── contract/          # Interfaces (repository, usecase)
├── internal/          # Use case implementations (private)
├── moduledto/         # Data Transfer Objects
└── exported.go        # Public API
```

## Prerequisites

Before running this project, ensure you have the following installed:

- **Go** 1.24 or higher ([Download](https://golang.org/dl/))
- **Docker** & **Docker Compose** ([Download](https://www.docker.com/get-started))
- **Make** (optional, for build automation)
- **Air** (optional, for hot reloading): `go install github.com/air-verse/air@latest`
- **PostgreSQL Client** (optional, for manual DB operations)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd iam-service
```

### 2. Set Up Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` and fill in the required values:

```env
# App Configuration
APP_NAME=iam-service
APP_ENV=development
APP_VERSION=1.0.0

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# PostgreSQL (when using Docker)
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=iam
POSTGRES_PASSWORD=iam_secret
POSTGRES_DB=iam_db
POSTGRES_SSL_MODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_secret
REDIS_DB=0

# Vault
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=vault_root_token

# MinIO
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minio_admin
MINIO_SECRET_KEY=minio_secret
MINIO_BUCKET=iam-storage
MINIO_USE_SSL=false

# JWT (Generate strong secrets for production!)
JWT_ACCESS_SECRET=your-super-secret-access-key-min-32-chars
JWT_REFRESH_SECRET=your-super-secret-refresh-key-min-32-chars
JWT_ISSUER=iam-service

# Email (Optional - for OTP emails)
EMAIL_PROVIDER=smtp
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your-email@gmail.com
EMAIL_SMTP_PASS=your-app-password
EMAIL_FROM_ADDRESS=noreply@example.com
EMAIL_FROM_NAME=IAM System
```

### 3. Start Infrastructure Services

Start PostgreSQL, Redis, Vault, and MinIO using Docker Compose:

```bash
cd deployment/docker
docker-compose up -d
```

Verify all services are running:

```bash
docker-compose ps
```

You should see:
- `iam-postgres` on port 5432
- `iam-redis` on port 6379
- `iam-vault` on port 8200
- `iam-minio` on ports 9000 (API) and 9001 (Console)

### 4. Install Go Dependencies

```bash
go mod download
```

## Running Locally

### Option 1: Using Air (Hot Reload - Recommended for Development)

Air automatically rebuilds and restarts your application when you make changes:

```bash
# Install Air if you haven't already
go install github.com/air-verse/air@latest

# Run with Air
air
```

The server will start on `http://localhost:8080` with hot reload enabled.

### Option 2: Using `go run`

```bash
go run cmd/http/main.go
```

### Option 3: Build and Run Binary

```bash
# Build the binary
go build -o bin/iam-service cmd/http/main.go

# Run the binary
./bin/iam-service
```

### Verify the Server is Running

```bash
# Health check
curl http://localhost:8080/api/v1/health/

# Should return: {"status":"healthy","timestamp":"..."}
```

## API Documentation

### Postman Collection

Import the comprehensive Postman collection located at:
```
doc/postman/iam-service-api-collection.json
```

The collection includes:
- ✅ All 16 endpoints with example requests
- ✅ Success and failure response examples
- ✅ Auto-script to save access tokens
- ✅ Environment variables setup

### Available Endpoints

#### Health Endpoints
- `GET /api/v1/health/` - Basic health check
- `GET /api/v1/health/ready` - Readiness probe
- `GET /api/v1/health/live` - Liveness probe

#### Authentication Endpoints
- `POST /api/v1/auth/register` - User self-registration
- `POST /api/v1/auth/register/special-account` - Create admin/approver accounts
- `POST /api/v1/auth/verify-otp` - Verify email OTP
- `POST /api/v1/auth/complete-profile` - Complete user profile
- `POST /api/v1/auth/resend-otp` - Resend OTP email
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/setup-pin` - Setup 6-digit PIN (requires JWT)
- `POST /api/v1/auth/request-password-reset` - Request password reset OTP
- `POST /api/v1/auth/reset-password` - Reset password with OTP

#### Role Management (Platform Admin Only)
- `POST /api/v1/roles/` - Create new role

#### User Management (Platform Admin Only)
- `POST /api/v1/users/` - Create new user with system role

### Quick Start API Testing

1. **Register a special account (Platform Admin)**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register/special-account \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "00000000-0000-0000-0000-000000000001",
    "email": "admin@example.com",
    "password": "AdminPass123!",
    "password_confirm": "AdminPass123!",
    "full_name": "System Admin",
    "phone": "+6281234567890",
    "role_code": "PLATFORM_ADMIN"
  }'
```

2. **Login**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "00000000-0000-0000-0000-000000000001",
    "email": "admin@example.com",
    "password": "AdminPass123!"
  }'
```

Save the `access_token` from the response for authenticated requests.

3. **Access protected endpoint**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/setup-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "pin": "147258",
    "pin_confirm": "147258"
  }'
```

## Database Migrations

Database migrations are located in the `migration/` directory and follow the naming pattern:
```
XXXXXX_description.up.sql    # Apply migration
XXXXXX_description.down.sql  # Rollback migration
```

### Available Migrations

1. `000001_extensions` - PostgreSQL extensions (uuid-ossp)
2. `000002_tenants` - Tenant management
3. `000003_users` - User, profile, credentials, security tables
4. `000004_auth_tokens` - Refresh tokens
5. `000005_saml_mfa` - SAML and MFA support
6. `000006_email_verification` - OTP verification
7. `000007_user_activation_tracking` - Activation workflow
8. `000008_roles_permissions` - RBAC system
9. `000009_user_roles` - User role assignments
10. `000010_refresh_token_family` - Token rotation
11. `000011_platform_admin` - Platform admin role

### Running Migrations

#### Manual Migration (Using psql)

```bash
# Connect to database
psql -h localhost -p 5432 -U iam -d iam_db

# Run migrations manually in order
\i migration/000001_extensions.up.sql
\i migration/000002_tenants.up.sql
# ... and so on
```

#### Using Migration Tool (migrate)

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run all migrations
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" up

# Rollback last migration
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" down 1

# Check migration version
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" version
```

## Development

### Code Structure Conventions

1. **Use Cases**: All business logic goes in `iam/*/iam/`
2. **DTOs**: Request/Response structures in `iam/*/dto/`
3. **Entities**: Database models in `entity/`
4. **Repositories** atau **Services**: Data access in `impl/`
5. **Routing**: HTTP routes in `delivery/http/fiber.go`
6. **Controllers**: Request handlers in `delivery/http/controller/`

### Adding a New Feature

1. Define entities in `entity/`
2. Create repository or service interface in module's `contract/`
3. Implement repository in `impl/`
4. Create DTOs in module's `dto/`
5. Implement use case in module's `internal/`
6. Expose use case in module's `exported.go`
7. Create controller in `delivery/http/controller/`
8. Add routes in `delivery/http/fiber.go`
9. Test with Postman collection

### Building for Production

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/iam-service cmd/http/main.go

# Build Docker image
docker build -f deployment/docker/Dockerfile -t iam-service:latest .

# Run Docker container
docker run -p 8080:8080 --env-file .env iam-service:latest
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./iam/auth/iam/...
```

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `APP_NAME` | Application name | `iam-service` |
| `APP_ENV` | Environment (development/staging/production) | `development` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `POSTGRES_HOST` | PostgreSQL host | `localhost` |
| `POSTGRES_PORT` | PostgreSQL port | `5432` |
| `POSTGRES_USER` | PostgreSQL username | `iam` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `iam_secret` |
| `POSTGRES_DB` | PostgreSQL database name | `iam_db` |
| `JWT_ACCESS_SECRET` | JWT access token secret (min 32 chars) | `your-secret-key` |
| `JWT_REFRESH_SECRET` | JWT refresh token secret (min 32 chars) | `your-secret-key` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `VAULT_ADDR` | Vault server address | `http://localhost:8200` |
| `VAULT_TOKEN` | Vault access token | `` |
| `MINIO_ENDPOINT` | MinIO endpoint | `localhost:9000` |
| `EMAIL_SMTP_HOST` | SMTP server host | `` |
| `EMAIL_SMTP_PORT` | SMTP server port | `587` |

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep iam-postgres

# Test database connection
psql -h localhost -p 5432 -U iam -d iam_db

# Check PostgreSQL logs
docker logs iam-postgres
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Migration Failures

```bash
# Reset database (WARNING: Deletes all data!)
docker-compose down -v
docker-compose up -d postgres
# Then re-run migrations
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/iam-new-feature`)
3. Commit your changes (`git commit -m 'Add new IAM feature'`)
4. Push to the branch (`git push origin feature/iam-new-feature`)
5. Open a Pull Request

## License

Should be apache or any other licenses.

## Support

For support, email abrahampurnomo144@gmail.com or create an issue in the repository.

---
