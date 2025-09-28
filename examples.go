package main

import (
	"fmt"
	"log"

	"elasticgaze/internal/models"
)

// Example usage of the configuration management functions
// This file demonstrates how to use the database functionality
// Note: This is for demonstration purposes - remove before production

func exampleUsage(app *App) {
	fmt.Println("=== ElasticGaze Configuration Management Examples ===")

	// Example 1: Create a new configuration
	fmt.Println("\n1. Creating a new configuration...")
	createReq := &models.CreateConfigRequest{
		ConnectionName:       "Production Cluster",
		EnvIndicatorColor:    "red",
		Host:                 "elasticsearch.prod.company.com",
		Port:                 "9200",
		SSLOrHTTPS:           true,
		AuthenticationMethod: "basic",
		Username:             models.StringPtr("admin"),
		Password:             models.StringPtr("secretpassword"),
		SetAsDefault:         true,
	}

	config, err := app.CreateConfig(createReq)
	if err != nil {
		log.Printf("Error creating config: %v", err)
		return
	}
	fmt.Printf("Created config with ID: %d, Name: %s\n", config.ID, config.ConnectionName)

	// Example 2: Create another configuration
	fmt.Println("\n2. Creating a development configuration...")
	devReq := &models.CreateConfigRequest{
		ConnectionName:       "Development Cluster",
		EnvIndicatorColor:    "green",
		Host:                 "localhost",
		Port:                 "9200",
		SSLOrHTTPS:           false,
		AuthenticationMethod: "none",
		SetAsDefault:         false,
	}

	devConfig, err := app.CreateConfig(devReq)
	if err != nil {
		log.Printf("Error creating dev config: %v", err)
		return
	}
	fmt.Printf("Created dev config with ID: %d, Name: %s\n", devConfig.ID, devConfig.ConnectionName)

	// Example 3: Get all configurations
	fmt.Println("\n3. Retrieving all configurations...")
	configs, err := app.GetAllConfigs()
	if err != nil {
		log.Printf("Error getting all configs: %v", err)
		return
	}

	for _, cfg := range configs {
		defaultStr := ""
		if cfg.SetAsDefault {
			defaultStr = " (DEFAULT)"
		}
		fmt.Printf("- ID: %d, Name: %s, Host: %s:%s%s\n",
			cfg.ID, cfg.ConnectionName, cfg.Host, cfg.Port, defaultStr)
	}

	// Example 4: Get configuration by ID
	fmt.Println("\n4. Getting configuration by ID...")
	retrievedConfig, err := app.GetConfigByID(config.ID)
	if err != nil {
		log.Printf("Error getting config by ID: %v", err)
		return
	}
	fmt.Printf("Retrieved config: %s at %s:%s\n",
		retrievedConfig.ConnectionName, retrievedConfig.Host, retrievedConfig.Port)

	// Example 5: Update a configuration
	fmt.Println("\n5. Updating configuration...")
	updateReq := &models.UpdateConfigRequest{
		Host:              models.StringPtr("new-elasticsearch.company.com"),
		Port:              models.StringPtr("9201"),
		EnvIndicatorColor: models.StringPtr("orange"),
	}

	updatedConfig, err := app.UpdateConfig(devConfig.ID, updateReq)
	if err != nil {
		log.Printf("Error updating config: %v", err)
		return
	}
	fmt.Printf("Updated config: %s now points to %s:%s\n",
		updatedConfig.ConnectionName, updatedConfig.Host, updatedConfig.Port)

	// Example 6: Get default configuration
	fmt.Println("\n6. Getting default configuration...")
	defaultConfig, err := app.GetDefaultConfig()
	if err != nil {
		log.Printf("Error getting default config: %v", err)
		return
	}
	fmt.Printf("Default config: %s (%s:%s)\n",
		defaultConfig.ConnectionName, defaultConfig.Host, defaultConfig.Port)

	// Example 7: Delete a configuration
	fmt.Println("\n7. Deleting a configuration...")
	err = app.DeleteConfig(devConfig.ID)
	if err != nil {
		log.Printf("Error deleting config: %v", err)
		return
	}
	fmt.Printf("Deleted config with ID: %d\n", devConfig.ID)

	// Verify deletion
	fmt.Println("\n8. Verifying deletion - getting all configs again...")
	finalConfigs, err := app.GetAllConfigs()
	if err != nil {
		log.Printf("Error getting final configs: %v", err)
		return
	}

	fmt.Printf("Remaining configurations: %d\n", len(finalConfigs))
	for _, cfg := range finalConfigs {
		defaultStr := ""
		if cfg.SetAsDefault {
			defaultStr = " (DEFAULT)"
		}
		fmt.Printf("- ID: %d, Name: %s, Host: %s:%s%s\n",
			cfg.ID, cfg.ConnectionName, cfg.Host, cfg.Port, defaultStr)
	}
}

// Uncomment the following code in main.go to test the database functionality:
/*
func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Initialize context (normally done by Wails)
	ctx := context.Background()
	app.startup(ctx)
	defer app.Close()

	// Run examples
	exampleUsage(app)

	fmt.Println("\n=== Database functionality test complete ===")
}
*/
