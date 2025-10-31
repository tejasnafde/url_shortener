package database

import (
	"context"
	"time"

	"url_shortener/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type sqliteDB struct {
	db *sqlx.DB
}

func NewSQLite(dsn string) (SQLExecutor, error) {
	conn, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		logger.Logger.WithError(err).Error("Failed to establish database connection")
		return nil, err
	}
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := conn.PingContext(ctx); err != nil {
		conn.Close()
		logger.Logger.WithError(err).Error("Failed to ping database connection")
		return nil, err
	}

	logger.Logger.Info("Successfully established SQLite connection")
	return &sqliteDB{db: conn}, nil
}

func (s *sqliteDB) Close() error { return s.db.Close() }

func (s *sqliteDB) ExecContext(ctx context.Context, query string, args ...any) (int64, error) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}
	logger.WithComponent("db").WithFields(map[string]any{"query": query, "args": args}).Debug("Exec query")
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

type rowsAdapter struct{ *sqlx.Rows }

func (r *rowsAdapter) Scan(dest ...any) error { return r.Rows.Scan(dest...) }

func (s *sqliteDB) QueryContext(ctx context.Context, query string, args ...any) (Rows, error) {
	logger.WithComponent("db").WithFields(map[string]any{"query": query, "args": args}).Debug("Query")
	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &rowsAdapter{Rows: rows}, nil
}

type rowAdapter struct{ *sqlx.Row }

func (r *rowAdapter) Scan(dest ...any) error { return r.Row.Scan(dest...) }

func (s *sqliteDB) QueryRowContext(ctx context.Context, query string, args ...any) RowScanner {
	logger.WithComponent("db").WithFields(map[string]any{"query": query, "args": args}).Debug("QueryRow")
	return &rowAdapter{Row: s.db.QueryRowxContext(ctx, query, args...)}
}
