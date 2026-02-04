package contract

import (
	"context"
	"time"

	"iam-service/entity"
	usercontract "iam-service/iam/user/contract"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.User, error)
	GetByEmailAnyTenant(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	EmailExistsInTenant(ctx context.Context, tenantID uuid.UUID, email string) (bool, error)
	List(ctx context.Context, filter *usercontract.UserListFilter) ([]*entity.User, int64, error)
	GetPendingApprovalUsers(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error)
}
type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
}
type UserCredentialsRepository interface {
	Create(ctx context.Context, credentials *entity.UserCredentials) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredentials, error)
	Update(ctx context.Context, credentials *entity.UserCredentials) error
}
type UserSecurityRepository interface {
	Create(ctx context.Context, security *entity.UserSecurity) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurity, error)
	Update(ctx context.Context, security *entity.UserSecurity) error
}
type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *entity.EmailVerification) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.EmailVerification, error)
	GetLatestByEmail(ctx context.Context, email string, otpType entity.OTPType) (*entity.EmailVerification, error)
	GetLatestByUserID(ctx context.Context, userID uuid.UUID, otpType entity.OTPType) (*entity.EmailVerification, error)
	MarkAsVerified(ctx context.Context, id uuid.UUID) error
	CountActiveOTPsByEmail(ctx context.Context, email string, otpType entity.OTPType) (int, error)
	DeleteExpiredByEmail(ctx context.Context, email string) error
}
type TenantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}
type UserActivationTrackingRepository interface {
	Create(ctx context.Context, tracking *entity.UserActivationTracking) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserActivationTracking, error)
	Update(ctx context.Context, tracking *entity.UserActivationTracking) error
}

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Role, error)
	GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Role, error)
}
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *entity.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID, reason string) error
	RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error
	RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error
}
type PINVerificationLogRepository interface {
	Create(ctx context.Context, log *entity.PINVerificationLog) error
	CountRecentFailures(ctx context.Context, userID uuid.UUID, since int) (int, error)
}

type RegistrationSessionStore interface {
	CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error
	GetRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) (*entity.RegistrationSession, error)
	UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error
	DeleteRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) error

	IncrementRegistrationAttempts(ctx context.Context, tenantID, sessionID uuid.UUID) (int, error)
	UpdateRegistrationOTP(ctx context.Context, tenantID, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error
	MarkRegistrationVerified(ctx context.Context, tenantID, sessionID uuid.UUID, tokenHash string) error

	LockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (bool, error)
	UnlockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string) error
	IsRegistrationEmailLocked(ctx context.Context, tenantID uuid.UUID, email string) (bool, error)

	IncrementRegistrationRateLimit(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (int64, error)
	GetRegistrationRateLimitCount(ctx context.Context, tenantID uuid.UUID, email string) (int64, error)
}
