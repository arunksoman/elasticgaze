package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

// NewConnection creates a new database connection
func NewConnection(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open SQLite database
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Initialize the database schema
	if err := db.initializeSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// GetConnection returns the underlying sql.DB connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// initializeSchema creates the necessary tables if they don't exist
func (db *DB) initializeSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS tbl_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		connection_name VARCHAR(255) NOT NULL UNIQUE,
		env_indicator_color VARCHAR(30) NOT NULL DEFAULT 'blue',
		host VARCHAR(255) NOT NULL,
		port VARCHAR(8) NOT NULL DEFAULT '9200',
		ssl_or_https BOOLEAN NOT NULL DEFAULT 0,
		authentication_method VARCHAR(255) NOT NULL DEFAULT 'none',
		username VARCHAR(255),
		password VARCHAR(255),
		set_as_default BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.conn.Exec(query); err != nil {
		return fmt.Errorf("failed to create tbl_config table: %w", err)
	}

	// Create trigger to update updated_at field
	triggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_tbl_config_updated_at 
	AFTER UPDATE ON tbl_config
	FOR EACH ROW
	BEGIN
		UPDATE tbl_config SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.conn.Exec(triggerQuery); err != nil {
		return fmt.Errorf("failed to create trigger: %w", err)
	}

	return nil
}
