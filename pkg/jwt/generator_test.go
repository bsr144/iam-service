package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAccessToken_JTI(t *testing.T) {
	config := &TokenConfig{
		SigningMethod: "HS256",
		AccessSecret:  "test-secret",
		AccessExpiry:  15 * time.Minute,
		Issuer:        "iam-service",
		Audience:      []string{"iam-service"},
	}

	userID := uuid.New()
	email := "test@example.com"
	sessionID := uuid.New()

	token, err := GenerateAccessToken(
		userID,
		email,
		nil,
		nil,
		[]string{"user"},
		[]string{"read:profile"},
		nil,
		sessionID,
		config,
	)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Parse token to verify JTI
	claims, err := ParseAccessToken(token, config)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify JTI is set and is a valid UUID
	assert.NotEmpty(t, claims.RegisteredClaims.ID, "JTI should be set")
	_, err = uuid.Parse(claims.RegisteredClaims.ID)
	assert.NoError(t, err, "JTI should be a valid UUID")
}

func TestGenerateAccessToken_JTI_Uniqueness(t *testing.T) {
	config := &TokenConfig{
		SigningMethod: "HS256",
		AccessSecret:  "test-secret",
		AccessExpiry:  15 * time.Minute,
		Issuer:        "iam-service",
		Audience:      []string{"iam-service"},
	}

	userID := uuid.New()
	email := "test@example.com"
	sessionID := uuid.New()

	// Generate two tokens
	token1, err := GenerateAccessToken(
		userID,
		email,
		nil,
		nil,
		[]string{"user"},
		[]string{"read:profile"},
		nil,
		sessionID,
		config,
	)
	require.NoError(t, err)

	token2, err := GenerateAccessToken(
		userID,
		email,
		nil,
		nil,
		[]string{"user"},
		[]string{"read:profile"},
		nil,
		sessionID,
		config,
	)
	require.NoError(t, err)

	// Parse both tokens
	claims1, err := ParseAccessToken(token1, config)
	require.NoError(t, err)

	claims2, err := ParseAccessToken(token2, config)
	require.NoError(t, err)

	// Verify JTIs are different
	assert.NotEqual(t, claims1.RegisteredClaims.ID, claims2.RegisteredClaims.ID,
		"JTI should be unique for each token")
}

func TestGenerateMultiTenantAccessToken_JTI(t *testing.T) {
	config := &TokenConfig{
		SigningMethod: "HS256",
		AccessSecret:  "test-secret",
		AccessExpiry:  15 * time.Minute,
		Issuer:        "iam-service",
		Audience:      []string{"iam-service"},
	}

	userID := uuid.New()
	email := "test@example.com"
	sessionID := uuid.New()
	tenantID := uuid.New()

	tenants := []TenantClaim{
		{
			TenantID: tenantID,
			Products: []ProductClaim{
				{
					ProductID:   uuid.New(),
					ProductCode: "APP1",
					Roles:       []string{"admin"},
					Permissions: []string{"read:all", "write:all"},
				},
			},
		},
	}

	token, err := GenerateMultiTenantAccessToken(
		userID,
		email,
		tenants,
		sessionID,
		config,
	)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Parse token to verify JTI
	claims, err := ParseMultiTenantAccessToken(token, config)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify JTI is set and is a valid UUID
	assert.NotEmpty(t, claims.RegisteredClaims.ID, "JTI should be set")
	_, err = uuid.Parse(claims.RegisteredClaims.ID)
	assert.NoError(t, err, "JTI should be a valid UUID")
}
