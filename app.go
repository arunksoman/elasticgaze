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
	collectionsService *service.CollectionsService
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

	// Initialize collections repository and service
	collectionsRepo := repository.NewCollectionsRepository(db.GetConnection())
	a.collectionsService = service.NewCollectionsService(collectionsRepo)

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

// Collections API Methods

// CreateCollection creates a new REST request collection
func (a *App) CreateCollection(req *models.CreateCollectionRequest) (*models.Collection, error) {
	runtime.LogInfof(a.ctx, "Creating new collection: %s", req.Name)
	collection, err := a.collectionsService.CreateCollection(req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to create collection: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully created collection with ID: %d", collection.ID)
	return collection, nil
}

// GetCollectionByID retrieves a collection by ID
func (a *App) GetCollectionByID(id int) (*models.Collection, error) {
	return a.collectionsService.GetCollectionByID(id)
}

// GetAllCollections retrieves all collections
func (a *App) GetAllCollections() ([]*models.Collection, error) {
	return a.collectionsService.GetAllCollections()
}

// UpdateCollection updates an existing collection
func (a *App) UpdateCollection(id int, req *models.UpdateCollectionRequest) (*models.Collection, error) {
	runtime.LogInfof(a.ctx, "Updating collection ID: %d", id)
	collection, err := a.collectionsService.UpdateCollection(id, req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to update collection: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully updated collection ID: %d", id)
	return collection, nil
}

// DeleteCollection deletes a collection by ID
func (a *App) DeleteCollection(id int) error {
	runtime.LogInfof(a.ctx, "Deleting collection ID: %d", id)
	err := a.collectionsService.DeleteCollection(id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to delete collection: %v", err)
		return err
	}
	runtime.LogInfof(a.ctx, "Successfully deleted collection ID: %d", id)
	return nil
}

// CreateFolder creates a new folder within a collection
func (a *App) CreateFolder(req *models.CreateFolderRequest) (*models.Folder, error) {
	runtime.LogInfof(a.ctx, "Creating new folder: %s in collection ID: %d", req.Name, req.CollectionID)
	folder, err := a.collectionsService.CreateFolder(req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to create folder: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully created folder with ID: %d", folder.ID)
	return folder, nil
}

// GetFolderByID retrieves a folder by ID
func (a *App) GetFolderByID(id int) (*models.Folder, error) {
	return a.collectionsService.GetFolderByID(id)
}

// GetFoldersByCollectionID retrieves all folders for a collection
func (a *App) GetFoldersByCollectionID(collectionID int) ([]*models.Folder, error) {
	return a.collectionsService.GetFoldersByCollectionID(collectionID)
}

// UpdateFolder updates an existing folder
func (a *App) UpdateFolder(id int, req *models.UpdateFolderRequest) (*models.Folder, error) {
	runtime.LogInfof(a.ctx, "Updating folder ID: %d", id)
	folder, err := a.collectionsService.UpdateFolder(id, req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to update folder: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully updated folder ID: %d", id)
	return folder, nil
}

// DeleteFolder deletes a folder by ID
func (a *App) DeleteFolder(id int) error {
	runtime.LogInfof(a.ctx, "Deleting folder ID: %d", id)
	err := a.collectionsService.DeleteFolder(id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to delete folder: %v", err)
		return err
	}
	runtime.LogInfof(a.ctx, "Successfully deleted folder ID: %d", id)
	return nil
}

// CreateRestRequest creates a new REST request within a collection
func (a *App) CreateRestRequest(req *models.CreateRequestRequest) (*models.Request, error) {
	runtime.LogInfof(a.ctx, "Creating new request: %s %s in collection ID: %d", req.Method, req.Name, req.CollectionID)
	request, err := a.collectionsService.CreateRequest(req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to create request: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully created request with ID: %d", request.ID)
	return request, nil
}

// GetRestRequestByID retrieves a REST request by ID
func (a *App) GetRestRequestByID(id int) (*models.Request, error) {
	return a.collectionsService.GetRequestByID(id)
}

// GetRestRequestsByCollectionID retrieves all requests for a collection
func (a *App) GetRestRequestsByCollectionID(collectionID int) ([]*models.Request, error) {
	return a.collectionsService.GetRequestsByCollectionID(collectionID)
}

// GetRestRequestsByFolderID retrieves all requests for a folder
func (a *App) GetRestRequestsByFolderID(folderID int) ([]*models.Request, error) {
	return a.collectionsService.GetRequestsByFolderID(folderID)
}

// UpdateRestRequest updates an existing REST request
func (a *App) UpdateRestRequest(id int, req *models.UpdateRequestRequest) (*models.Request, error) {
	runtime.LogInfof(a.ctx, "Updating request ID: %d", id)
	request, err := a.collectionsService.UpdateRequest(id, req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to update request: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Successfully updated request ID: %d", id)
	return request, nil
}

// DeleteRestRequest deletes a REST request by ID
func (a *App) DeleteRestRequest(id int) error {
	runtime.LogInfof(a.ctx, "Deleting request ID: %d", id)
	err := a.collectionsService.DeleteRequest(id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to delete request: %v", err)
		return err
	}
	runtime.LogInfof(a.ctx, "Successfully deleted request ID: %d", id)
	return nil
}

// GetCollectionTree retrieves the tree structure for a specific collection
func (a *App) GetCollectionTree(collectionID int) (*models.CollectionTreeNode, error) {
	return a.collectionsService.GetCollectionTree(collectionID)
}

// GetAllCollectionTrees retrieves the tree structure for all collections
func (a *App) GetAllCollectionTrees() ([]*models.CollectionTreeNode, error) {
	return a.collectionsService.GetAllCollectionTrees()
}

// EnsureDefaultCollection ensures a default collection exists and returns it
func (a *App) EnsureDefaultCollection() (*models.Collection, error) {
	runtime.LogInfo(a.ctx, "Ensuring default collection exists")
	collection, err := a.collectionsService.EnsureDefaultCollection()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to ensure default collection: %v", err)
		return nil, err
	}
	runtime.LogInfof(a.ctx, "Default collection ensured with ID: %d", collection.ID)
	return collection, nil
}
