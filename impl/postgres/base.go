package postgres

import (
	"context"

	"gorm.io/gorm"
)

// baseRepository provides common functionality for all PostgreSQL repositories.
// All repositories should embed this struct to gain transaction-aware database access.
//
// Usage:
//
//	type userRepository struct {
//	    baseRepository
//	}
//
//	func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
//	    return r.getDB(ctx).Create(user).Error
//	}
type baseRepository struct {
	db *gorm.DB
}

// getDB returns the transaction from context if available,
// otherwise returns the default database connection with context.
//
// This is the KEY method that enables transaction participation.
// When a usecase wraps operations in TxManager.WithTransaction(),
// the transaction is stored in context and all repository operations
// automatically use that transaction.
//
// Example flow:
//  1. Usecase calls TxManager.WithTransaction(ctx, fn)
//  2. TransactionManager begins tx and stores in context via withTx()
//  3. Usecase calls repo.Create(txCtx, entity)
//  4. Repository calls r.getDB(txCtx) which returns the tx from context
//  5. All operations within the transaction share the same tx
func (r *baseRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := getTx(ctx); ok {
		return tx
	}
	return r.db.WithContext(ctx)
}
