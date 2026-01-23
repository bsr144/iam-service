package internal

import (
	"context"
	"encoding/json"
	"iam-service/entity"
	"iam-service/internal/auth/authdto"
	"iam-service/pkg/errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (uc *usecase) RegisterSpecialAccount(ctx context.Context, req *authdto.RegisterSpecialAccountRequest) (*authdto.RegisterSpecialAccountResponse, error) {
	tenantExists, err := uc.TenantRepo.Exists(ctx, req.TenantID)
	if err != nil {
		return nil, errors.ErrInternal("failed to verify tenant").WithError(err)
	}
	if !tenantExists {
		return nil, errors.ErrTenantNotFound()
	}

	emailExists, err := uc.UserRepo.EmailExistsInTenant(ctx, req.TenantID, req.Email)
	if err != nil {
		return nil, errors.ErrInternal("failed to check email").WithError(err)
	}
	if emailExists {
		return nil, errors.ErrUserAlreadyExists()
	}

	if err := uc.validatePassword(req.Password); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.ErrInternal("failed to hash password").WithError(err)
	}

	var response *authdto.RegisterSpecialAccountResponse
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

		role, err := uc.RoleRepo.GetByCode(ctx, req.TenantID, req.UserType)
		if err != nil {
			return err
		}

		if role == nil {
			return errors.ErrRoleNotFound()
		}

		userRoleID, err := uuid.NewV7()
		if err != nil {
			return err
		}
		userRole := &entity.UserRole{
			UserRoleID: userRoleID,
			UserID:     userID,
			RoleID:     role.RoleID,
			CreatedAt:  now,
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

		response = &authdto.RegisterSpecialAccountResponse{
			UserID: userID,
			Email:  req.Email,
		}

		return nil
	})

	if err != nil {
		return nil, errors.ErrInternal("failed to create user").WithError(err)
	}

	return response, nil
}
