package postgres

import (
	"context"

	"iam-service/entity"
	membercontract "iam-service/saving/member/contract"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type productRegistrationConfigRepository struct {
	baseRepository
}

func NewProductRegistrationConfigRepository(db *gorm.DB) membercontract.ProductRegistrationConfigRepository {
	return &productRegistrationConfigRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *productRegistrationConfigRepository) GetByApplicationAndType(ctx context.Context, applicationID uuid.UUID, regType string) (*entity.ProductRegistrationConfig, error) {
	var config entity.ProductRegistrationConfig
	err := r.getDB(ctx).
		Where("application_id = ? AND registration_type = ?", applicationID, regType).
		First(&config).Error
	if err != nil {
		return nil, translateError(err, "product registration config")
	}
	return &config, nil
}
