package postgres

import (
	"context"
	"fmt"
	"strings"

	"iam-service/entity"
	"iam-service/saving/participant/contract"
	"iam-service/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var allowedParticipantSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"full_name":  true,
	"status":     true,
}

func escapeILIKE(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

type participantRepository struct {
	baseRepository
}

func NewParticipantRepository(db *gorm.DB) contract.ParticipantRepository {
	return &participantRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *participantRepository) Create(ctx context.Context, participant *entity.Participant) error {
	if err := r.getDB(ctx).Create(participant).Error; err != nil {
		return translateError(err, "participant")
	}
	return nil
}

func (r *participantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	var participant entity.Participant
	err := r.getDB(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&participant).Error
	if err != nil {
		return nil, translateError(err, "participant")
	}
	return &participant, nil
}

func (r *participantRepository) Update(ctx context.Context, participant *entity.Participant) error {
	oldVersion := participant.Version
	participant.Version = oldVersion + 1

	result := r.getDB(ctx).Where("version = ?", oldVersion).Save(participant)
	if result.Error != nil {
		participant.Version = oldVersion
		return translateError(result.Error, "participant")
	}
	if result.RowsAffected == 0 {
		participant.Version = oldVersion
		return errors.ErrConflict("participant was modified by another request")
	}
	return nil
}

func (r *participantRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Model(&entity.Participant{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return translateError(err, "participant")
	}
	return nil
}

func (r *participantRepository) List(ctx context.Context, filter *contract.ParticipantFilter) ([]*entity.Participant, int64, error) {
	var participants []*entity.Participant
	var total int64

	query := r.getDB(ctx).Model(&entity.Participant{}).
		Where("tenant_id = ? AND deleted_at IS NULL", filter.TenantID)

	if filter.ApplicationID != nil {
		query = query.Where("application_id = ?", *filter.ApplicationID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != "" {
		search := "%" + escapeILIKE(filter.Search) + "%"
		query = query.Where(
			"full_name ILIKE ? OR ktp_number ILIKE ? OR employee_number ILIKE ? OR phone_number ILIKE ?",
			search, search, search, search,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, translateError(err, "participant")
	}

	sortBy := "created_at"
	if filter.SortBy != "" && allowedParticipantSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}
	sortOrder := "desc"
	if filter.SortOrder == "asc" {
		sortOrder = "asc"
	}

	offset := (filter.Page - 1) * filter.PerPage
	err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&participants).Error

	if err != nil {
		return nil, 0, translateError(err, "participant")
	}

	return participants, total, nil
}
