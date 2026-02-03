package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"iam-service/entity"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"
	jwtpkg "iam-service/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (uc *usecase) Login(ctx context.Context, req *authdto.LoginRequest) (*authdto.LoginResponse, error) {
	user, err := uc.UserRepo.GetByEmail(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to get user").WithError(err)
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials()
	}

	if !user.IsActive {
		return nil, errors.ErrUserSuspended()
	}

	credentials, err := uc.UserCredentialsRepo.GetByUserID(ctx, user.UserID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get credentials").WithError(err)
	}
	if credentials == nil || credentials.PasswordHash == nil {
		return nil, errors.ErrInvalidCredentials()
	}

	err = bcrypt.CompareHashAndPassword([]byte(*credentials.PasswordHash), []byte(req.Password))
	if err != nil {

		if err := uc.logFailedLogin(ctx, user.UserID); err != nil {

		}
		return nil, errors.ErrInvalidCredentials()
	}

	profile, err := uc.UserProfileRepo.GetByUserID(ctx, user.UserID)
	if err != nil {
		return nil, errors.ErrInternal("failed to get profile").WithError(err)
	}
	if profile == nil {
		return nil, errors.ErrProfileIncomplete()
	}

	var productID *uuid.UUID
	if req.ProductCode != nil {
		var product entity.Product
		err := uc.DB.Where("tenant_id = ? AND code = ? AND is_active = ? AND deleted_at IS NULL",
			req.TenantID, *req.ProductCode, true).First(&product).Error
		if err != nil {
			return nil, errors.ErrBadRequest("invalid product code")
		}
		productID = &product.ProductID
	} else if req.ProductID != nil {

		var product entity.Product
		err := uc.DB.Where("product_id = ? AND tenant_id = ? AND is_active = ? AND deleted_at IS NULL",
			*req.ProductID, req.TenantID, true).First(&product).Error
		if err != nil {
			return nil, errors.ErrBadRequest("invalid product id")
		}
		productID = req.ProductID
	}

	roles, err := uc.getUserRoles(ctx, user.UserID, productID)
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

	sessionID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate session ID").WithError(err)
	}

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
		user.UserID,
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

	refreshToken, err := jwtpkg.GenerateRefreshToken(user.UserID, sessionID, tokenConfig)
	if err != nil {
		return nil, errors.ErrInternal("failed to generate refresh token").WithError(err)
	}

	refreshTokenHash := hashToken(refreshToken)

	tokenFamilyID, err := uuid.NewV7()
	if err != nil {
		return nil, errors.ErrInternal("failed to generate token family ID").WithError(err)
	}

	now := time.Now()
	tenantID := req.TenantID
	refreshTokenEntity := &entity.RefreshToken{
		RefreshTokenID: uuid.New(),
		TenantID:       tenantID,
		UserID:         user.UserID,
		TokenHash:      refreshTokenHash,
		TokenFamily:    tokenFamilyID,
		ExpiresAt:      now.Add(uc.Config.JWT.RefreshExpiry),
		CreatedAt:      now,
	}

	if err := uc.DB.Create(refreshTokenEntity).Error; err != nil {
		return nil, errors.ErrInternal("failed to store refresh token").WithError(err)
	}

	if err := uc.updateLastLogin(ctx, user.UserID); err != nil {
		return nil, errors.ErrInternal("failed to update last login").WithError(err)
	}

	response := &authdto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(uc.Config.JWT.AccessExpiry.Seconds()),
		TokenType:    "Bearer",
		User: authdto.UserResponse{
			ID:         user.UserID,
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
	var userRoles []entity.UserRole
	now := time.Now()

	query := uc.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Where("effective_from <= ?", now).
		Where("effective_to IS NULL OR effective_to > ?", now)

	if productID != nil {

		query = query.Where("product_id = ? OR product_id IS NULL", *productID)
	}

	err := query.Find(&userRoles).Error
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

	var roles []*entity.Role
	err = uc.DB.Where("role_id IN ? AND is_active = ?", roleIDs, true).Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (uc *usecase) resolvePermissionsFromRoles(ctx context.Context, roles []*entity.Role) ([]string, error) {
	if len(roles) == 0 {
		return []string{}, nil
	}

	roleIDs := make([]uuid.UUID, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.RoleID
	}

	type permissionResult struct {
		Code string
	}

	var results []permissionResult
	err := uc.DB.Raw(`
		SELECT DISTINCT p.code
		FROM role_permissions rp
		INNER JOIN permissions p ON p.permission_id = rp.permission_id
		WHERE rp.role_id IN ? AND p.deleted_at IS NULL
		ORDER BY p.code
	`, roleIDs).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	permissions := make([]string, len(results))
	for i, r := range results {
		permissions[i] = r.Code
	}

	return permissions, nil
}

func (uc *usecase) updateLastLogin(ctx context.Context, userID uuid.UUID) error {
	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil || security == nil {
		return err
	}

	now := time.Now()
	security.LastLoginAt = &now
	security.FailedLoginAttempts = 0

	return uc.UserSecurityRepo.Update(ctx, security)
}

func (uc *usecase) logFailedLogin(ctx context.Context, userID uuid.UUID) error {
	security, err := uc.UserSecurityRepo.GetByUserID(ctx, userID)
	if err != nil || security == nil {
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
