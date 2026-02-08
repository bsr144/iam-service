package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"iam-service/config"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompleteRegistration(t *testing.T) {
	tenantID := uuid.New()
	registrationID := uuid.New()
	email := "test@example.com"
	ipAddress := "192.168.1.1"
	userAgent := "Mozilla/5.0"
	jwtSecret := "test-secret-key-for-testing-purposes"

	generateValidToken := func() (string, string) {
		claims := jwt.MapClaims{
			"registration_id": registrationID.String(),
			"tenant_id":       tenantID.String(),
			"email":           email,
			"purpose":         RegistrationCompleteTokenPurpose,
			"exp":             time.Now().Add(15 * time.Minute).Unix(),
			"iat":             time.Now().Unix(),
			"jti":             uuid.New().String(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(jwtSecret))
		hash := sha256.Sum256([]byte(tokenString))
		tokenHash := hex.EncodeToString(hash[:])
		return tokenString, tokenHash
	}

	validReq := &authdto.CompleteRegistrationRequest{
		Password:             "SecureP@ssw0rd!",
		PasswordConfirmation: "SecureP@ssw0rd!",
		FirstName:            "John",
		LastName:             "Doe",
	}

	tests := []struct {
		name          string
		req           *authdto.CompleteRegistrationRequest
		setupToken    func() string
		setupMocks    func(*MockTransactionManager, *MockRegistrationSessionStore, *MockUserRepository, *MockUserProfileRepository, *MockUserCredentialsRepository, *MockUserSecurityRepository, *MockUserActivationTrackingRepository, *MockEmailService, string)
		expectedError string
		expectedCode  string
	}{
		{
			name: "success - registration completed",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					TenantID:              tenantID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				credentialsRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserCredentials")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurity")).Return(nil)
				trackingRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserActivationTracking")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, tenantID, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, tenantID, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John").Return(nil)
			},
		},
		{
			name: "error - invalid token format",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "invalid",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - expired token",
			req:  validReq,
			setupToken: func() string {
				claims := jwt.MapClaims{
					"registration_id": registrationID.String(),
					"tenant_id":       tenantID.String(),
					"email":           email,
					"purpose":         RegistrationCompleteTokenPurpose,
					"exp":             time.Now().Add(-1 * time.Minute).Unix(),
					"iat":             time.Now().Add(-16 * time.Minute).Unix(),
					"jti":             uuid.New().String(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(jwtSecret))
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "invalid or expired",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - token purpose mismatch",
			req:  validReq,
			setupToken: func() string {
				claims := jwt.MapClaims{
					"registration_id": registrationID.String(),
					"tenant_id":       tenantID.String(),
					"email":           email,
					"purpose":         "wrong_purpose",
					"exp":             time.Now().Add(15 * time.Minute).Unix(),
					"iat":             time.Now().Unix(),
					"jti":             uuid.New().String(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(jwtSecret))
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {

			},
			expectedError: "not a registration completion token",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - registration not found",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(nil, errors.ErrNotFound("registration not found"))
			},
			expectedError: "not found",
			expectedCode:  errors.CodeNotFound,
		},
		{
			name: "error - session expired",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					TenantID:  tenantID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusVerified,
					ExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - email not verified",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:        registrationID,
					TenantID:  tenantID,
					Email:     email,
					Status:    entity.RegistrationSessionStatusPendingVerification,
					ExpiresAt: time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
			},
			expectedError: "not been verified",
			expectedCode:  errors.CodeForbidden,
		},
		{
			name: "error - token already used (hash mismatch)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				wrongHash := "different-hash-value"
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					TenantID:              tenantID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &wrongHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
			},
			expectedError: "already been used",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - weak password",
			req: &authdto.CompleteRegistrationRequest{
				Password:             "weak",
				PasswordConfirmation: "weak",
				FirstName:            "John",
				LastName:             "Doe",
			},
			setupToken: func() string {
				tokenString, tokenHash := generateValidToken()

				_ = tokenHash
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					TenantID:              tenantID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
			},
			expectedError: "Password",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - email already registered (race condition)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					TenantID:              tenantID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(true, nil)
			},
			expectedError: "already been registered",
			expectedCode:  errors.CodeConflict,
		},
		{
			name: "error - transaction rollback on user creation failure",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					TenantID:              tenantID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusVerified,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, tenantID, registrationID).Return(session, nil)
				userRepo.On("EmailExistsInTenant", mock.Anything, tenantID, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.ErrInternal("database error"))
			},
			expectedError: "failed to create user",
			expectedCode:  errors.CodeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := NewMockTransactionManager()
			redis := new(MockRegistrationSessionStore)
			userRepo := new(MockUserRepository)
			profileRepo := new(MockUserProfileRepository)
			credentialsRepo := new(MockUserCredentialsRepository)
			securityRepo := new(MockUserSecurityRepository)
			trackingRepo := new(MockUserActivationTrackingRepository)
			emailSvc := new(MockEmailService)

			token := tt.setupToken()
			hash := sha256.Sum256([]byte(token))
			tokenHash := hex.EncodeToString(hash[:])

			tt.setupMocks(txManager, redis, userRepo, profileRepo, credentialsRepo, securityRepo, trackingRepo, emailSvc, tokenHash)

			uc := &usecase{
				TxManager: txManager,
				Config: &config.Config{
					JWT: config.JWTConfig{
						AccessSecret: jwtSecret,
					},
				},
				Redis:                      redis,
				UserRepo:                   userRepo,
				UserProfileRepo:            profileRepo,
				UserCredentialsRepo:        credentialsRepo,
				UserSecurityRepo:           securityRepo,
				UserActivationTrackingRepo: trackingRepo,
				EmailService:               emailSvc,
			}

			ctx := context.Background()
			resp, err := uc.CompleteRegistration(ctx, tenantID, registrationID, token, tt.req, ipAddress, userAgent)

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
				assert.NotEqual(t, uuid.Nil, resp.UserID)
				assert.Equal(t, email, resp.Email)
			}

			redis.AssertExpectations(t)
			userRepo.AssertExpectations(t)
			profileRepo.AssertExpectations(t)
			credentialsRepo.AssertExpectations(t)
			securityRepo.AssertExpectations(t)
			trackingRepo.AssertExpectations(t)
		})
	}
}
