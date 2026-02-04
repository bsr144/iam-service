package internal

import (
	"context"
	"time"

	"iam-service/entity"
	usercontract "iam-service/iam/user/contract"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockTenantRepository is a mock implementation of TenantRepository
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

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
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

// MockRegistrationSessionStore is a mock implementation of RegistrationSessionStore
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

// MockEmailService is a mock implementation of EmailService
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
