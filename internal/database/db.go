package database

import "context"

type RowScanner interface {
	Scan(dest ...any) error
}

type Rows interface {
	Next() bool
	Close() error
	Err() error
	Scan(dest ...any) error
}

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (rowsAffected int64, err error)
	QueryContext(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) RowScanner
	Close() error
}
