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

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx                context.Context
	db                 *database.DB
	configService      *service.ConfigService
	esService          *service.ElasticsearchService
	monacoCacheService *service.MonacoCacheService
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
		runtime.LogErrorf(ctx, "Failed to initialize database: %v", err)
		fmt.Printf("Failed to initialize database: %v\n", err)
		// You might want to handle this more gracefully
		os.Exit(1)
	}

	// Logger is already initialized in main.go, just log that we're ready
	runtime.LogInfo(ctx, "ElasticGaze application startup completed successfully")
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

	// Initialize Monaco cache service
	a.monacoCacheService = service.NewMonacoCacheService(elasticGazeDir)

	return nil
}

// Close closes the database connection
func (a *App) Close() error {
	runtime.LogInfo(a.ctx, "Closing application and database connection")
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
	runtime.LogInfof(a.ctx, "Creating new Elasticsearch configuration: %s", req.ConnectionName)
	config, err := a.configService.CreateConfig(req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to create configuration: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully created configuration with ID: %d", config.ID)
	return config, nil
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
	runtime.LogInfof(a.ctx, "Testing Elasticsearch connection to %s:%s", req.Host, req.Port)
	response, err := a.esService.TestConnection(req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Connection test failed: %v", err)
		return response, err
	}
	if response.Success {
		runtime.LogInfo(a.ctx, "Connection test successful")
	} else {
		runtime.LogWarningf(a.ctx, "Connection test failed: %s", response.Message)
	}
	return response, nil
}

// HasDefaultConfig checks if there is a default connection configured
func (a *App) HasDefaultConfig() (bool, error) {
	return a.configService.HasDefaultConfig()
}

// TestDefaultConnection tests if the default connection exists and is working
func (a *App) TestDefaultConnection() (*models.TestConnectionResponse, error) {
	// First check if default config exists
	defaultConfig, err := a.configService.GetDefaultConfig()
	if err != nil {
		return &models.TestConnectionResponse{
			Success:      false,
			Message:      "No default connection configured",
			ErrorDetails: "No default connection found",
			ErrorCode:    "NO_DEFAULT_CONNECTION",
		}, nil
	}

	// Convert config to test request
	testReq := &models.TestConnectionRequest{
		Host:                 defaultConfig.Host,
		Port:                 defaultConfig.Port,
		SSLOrHTTPS:           defaultConfig.SSLOrHTTPS,
		AuthenticationMethod: defaultConfig.AuthenticationMethod,
		Username:             defaultConfig.Username,
		Password:             defaultConfig.Password,
	}

	// Test the connection
	return a.esService.TestConnection(testReq)
}

// GetClusterDashboardData retrieves dashboard data for the default cluster
func (a *App) GetClusterDashboardData() (*models.ProcessedDashboardData, error) {
	// Get default config
	defaultConfig, err := a.configService.GetDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("no default connection configured: %w", err)
	}

	// Fetch cluster dashboard data
	return a.esService.GetClusterDashboardData(defaultConfig)
}

// GetClusterDashboardDataByConfig retrieves dashboard data for a specific cluster configuration
func (a *App) GetClusterDashboardDataByConfig(configID int) (*models.ProcessedDashboardData, error) {
	// Get all configs to find the one with the specified ID
	configs, err := a.configService.GetAllConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to get configurations: %w", err)
	}

	var selectedConfig *models.Config
	for _, config := range configs {
		if config.ID == configID {
			selectedConfig = config
			break
		}
	}

	if selectedConfig == nil {
		return nil, fmt.Errorf("configuration with ID %d not found", configID)
	}

	// Fetch cluster dashboard data
	return a.esService.GetClusterDashboardData(selectedConfig)
}

// GetClusterHealthForAllConfigs retrieves cluster health for all configurations
func (a *App) GetClusterHealthForAllConfigs() (map[string]string, error) {
	configs, err := a.configService.GetAllConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to get configurations: %w", err)
	}

	healthMap := make(map[string]string)

	for _, config := range configs {
		// Convert config to test request
		testReq := &models.TestConnectionRequest{
			Host:                 config.Host,
			Port:                 config.Port,
			SSLOrHTTPS:           config.SSLOrHTTPS,
			AuthenticationMethod: config.AuthenticationMethod,
			Username:             config.Username,
			Password:             config.Password,
		}

		// Get cluster health
		health, err := a.esService.GetClusterHealthByConfig(testReq)
		if err != nil {
			healthMap[config.ConnectionName] = "red" // Default to red on error
		} else {
			healthMap[config.ConnectionName] = health.Status
		}
	}

	return healthMap, nil
}

// ExecuteElasticsearchRequest executes a generic REST request to the default Elasticsearch cluster
func (a *App) ExecuteElasticsearchRequest(req *models.ElasticsearchRestRequest) (*models.ElasticsearchRestResponse, error) {
	// Get default config
	defaultConfig, err := a.configService.GetDefaultConfig()
	if err != nil {
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   500,
			ErrorDetails: "No default connection configured",
			ErrorCode:    "NO_DEFAULT_CONNECTION",
		}, nil
	}

	// Execute the request
	return a.esService.ExecuteRestRequest(defaultConfig, req)
}

// Monaco Cache API Methods

// GetMonacoCacheInfo retrieves cache information for a specific Monaco Editor version
func (a *App) GetMonacoCacheInfo(version string) (*service.CacheInfo, error) {
	runtime.LogInfof(a.ctx, "Getting Monaco cache info for version: %s", version)
	return a.monacoCacheService.GetCacheInfo(version)
}

// WriteMonacoCache stores Monaco Editor data in cache
func (a *App) WriteMonacoCache(version string, data string) error {
	runtime.LogInfof(a.ctx, "Writing Monaco cache for version: %s", version)
	return a.monacoCacheService.WriteCache(version, []byte(data))
}

// ReadMonacoCache retrieves Monaco Editor data from cache
func (a *App) ReadMonacoCache(version string) (string, error) {
	runtime.LogInfof(a.ctx, "Reading Monaco cache for version: %s", version)
	data, err := a.monacoCacheService.ReadCache(version)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// InvalidateMonacoCache removes cache for a specific version
func (a *App) InvalidateMonacoCache(version string) error {
	runtime.LogInfof(a.ctx, "Invalidating Monaco cache for version: %s", version)
	return a.monacoCacheService.InvalidateCache(version)
}

// ClearAllMonacoCache removes all Monaco Editor cache
func (a *App) ClearAllMonacoCache() error {
	runtime.LogInfo(a.ctx, "Clearing all Monaco cache")
	return a.monacoCacheService.ClearAllCache()
}

// GetMonacoCacheSize returns the total size of Monaco Editor cache
func (a *App) GetMonacoCacheSize() (int64, error) {
	return a.monacoCacheService.GetCacheSize()
}
