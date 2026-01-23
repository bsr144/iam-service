package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// RoleAssignmentQueueStatus represents the state of a role assignment in the queue
type RoleAssignmentQueueStatus string

const (
	RoleAssignmentQueueStatusPending    RoleAssignmentQueueStatus = "pending"
	RoleAssignmentQueueStatusProcessing RoleAssignmentQueueStatus = "processing"
	RoleAssignmentQueueStatusCompleted  RoleAssignmentQueueStatus = "completed"
	RoleAssignmentQueueStatusFailed     RoleAssignmentQueueStatus = "failed"
	RoleAssignmentQueueStatusCancelled  RoleAssignmentQueueStatus = "cancelled"
)

// RoleAssignmentQueue tracks admin role assignments before user activation
// Supports bulk operations: admin can assign roles to many users at once
// Provides audit trail and progress tracking
type RoleAssignmentQueue struct {
	QueueID  uuid.UUID `json:"queue_id" gorm:"column:queue_id;primaryKey" db:"queue_id"`

	// Target user
	UserID   uuid.UUID `json:"user_id" gorm:"column:user_id;not null" db:"user_id"`
	TenantID uuid.UUID `json:"tenant_id" gorm:"column:tenant_id;not null" db:"tenant_id"`

	// Role to assign
	RoleID    uuid.UUID  `json:"role_id" gorm:"column:role_id;not null" db:"role_id"`
	ProductID *uuid.UUID `json:"product_id,omitempty" gorm:"column:product_id" db:"product_id"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty" gorm:"column:branch_id" db:"branch_id"`

	// Effective period
	EffectiveFrom time.Time  `json:"effective_from" gorm:"column:effective_from;default:CURRENT_TIMESTAMP" db:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty" gorm:"column:effective_to" db:"effective_to"`

	// Queue management
	Status RoleAssignmentQueueStatus `json:"status" gorm:"column:status;not null;default:'pending'" db:"status"`

	// Admin context (who initiated the assignment)
	AssignedBy uuid.UUID `json:"assigned_by" gorm:"column:assigned_by;not null" db:"assigned_by"`
	AssignedAt time.Time `json:"assigned_at" gorm:"column:assigned_at;not null;default:CURRENT_TIMESTAMP" db:"assigned_at"`

	// Bulk operation tracking
	BatchID       *uuid.UUID `json:"batch_id,omitempty" gorm:"column:batch_id" db:"batch_id"`       // Groups assignments together
	BatchTotal    *int       `json:"batch_total,omitempty" gorm:"column:batch_total" db:"batch_total"` // Total in batch
	BatchSequence *int       `json:"batch_sequence,omitempty" gorm:"column:batch_sequence" db:"batch_sequence"` // Position in batch

	// Processing tracking
	ProcessedAt         *time.Time `json:"processed_at,omitempty" gorm:"column:processed_at" db:"processed_at"`
	ProcessingStartedAt *time.Time `json:"processing_started_at,omitempty" gorm:"column:processing_started_at" db:"processing_started_at"`
	FailureReason       *string    `json:"failure_reason,omitempty" gorm:"column:failure_reason" db:"failure_reason"`
	RetryCount          int        `json:"retry_count" gorm:"column:retry_count;default:0" db:"retry_count"`

	// Result (links to created user_role after successful processing)
	UserRoleID *uuid.UUID `json:"user_role_id,omitempty" gorm:"column:user_role_id" db:"user_role_id"`

	// Metadata
	Metadata json.RawMessage `json:"metadata,omitempty" gorm:"column:metadata;type:jsonb;default:'{}'" db:"metadata"`
}

func (RoleAssignmentQueue) TableName() string {
	return "role_assignments_queue"
}

// IsPending checks if the assignment is still pending
func (raq *RoleAssignmentQueue) IsPending() bool {
	return raq.Status == RoleAssignmentQueueStatusPending
}

// IsProcessing checks if the assignment is currently being processed
func (raq *RoleAssignmentQueue) IsProcessing() bool {
	return raq.Status == RoleAssignmentQueueStatusProcessing
}

// IsCompleted checks if the assignment has been completed successfully
func (raq *RoleAssignmentQueue) IsCompleted() bool {
	return raq.Status == RoleAssignmentQueueStatusCompleted
}

// IsFailed checks if the assignment has failed
func (raq *RoleAssignmentQueue) IsFailed() bool {
	return raq.Status == RoleAssignmentQueueStatusFailed
}

// IsCancelled checks if the assignment has been cancelled
func (raq *RoleAssignmentQueue) IsCancelled() bool {
	return raq.Status == RoleAssignmentQueueStatusCancelled
}

// CanBeProcessed checks if the assignment can be processed
func (raq *RoleAssignmentQueue) CanBeProcessed() bool {
	return raq.Status == RoleAssignmentQueueStatusPending || raq.Status == RoleAssignmentQueueStatusFailed
}

// CanBeRetried checks if the assignment can be retried after failure
func (raq *RoleAssignmentQueue) CanBeRetried() bool {
	return raq.Status == RoleAssignmentQueueStatusFailed
}

// IsPartOfBatch checks if this assignment is part of a bulk operation
func (raq *RoleAssignmentQueue) IsPartOfBatch() bool {
	return raq.BatchID != nil
}

// GetBatchProgress returns the progress percentage if part of a batch
func (raq *RoleAssignmentQueue) GetBatchProgress() *float64 {
	if !raq.IsPartOfBatch() || raq.BatchTotal == nil || raq.BatchSequence == nil {
		return nil
	}
	if *raq.BatchTotal == 0 {
		return nil
	}
	progress := float64(*raq.BatchSequence) / float64(*raq.BatchTotal) * 100
	return &progress
}
