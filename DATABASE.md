# ElasticGaze - Database Layer

This document describes the database layer implementation for ElasticGaze, a Wails application for managing Elasticsearch connections.

## Architecture

The application follows a clean architecture pattern with the following layers:

- **Models** (`internal/models`): Data structures and validation
- **Repository** (`internal/repository`): Data access layer
- **Service** (`internal/service`): Business logic layer
- **Database** (`internal/database`): Database connection and initialization

## Database Schema

### tbl_config Table

```sql
CREATE TABLE tbl_config (
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
);
```

## Available Functions

The following functions are available through the Wails app interface:

### Configuration Management

1. **CreateConfig(req *models.CreateConfigRequest) (*models.Config, error)**
   - Creates a new Elasticsearch connection configuration
   - Validates input data
   - Ensures only one default configuration exists

2. **GetConfigByID(id int) (*models.Config, error)**
   - Retrieves a specific configuration by its ID

3. **GetAllConfigs() ([]*models.Config, error)**
   - Retrieves all configurations ordered by creation date (newest first)

4. **GetDefaultConfig() (*models.Config, error)**
   - Retrieves the configuration marked as default

5. **UpdateConfig(id int, req *models.UpdateConfigRequest) (*models.Config, error)**
   - Updates an existing configuration
   - Supports partial updates (only provided fields are updated)

6. **DeleteConfig(id int) error**
   - Deletes a configuration by ID

## Usage Examples

### Creating a Configuration

```go
// Create a new configuration
req := &models.CreateConfigRequest{
    ConnectionName:       "Production Cluster",
    EnvIndicatorColor:    "red",
    Host:                 "elasticsearch.prod.company.com",
    Port:                 "9200",
    SSLOrHTTPS:          true,
    AuthenticationMethod: "basic",
    Username:             stringPtr("admin"),
    Password:             stringPtr("secret"),
    SetAsDefault:         true,
}

config, err := app.CreateConfig(req)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created config with ID: %d\n", config.ID)
```

### Retrieving All Configurations

```go
configs, err := app.GetAllConfigs()
if err != nil {
    log.Fatal(err)
}

for _, config := range configs {
    fmt.Printf("Config: %s (%s:%s)\n", config.ConnectionName, config.Host, config.Port)
}
```

### Updating a Configuration

```go
// Update only the host and port
updateReq := &models.UpdateConfigRequest{
    Host: stringPtr("new-elasticsearch.company.com"),
    Port: stringPtr("9201"),
}

config, err := app.UpdateConfig(configID, updateReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Updated config: %+v\n", config)
```

## Database Location

The SQLite database is stored in:
- **Windows**: `%APPDATA%/elasticgaze/elasticgaze.db`
- **macOS**: `~/Library/Application Support/elasticgaze/elasticgaze.db`
- **Linux**: `~/.config/elasticgaze/elasticgaze.db`

## Features

- **Automatic Schema Migration**: Database tables are created automatically on first run
- **Input Validation**: All requests are validated before database operations
- **Default Management**: Only one configuration can be marked as default at a time
- **Timestamps**: Automatic creation and update timestamps
- **Error Handling**: Comprehensive error messages for debugging
- **Clean Architecture**: Separation of concerns with distinct layers

## Dependencies

- **modernc.org/sqlite**: Pure Go SQLite driver
- **Standard Library**: Uses only Go standard library for database operations

## Development

To extend the database functionality:

1. Add new fields to the model in `internal/models/config.go`
2. Update the database schema in `internal/database/connection.go`
3. Add repository methods in `internal/repository/config_repository.go`
4. Add business logic in `internal/service/config_service.go`
5. Expose API methods in `app.go`

## Testing

For testing, you can create a temporary database:

```go
func createTestDB() (*database.DB, error) {
    return database.NewConnection(":memory:")
}
```

This creates an in-memory SQLite database perfect for unit tests.