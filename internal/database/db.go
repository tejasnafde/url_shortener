package database

import (
	// "database/sql"
	"context"
	"time"
	// "github.com/jmoiron/sqlx" temporarily commented
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