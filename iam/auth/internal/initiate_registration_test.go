package internal

import (
	"context"
	"testing"
	"time"

	"iam-service/config"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitiateRegistration(t *testing.T) {
	tenantID := uuid.New()
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"

	tests := []struct {
		name          string
		req           *authdto.InitiateRegistrationRequest
		setupMocks    func(*MockTenantRepository, *MockUserRepository, *MockRegistrationSessionStore, *MockEmailService)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - registration initiated",
			req:  &authdto.InitiateRegistrationRequest{Email: email},
			setupMocks: func(tenantRepo *MockTenantRepository, userRepo *MockUserRepository, redis *MockRegistrationSessionStore, emailSvc *MockEmailService) {
				tenantRepo.On("Exists", mock.Anything, tenantID).Return(true, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(false, nil)
				redis.On("IncrementRegistrationRateLimit", mock.Anything, tenantID, email, mock.Anything).Return(int64(1), nil)
				redis.On("IsRegistrationEmailLocked", mock.Anything, tenantID, email).Return(false, nil)
				redis.On("LockRegistrationEmail", mock.Anything, tenantID, email, mock.Anything).Return(true, nil)
				redis.On("CreateRegistrationSession", mock.Anything, mock.AnythingOfType("*entity.RegistrationSession"), mock.Anything).Return(nil)
				emailSvc.On("SendOTP", mock.Anything, email, mock.AnythingOfType("string"), RegistrationOTPExpiryMinutes).Return(nil)
			},
		},
		{
			name: "error - tenant not found",
			req:  &authdto.InitiateRegistrationRequest{Email: email},
			setupMocks: func(tenantRepo *MockTenantRepository, userRepo *MockUserRepository, redis *MockRegistrationSessionStore, emailSvc *MockEmailService) {
				tenantRepo.On("Exists", mock.Anything, tenantID).Return(false, nil)
			},
			expectedError: "Tenant not found",
			expectedCode:  errors.CodeTenantNotFound,
		},
		{
			name: "error - email already exists",
			req:  &authdto.InitiateRegistrationRequest{Email: email},
			setupMocks: func(tenantRepo *MockTenantRepository, userRepo *MockUserRepository, redis *MockRegistrationSessionStore, emailSvc *MockEmailService) {
				tenantRepo.On("Exists", mock.Anything, tenantID).Return(true, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(true, nil)
			},
			expectedError: "already exists",
			expectedCode:  errors.CodeUserAlreadyExists,
		},
		{
			name: "error - rate limit exceeded",
			req:  &authdto.InitiateRegistrationRequest{Email: email},
			setupMocks: func(tenantRepo *MockTenantRepository, userRepo *MockUserRepository, redis *MockRegistrationSessionStore, emailSvc *MockEmailService) {
				tenantRepo.On("Exists", mock.Anything, tenantID).Return(true, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(false, nil)
				redis.On("IncrementRegistrationRateLimit", mock.Anything, tenantID, email, mock.Anything).Return(int64(RegistrationRateLimitPerHour+1), nil)
			},
			expectedError: "Too many registration attempts",
			expectedCode:  errors.CodeTooManyRequests,
		},
		{
			name: "error - email already locked (registration in progress)",
			req:  &authdto.InitiateRegistrationRequest{Email: email},
			setupMocks: func(tenantRepo *MockTenantRepository, userRepo *MockUserRepository, redis *MockRegistrationSessionStore, emailSvc *MockEmailService) {
				tenantRepo.On("Exists", mock.Anything, tenantID).Return(true, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(false, nil)
				redis.On("IncrementRegistrationRateLimit", mock.Anything, tenantID, email, mock.Anything).Return(int64(1), nil)
				redis.On("IsRegistrationEmailLocked", mock.Anything, tenantID, email).Return(true, nil)
			},
			expectedError: "An active registration already exists",
			expectedCode:  errors.CodeConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tenantRepo := new(MockTenantRepository)
			userRepo := new(MockUserRepository)
			redis := new(MockRegistrationSessionStore)
			emailSvc := new(MockEmailService)

			tt.setupMocks(tenantRepo, userRepo, redis, emailSvc)

			// Create usecase
			uc := &usecase{
				Config:       &config.Config{},
				TenantRepo:   tenantRepo,
				UserRepo:     userRepo,
				Redis:        redis,
				EmailService: emailSvc,
			}

			// Execute
			ctx := context.Background()
			resp, err := uc.InitiateRegistration(ctx, tenantID, tt.req, ipAddress, userAgent)

			// Assert
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				if tt.expectedCode != "" {
					appErr := errors.GetAppError(err)
					require.NotNil(t, appErr, "Expected AppError but got: %v", err)
					assert.Equal(t, tt.expectedCode, appErr.Code)
				}
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, email, resp.Email)
				assert.Equal(t, string(entity.RegistrationSessionStatusPendingVerification), resp.Status)
				assert.NotEmpty(t, resp.RegistrationID)
				assert.True(t, resp.ExpiresAt.After(time.Now()))
			}

			// Verify mocks
			tenantRepo.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			redis.AssertExpectations(t)
		})
	}
}
