package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"elasticgaze/internal/logging"

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

	// Open SQLite database with connection parameters for better concurrency
	conn, err := sql.Open("sqlite", dbPath+"?_busy_timeout=10000&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=true")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for better concurrency handling
	conn.SetMaxOpenConns(25)   // Limit concurrent connections
	conn.SetMaxIdleConns(25)   // Keep connections alive for reuse
	conn.SetConnMaxLifetime(0) // No maximum lifetime

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}

	// Configure SQLite pragmas for multi-instance safety
	if err := db.configurePragmas(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to configure database pragmas: %w", err)
	}

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
		// Execute checkpoint to flush WAL to main database file
		if _, err := db.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
			// Log warning but don't fail the close operation
			logging.Warnf("Failed to checkpoint WAL: %v", err)
		}

		err := db.conn.Close()
		db.conn = nil
		return err
	}
	return nil
}

// GetConnection returns the underlying sql.DB connection
func (db *DB) GetConnection() *sql.DB {
	return db.conn
}

// configurePragmas sets SQLite pragmas for better multi-instance handling
func (db *DB) configurePragmas() error {
	pragmas := []string{
		"PRAGMA busy_timeout = 10000", // 10 second timeout for locked database
		// Note: journal_mode=WAL is already set in connection string
		"PRAGMA synchronous = NORMAL",      // Good balance of safety and performance
		"PRAGMA cache_size = 1000",         // Cache size in pages
		"PRAGMA temp_store = memory",       // Store temporary tables and indices in memory
		"PRAGMA mmap_size = 268435456",     // 256MB memory-mapped I/O
		"PRAGMA wal_autocheckpoint = 1000", // Checkpoint WAL file after 1000 pages
	}

	for _, pragma := range pragmas {
		if _, err := db.conn.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute pragma '%s': %w", pragma, err)
		}
	}

	return nil
}

// initializeSchema creates the necessary tables if they don't exist
func (db *DB) initializeSchema() error {
	// Create tbl_config table
	configQuery := `
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

	if _, err := db.conn.Exec(configQuery); err != nil {
		return fmt.Errorf("failed to create tbl_config table: %w", err)
	}

	// Create collections table
	collectionsQuery := `
	CREATE TABLE IF NOT EXISTS tbl_collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.conn.Exec(collectionsQuery); err != nil {
		return fmt.Errorf("failed to create tbl_collections table: %w", err)
	}

	// Create folders table
	foldersQuery := `
	CREATE TABLE IF NOT EXISTS tbl_folders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		parent_folder_id INTEGER,
		collection_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (parent_folder_id) REFERENCES tbl_folders(id) ON DELETE CASCADE,
		FOREIGN KEY (collection_id) REFERENCES tbl_collections(id) ON DELETE CASCADE
	);`

	if _, err := db.conn.Exec(foldersQuery); err != nil {
		return fmt.Errorf("failed to create tbl_folders table: %w", err)
	}

	// Create requests table
	requestsQuery := `
	CREATE TABLE IF NOT EXISTS tbl_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		method VARCHAR(10) NOT NULL DEFAULT 'GET',
		url TEXT NOT NULL,
		body TEXT,
		description TEXT,
		folder_id INTEGER,
		collection_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (folder_id) REFERENCES tbl_folders(id) ON DELETE CASCADE,
		FOREIGN KEY (collection_id) REFERENCES tbl_collections(id) ON DELETE CASCADE
	);`

	if _, err := db.conn.Exec(requestsQuery); err != nil {
		return fmt.Errorf("failed to create tbl_requests table: %w", err)
	}

	// Create trigger to update updated_at field for tbl_config
	configTriggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_tbl_config_updated_at 
	AFTER UPDATE ON tbl_config
	FOR EACH ROW
	BEGIN
		UPDATE tbl_config SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.conn.Exec(configTriggerQuery); err != nil {
		return fmt.Errorf("failed to create config trigger: %w", err)
	}

	// Create trigger to update updated_at field for tbl_collections
	collectionsTriggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_tbl_collections_updated_at 
	AFTER UPDATE ON tbl_collections
	FOR EACH ROW
	BEGIN
		UPDATE tbl_collections SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.conn.Exec(collectionsTriggerQuery); err != nil {
		return fmt.Errorf("failed to create collections trigger: %w", err)
	}

	// Create trigger to update updated_at field for tbl_folders
	foldersTriggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_tbl_folders_updated_at 
	AFTER UPDATE ON tbl_folders
	FOR EACH ROW
	BEGIN
		UPDATE tbl_folders SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.conn.Exec(foldersTriggerQuery); err != nil {
		return fmt.Errorf("failed to create folders trigger: %w", err)
	}

	// Create trigger to update updated_at field for tbl_requests
	requestsTriggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_tbl_requests_updated_at 
	AFTER UPDATE ON tbl_requests
	FOR EACH ROW
	BEGIN
		UPDATE tbl_requests SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;`

	if _, err := db.conn.Exec(requestsTriggerQuery); err != nil {
		return fmt.Errorf("failed to create requests trigger: %w", err)
	}

	return nil
}

// ExecuteWithRetry executes a function with retry logic for database busy errors
func (db *DB) ExecuteWithRetry(operation func() error, maxRetries int) error {
	var lastErr error

	for i := 0; i <= maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		// Check if it's a database busy error
		if err.Error() == "database is locked" || err.Error() == "database is busy" {
			lastErr = err
			if i < maxRetries {
				// Wait with exponential backoff
				waitTime := time.Duration(50*(i+1)) * time.Millisecond
				time.Sleep(waitTime)
				continue
			}
		} else {
			// If it's not a busy/lock error, return immediately
			return err
		}
	}

	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}
