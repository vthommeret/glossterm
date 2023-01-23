package base

import "context"

// Tx is a common interface implemented by all transactions.
type Tx interface {
	// Commit applies all changes made in the transaction.
	Commit(ctx context.Context) error
	// Close rolls back the transaction.
	// Committed transactions will not be affected by calling Close.
	Close() error
}

// IteratorContext is a common interface implemented by all iterators that are not bound to the context.
type IteratorContext interface {
	// Next advances an iterator.
	Next(ctx context.Context) bool
	// Err returns a last encountered error.
	Err() error
	// Close frees resources.
	Close() error
}

// Iterator is a common interface implemented by all iterators that does not require a context.
type Iterator interface {
	// Next advances an iterator.
	Next() bool
	// Err returns a last encountered error.
	Err() error
	// Close frees resources.
	Close() error
}
