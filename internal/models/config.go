package models

// Config represents an Elasticsearch connection configuration
type Config struct {
	ID                   int     `json:"id" db:"id"`
	ConnectionName       string  `json:"connection_name" db:"connection_name"`
	EnvIndicatorColor    string  `json:"env_indicator_color" db:"env_indicator_color"`
	Host                 string  `json:"host" db:"host"`
	Port                 string  `json:"port" db:"port"`
	SSLOrHTTPS           bool    `json:"ssl_or_https" db:"ssl_or_https"`
	AuthenticationMethod string  `json:"authentication_method" db:"authentication_method"`
	Username             *string `json:"username,omitempty" db:"username"`
	Password             *string `json:"password,omitempty" db:"password"`
	SetAsDefault         bool    `json:"set_as_default" db:"set_as_default"`
	CreatedAt            string  `json:"created_at" db:"created_at"`
	UpdatedAt            string  `json:"updated_at" db:"updated_at"`
}

// CreateConfigRequest represents the request payload for creating a new config
type CreateConfigRequest struct {
	ConnectionName       string  `json:"connection_name" validate:"required"`
	EnvIndicatorColor    string  `json:"env_indicator_color"`
	Host                 string  `json:"host" validate:"required"`
	Port                 string  `json:"port"`
	SSLOrHTTPS           bool    `json:"ssl_or_https"`
	AuthenticationMethod string  `json:"authentication_method"`
	Username             *string `json:"username,omitempty"`
	Password             *string `json:"password,omitempty"`
	SetAsDefault         bool    `json:"set_as_default"`
}

// UpdateConfigRequest represents the request payload for updating an existing config
type UpdateConfigRequest struct {
	ConnectionName       *string `json:"connection_name,omitempty"`
	EnvIndicatorColor    *string `json:"env_indicator_color,omitempty"`
	Host                 *string `json:"host,omitempty"`
	Port                 *string `json:"port,omitempty"`
	SSLOrHTTPS           *bool   `json:"ssl_or_https,omitempty"`
	AuthenticationMethod *string `json:"authentication_method,omitempty"`
	Username             *string `json:"username,omitempty"`
	Password             *string `json:"password,omitempty"`
	SetAsDefault         *bool   `json:"set_as_default,omitempty"`
}

// Validate performs basic validation on the CreateConfigRequest
func (c *CreateConfigRequest) Validate() error {
	if c.ConnectionName == "" {
		return ErrConnectionNameRequired
	}
	if c.Host == "" {
		return ErrHostRequired
	}
	return nil
}

// Common validation errors
var (
	ErrConnectionNameRequired = &ValidationError{Field: "connection_name", Message: "connection name is required"}
	ErrHostRequired           = &ValidationError{Field: "host", Message: "host is required"}
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
