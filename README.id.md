# IAM Service - Sistem Identity & Access Management

[![en](https://img.shields.io/badge/lang-en-blue.svg)](README.md) [![id](https://img.shields.io/badge/lang-id-red.svg)](README.id.md)

Sistem Identity & Access Management (IAM) yang komprehensif, dibuat menggunakan Go dan teknologi cloud-native modern. Service ini bersifat application-agnostic dan dapat diintegrasikan dengan sistem apapun yang membutuhkan autentikasi dan otorisasi.

## Daftar Isi

- [Tentang Project](#tentang-project)
- [Fitur-Fitur Saat Ini](#fitur-fitur-saat-ini)
- [Teknologi yang Dipake](#teknologi-yang-dipake)
- [Struktur Project](#struktur-project)
- [Yang Perlu Disiapkan](#yang-perlu-disiapkan)
- [Cara Mulai](#cara-mulai)
- [Jalanin di Lokal](#jalanin-di-lokal)
- [Dokumentasi API](#dokumentasi-api)
- [Migrasi Database](#migrasi-database)
- [Development](#development)
- [Environment Variables](#environment-variables)

## Tentang Project

IAM Service adalah sistem Identity & Access Management multi-tenant yang menyediakan autentikasi aman, kontrol akses berbasis role, dan kemampuan manajemen user yang komprehensif. Service ini dirancang untuk dapat digunakan oleh aplikasi apapun yang membutuhkan manajemen identitas.

## Fitur-Fitur Saat Ini

### Autentikasi & Otorisasi
- 🔐 Registrasi user dengan verifikasi email (OTP)
- 🔑 Login aman pake JWT (access & refresh tokens)
- 📱 Setup PIN 6 digit dengan validasi keamanan
- 🔄 Reset password pake verifikasi OTP
- 🚪 Logout dengan revoke token
- 👥 Support multi-tenant
- 🎭 Role-based access control (RBAC)
- 🛡️ Manajemen platform admin

### Manajemen User
- 📝 Workflow lengkapin profil
- 👤 Registrasi akun spesial (admin/approver)
- 🔒 Fitur keamanan akun (tracking login attempts, kunci akun)
- 📊 Tracking aktivasi user

### Fitur Keamanan
- 🔐 Validasi password
- 🔄 Tracking history password (5 password terakhir)
- 🔢 Validasi PIN (mencegah pola PIN yang lemah)
- 🚫 Rate limiting untuk request OTP
- 🔒 Hashing password pake Bcrypt
- 🎫 Autentikasi berbasis JWT
- 🔑 Integrasi HashiCorp Vault untuk secret management

## Teknologi yang Dipake

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

### Library & Tools
- **JWT**: golang-jwt/jwt v5
- **Validation**: go-playground/validator v10
- **Configuration**: Viper
- **Logging**: Uber Zap
- **Live Reload**: Air (development)

## Struktur Project

```
iam-service/
├── cmd/                          # Entry point aplikasi
│   ├── http/                     # Entry point HTTP server
│   ├── grpc/                     # gRPC server (kedepannya)
│   ├── graphql/                  # GraphQL server (kedepannya)
│   ├── websocket/                # WebSocket server (kedepannya)
│   ├── worker/                   # Background workers
│   └── migrate/                  # Runner migrasi database
├── config/                       # Manajemen konfigurasi
├── delivery/                     # Delivery layer (controllers/handlers)
│   ├── http/                     # HTTP controllers & routing
│   └── grpc/                     # gRPC handlers
├── entity/                       # Domain entities/models
├── iam/                          # Modul Identity & Access Management
│   ├── auth/                     # Use cases autentikasi
│   ├── user/                     # Use cases manajemen user
│   ├── role/                     # Use cases manajemen role
│   └── health/                   # Use cases health check
├── impl/                         # Implementasi infrastructure
│   ├── postgres/                 # Repository PostgreSQL
│   ├── redis/                    # Implementasi cache Redis
│   ├── hashivault/               # Integrasi Vault
│   ├── minio/                    # Implementasi storage MinIO
│   └── mailer/                   # Implementasi service email
├── pkg/                          # Package yang dishare
│   ├── jwt/                      # Utilities JWT
│   ├── errors/                   # Custom error types
│   ├── logger/                   # Utilities logging
│   ├── validator/                # Helper validasi
│   └── utils/                    # Utilities umum
├── infrastructure/               # Setup infrastructure
├── migration/                    # File migrasi database
├── deployment/                   # Konfigurasi deployment
│   ├── docker/                   # Docker & Docker Compose
│   └── k8s/                      # Manifest Kubernetes
├── doc/                          # Dokumentasi
│   ├── postman/                  # Koleksi API Postman
│   └── openapi/                  # Spec OpenAPI/Swagger
└── script/                       # Script utilities
```

### Pola Arsitektur

Project ini ngikutin prinsip **Clean Architecture**:

- **Entity Layer**: Model bisnis inti (`entity/`)
- **Use Case Layer**: Logika bisnis (`iam/*/iam/`)
- **Repository Layer**: Abstraksi akses data (`iam/*/contract/`)
- **Delivery Layer**: HTTP handlers (`delivery/http/`)
- **Infrastructure Layer**: Service eksternal (`impl/`, `infrastructure/`)

Setiap modul (auth, user, role) ngikutin struktur ini:
```
iam/module/
├── contract/          # Interface (repository, usecase)
├── internal/          # Implementasi use case (private)
├── moduledto/         # Data Transfer Objects
└── exported.go        # Public API
```

## Yang Perlu Disiapkan

Sebelum jalanin project ini, pastikan sudah install:

- **Go** 1.24 atau lebih tinggi ([Download](https://golang.org/dl/))
- **Docker** & **Docker Compose** ([Download](https://www.docker.com/get-started))
- **Make** (opsional, buat automasi build)
- **Air** (opsional, buat hot reloading): `go install github.com/air-verse/air@latest`
- **PostgreSQL Client** (opsional, buat operasi DB secara manual)

## Cara Mulai

### 1. Clone Repository

```bash
git clone <repository-url>
cd iam-service
```

### 2. Setup Environment Variables

Copy file example environment dan konfigurasi:

```bash
cp .env.example .env
```

Edit `.env` dan isi nilai yang diperluin:

```env
# Konfigurasi App
APP_NAME=iam-service
APP_ENV=development
APP_VERSION=1.0.0

# Konfigurasi Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# PostgreSQL (kalo pake Docker)
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

# JWT (Bikin secret yang kuat buat production!)
JWT_ACCESS_SECRET=secret-key-mu-yang-super-rahasia-minimal-32-karakter
JWT_REFRESH_SECRET=secret-key-mu-yang-super-rahasia-minimal-32-karakter
JWT_ISSUER=iam-service

# Email (Opsional - buat kirim OTP)
EMAIL_PROVIDER=smtp
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=email-kamu@gmail.com
EMAIL_SMTP_PASS=password-app-kamu
EMAIL_FROM_ADDRESS=noreply@example.com
EMAIL_FROM_NAME=IAM System
```

### 3. Jalanin Service Infrastructure

Start PostgreSQL, Redis, Vault, dan MinIO pake Docker Compose:

```bash
cd deployment/docker
docker-compose up -d
```

Cek semua service udah jalan:

```bash
docker-compose ps
```

Harusnya muncul:
- `iam-postgres` di port 5432
- `iam-redis` di port 6379
- `iam-vault` di port 8200
- `iam-minio` di port 9000 (API) dan 9001 (Console)

### 4. Install Dependencies Go

```bash
go mod download
```

## Jalanin di Lokal

### Opsi 1: Pake Air (Hot Reload - Recommended buat Development)

Air otomatis rebuild dan restart aplikasi kamu kalo ada perubahan:

```bash
# Install Air kalo belum punya
go install github.com/air-verse/air@latest

# Jalanin pake Air
air
```

Server bakal jalan di `http://localhost:8080` dengan hot reload aktif.

### Opsi 2: Pake `go run`

```bash
go run cmd/http/main.go
```

### Opsi 3: Build dan Run Binary

```bash
# Build binary
go build -o bin/iam-service cmd/http/main.go

# Jalanin binary
./bin/iam-service
```

### Cek Server Udah Jalan

```bash
# Health check
curl http://localhost:8080/api/v1/health/

# Harusnya return: {"status":"healthy","timestamp":"..."}
```

## Dokumentasi API

### Koleksi Postman

Import koleksi Postman yang lengkap di:
```
doc/postman/iam-service-api-collection.json
```

Koleksi ini mencakup:
- ✅ Semua 16 endpoint dengan contoh request
- ✅ Contoh response sukses dan gagal
- ✅ Auto-script buat save access token
- ✅ Setup environment variables

### Endpoint yang Tersedia

#### Endpoint Health
- `GET /api/v1/health/` - Health check dasar
- `GET /api/v1/health/ready` - Readiness probe
- `GET /api/v1/health/live` - Liveness probe

#### Endpoint Autentikasi
- `POST /api/v1/auth/register` - Registrasi user sendiri
- `POST /api/v1/auth/register/special-account` - Bikin akun admin/approver
- `POST /api/v1/auth/verify-otp` - Verifikasi OTP email
- `POST /api/v1/auth/complete-profile` - Lengkapin profil user
- `POST /api/v1/auth/resend-otp` - Kirim ulang OTP email
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/logout` - Logout user
- `POST /api/v1/auth/setup-pin` - Setup PIN 6 digit (butuh JWT)
- `POST /api/v1/auth/request-password-reset` - Request OTP reset password
- `POST /api/v1/auth/reset-password` - Reset password pake OTP

#### Manajemen Role (Khusus Platform Admin)
- `POST /api/v1/roles/` - Bikin role baru

#### Manajemen User (Khusus Platform Admin)
- `POST /api/v1/users/` - Bikin user baru dengan system role

### Quick Start Testing API

1. **Daftar akun spesial (Platform Admin)**:
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

Simpen `access_token` dari response buat request yang butuh autentikasi.

3. **Akses endpoint yang dilindungi**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/setup-pin \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ACCESS_TOKEN_KAMU" \
  -d '{
    "pin": "147258",
    "pin_confirm": "147258"
  }'
```

## Migrasi Database

File migrasi database ada di folder `migration/` dan ngikutin pola penamaan:
```
XXXXXX_deskripsi.up.sql    # Apply migrasi
XXXXXX_deskripsi.down.sql  # Rollback migrasi
```

### Migrasi yang Tersedia

1. `000001_extensions` - Extension PostgreSQL (uuid-ossp)
2. `000002_tenants` - Manajemen tenant
3. `000003_users` - Tabel user, profile, credentials, security
4. `000004_auth_tokens` - Refresh tokens
5. `000005_saml_mfa` - Support SAML dan MFA
6. `000006_email_verification` - Verifikasi OTP
7. `000007_user_activation_tracking` - Workflow aktivasi
8. `000008_roles_permissions` - Sistem RBAC
9. `000009_user_roles` - Assignment role user
10. `000010_refresh_token_family` - Rotasi token
11. `000011_platform_admin` - Role platform admin

### Jalanin Migrasi

#### Migrasi Manual (Pake psql)

```bash
# Connect ke database
psql -h localhost -p 5432 -U iam -d iam_db

# Jalanin migrasi manual secara berurutan
\i migration/000001_extensions.up.sql
\i migration/000002_tenants.up.sql
# ... dan seterusnya
```

#### Pake Migration Tool (migrate)

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Jalanin semua migrasi
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" up

# Rollback migrasi terakhir
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" down 1

# Cek versi migrasi
migrate -path migration -database "postgresql://iam:iam_secret@localhost:5432/iam_db?sslmode=disable" version
```

## Development

### Konvensi Struktur Code

1. **Use Cases**: Semua logika bisnis ada di `iam/*/iam/`
2. **DTOs**: Struktur Request/Response di `iam/*/dto/`
3. **Entities**: Model database di `entity/`
4. **Repositories** atau **Services**: Akses data di `impl/`
5. **Routing**: Route HTTP di `delivery/http/fiber.go`
6. **Controllers**: Request handlers di `delivery/http/controller/`

### Nambah Fitur Baru

1. Define entities di `entity/`
2. Bikin repository atau service interface di `contract/` modul
3. Implementasi repository atau service di `impl/`
4. Bikin DTOs di `dto/` modul
5. Implementasi use case di `internal/` modul
6. Expose use case di `exported.go` modul
7. Bikin controller di `delivery/http/controller/`
8. Tambahin routes di `delivery/http/fiber.go`
9. Testing pake koleksi Postman

### Build buat Production

```bash
# Build binary yang optimized
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/iam-service cmd/http/main.go

# Build Docker image
docker build -f deployment/docker/Dockerfile -t iam-service:latest .

# Jalanin Docker container
docker run -p 8080:8080 --env-file .env iam-service:latest
```

### Jalanin Tests

```bash
# Jalanin semua tests
go test ./...

# Jalanin tests dengan coverage
go test -cover ./...

# Jalanin tests untuk package tertentu
go test ./iam/auth/iam/...
```

## Environment Variables

### Variable yang Wajib

| Variable | Deskripsi | Contoh |
|----------|-----------|--------|
| `APP_NAME` | Nama aplikasi | `iam-service` |
| `APP_ENV` | Environment (development/staging/production) | `development` |
| `SERVER_PORT` | Port HTTP server | `8080` |
| `POSTGRES_HOST` | Host PostgreSQL | `localhost` |
| `POSTGRES_PORT` | Port PostgreSQL | `5432` |
| `POSTGRES_USER` | Username PostgreSQL | `iam` |
| `POSTGRES_PASSWORD` | Password PostgreSQL | `iam_secret` |
| `POSTGRES_DB` | Nama database PostgreSQL | `iam_db` |
| `JWT_ACCESS_SECRET` | Secret JWT access token (min 32 karakter) | `secret-key-kamu` |
| `JWT_REFRESH_SECRET` | Secret JWT refresh token (min 32 karakter) | `secret-key-kamu` |

### Variable Opsional

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `REDIS_HOST` | Host Redis | `localhost` |
| `REDIS_PORT` | Port Redis | `6379` |
| `REDIS_PASSWORD` | Password Redis | `` |
| `VAULT_ADDR` | Address server Vault | `http://localhost:8200` |
| `VAULT_TOKEN` | Token akses Vault | `` |
| `MINIO_ENDPOINT` | Endpoint MinIO | `localhost:9000` |
| `EMAIL_SMTP_HOST` | Host server SMTP | `` |
| `EMAIL_SMTP_PORT` | Port server SMTP | `587` |

## Troubleshooting

### Masalah Koneksi Database

```bash
# Cek PostgreSQL lagi jalan atau nggak
docker ps | grep iam-postgres

# Test koneksi database
psql -h localhost -p 5432 -U iam -d iam_db

# Cek logs PostgreSQL
docker logs iam-postgres
```

### Port Udah Kepake

```bash
# Cari process yang pake port 8080
lsof -i :8080

# Kill process nya
kill -9 <PID>
```

### Migrasi Gagal

```bash
# Reset database (PERINGATAN: Hapus semua data!)
docker-compose down -v
docker-compose up -d postgres
# Terus jalanin ulang migrasi
```

## Contributing

1. Fork repository ini
2. Bikin feature branch (`git checkout -b feature/iam-fitur-baru`)
3. Commit perubahan (`git commit -m 'Nambah fitur baru di IAM'`)
4. Push ke branch (`git push origin feature/iam-fitur-baru`)
5. Bikin Pull Request

## License

Harusnya pake Apache atau license lainnya yang sesuai.

## Support

Kalo butuh bantuan, email ke abrahampurnomo144@gmail.com atau bikin issue di repository.

---
