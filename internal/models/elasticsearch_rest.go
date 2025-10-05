package models

// ElasticsearchRestRequest represents a generic REST request to Elasticsearch
type ElasticsearchRestRequest struct {
	Method   string  `json:"method" validate:"required"`   // HTTP method (GET, POST, PUT, DELETE, etc.)
	Endpoint string  `json:"endpoint" validate:"required"` // Elasticsearch API endpoint (e.g., "_search", "_cat/indices")
	Body     *string `json:"body,omitempty"`               // Request body (JSON string, optional)
}

// ElasticsearchRestResponse represents the response from an Elasticsearch REST request
type ElasticsearchRestResponse struct {
	Success      bool   `json:"success"`
	StatusCode   int    `json:"status_code"`
	Response     string `json:"response"` // The actual Elasticsearch response as JSON string
	ErrorDetails string `json:"error_details,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
}

// Validate performs basic validation on the ElasticsearchRestRequest
func (e *ElasticsearchRestRequest) Validate() error {
	if e.Method == "" {
		return ErrMethodRequired
	}
	if e.Endpoint == "" {
		return ErrEndpointRequired
	}
	return nil
}
