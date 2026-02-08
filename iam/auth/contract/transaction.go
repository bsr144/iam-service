package contract

import "context"

// TransactionManager handles database transaction lifecycle.
// This interface is defined by the consumer (usecase layer) following
// the Go idiom "accept interfaces, return structs".
//
// The implementation (e.g., postgres.GormTransactionManager) satisfies
// this interface via Go's structural typing - no explicit declaration needed.
type TransactionManager interface {
	// WithTransaction executes fn within a database transaction.
	// If fn returns an error, the transaction is rolled back.
	// If fn returns nil, the transaction is committed.
	// Supports nested transactions via savepoints.
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
