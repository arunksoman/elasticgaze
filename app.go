package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"elasticgaze/internal/database"
	"elasticgaze/internal/models"
	"elasticgaze/internal/repository"
	"elasticgaze/internal/service"
)

// App struct
type App struct {
	ctx           context.Context
	db            *database.DB
	configService *service.ConfigService
	esService     *service.ElasticsearchService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := a.initDatabase(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		// You might want to handle this more gracefully
		os.Exit(1)
	}
}

// initDatabase initializes the SQLite database and services
func (a *App) initDatabase() error {
	// Get the application data directory
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %w", err)
	}

	// Create elasticgaze directory if it doesn't exist
	elasticGazeDir := filepath.Join(appDataDir, "elasticgaze")
	if err := os.MkdirAll(elasticGazeDir, 0755); err != nil {
		return fmt.Errorf("failed to create elasticgaze directory: %w", err)
	}

	// Database path
	dbPath := filepath.Join(elasticGazeDir, "elasticgaze.db")

	// Initialize database connection
	db, err := database.NewConnection(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	a.db = db

	// Initialize repository and service layers
	configRepo := repository.NewConfigRepository(db.GetConnection())
	a.configService = service.NewConfigService(configRepo)
	a.esService = service.NewElasticsearchService()

	return nil
}

// Close closes the database connection
func (a *App) Close() error {
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// Configuration Management API Methods

// CreateConfig creates a new Elasticsearch connection configuration
func (a *App) CreateConfig(req *models.CreateConfigRequest) (*models.Config, error) {
	return a.configService.CreateConfig(req)
}

// GetConfigByID retrieves a configuration by ID
func (a *App) GetConfigByID(id int) (*models.Config, error) {
	return a.configService.GetConfigByID(id)
}

// GetAllConfigs retrieves all configurations
func (a *App) GetAllConfigs() ([]*models.Config, error) {
	return a.configService.GetAllConfigs()
}

// GetDefaultConfig retrieves the default configuration
func (a *App) GetDefaultConfig() (*models.Config, error) {
	return a.configService.GetDefaultConfig()
}

// UpdateConfig updates an existing configuration
func (a *App) UpdateConfig(id int, req *models.UpdateConfigRequest) (*models.Config, error) {
	return a.configService.UpdateConfig(id, req)
}

// DeleteConfig deletes a configuration by ID
func (a *App) DeleteConfig(id int) error {
	return a.configService.DeleteConfig(id)
}

// TestConnection tests an Elasticsearch connection
func (a *App) TestConnection(req *models.TestConnectionRequest) (*models.TestConnectionResponse, error) {
	return a.esService.TestConnection(req)
}
