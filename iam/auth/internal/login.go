package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrInvalidCredentials()
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.ErrUserSuspended()
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, user.ID)
	if err != nil {

		if errors.IsNotFound(err) {
			return nil, errors.ErrInvalidCredentials()
		}
		return nil, err
	}
	if credentials.PasswordHash == nil {
		return nil, errors.ErrInvalidCredentials()
	}

	err = bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(req.Password))
	if err != nil {

		if err := uc.logFailedLogin(ctx, user.ID); err != nil {

		}
		return nil, errors.ErrInvalidCredentials()
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.ErrProfileIncomplete()
		}
		return nil, err
	}

	var productID *uuid.UUID
	if req.ProductCode != nil {
		product, err := uc.ProductRepo.GetByCodeAndTenant(ctx, req.TenantID, *req.ProductCode)
		if err != nil {
			return nil, errors.ErrBadRequest("invalid product code")
		}
		productID = &product.ID
	} else if req.ProductID != nil {
		product, err := uc.ProductRepo.GetByIDAndTenant(ctx, *req.ProductID, req.TenantID)
		if err != nil {
			return nil, errors.ErrBadRequest("invalid product id")
		}
		productID = &product.ID
	}

	roles, err := uc.getUserRoles(ctx, user.ID, productID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user roles").WithError(err)
	}

	if productID != nil && len(roles) == 0 {
		return nil, errors.ErrForbidden("no access to the specified product")
	}

	roleCodes := make([]string, len(roles))
	for i, role := range roles {
		roleCodes[i] = role.Code
	}

	permissions, err := uc.resolvePermissionsFromRoles(ctx, roles)
	if err != nil {
		return nil, errors.ErrInternal("failed to resolve permissions").WithError(err)
	}

	sessionID := uuid.New()

	tokenConfig := &jwtpkg.TokenConfig{
		SigningMethod: uc.Config.JWT.SigningMethod,
		AccessSecret:  uc.Config.JWT.AccessSecret,
		RefreshSecret: uc.Config.JWT.RefreshSecret,
		AccessExpiry:  uc.Config.JWT.AccessExpiry,
		RefreshExpiry: uc.Config.JWT.RefreshExpiry,
		Issuer:        uc.Config.JWT.Issuer,
		Audience:      uc.Config.JWT.Audience,
	}

	if uc.Config.JWT.SigningMethod == "RS256" {
		privateKey, err := jwtpkg.LoadPrivateKeyFromFile(uc.Config.JWT.PrivateKeyPath)
		if err != nil {
			return nil, errors.ErrInternal("failed to load private key").WithError(err)
		}
		publicKey, err := jwtpkg.LoadPublicKeyFromFile(uc.Config.JWT.PublicKeyPath)
		if err != nil {
			return nil, errors.ErrInternal("failed to load public key").WithError(err)
		}
		tokenConfig.PrivateKey = privateKey
		tokenConfig.PublicKey = publicKey
	}

	accessToken, err := jwtpkg.GenerateAccessToken(
		user.ID,
		user.Email,
		user.TenantID,
		productID,
		roleCodes,
		permissions,
		user.BranchID,
		sessionID,
		tokenConfig,
	)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate access token").WithError(err)
	}

	refreshToken, err := jwtpkg.GenerateRefreshToken(user.ID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)

	now := time.Now()
	refreshTokenEntity := &entity.RefreshToken{
		TenantID:  req.TenantID,
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt: now,
	}

	if err := uc.RefreshTokenRepo.Create(ctx, refreshTokenEntity); err != nil {
		return nil, errors.ErrInternal("failed to store refresh token").WithError(err)
	}

	if err := uc.updateLastLogin(ctx, user.ID); err != nil {
		return nil, errors.ErrInternal("failed to update last login").WithError(err)
	}

	response := &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		User: authdto.UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			FullName:   profile.FullName(),
			TenantID:   user.TenantID,
			ProductID:  productID,
			BranchID:   user.BranchID,
			Roles:      roleCodes,
			MFAEnabled: credentials.MFAEnabled,
		},
	}

	return response, nil
}

func (uc *usecase) getUserRoles(ctx context.Context, userID uuid.UUID, productID *uuid.UUID) ([]*entity.Role, error) {
	userRoles, err := uc.UserRoleRepo.ListActiveByUserID(ctx, userID, productID)
	if err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []*entity.Role{}, nil
	}

	roleIDs := make([]uuid.UUID, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	return uc.RoleRepo.GetByIDs(ctx, roleIDs)
}

func (uc *usecase) resolvePermissionsFromRoles(ctx context.Context, roles []*entity.Role) ([]string, error) {
	if len(roles) == 0 {
		return []string{}, nil
	}

	roleIDs := make([]uuid.UUID, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	return uc.PermissionRepo.GetCodesByRoleIDs(ctx, roleIDs)
}

func (uc *usecase) updateLastLogin(ctx context.Context, userID uuid.UUID) error {
	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	now := time.Now()
	security.LastLoginAt = &now
	security.FailedLoginAttempts = 0

	return uc.UserSecurityRepo.Update(ctx, security)
}

func (uc *usecase) logFailedLogin(ctx context.Context, userID uuid.UUID) error {
	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	security.FailedLoginAttempts++

	if security.FailedLoginAttempts >= 5 {
		lockUntil := time.Now().Add(15 * time.Minute)
		security.LockedUntil = &lockUntil
	}

	return uc.UserSecurityRepo.Update(ctx, security)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
