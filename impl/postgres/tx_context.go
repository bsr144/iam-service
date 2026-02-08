package postgres

import (
	"context"

	"gorm.io/gorm"
)

// txKey is the context key for storing the GORM transaction.
// Using unexported struct type prevents key collisions with other packages.
type txKey struct{}

// withTx stores the GORM transaction in context.
// This is called by GormTransactionManager when starting a transaction.
func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// getTx retrieves the GORM transaction from context.
// Returns (nil, false) if no transaction is active.
// This is called by repositories via baseRepository.getDB().
func getTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

// hasTx checks if a transaction exists in context.
// Useful for debugging or conditional logic.
func hasTx(ctx context.Context) bool {
	_, ok := getTx(ctx)
	return ok
}
