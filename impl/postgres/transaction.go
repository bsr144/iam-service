package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// GormTransactionManager implements transaction management using GORM.
// It satisfies any consumer-defined TransactionManager interface with
// the same method signature via Go's structural typing.
//
// Usage in delivery layer (composition root):
//
//	txManager := postgres.NewTransactionManager(db)
//	authUsecase := authinternal.NewUsecase(txManager, ...)
//
// The usecase receives it as contract.TransactionManager (consumer-defined interface).
type GormTransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new GORM-based transaction manager.
// Returns concrete struct (Go idiom: "return structs, accept interfaces").
func NewTransactionManager(db *gorm.DB) *GormTransactionManager {
	return &GormTransactionManager{db: db}
}

// WithTransaction executes fn within a database transaction.
//
// Behavior:
//   - If fn returns nil, the transaction is committed
//   - If fn returns an error, the transaction is rolled back
//   - Supports nested transactions via GORM's savepoint mechanism
//
// The transaction is stored in context and can be retrieved by repositories
// via baseRepository.getDB(ctx). This enables all repository operations
// within fn to participate in the same transaction.
//
// Example:
//
//	err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
//	    if err := userRepo.Create(txCtx, user); err != nil {
//	        return err  // triggers rollback
//	    }
//	    if err := profileRepo.Create(txCtx, profile); err != nil {
//	        return err  // triggers rollback
//	    }
//	    return nil  // triggers commit
//	})
func (m *GormTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// Check for existing transaction (nested transaction support)
	if existingTx, ok := getTx(ctx); ok {
		// Use GORM's nested transaction (creates savepoint)
		return existingTx.Transaction(func(nestedTx *gorm.DB) error {
			nestedCtx := withTx(ctx, nestedTx)
			return fn(nestedCtx)
		})
	}

	// Start new transaction
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := withTx(ctx, tx)
		if err := fn(txCtx); err != nil {
			return fmt.Errorf("transaction failed: %w", err)
		}
		return nil
	})
}
