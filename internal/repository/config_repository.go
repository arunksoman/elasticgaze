package repository

import (
	"database/sql"
	"fmt"

	"elasticgaze/internal/models"
)

// ConfigRepository handles database operations for configuration
type ConfigRepository struct {
	db *sql.DB
}

// NewConfigRepository creates a new config repository
func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// Create creates a new configuration entry
func (r *ConfigRepository) Create(req *models.CreateConfigRequest) (*models.Config, error) {
	// If this config is being set as default, unset all other defaults first
	if req.SetAsDefault {
		if err := r.unsetAllDefaults(); err != nil {
			return nil, fmt.Errorf("failed to unset other defaults: %w", err)
		}
	}

	// Set defaults
	if req.EnvIndicatorColor == "" {
		req.EnvIndicatorColor = "blue"
	}
	if req.Port == "" {
		req.Port = "9200"
	}
	if req.AuthenticationMethod == "" {
		req.AuthenticationMethod = "none"
	}

	query := `
		INSERT INTO tbl_config (
			connection_name, env_indicator_color, host, port, ssl_or_https,
			authentication_method, username, password, set_as_default
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
	`

	var config models.Config
	err := r.db.QueryRow(query,
		req.ConnectionName,
		req.EnvIndicatorColor,
		req.Host,
		req.Port,
		req.SSLOrHTTPS,
		req.AuthenticationMethod,
		req.Username,
		req.Password,
		req.SetAsDefault,
	).Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	// Fill the config struct
	config.ConnectionName = req.ConnectionName
	config.EnvIndicatorColor = req.EnvIndicatorColor
	config.Host = req.Host
	config.Port = req.Port
	config.SSLOrHTTPS = req.SSLOrHTTPS
	config.AuthenticationMethod = req.AuthenticationMethod
	config.Username = req.Username
	config.Password = req.Password
	config.SetAsDefault = req.SetAsDefault

	return &config, nil
}

// GetByID retrieves a configuration by ID
func (r *ConfigRepository) GetByID(id int) (*models.Config, error) {
	query := `
		SELECT id, connection_name, env_indicator_color, host, port, ssl_or_https,
		       authentication_method, username, password, set_as_default, created_at, updated_at
		FROM tbl_config
		WHERE id = ?
	`

	var config models.Config
	err := r.db.QueryRow(query, id).Scan(
		&config.ID,
		&config.ConnectionName,
		&config.EnvIndicatorColor,
		&config.Host,
		&config.Port,
		&config.SSLOrHTTPS,
		&config.AuthenticationMethod,
		&config.Username,
		&config.Password,
		&config.SetAsDefault,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("config with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get config by ID: %w", err)
	}

	return &config, nil
}

// GetAll retrieves all configurations
func (r *ConfigRepository) GetAll() ([]*models.Config, error) {
	query := `
		SELECT id, connection_name, env_indicator_color, host, port, ssl_or_https,
		       authentication_method, username, password, set_as_default, created_at, updated_at
		FROM tbl_config
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.Config
	for rows.Next() {
		var config models.Config
		err := rows.Scan(
			&config.ID,
			&config.ConnectionName,
			&config.EnvIndicatorColor,
			&config.Host,
			&config.Port,
			&config.SSLOrHTTPS,
			&config.AuthenticationMethod,
			&config.Username,
			&config.Password,
			&config.SetAsDefault,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan config row: %w", err)
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating config rows: %w", err)
	}

	return configs, nil
}

// GetDefault retrieves the default configuration
func (r *ConfigRepository) GetDefault() (*models.Config, error) {
	query := `
		SELECT id, connection_name, env_indicator_color, host, port, ssl_or_https,
		       authentication_method, username, password, set_as_default, created_at, updated_at
		FROM tbl_config
		WHERE set_as_default = 1
		LIMIT 1
	`

	var config models.Config
	err := r.db.QueryRow(query).Scan(
		&config.ID,
		&config.ConnectionName,
		&config.EnvIndicatorColor,
		&config.Host,
		&config.Port,
		&config.SSLOrHTTPS,
		&config.AuthenticationMethod,
		&config.Username,
		&config.Password,
		&config.SetAsDefault,
		&config.CreatedAt,
		&config.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no default configuration found")
		}
		return nil, fmt.Errorf("failed to get default config: %w", err)
	}

	return &config, nil
}

// Update updates an existing configuration
func (r *ConfigRepository) Update(id int, req *models.UpdateConfigRequest) (*models.Config, error) {
	// If this config is being set as default, unset all other defaults first
	if req.SetAsDefault != nil && *req.SetAsDefault {
		if err := r.unsetAllDefaults(); err != nil {
			return nil, fmt.Errorf("failed to unset other defaults: %w", err)
		}
	}

	// Build dynamic query based on provided fields
	query := "UPDATE tbl_config SET "
	args := []interface{}{}
	setParts := []string{}

	if req.ConnectionName != nil {
		setParts = append(setParts, "connection_name = ?")
		args = append(args, *req.ConnectionName)
	}
	if req.EnvIndicatorColor != nil {
		setParts = append(setParts, "env_indicator_color = ?")
		args = append(args, *req.EnvIndicatorColor)
	}
	if req.Host != nil {
		setParts = append(setParts, "host = ?")
		args = append(args, *req.Host)
	}
	if req.Port != nil {
		setParts = append(setParts, "port = ?")
		args = append(args, *req.Port)
	}
	if req.SSLOrHTTPS != nil {
		setParts = append(setParts, "ssl_or_https = ?")
		args = append(args, *req.SSLOrHTTPS)
	}
	if req.AuthenticationMethod != nil {
		setParts = append(setParts, "authentication_method = ?")
		args = append(args, *req.AuthenticationMethod)
	}
	if req.Username != nil {
		setParts = append(setParts, "username = ?")
		args = append(args, *req.Username)
	}
	if req.Password != nil {
		setParts = append(setParts, "password = ?")
		args = append(args, *req.Password)
	}
	if req.SetAsDefault != nil {
		setParts = append(setParts, "set_as_default = ?")
		args = append(args, *req.SetAsDefault)
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields provided for update")
	}

	// Join all SET parts with commas
	query += setParts[0]
	for i := 1; i < len(setParts); i++ {
		query += ", " + setParts[i]
	}
	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	// Return the updated config
	return r.GetByID(id)
}

// Delete deletes a configuration by ID
func (r *ConfigRepository) Delete(id int) error {
	query := "DELETE FROM tbl_config WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("config with ID %d not found", id)
	}

	return nil
}

// HasDefaultConfig checks if there is already a default configuration
func (r *ConfigRepository) HasDefaultConfig() (bool, error) {
	query := "SELECT COUNT(*) FROM tbl_config WHERE set_as_default = 1"
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check for default config: %w", err)
	}
	return count > 0, nil
}

// unsetAllDefaults removes default flag from all configurations
func (r *ConfigRepository) unsetAllDefaults() error {
	query := "UPDATE tbl_config SET set_as_default = 0 WHERE set_as_default = 1"
	_, err := r.db.Exec(query)
	return err
}
