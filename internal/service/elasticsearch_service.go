package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"elasticgaze/internal/models"
)

// ElasticsearchService handles Elasticsearch connection testing
type ElasticsearchService struct {
	client *http.Client
}

// NewElasticsearchService creates a new Elasticsearch service
func NewElasticsearchService() *ElasticsearchService {
	// Create HTTP client with timeout and SSL configuration
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // For development - in production, you might want to make this configurable
			},
		},
	}

	return &ElasticsearchService{
		client: client,
	}
}

// TestConnection tests the connection to an Elasticsearch cluster
func (s *ElasticsearchService) TestConnection(req *models.TestConnectionRequest) (*models.TestConnectionResponse, error) {
	log.Printf("🔍 Testing Elasticsearch connection to %s:%s (SSL: %v, Auth: %s)",
		req.Host, req.Port, req.SSLOrHTTPS, req.AuthenticationMethod)

	// Validate the request
	if err := req.Validate(); err != nil {
		log.Printf("❌ Connection test validation failed: %v", err)
		return &models.TestConnectionResponse{
			Success:      false,
			Message:      "Validation failed",
			ErrorDetails: err.Error(),
			ErrorCode:    "VALIDATION_ERROR",
		}, nil
	}

	// Build the URL
	protocol := "http"
	if req.SSLOrHTTPS {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s:%s/", protocol, req.Host, req.Port)
	log.Printf("🌐 Connection URL: %s", url)

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("❌ Failed to create HTTP request: %v", err)
		return &models.TestConnectionResponse{
			Success:      false,
			Message:      "Failed to create request",
			ErrorDetails: fmt.Sprintf("HTTP request creation failed: %v", err),
			ErrorCode:    "REQUEST_CREATION_ERROR",
		}, nil
	}

	// Add authentication
	if err := s.addAuthentication(httpReq, req); err != nil {
		log.Printf("❌ Authentication setup failed: %v", err)
		return &models.TestConnectionResponse{
			Success:      false,
			Message:      "Authentication setup failed",
			ErrorDetails: fmt.Sprintf("Authentication error: %v", err),
			ErrorCode:    "AUTH_ERROR",
		}, nil
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "ElasticGaze/1.0")

	log.Printf("🔐 Authentication method: %s", req.AuthenticationMethod)
	if req.AuthenticationMethod == "basic" && req.Username != nil {
		log.Printf("👤 Username: %s", *req.Username)
	}

	// Make the request
	log.Printf("🚀 Making HTTP request...")
	start := time.Now()
	resp, err := s.client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ HTTP request failed after %v: %v", duration, err)
		errorDetails := fmt.Sprintf("Connection failed after %v\nURL: %s\nError: %v", duration, url, err)

		// Check for specific error types
		errorCode := "CONNECTION_ERROR"
		errorMessage := "Connection failed"

		if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
			errorCode = "TIMEOUT_ERROR"
			errorMessage = "Connection timeout"
		}

		return &models.TestConnectionResponse{
			Success:      false,
			Message:      errorMessage,
			ErrorDetails: errorDetails,
			ErrorCode:    errorCode,
		}, nil
	}
	defer resp.Body.Close()

	log.Printf("✅ HTTP response received after %v - Status: %d %s", duration, resp.StatusCode, resp.Status)

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response body: %v", err)
		return &models.TestConnectionResponse{
			Success:      false,
			Message:      "Failed to read response",
			ErrorDetails: fmt.Sprintf("Response reading failed: %v", err),
			ErrorCode:    "RESPONSE_READ_ERROR",
		}, nil
	}

	log.Printf("📄 Response body length: %d bytes", len(body))

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ HTTP error status %d: %s", resp.StatusCode, string(body))
		errorDetails := fmt.Sprintf("HTTP %d %s\nURL: %s\nResponse: %s",
			resp.StatusCode, resp.Status, url, string(body))

		errorCode := fmt.Sprintf("HTTP_%d", resp.StatusCode)
		errorMessage := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)

		// Handle specific HTTP status codes
		switch resp.StatusCode {
		case 401:
			errorMessage = "Authentication failed"
			errorCode = "AUTH_FAILED"
		case 403:
			errorMessage = "Access forbidden"
			errorCode = "ACCESS_FORBIDDEN"
		case 404:
			errorMessage = "Elasticsearch not found at this URL"
			errorCode = "NOT_FOUND"
		case 500:
			errorMessage = "Elasticsearch server error"
			errorCode = "SERVER_ERROR"
		}

		return &models.TestConnectionResponse{
			Success:      false,
			Message:      errorMessage,
			ErrorDetails: errorDetails,
			ErrorCode:    errorCode,
		}, nil
	}

	// Parse Elasticsearch info response
	var esInfo struct {
		Name        string `json:"name"`
		ClusterName string `json:"cluster_name"`
		ClusterUUID string `json:"cluster_uuid"`
		Version     struct {
			Number        string `json:"number"`
			BuildFlavor   string `json:"build_flavor"`
			BuildType     string `json:"build_type"`
			BuildHash     string `json:"build_hash"`
			BuildDate     string `json:"build_date"`
			BuildSnapshot bool   `json:"build_snapshot"`
		} `json:"version"`
		Tagline string `json:"tagline"`
	}

	if err := json.Unmarshal(body, &esInfo); err != nil {
		log.Printf("⚠️ Response parsing failed (connection still successful): %v", err)
		log.Printf("📄 Raw response: %s", string(body))
		// Connection succeeded but response parsing failed
		return &models.TestConnectionResponse{
			Success:   true,
			Message:   "Connection successful (unable to parse cluster info)",
			ErrorCode: "PARSE_WARNING",
		}, nil
	}

	// Success!
	log.Printf("🎉 Connection test successful!")
	log.Printf("🏷️  Cluster Name: %s", esInfo.ClusterName)
	log.Printf("🏷️  Cluster UUID: %s", esInfo.ClusterUUID)
	log.Printf("📦 Elasticsearch Version: %s (%s)", esInfo.Version.Number, esInfo.Version.BuildFlavor)
	log.Printf("🏗️  Build: %s (%s)", esInfo.Version.BuildHash[:8], esInfo.Version.BuildDate)

	return &models.TestConnectionResponse{
		Success:     true,
		Message:     "Connection successful",
		ClusterName: esInfo.ClusterName,
		Version:     esInfo.Version.Number,
	}, nil
}

// addAuthentication adds authentication headers to the HTTP request
func (s *ElasticsearchService) addAuthentication(req *http.Request, connReq *models.TestConnectionRequest) error {
	switch connReq.AuthenticationMethod {
	case "basic":
		if connReq.Username == nil || connReq.Password == nil {
			return fmt.Errorf("username and password required for basic authentication")
		}
		req.SetBasicAuth(*connReq.Username, *connReq.Password)

	case "apikey":
		if connReq.APIKey == nil {
			return fmt.Errorf("API key required for API key authentication")
		}
		req.Header.Set("Authorization", "ApiKey "+*connReq.APIKey)

	case "none":
		// No authentication needed

	default:
		return fmt.Errorf("unsupported authentication method: %s", connReq.AuthenticationMethod)
	}

	return nil
}
