package internal

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"iam-service/entity"
	"iam-service/iam/user/userdto"
	"iam-service/impl/postgres"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (uc *usecase) Create(ctx context.Context, req *userdto.CreateRequest) (*userdto.CreateResponse, error) {
	tenantExists, err := uc.TenantRepo.Exists(ctx, req.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenantExists {
		return nil, errors.ErrTenantNotFound()
	}

	role, err := uc.RoleRepo.GetByCode(ctx, req.TenantID, req.RoleCode)
	if err != nil {
		if stderrors.Is(err, postgres.ErrRecordNotFound) {
			return nil, errors.ErrRoleNotFound()
		}
		return nil, errors.ErrInternal("failed to get role").WithError(err)
	}

	if !role.IsSystem {
		return nil, errors.ErrBadRequest("Only system roles can be assigned through this endpoint")
	}

	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrUserAlreadyExists()
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}

	var response *userdto.CreateResponse
	now := time.Now()
	passwordHashStr := string(passwordHash)

	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		userID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		user := &entity.User{
			UserID:        userID,
			TenantID:      &req.TenantID,
			BranchID:      req.BranchID,
			Email:         req.Email,
			EmailVerified: true,
			IsActive:      true,
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userCredentialID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		credentials := &entity.UserCredentials{
			UserCredentialID: userCredentialID,
			UserID:           userID,
			PasswordHash:     &passwordHashStr,
			PasswordHistory:  json.RawMessage("[]"),
			PINHistory:       json.RawMessage("[]"),
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := tx.Create(credentials).Error; err != nil {
			return err
		}
		userProfileID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		profile := &entity.UserProfile{
			UserProfileID: userProfileID,
			UserID:        userID,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		userRoleID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		userRole := &entity.UserRole{
			UserRoleID:    userRoleID,
			UserID:        userID,
			RoleID:        role.RoleID,
			BranchID:      req.BranchID,
			EffectiveFrom: now,
			CreatedAt:     now,
		}
		if err := tx.Create(userRole).Error; err != nil {
			return err
		}
		userSecurityID, err := uuid.NewV7()
		if err != nil {
			return err
		}

		security := &entity.UserSecurity{
			UserSecurityID: userSecurityID,
			UserID:         userID,
			Metadata:       json.RawMessage("{}"),
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := tx.Create(security).Error; err != nil {
			return err
		}
		tracking := entity.NewUserActivationTracking(userID, &req.TenantID)
		if err := tracking.MarkUserCreatedBySystem(); err != nil {
			return err
		}
		if err := tx.Create(tracking).Error; err != nil {
			return err
		}

		response = &userdto.CreateResponse{
			UserID:   userID,
			Email:    req.Email,
			FullName: req.FirstName + " " + req.LastName,
			RoleCode: req.RoleCode,
			TenantID: req.TenantID,
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	return response, nil
}
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.ErrValidation("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.ErrValidation("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.ErrValidation("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.ErrValidation("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.ErrValidation("password must contain at least one special character")
	}

	return nil
}
