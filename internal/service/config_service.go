package service

import (
	"fmt"

	"elasticgaze/internal/models"
	"elasticgaze/internal/repository"
)

// ConfigService handles business logic for configuration operations
type ConfigService struct {
	repo *repository.ConfigRepository
}

// NewConfigService creates a new config service
func NewConfigService(repo *repository.ConfigRepository) *ConfigService {
	return &ConfigService{repo: repo}
}

// CreateConfig creates a new configuration
func (s *ConfigService) CreateConfig(req *models.CreateConfigRequest) (*models.Config, error) {
	// Validate the request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// Check if trying to create a default connection when one already exists
	if req.SetAsDefault {
		hasDefault, err := s.repo.HasDefaultConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to check for existing default: %w", err)
		}
		if hasDefault {
			return nil, models.ErrMultipleDefaultsNotAllowed
		}
	}

	// Create the configuration
	config, err := s.repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return config, nil
}

// GetConfigByID retrieves a configuration by ID
func (s *ConfigService) GetConfigByID(id int) (*models.Config, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID: must be greater than 0")
	}

	config, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}

// GetAllConfigs retrieves all configurations
func (s *ConfigService) GetAllConfigs() ([]*models.Config, error) {
	configs, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all configs: %w", err)
	}

	return configs, nil
}

// GetDefaultConfig retrieves the default configuration
func (s *ConfigService) GetDefaultConfig() (*models.Config, error) {
	config, err := s.repo.GetDefault()
	if err != nil {
		return nil, fmt.Errorf("failed to get default config: %w", err)
	}

	return config, nil
}

// UpdateConfig updates an existing configuration
func (s *ConfigService) UpdateConfig(id int, req *models.UpdateConfigRequest) (*models.Config, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid ID: must be greater than 0")
	}

	// Check if config exists
	existingConfig, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("config not found: %w", err)
	}

	// Check if trying to set as default when another default already exists
	if req.SetAsDefault != nil && *req.SetAsDefault {
		// Only validate if this config is not already the default
		if !existingConfig.SetAsDefault {
			hasDefault, err := s.repo.HasDefaultConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to check for existing default: %w", err)
			}
			if hasDefault {
				return nil, models.ErrMultipleDefaultsNotAllowed
			}
		}
	}

	config, err := s.repo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	return config, nil
}

// DeleteConfig deletes a configuration by ID
func (s *ConfigService) DeleteConfig(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid ID: must be greater than 0")
	}

	// Check if this is the default config and if there are other configs
	config, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("config not found: %w", err)
	}

	// If this is the default config, we might want to handle this specially
	if config.SetAsDefault {
		// You could either prevent deletion or automatically set another config as default
		// For now, we'll allow deletion but you might want to add business logic here
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	return nil
}
