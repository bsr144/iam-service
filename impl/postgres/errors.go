package postgres

import "errors"

// Sentinel errors for expected database conditions.
// These errors are used by repositories to indicate specific database states
// that the usecase layer can handle appropriately.
var (
	// ErrRecordNotFound indicates the requested record does not exist.
	ErrRecordNotFound = errors.New("record not found")

	// ErrDuplicateEntry indicates a unique constraint violation.
	ErrDuplicateEntry = errors.New("duplicate entry")

	// ErrForeignKeyFailed indicates a foreign key constraint violation.
	ErrForeignKeyFailed = errors.New("foreign key constraint failed")

	// ErrDeadlock indicates a database deadlock was detected.
	ErrDeadlock = errors.New("database deadlock")
)
