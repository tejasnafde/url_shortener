package database

import (
	// "database/sql"
	"context"
	"time"
	"github.com/jmoiron/sqlx" 
	_ "github.com/mattn/go-sqlite3"
	"url_shortener/pkg/logger"
)

// TODO: Implement generic database utilities (equivalent to your Python DBUtil class)
//
// What to implement:
// 1. Database connection management
//    - Connect to SQLite database
//    - Connection pooling and configuration
//    - Graceful connection closing
//
// 2. Generic query execution methods
//    - ExecuteQuery(query string, args ...interface{}) ([]map[string]interface{}, error)
//    - ExecuteQueryWithoutOutput(query string, args ...interface{}) (map[string]interface{}, error)
//
//
// Enterprise patterns to follow:
// - Use dependency injection for database connection
// - Return interfaces, not concrete types (for testability)
// - Proper error handling and logging
// - Connection health checks
// - Query timeout configuration

// DB represents a database connection
type DB struct {
	conn *sqlx.DB
}


func GetConnection(dsn string) (*DB, error) {
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
	db := &DB{conn: conn}
	logger.Logger.Info("Successfully Established database connection, returning conn obj")
	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) ExecuteQuery(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx , cancel = context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	}
	// if ctx.Err() == context.DeadlineExceeded {
	// 	logger.Logger.Warn("Query timed out")
	// 	return nil, ctx.Err()
	// }
	logger.WithComponent("db").
	WithFields(map[string]any{
		"query": query,
		"args":  args,
	}).Debug("Executing Query")
	rows, err := db.conn.QueryxContext(ctx, query, args...)
	if err != nil {
		if err.Is(err, context.DeadlineExceeded) || err.Is(err, context.Canceled) {
			logger.WithComponent("db").WithError(err).
				WithField("query", query).
				Warn("Query canceled or timed out")
		} else {
			logger.WithComponent("db").WithError(err).
				WithField("query", query).
				Error("Query failed")
		}
		return nil, err
	}
	defer rows.Close()
	var out []map[string]any
	for rows.Next() {
		m := map[string]any{}
		if err := rows.MapScan(m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	logger.WithComponent("db").
	WithFields(map[string]any{
		"rows":  len(out),
		"query": query,
	}).
	Debug("Query succeeded")
	return out, rows.Err()
}