package internal

import (
	"context"
	"time"

	"iam-service/entity"
	usercontract "iam-service/iam/user/contract"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}

func NewMockTransactionManager() *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	return m
}

type MockTenantRepository struct {
	mock.Mock
}

func (m *MockTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) GetBySlug(ctx context.Context, slug string) (*entity.Tenant, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Tenant), args.Error(1)
}

func (m *MockTenantRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	// Simulate PostgreSQL default UUID generation
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.User, error) {
	args := m.Called(ctx, tenantID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmailAnyTenant(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) EmailExistsInTenant(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	args := m.Called(ctx, tenantID, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, filter *usercontract.UserListFilter) ([]*entity.User, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetPendingApprovalUsers(ctx context.Context, tenantID uuid.UUID) ([]*entity.User, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

type MockRegistrationSessionStore struct {
	mock.Mock
}

func (m *MockRegistrationSessionStore) CreateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) GetRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) (*entity.RegistrationSession, error) {
	args := m.Called(ctx, tenantID, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RegistrationSession), args.Error(1)
}

func (m *MockRegistrationSessionStore) UpdateRegistrationSession(ctx context.Context, session *entity.RegistrationSession, ttl time.Duration) error {
	args := m.Called(ctx, session, ttl)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) DeleteRegistrationSession(ctx context.Context, tenantID, sessionID uuid.UUID) error {
	args := m.Called(ctx, tenantID, sessionID)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) IncrementRegistrationAttempts(ctx context.Context, tenantID, sessionID uuid.UUID) (int, error) {
	args := m.Called(ctx, tenantID, sessionID)
	return args.Int(0), args.Error(1)
}

func (m *MockRegistrationSessionStore) UpdateRegistrationOTP(ctx context.Context, tenantID, sessionID uuid.UUID, otpHash string, expiresAt time.Time) error {
	args := m.Called(ctx, tenantID, sessionID, otpHash, expiresAt)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) MarkRegistrationVerified(ctx context.Context, tenantID, sessionID uuid.UUID, tokenHash string) error {
	args := m.Called(ctx, tenantID, sessionID, tokenHash)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) LockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (bool, error) {
	args := m.Called(ctx, tenantID, email, ttl)
	return args.Bool(0), args.Error(1)
}

func (m *MockRegistrationSessionStore) UnlockRegistrationEmail(ctx context.Context, tenantID uuid.UUID, email string) error {
	args := m.Called(ctx, tenantID, email)
	return args.Error(0)
}

func (m *MockRegistrationSessionStore) IsRegistrationEmailLocked(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	args := m.Called(ctx, tenantID, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockRegistrationSessionStore) IncrementRegistrationRateLimit(ctx context.Context, tenantID uuid.UUID, email string, ttl time.Duration) (int64, error) {
	args := m.Called(ctx, tenantID, email, ttl)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRegistrationSessionStore) GetRegistrationRateLimitCount(ctx context.Context, tenantID uuid.UUID, email string) (int64, error) {
	args := m.Called(ctx, tenantID, email)
	return args.Get(0).(int64), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendOTP(ctx context.Context, email, otp string, expiryMinutes int) error {
	args := m.Called(ctx, email, otp, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendWelcome(ctx context.Context, email, firstName string) error {
	args := m.Called(ctx, email, firstName)
	return args.Error(0)
}

func (m *MockEmailService) SendPasswordReset(ctx context.Context, email, token string, expiryMinutes int) error {
	args := m.Called(ctx, email, token, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendPINReset(ctx context.Context, email, otp string, expiryMinutes int) error {
	args := m.Called(ctx, email, otp, expiryMinutes)
	return args.Error(0)
}

func (m *MockEmailService) SendAdminInvitation(ctx context.Context, email, token string, expiryMinutes int) error {
	args := m.Called(ctx, email, token, expiryMinutes)
	return args.Error(0)
}

type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	// Simulate PostgreSQL default UUID generation
	if profile.UserProfileID == uuid.Nil {
		profile.UserProfileID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

type MockUserCredentialsRepository struct {
	mock.Mock
}

func (m *MockUserCredentialsRepository) Create(ctx context.Context, credentials *entity.UserCredentials) error {
	args := m.Called(ctx, credentials)
	// Simulate PostgreSQL default UUID generation
	if credentials.UserCredentialID == uuid.Nil {
		credentials.UserCredentialID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserCredentialsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserCredentials, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserCredentials), args.Error(1)
}

func (m *MockUserCredentialsRepository) Update(ctx context.Context, credentials *entity.UserCredentials) error {
	args := m.Called(ctx, credentials)
	return args.Error(0)
}

type MockUserSecurityRepository struct {
	mock.Mock
}

func (m *MockUserSecurityRepository) Create(ctx context.Context, security *entity.UserSecurity) error {
	args := m.Called(ctx, security)
	// Simulate PostgreSQL default UUID generation
	if security.UserSecurityID == uuid.Nil {
		security.UserSecurityID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserSecurityRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserSecurity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSecurity), args.Error(1)
}

func (m *MockUserSecurityRepository) Update(ctx context.Context, security *entity.UserSecurity) error {
	args := m.Called(ctx, security)
	return args.Error(0)
}

type MockUserActivationTrackingRepository struct {
	mock.Mock
}

func (m *MockUserActivationTrackingRepository) Create(ctx context.Context, tracking *entity.UserActivationTracking) error {
	args := m.Called(ctx, tracking)
	// Simulate PostgreSQL default UUID generation
	if tracking.UserActivationTrackingID == uuid.Nil {
		tracking.UserActivationTrackingID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserActivationTrackingRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserActivationTracking, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserActivationTracking), args.Error(1)
}

func (m *MockUserActivationTrackingRepository) Update(ctx context.Context, tracking *entity.UserActivationTracking) error {
	args := m.Called(ctx, tracking)
	return args.Error(0)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	// Simulate PostgreSQL default UUID generation
	if role.RoleID == uuid.Nil {
		role.RoleID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*entity.Role, error) {
	args := m.Called(ctx, tenantID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Role, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entity.Role, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Role), args.Error(1)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, token *entity.RefreshToken) error {
	args := m.Called(ctx, token)
	// Simulate PostgreSQL default UUID generation
	if token.RefreshTokenID == uuid.Nil {
		token.RefreshTokenID = uuid.New()
	}
	if token.TokenFamily == uuid.Nil {
		token.TokenFamily = uuid.New()
	}
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	args := m.Called(ctx, id, reason)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeAllByUserID(ctx context.Context, userID uuid.UUID, reason string) error {
	args := m.Called(ctx, userID, reason)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeByFamily(ctx context.Context, tokenFamily uuid.UUID, reason string) error {
	args := m.Called(ctx, tokenFamily, reason)
	return args.Error(0)
}

type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ctx context.Context, userRole *entity.UserRole) error {
	args := m.Called(ctx, userRole)
	// Simulate PostgreSQL default UUID generation
	if userRole.UserRoleID == uuid.Nil {
		userRole.UserRoleID = uuid.New()
	}
	return args.Error(0)
}

func (m *MockUserRoleRepository) ListActiveByUserID(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]entity.UserRole, error) {
	args := m.Called(ctx, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.UserRole), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByCodeAndTenant(ctx context.Context, tenantID uuid.UUID, code string) (*entity.Product, error) {
	args := m.Called(ctx, tenantID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetByIDAndTenant(ctx context.Context, productID, tenantID uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, productID, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetCodesByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]string, error) {
	args := m.Called(ctx, roleIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}
