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

func TestCompleteProfileRegistration(t *testing.T) {
	registrationID := uuid.New()
	email := "test@example.com"
	jwtSecret := "test-secret-key-for-testing-purposes"
	passwordHash := "$2a$10$hashedpassword"

	generateValidToken := func() (string, string) {
		claims := jwt.MapClaims{
			"registration_id": registrationID.String(),
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

	validReq := &authdto.CompleteProfileRegistrationRequest{
		FullName:      "John Michael Smith",
		PhoneNumber:   "+6281234567890",
		DateOfBirth:   "1990-01-15",
		Gender:        "male",
		MaritalStatus: "single",
		Address:       "Jl. Sudirman No. 123, Jakarta Pusat",
		PlaceOfBirth:  "Jakarta",
	}

	tests := []struct {
		name          string
		req           *authdto.CompleteProfileRegistrationRequest
		setupToken    func() string
		setupMocks    func(*MockTransactionManager, *MockRegistrationSessionStore, *MockUserRepository, *MockUserProfileRepository, *MockUserCredentialsRepository, *MockUserSecurityRepository, *MockUserActivationTrackingRepository, *MockEmailService, *MockRefreshTokenRepository, string)
		expectedError string
		expectedCode  string
		validateResp  func(*testing.T, *authdto.CompleteProfileRegistrationResponse)
	}{
		{
			name: "success - new flow with PASSWORD_SET status",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				credentialsRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserCredentials")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurity")).Return(nil)
				trackingRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserActivationTracking")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John Michael").Return(nil)
			},
			validateResp: func(t *testing.T, resp *authdto.CompleteProfileRegistrationResponse) {
				assert.Equal(t, email, resp.Email)
				assert.NotEmpty(t, resp.AccessToken)
				assert.NotEmpty(t, resp.RefreshToken)
				assert.Equal(t, "Bearer", resp.TokenType)
				assert.Greater(t, resp.ExpiresIn, 0)
				assert.Equal(t, "John Michael", resp.Profile.FirstName)
				assert.Equal(t, "Smith", resp.Profile.LastName)
			},
		},
		{
			name: "error - invalid token",
			req:  validReq,
			setupToken: func() string {
				return "invalid-token"
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
			},
			expectedError: "invalid",
			expectedCode:  errors.CodeUnauthorized,
		},
		{
			name: "error - expired session",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(-1 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "expired",
			expectedCode:  "REGISTRATION_EXPIRED",
		},
		{
			name: "error - wrong status PENDING_VERIFICATION",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPendingVerification,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "not ready",
			expectedCode:  errors.CodeForbidden,
		},
		{
			name: "error - missing full_name",
			req: &authdto.CompleteProfileRegistrationRequest{
				FullName:      "",
				PhoneNumber:   "+6281234567890",
				DateOfBirth:   "1990-01-15",
				Gender:        "male",
				MaritalStatus: "single",
				Address:       "Jl. Sudirman No. 123",
				PlaceOfBirth:  "Jakarta",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "full_name",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - age under 18",
			req: &authdto.CompleteProfileRegistrationRequest{
				FullName:      "Young Person",
				PhoneNumber:   "+6281234567890",
				DateOfBirth:   time.Now().AddDate(-17, 0, 0).Format("2006-01-02"),
				Gender:        "male",
				MaritalStatus: "single",
				Address:       "Jl. Sudirman No. 123",
				PlaceOfBirth:  "Jakarta",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "18 years",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - invalid gender value",
			req: &authdto.CompleteProfileRegistrationRequest{
				FullName:      "John Doe",
				PhoneNumber:   "+6281234567890",
				DateOfBirth:   "1990-01-15",
				Gender:        "unknown",
				MaritalStatus: "single",
				Address:       "Jl. Sudirman No. 123",
				PlaceOfBirth:  "Jakarta",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
			},
			expectedError: "gender",
			expectedCode:  errors.CodeValidation,
		},
		{
			name: "error - email already registered (race condition)",
			req:  validReq,
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(true, nil)
			},
			expectedError: "already been registered",
			expectedCode:  errors.CodeConflict,
		},
		{
			name: "success - name splitting: John Smith",
			req: &authdto.CompleteProfileRegistrationRequest{
				FullName:      "John Smith",
				PhoneNumber:   "+6281234567890",
				DateOfBirth:   "1990-01-15",
				Gender:        "male",
				MaritalStatus: "single",
				Address:       "Jl. Sudirman No. 123",
				PlaceOfBirth:  "Jakarta",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				credentialsRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserCredentials")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurity")).Return(nil)
				trackingRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserActivationTracking")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "John").Return(nil)
			},
			validateResp: func(t *testing.T, resp *authdto.CompleteProfileRegistrationResponse) {
				assert.Equal(t, "John", resp.Profile.FirstName)
				assert.Equal(t, "Smith", resp.Profile.LastName)
			},
		},
		{
			name: "success - name splitting: Madonna (single name)",
			req: &authdto.CompleteProfileRegistrationRequest{
				FullName:      "Madonna",
				PhoneNumber:   "+6281234567890",
				DateOfBirth:   "1990-01-15",
				Gender:        "female",
				MaritalStatus: "single",
				Address:       "Jl. Sudirman No. 123",
				PlaceOfBirth:  "Jakarta",
			},
			setupToken: func() string {
				tokenString, _ := generateValidToken()
				return tokenString
			},
			setupMocks: func(txManager *MockTransactionManager, redis *MockRegistrationSessionStore, userRepo *MockUserRepository, profileRepo *MockUserProfileRepository, credentialsRepo *MockUserCredentialsRepository, securityRepo *MockUserSecurityRepository, trackingRepo *MockUserActivationTrackingRepository, emailSvc *MockEmailService, refreshTokenRepo *MockRefreshTokenRepository, tokenHash string) {
				session := &entity.RegistrationSession{
					ID:                    registrationID,
					Email:                 email,
					Status:                entity.RegistrationSessionStatusPasswordSet,
					RegistrationTokenHash: &tokenHash,
					ExpiresAt:             time.Now().Add(10 * time.Minute),
				}
				redis.On("GetRegistrationSession", mock.Anything, registrationID).Return(session, nil)
				redis.On("GetRegistrationPasswordHash", mock.Anything, registrationID).Return(passwordHash, nil)
				userRepo.On("EmailExists", mock.Anything, email).Return(false, nil)
				userRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)
				credentialsRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserCredentials")).Return(nil)
				profileRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserProfile")).Return(nil)
				securityRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserSecurity")).Return(nil)
				trackingRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.UserActivationTracking")).Return(nil)
				refreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				redis.On("DeleteRegistrationSession", mock.Anything, registrationID).Return(nil)
				redis.On("UnlockRegistrationEmail", mock.Anything, email).Return(nil)
				emailSvc.On("SendWelcome", mock.Anything, email, "Madonna").Return(nil)
			},
			validateResp: func(t *testing.T, resp *authdto.CompleteProfileRegistrationResponse) {
				assert.Equal(t, "Madonna", resp.Profile.FirstName)
				assert.Equal(t, "", resp.Profile.LastName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txManager := NewMockTransactionManager()
			redis := &MockRegistrationSessionStore{}
			userRepo := &MockUserRepository{}
			profileRepo := &MockUserProfileRepository{}
			credentialsRepo := &MockUserCredentialsRepository{}
			securityRepo := &MockUserSecurityRepository{}
			trackingRepo := &MockUserActivationTrackingRepository{}
			emailSvc := &MockEmailService{}
			refreshTokenRepo := &MockRefreshTokenRepository{}

			tokenString := tt.setupToken()
			_, tokenHash := generateValidToken()
			if tokenString != "invalid-token" {
				hash := sha256.Sum256([]byte(tokenString))
				tokenHash = hex.EncodeToString(hash[:])
			}
			tt.setupMocks(txManager, redis, userRepo, profileRepo, credentialsRepo, securityRepo, trackingRepo, emailSvc, refreshTokenRepo, tokenHash)

			cfg := &config.Config{
				JWT: config.JWTConfig{
					AccessSecret:  jwtSecret,
					RefreshSecret: "refresh-secret",
					SigningMethod: "HS256",
					AccessExpiry:  3600 * time.Second,
					RefreshExpiry: 86400 * time.Second,
					Issuer:        "iam-service",
					Audience:      []string{"iam-api"},
				},
			}

			uc := &usecase{
				Config:                     cfg,
				TxManager:                  txManager,
				Redis:                      redis,
				UserRepo:                   userRepo,
				UserProfileRepo:            profileRepo,
				UserCredentialsRepo:        credentialsRepo,
				UserSecurityRepo:           securityRepo,
				UserActivationTrackingRepo: trackingRepo,
				EmailService:               emailSvc,
				RefreshTokenRepo:           refreshTokenRepo,
			}

			tt.req.RegistrationID = registrationID
			tt.req.RegistrationToken = tokenString

			response, err := uc.CompleteProfileRegistration(
				context.Background(),
				tt.req,
			)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				appErr, ok := err.(*errors.AppError)
				require.True(t, ok, "Error should be AppError")
				assert.Equal(t, tt.expectedCode, appErr.Code)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, response)

			if tt.validateResp != nil {
				tt.validateResp(t, response)
			}
		})
	}
}
