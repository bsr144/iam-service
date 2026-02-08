package testutil

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockTransactionManager implements any TransactionManager interface for testing.
// Works with any consumer-defined TransactionManager via Go's structural typing.
type MockTransactionManager struct {
	mock.Mock
}

// WithTransaction executes the function directly without a real transaction.
func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)

	// If mock returns nil, execute the function
	if args.Get(0) == nil {
		return fn(ctx)
	}
	return args.Error(0)
}

// NewMockTransactionManager creates a mock that executes functions directly.
// Use this for unit tests where you want the transaction callback to execute.
func NewMockTransactionManager() *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(nil)
	return m
}

// NewFailingTransactionManager creates a mock that simulates transaction failure.
// The transaction callback will NOT be executed, and the provided error is returned.
func NewFailingTransactionManager(err error) *MockTransactionManager {
	m := &MockTransactionManager{}
	m.On("WithTransaction", mock.Anything, mock.Anything).Return(err)
	return m
}

// PassthroughTransactionManager is a simple mock that always executes the function.
// Use when you don't need mock assertions and just want transactions to pass through.
type PassthroughTransactionManager struct{}

// WithTransaction executes fn directly without any transaction wrapping.
func (m *PassthroughTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// NewPassthroughTransactionManager creates a passthrough transaction manager.
func NewPassthroughTransactionManager() *PassthroughTransactionManager {
	return &PassthroughTransactionManager{}
}
