# IAM System - Detailed Task Breakdown
# Supplementary Implementation Guide

---

## Document Control

| Attribute | Value |
|-----------|-------|
| **Document Title** | IAM System Detailed Task Breakdown |
| **Version** | 1.0 |
| **Created Date** | 2026-02-03 |
| **Parent Document** | IAM-Implementation-Plan-V2.md |

---

## 1. Masterdata Service Integration

### 1.1 Overview

The IAM system depends on the Masterdata Service for validating reference data fields in user profiles:
- Gender
- Marital Status
- Country, Province, City (for address validation)
- Other tenant-configurable reference data

### 1.2 Integration Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   IAM Service   │────▶│  Redis Cache    │────▶│   Masterdata    │
│                 │     │   (1 hour TTL)  │     │    Service      │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │   1. Check cache      │                       │
        │◀──────────────────────│                       │
        │                       │                       │
        │   2. Cache miss       │   3. Fetch from API   │
        │──────────────────────▶│──────────────────────▶│
        │                       │                       │
        │   4. Store in cache   │   5. Return data      │
        │◀──────────────────────│◀──────────────────────│
```

### 1.3 Detailed Tasks

#### Task 1.1.15: Set up Masterdata Service Stub (8 hours)

**Objective:** Create a mock Masterdata API for development and testing.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.1.15.1 | Create Masterdata mock server (Docker container) | 2 | Use WireMock or custom Go service |
| 1.1.15.2 | Define mock endpoints | 1 | GET /categories, GET /items |
| 1.1.15.3 | Create mock data (gender, marital status, countries, provinces, cities) | 2 | JSON fixtures |
| 1.1.15.4 | Add to Docker Compose | 1 | Include in local dev environment |
| 1.1.15.5 | Document mock service usage | 1 | README |
| 1.1.15.6 | Create toggle for mock vs real service | 1 | Environment variable |

**Mock Endpoints:**

```yaml
# Mock Masterdata API
GET /api/masterdata/categories
GET /api/masterdata/categories/{code}/items
GET /api/masterdata/items/{id}
GET /api/masterdata/items/validate?category={code}&value={value}
```

**Mock Data Example (gender):**

```json
{
  "category": "GENDER",
  "items": [
    { "code": "M", "name": "Male", "status": "ACTIVE" },
    { "code": "F", "name": "Female", "status": "ACTIVE" },
    { "code": "O", "name": "Other", "status": "ACTIVE" }
  ]
}
```

**Acceptance Criteria:**
- [ ] Mock server starts with `make dev`
- [ ] All required categories have mock data
- [ ] Toggle between mock and real service works
- [ ] Documentation complete

---

#### Task 1.2.15: Implement Masterdata Client (12 hours)

**Objective:** Create a client to interact with Masterdata service with Redis caching.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.2.15.1 | Define Masterdata client interface | 1 | Go interface |
| 1.2.15.2 | Implement HTTP client with retry logic | 2 | Exponential backoff |
| 1.2.15.3 | Implement Redis caching layer | 3 | 1 hour TTL |
| 1.2.15.4 | Implement graceful degradation | 2 | Use stale cache on failure |
| 1.2.15.5 | Create validation functions | 2 | ValidateGender, ValidateMaritalStatus, etc. |
| 1.2.15.6 | Write unit tests | 2 | Mock HTTP responses |

**Interface Definition:**

```go
// internal/masterdata/client.go

type MasterdataClient interface {
    // GetCategories returns all available categories
    GetCategories(ctx context.Context) ([]Category, error)
    
    // GetItems returns items for a category
    GetItems(ctx context.Context, categoryCode string) ([]Item, error)
    
    // ValidateItem checks if a value is valid for a category
    ValidateItem(ctx context.Context, categoryCode, value string) (bool, error)
    
    // GetItem returns a single item by ID
    GetItem(ctx context.Context, itemID string) (*Item, error)
}

type Category struct {
    ID          string `json:"id"`
    Code        string `json:"code"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type Item struct {
    ID           string                 `json:"id"`
    CategoryCode string                 `json:"category_code"`
    Code         string                 `json:"code"`
    Name         string                 `json:"name"`
    Status       string                 `json:"status"`
    Metadata     map[string]interface{} `json:"metadata"`
}
```

**Redis Cache Keys:**

```
masterdata:categories                    -> list of all categories (TTL: 1 hour)
masterdata:category:{code}:items         -> items for a category (TTL: 1 hour)
masterdata:item:{id}                     -> single item (TTL: 1 hour)
```

**Graceful Degradation Logic:**

```go
func (c *client) GetItems(ctx context.Context, categoryCode string) ([]Item, error) {
    // 1. Try cache first
    items, err := c.cache.GetItems(ctx, categoryCode)
    if err == nil {
        return items, nil
    }
    
    // 2. Try API
    items, err = c.api.GetItems(ctx, categoryCode)
    if err != nil {
        // 3. Fallback: try stale cache (ignore TTL)
        items, staleErr := c.cache.GetItemsStale(ctx, categoryCode)
        if staleErr == nil {
            c.logger.Warn("using stale cache for masterdata", 
                zap.String("category", categoryCode),
                zap.Error(err))
            return items, nil
        }
        return nil, fmt.Errorf("masterdata unavailable: %w", err)
    }
    
    // 4. Update cache
    c.cache.SetItems(ctx, categoryCode, items)
    return items, nil
}
```

**Acceptance Criteria:**
- [ ] Client interface defined and implemented
- [ ] Redis caching works (1 hour TTL)
- [ ] Graceful degradation: stale cache used when API fails
- [ ] Retry logic with exponential backoff
- [ ] Unit tests with ≥80% coverage

---

#### Task 2.2.13: Validate Profile Fields via Masterdata (4 hours)

**Objective:** Validate user profile fields (gender, marital status) against Masterdata.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 2.2.13.1 | Add validation to user profile update | 1 | Use Masterdata client |
| 2.2.13.2 | Add validation to admin user creation | 1 | Same validation |
| 2.2.13.3 | Return meaningful error messages | 1 | "Invalid gender code" |
| 2.2.13.4 | Write integration tests | 1 | Test with mock Masterdata |

**Validation Logic:**

```go
// internal/user/service.go

func (s *Service) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) error {
    // Validate gender if provided
    if req.Gender != "" {
        valid, err := s.masterdata.ValidateItem(ctx, "GENDER", req.Gender)
        if err != nil {
            return fmt.Errorf("failed to validate gender: %w", err)
        }
        if !valid {
            return NewValidationError("gender", "invalid gender code")
        }
    }
    
    // Validate marital status if provided
    if req.MaritalStatus != "" {
        valid, err := s.masterdata.ValidateItem(ctx, "MARITAL_STATUS", req.MaritalStatus)
        if err != nil {
            return fmt.Errorf("failed to validate marital status: %w", err)
        }
        if !valid {
            return NewValidationError("marital_status", "invalid marital status code")
        }
    }
    
    // Continue with update...
}
```

**Error Response Example:**

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      {
        "field": "gender",
        "message": "Invalid gender code. Valid values: M, F, O"
      }
    ]
  }
}
```

**Acceptance Criteria:**
- [ ] Profile update validates against Masterdata
- [ ] Admin user creation validates against Masterdata
- [ ] Meaningful error messages returned
- [ ] Validation is skipped if Masterdata unavailable (graceful)

---

## 2. Email Service Setup

### 2.1 Overview

The IAM system requires email functionality for:
- Email OTP verification during registration
- Welcome emails after approval
- Password reset emails
- Account lockout notifications

### 2.2 Email Architecture

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   IAM Service   │────▶│  Email Queue    │────▶│   SMTP Server   │
│                 │     │    (Redis)      │     │  (or SendGrid)  │
└─────────────────┘     └─────────────────┘     └─────────────────┘
        │                       │                       │
        │   1. Queue email      │                       │
        │──────────────────────▶│                       │
        │                       │                       │
        │                       │   2. Worker sends     │
        │                       │──────────────────────▶│
        │                       │                       │
        │                       │   3. Delivery status  │
        │                       │◀──────────────────────│
```

### 2.3 Detailed Tasks

#### Task 1.1.13: Set up SMTP Service (4 hours)

**Objective:** Configure email infrastructure for development and production.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.1.13.1 | Add Mailhog to Docker Compose | 1 | Local email capture |
| 1.1.13.2 | Configure SMTP settings in Vault | 1 | Production credentials |
| 1.1.13.3 | Create SMTP configuration struct | 1 | Go config |
| 1.1.13.4 | Document email setup | 1 | README section |

**Docker Compose (development):**

```yaml
# docker-compose.yml
services:
  mailhog:
    image: mailhog/mailhog:v1.0.1
    ports:
      - "1025:1025"   # SMTP port
      - "8025:8025"   # Web UI
    networks:
      - iam-network
```

**SMTP Configuration:**

```go
// internal/config/email.go

type EmailConfig struct {
    // SMTP settings
    Host     string `env:"SMTP_HOST" envDefault:"localhost"`
    Port     int    `env:"SMTP_PORT" envDefault:"1025"`
    Username string `env:"SMTP_USERNAME"`
    Password string `env:"SMTP_PASSWORD"`
    
    // TLS settings
    UseTLS      bool `env:"SMTP_USE_TLS" envDefault:"false"`
    SkipVerify  bool `env:"SMTP_SKIP_VERIFY" envDefault:"false"`
    
    // Sender settings
    FromAddress string `env:"SMTP_FROM_ADDRESS" envDefault:"noreply@iam.local"`
    FromName    string `env:"SMTP_FROM_NAME" envDefault:"IAM System"`
    
    // Queue settings
    QueueSize    int           `env:"EMAIL_QUEUE_SIZE" envDefault:"1000"`
    WorkerCount  int           `env:"EMAIL_WORKER_COUNT" envDefault:"3"`
    RetryCount   int           `env:"EMAIL_RETRY_COUNT" envDefault:"3"`
    RetryDelay   time.Duration `env:"EMAIL_RETRY_DELAY" envDefault:"1m"`
}
```

**Vault Secrets (production):**

```bash
# Store SMTP credentials in Vault
vault kv put secret/iam/smtp \
    host="smtp.sendgrid.net" \
    port="587" \
    username="apikey" \
    password="SG.xxxxxx" \
    use_tls="true"
```

**Acceptance Criteria:**
- [ ] Mailhog running locally (http://localhost:8025)
- [ ] SMTP config loaded from environment/Vault
- [ ] Production SMTP credentials documented (not committed)

---

#### Task 1.1.14: Create Email Templates (8 hours)

**Objective:** Create HTML email templates for all email types.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.1.14.1 | Design email base template | 2 | Header, footer, styling |
| 1.1.14.2 | Create OTP verification template | 1.5 | 6-digit code display |
| 1.1.14.3 | Create welcome email template | 1 | Post-approval |
| 1.1.14.4 | Create password reset template | 1.5 | Reset link |
| 1.1.14.5 | Create account locked template | 1 | Unlock instructions |
| 1.1.14.6 | Test templates across email clients | 1 | Outlook, Gmail, etc. |

**Template Directory Structure:**

```
templates/
├── email/
│   ├── base.html              # Base layout
│   ├── otp_verification.html  # Email OTP
│   ├── welcome.html           # Welcome after approval
│   ├── password_reset.html    # Password reset link
│   ├── account_locked.html    # Lockout notification
│   └── password_changed.html  # Password change confirmation
```

**OTP Template Example:**

```html
<!-- templates/email/otp_verification.html -->
{{define "subject"}}Verify your email - {{.TenantName}}{{end}}

{{define "body"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
    <style>
        /* Inline styles for email compatibility */
        .container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
        .otp-code { font-size: 32px; font-weight: bold; letter-spacing: 8px; 
                    color: #2563eb; padding: 20px; background: #f3f4f6; 
                    text-align: center; border-radius: 8px; }
        .footer { color: #6b7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Verify Your Email</h2>
        <p>Hello {{.FirstName}},</p>
        <p>Your verification code is:</p>
        <div class="otp-code">{{.OTPCode}}</div>
        <p>This code expires in <strong>{{.ExpiresInMinutes}} minutes</strong>.</p>
        <p>If you didn't request this code, please ignore this email.</p>
        <div class="footer">
            <p>This is an automated message from {{.TenantName}}.</p>
            <p>© {{.Year}} {{.TenantName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
{{end}}
```

**Template Data Structures:**

```go
// internal/email/templates.go

type OTPEmailData struct {
    FirstName        string
    OTPCode          string
    ExpiresInMinutes int
    TenantName       string
    Year             int
}

type WelcomeEmailData struct {
    FirstName  string
    Email      string
    LoginURL   string
    TenantName string
}

type PasswordResetEmailData struct {
    FirstName        string
    ResetLink        string
    ExpiresInMinutes int
    TenantName       string
    IPAddress        string
}
```

**Acceptance Criteria:**
- [ ] All templates created and tested
- [ ] Templates render correctly in Mailhog
- [ ] Templates work in major email clients (Outlook, Gmail)
- [ ] Template data structures defined

---

#### Task 1.2.16: Implement Email Service Wrapper (8 hours)

**Objective:** Create email service with template rendering and async sending.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.2.16.1 | Create email service interface | 1 | Go interface |
| 1.2.16.2 | Implement template rendering | 2 | Go html/template |
| 1.2.16.3 | Implement SMTP sender | 2 | net/smtp |
| 1.2.16.4 | Implement async queue (Redis) | 2 | Background workers |
| 1.2.16.5 | Write unit tests | 1 | Mock SMTP |

**Email Service Interface:**

```go
// internal/email/service.go

type EmailService interface {
    // SendOTP sends email verification OTP
    SendOTP(ctx context.Context, to string, data OTPEmailData) error
    
    // SendWelcome sends welcome email after approval
    SendWelcome(ctx context.Context, to string, data WelcomeEmailData) error
    
    // SendPasswordReset sends password reset link
    SendPasswordReset(ctx context.Context, to string, data PasswordResetEmailData) error
    
    // SendAccountLocked sends lockout notification
    SendAccountLocked(ctx context.Context, to string, data AccountLockedEmailData) error
    
    // SendPasswordChanged sends password change confirmation
    SendPasswordChanged(ctx context.Context, to string, data PasswordChangedEmailData) error
}

type Email struct {
    To       string
    Subject  string
    HTMLBody string
    TextBody string
}
```

**Async Queue Implementation:**

```go
// internal/email/queue.go

type EmailQueue struct {
    redis  *redis.Client
    logger *zap.Logger
}

func (q *EmailQueue) Enqueue(ctx context.Context, email Email) error {
    data, err := json.Marshal(email)
    if err != nil {
        return err
    }
    return q.redis.RPush(ctx, "email:queue", data).Err()
}

func (q *EmailQueue) StartWorkers(ctx context.Context, workerCount int, sender EmailSender) {
    for i := 0; i < workerCount; i++ {
        go q.worker(ctx, i, sender)
    }
}

func (q *EmailQueue) worker(ctx context.Context, id int, sender EmailSender) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            data, err := q.redis.BLPop(ctx, time.Second, "email:queue").Result()
            if err != nil {
                continue
            }
            
            var email Email
            if err := json.Unmarshal([]byte(data[1]), &email); err != nil {
                q.logger.Error("failed to unmarshal email", zap.Error(err))
                continue
            }
            
            if err := sender.Send(ctx, email); err != nil {
                q.logger.Error("failed to send email", 
                    zap.String("to", email.To),
                    zap.Error(err))
                // Retry logic...
            }
        }
    }
}
```

**Acceptance Criteria:**
- [ ] Email service interface implemented
- [ ] Template rendering works correctly
- [ ] Async queue with background workers
- [ ] Retry logic for failed sends
- [ ] Unit tests with ≥80% coverage

---

## 3. Token Blacklist Implementation

### 3.1 Overview

Token blacklist ensures that logged-out tokens cannot be reused. When a user logs out, their access token is added to a Redis blacklist until it expires naturally.

### 3.2 Architecture

```
┌─────────────────┐
│   POST /logout  │
│   (with token)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│  Extract token  │────▶│  Add to Redis   │
│  (JWT decode)   │     │  blacklist      │
└────────┬────────┘     │  (TTL = exp)    │
         │              └─────────────────┘
         ▼
┌─────────────────┐
│   Auth Middleware   │
│   (every request)   │
└────────┬────────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐
│  Check blacklist │────▶│  Redis EXISTS   │
│  before validation│     │  O(1) lookup    │
└─────────────────┘     └─────────────────┘
```

### 3.3 Detailed Tasks

#### Task 2.1.9: Implement Token Blacklist in Redis (8 hours)

**Objective:** Create token blacklist for secure logout.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 2.1.9.1 | Design Redis key structure | 1 | `token:blacklist:{jti}` |
| 2.1.9.2 | Implement blacklist repository | 2 | Add, Check, Cleanup |
| 2.1.9.3 | Integrate with logout endpoint | 1 | Extract JTI, add to blacklist |
| 2.1.9.4 | Integrate with auth middleware | 2 | Check blacklist on every request |
| 2.1.9.5 | Write unit tests | 2 | Test expiry, concurrent access |

**Redis Key Structure:**

```
token:blacklist:{jti}    -> "1" (value doesn't matter)
                         -> TTL = token expiry time - current time

Example:
token:blacklist:550e8400-e29b-41d4-a716-446655440000 -> "1" (TTL: 3600s)
```

**Blacklist Repository:**

```go
// internal/auth/blacklist.go

type TokenBlacklist interface {
    // Add adds a token to the blacklist
    Add(ctx context.Context, jti string, expiresAt time.Time) error
    
    // IsBlacklisted checks if a token is blacklisted
    IsBlacklisted(ctx context.Context, jti string) (bool, error)
}

type redisBlacklist struct {
    client *redis.Client
    logger *zap.Logger
}

func (b *redisBlacklist) Add(ctx context.Context, jti string, expiresAt time.Time) error {
    ttl := time.Until(expiresAt)
    if ttl <= 0 {
        // Token already expired, no need to blacklist
        return nil
    }
    
    key := fmt.Sprintf("token:blacklist:%s", jti)
    return b.client.Set(ctx, key, "1", ttl).Err()
}

func (b *redisBlacklist) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
    key := fmt.Sprintf("token:blacklist:%s", jti)
    exists, err := b.client.Exists(ctx, key).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, nil
}
```

**Logout Integration:**

```go
// internal/auth/handler.go

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
    // Extract token from header
    token := c.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")
    
    // Decode JWT to get JTI and expiry
    claims, err := h.jwtService.Decode(token)
    if err != nil {
        // Invalid token, but still return success (idempotent)
        return c.SendStatus(fiber.StatusNoContent)
    }
    
    // Add to blacklist
    if err := h.blacklist.Add(c.Context(), claims.JTI, claims.ExpiresAt); err != nil {
        h.logger.Error("failed to blacklist token", zap.Error(err))
        // Don't fail the request, log for monitoring
    }
    
    // Also invalidate refresh token family (optional, for enhanced security)
    if claims.RefreshTokenID != "" {
        h.refreshTokenRepo.Revoke(c.Context(), claims.RefreshTokenID)
    }
    
    return c.SendStatus(fiber.StatusNoContent)
}
```

**Auth Middleware Integration:**

```go
// internal/middleware/auth.go

func (m *AuthMiddleware) Handle(c *fiber.Ctx) error {
    token := c.Get("Authorization")
    token = strings.TrimPrefix(token, "Bearer ")
    
    if token == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Code:    "UNAUTHORIZED",
            Message: "Missing authorization token",
        })
    }
    
    // Decode and validate JWT
    claims, err := m.jwtService.Validate(token)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Code:    "INVALID_TOKEN",
            Message: "Invalid or expired token",
        })
    }
    
    // Check blacklist (CRITICAL: must happen after validation)
    blacklisted, err := m.blacklist.IsBlacklisted(c.Context(), claims.JTI)
    if err != nil {
        m.logger.Error("failed to check blacklist", zap.Error(err))
        // Fail open or closed based on security requirements
        // Recommendation: fail closed (deny access)
        return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
            Code:    "INTERNAL_ERROR",
            Message: "Authentication check failed",
        })
    }
    
    if blacklisted {
        return c.Status(fiber.StatusUnauthorized).JSON(ErrorResponse{
            Code:    "TOKEN_REVOKED",
            Message: "Token has been revoked",
        })
    }
    
    // Set user context
    c.Locals("user", claims)
    return c.Next()
}
```

**Acceptance Criteria:**
- [ ] Blacklist repository implemented
- [ ] Logout adds token to blacklist
- [ ] Auth middleware checks blacklist
- [ ] TTL matches token expiry (auto-cleanup)
- [ ] Unit tests with ≥80% coverage

---

## 4. Password History Implementation

### 4.1 Overview

Password history prevents users from reusing recent passwords. The system stores hashed versions of the last N passwords (configurable, default 5).

### 4.2 Detailed Tasks

#### Task 2.2.7: Implement Password History Check (4 hours)

**Objective:** Prevent password reuse by checking against history.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 2.2.7.1 | Implement password history repository | 1 | CRUD operations |
| 2.2.7.2 | Integrate with password change flow | 1 | Check before update |
| 2.2.7.3 | Integrate with password reset flow | 1 | Check before reset |
| 2.2.7.4 | Write unit tests | 1 | Test history check |

**Repository Implementation:**

```go
// internal/user/password_history.go

type PasswordHistoryRepository interface {
    // Add adds a password hash to history
    Add(ctx context.Context, userID string, passwordHash string) error
    
    // GetRecent returns the most recent N password hashes
    GetRecent(ctx context.Context, userID string, count int) ([]string, error)
    
    // IsReused checks if a password matches any in history
    IsReused(ctx context.Context, userID string, password string, count int) (bool, error)
}

func (r *passwordHistoryRepo) IsReused(ctx context.Context, userID string, password string, count int) (bool, error) {
    hashes, err := r.GetRecent(ctx, userID, count)
    if err != nil {
        return false, err
    }
    
    for _, hash := range hashes {
        if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil {
            return true, nil // Password matches a historical hash
        }
    }
    
    return false, nil
}
```

**Integration with Password Change:**

```go
// internal/auth/service.go

func (s *AuthService) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
    // Verify current password
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // ... verify current password ...
    
    // Check password history (configurable count from tenant settings)
    historyCount := s.getTenantPasswordHistoryCount(ctx, user.TenantID)
    reused, err := s.passwordHistory.IsReused(ctx, userID, req.NewPassword, historyCount)
    if err != nil {
        return err
    }
    if reused {
        return NewValidationError("new_password", 
            fmt.Sprintf("Cannot reuse any of your last %d passwords", historyCount))
    }
    
    // Hash new password
    newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    
    // Update password
    if err := s.authMethodRepo.UpdatePassword(ctx, userID, string(newHash)); err != nil {
        return err
    }
    
    // Add old password to history
    if err := s.passwordHistory.Add(ctx, userID, user.PasswordHash); err != nil {
        s.logger.Error("failed to add password to history", zap.Error(err))
        // Non-critical, continue
    }
    
    return nil
}
```

**Acceptance Criteria:**
- [ ] Password history stored in database
- [ ] History check integrated with password change
- [ ] History check integrated with password reset
- [ ] Configurable history count (default 5)

---

#### Task 2.2.8: Implement Password History Pruning (4 hours)

**Objective:** Automatically clean up old password history records.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 2.2.8.1 | Implement pruning function | 1 | Keep only last N |
| 2.2.8.2 | Add pruning to password change flow | 1 | Prune after adding |
| 2.2.8.3 | Create background cleanup job (optional) | 1 | Periodic cleanup |
| 2.2.8.4 | Write unit tests | 1 | Test pruning logic |

**Pruning Implementation:**

```go
// internal/user/password_history.go

func (r *passwordHistoryRepo) Add(ctx context.Context, userID string, passwordHash string) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    // Insert new record
    _, err = tx.ExecContext(ctx, `
        INSERT INTO password_history (id, user_id, password_hash, created_at)
        VALUES ($1, $2, $3, NOW())
    `, uuid.New(), userID, passwordHash)
    if err != nil {
        return err
    }
    
    // Prune old records (keep only last N)
    // Using tenant setting for history count
    _, err = tx.ExecContext(ctx, `
        DELETE FROM password_history
        WHERE user_id = $1
        AND id NOT IN (
            SELECT id FROM password_history
            WHERE user_id = $1
            ORDER BY created_at DESC
            LIMIT $2
        )
    `, userID, maxHistoryCount)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

**Acceptance Criteria:**
- [ ] Pruning happens automatically on add
- [ ] Only configured number of records retained
- [ ] No orphan records in database

---

## 5. API Contract Management (OpenAPI + Postman)

### 5.1 Overview

API contracts are defined using OpenAPI specification and tested using Postman collections. This enables frontend teams to start integration early.

### 5.2 Detailed Tasks

#### Task 1.2.13: Create OpenAPI Specification Draft (16 hours)

**Objective:** Create comprehensive OpenAPI specification for all endpoints.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.2.13.1 | Define info, servers, security schemes | 1 | Base setup |
| 1.2.13.2 | Define authentication endpoints | 3 | /auth/* |
| 1.2.13.3 | Define user management endpoints | 3 | /users/* |
| 1.2.13.4 | Define authorization endpoints | 3 | /roles/*, /permissions/* |
| 1.2.13.5 | Define organization endpoints | 2 | /tenants/*, /branches/* |
| 1.2.13.6 | Define audit endpoints | 2 | /audit/* |
| 1.2.13.7 | Validate with Spectral | 1 | Lint OpenAPI |
| 1.2.13.8 | Document error responses | 1 | All error codes |

**OpenAPI Structure:**

```yaml
openapi: 3.1.0
info:
  title: IAM System API
  version: 1.0.0
  description: Generic Multi-Tenant Identity and Access Management System

servers:
  - url: http://localhost:3000/api/iam
    description: Development
  - url: https://staging.iam.example.com/api/iam
    description: Staging
  - url: https://iam.example.com/api/iam
    description: Production

security:
  - bearerAuth: []

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
  
  schemas:
    Error:
      type: object
      properties:
        error:
          type: object
          properties:
            code:
              type: string
            message:
              type: string
            details:
              type: array
              items:
                type: object
    
    # ... more schemas ...

paths:
  /auth/login:
    post:
      # ... endpoint definition ...
```

**Acceptance Criteria:**
- [ ] All endpoints documented
- [ ] Request/response schemas defined
- [ ] Error responses documented
- [ ] Spectral validation passes
- [ ] Security schemes defined

---

#### Task 1.2.14: Create Postman Collection (12 hours)

**Objective:** Create Postman collection with all endpoints and tests.

**Subtasks:**

| ID | Subtask | Est. Hours | Notes |
|----|---------|------------|-------|
| 1.2.14.1 | Create collection structure | 1 | Folders by feature |
| 1.2.14.2 | Add authentication endpoints | 2 | With pre-request scripts |
| 1.2.14.3 | Add user management endpoints | 2 | With tests |
| 1.2.14.4 | Add authorization endpoints | 2 | With tests |
| 1.2.14.5 | Add organization endpoints | 2 | With tests |
| 1.2.14.6 | Add audit endpoints | 1 | With tests |
| 1.2.14.7 | Create environment files | 1 | Dev, Staging, Prod |
| 1.2.14.8 | Add collection-level scripts | 1 | Auth token management |

**Collection Structure:**

```
IAM API v1.0
├── Auth
│   ├── Login
│   ├── Verify PIN
│   ├── Refresh Token
│   ├── Logout
│   ├── Register
│   ├── Verify Email
│   ├── Setup PIN
│   ├── Change Password
│   ├── Forgot Password
│   ├── Reset Password
│   └── OAuth
│       ├── Google Login
│       └── Google Callback
├── Users
│   ├── Create User
│   ├── Get User
│   ├── Update User
│   ├── Delete User
│   ├── List Users
│   ├── Approve User
│   ├── Reject User
│   ├── Unlock User
│   └── Profile
│       ├── Get My Profile
│       └── Update My Profile
├── Authorization
│   ├── Applications
│   ├── Roles
│   ├── Permissions
│   ├── Role Assignments
│   └── Permission Check
├── Organization
│   ├── Tenants
│   ├── Branches
│   └── User Branches
└── Audit
    ├── Query Logs
    └── Export
```

**Environment Variables:**

```json
{
  "id": "iam-dev",
  "name": "IAM - Development",
  "values": [
    { "key": "baseUrl", "value": "http://localhost:3000/api/iam" },
    { "key": "tenantId", "value": "test-tenant" },
    { "key": "adminEmail", "value": "admin@test.local" },
    { "key": "adminPassword", "value": "Test@123" },
    { "key": "accessToken", "value": "" },
    { "key": "refreshToken", "value": "" }
  ]
}
```

**Pre-request Script (Collection Level):**

```javascript
// Auto-refresh token if expired
const accessToken = pm.environment.get("accessToken");
const tokenExpiry = pm.environment.get("tokenExpiry");

if (accessToken && tokenExpiry) {
    const now = Date.now();
    if (now > parseInt(tokenExpiry) - 60000) { // 1 minute buffer
        // Token expired or about to expire, refresh
        const refreshToken = pm.environment.get("refreshToken");
        if (refreshToken) {
            pm.sendRequest({
                url: pm.environment.get("baseUrl") + "/auth/refresh",
                method: "POST",
                header: { "Content-Type": "application/json" },
                body: { mode: "raw", raw: JSON.stringify({ refresh_token: refreshToken }) }
            }, (err, res) => {
                if (!err && res.code === 200) {
                    const body = res.json();
                    pm.environment.set("accessToken", body.data.access_token);
                    pm.environment.set("refreshToken", body.data.refresh_token);
                    pm.environment.set("tokenExpiry", Date.now() + (body.data.expires_in * 1000));
                }
            });
        }
    }
}
```

**Test Script Example (Login):**

```javascript
pm.test("Status code is 200", () => {
    pm.response.to.have.status(200);
});

pm.test("Response has access_token", () => {
    const jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property("access_token");
    pm.expect(jsonData.data).to.have.property("refresh_token");
    pm.expect(jsonData.data).to.have.property("expires_in");
    
    // Save tokens to environment
    pm.environment.set("accessToken", jsonData.data.access_token);
    pm.environment.set("refreshToken", jsonData.data.refresh_token);
    pm.environment.set("tokenExpiry", Date.now() + (jsonData.data.expires_in * 1000));
});

pm.test("Token is valid JWT", () => {
    const jsonData = pm.response.json();
    const token = jsonData.data.access_token;
    const parts = token.split(".");
    pm.expect(parts.length).to.equal(3);
});
```

**Acceptance Criteria:**
- [ ] All endpoints in collection
- [ ] Environment files for Dev, Staging, Prod
- [ ] Tests for every endpoint
- [ ] Pre-request scripts for auth
- [ ] Collection runs end-to-end (Newman)

---

## 6. Environment Configuration

### 6.1 Overview

Proper environment configuration ensures consistent deployments across development, staging, and production.

### 6.2 Environment Files

#### Development (.env.development)

```bash
# Application
APP_ENV=development
APP_PORT=3000
APP_LOG_LEVEL=debug

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=iam_dev
DB_USER=iam
DB_PASSWORD=iam_dev_password
DB_SSL_MODE=disable
DB_MAX_CONNECTIONS=10

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# OpenSearch
OPENSEARCH_URL=http://localhost:9200
OPENSEARCH_USERNAME=admin
OPENSEARCH_PASSWORD=admin

# Vault (dev mode)
VAULT_ADDR=http://localhost:8200
VAULT_TOKEN=dev-root-token

# SMTP (Mailhog)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USE_TLS=false

# Masterdata
MASTERDATA_URL=http://localhost:3001/api/masterdata
MASTERDATA_USE_MOCK=true

# JWT
JWT_ISSUER=iam-dev
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=7d
```

#### Staging (.env.staging)

```bash
# Application
APP_ENV=staging
APP_PORT=3000
APP_LOG_LEVEL=info

# Database (from Vault)
DB_HOST_FROM_VAULT=true
DB_SSL_MODE=require

# Redis (from Vault)
REDIS_HOST_FROM_VAULT=true

# OpenSearch (from Vault)
OPENSEARCH_URL_FROM_VAULT=true

# Vault
VAULT_ADDR=https://vault.staging.example.com
VAULT_ROLE=iam-staging

# SMTP (from Vault)
SMTP_FROM_VAULT=true

# Masterdata
MASTERDATA_URL=https://masterdata.staging.example.com/api/masterdata
MASTERDATA_USE_MOCK=false

# JWT
JWT_ISSUER=iam-staging
JWT_ACCESS_TOKEN_TTL=15m
JWT_REFRESH_TOKEN_TTL=7d
```

---

## Document Sign-Off

| Role | Name | Signature | Date |
|------|------|-----------|------|
| Tech Lead | | | |
| Backend Lead | | | |
| DevOps Lead | | | |

---

**End of Detailed Task Breakdown**
