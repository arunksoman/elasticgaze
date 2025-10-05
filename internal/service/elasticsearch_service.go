package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	url := s.buildURL(req, "/")
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

// GetClusterDashboardData fetches all cluster data needed for the dashboard
func (s *ElasticsearchService) GetClusterDashboardData(config *models.Config) (*models.ProcessedDashboardData, error) {
	log.Printf("🔍 Fetching cluster dashboard data for %s", config.ConnectionName)

	// Create test connection request from config
	testReq := &models.TestConnectionRequest{
		Host:                 config.Host,
		Port:                 config.Port,
		SSLOrHTTPS:           config.SSLOrHTTPS,
		AuthenticationMethod: config.AuthenticationMethod,
		Username:             config.Username,
		Password:             config.Password,
	}

	// Get cluster info
	clusterInfo, err := s.getClusterInfo(testReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}

	// Get cluster health
	clusterHealth, err := s.getClusterHealth(testReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster health: %w", err)
	}

	// Get nodes info
	nodesInfo, err := s.getNodesInfo(testReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes info: %w", err)
	}

	// Get indices stats
	indicesStats, err := s.getIndicesStats(testReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get indices stats: %w", err)
	}

	// Process the data
	processedData := s.processClusterData(clusterInfo, clusterHealth, nodesInfo, indicesStats)

	log.Printf("✅ Successfully fetched cluster dashboard data")
	return processedData, nil
}

// getClusterInfo fetches cluster information
func (s *ElasticsearchService) getClusterInfo(connReq *models.TestConnectionRequest) (*models.ClusterInfo, error) {
	url := s.buildURL(connReq, "/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := s.addAuthentication(req, connReq); err != nil {
		return nil, fmt.Errorf("failed to add authentication: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var clusterInfo models.ClusterInfo
	if err := json.Unmarshal(body, &clusterInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &clusterInfo, nil
}

// getClusterHealth fetches cluster health
func (s *ElasticsearchService) getClusterHealth(connReq *models.TestConnectionRequest) (*models.ClusterHealth, error) {
	url := s.buildURL(connReq, "/_cluster/health")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := s.addAuthentication(req, connReq); err != nil {
		return nil, fmt.Errorf("failed to add authentication: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var clusterHealth models.ClusterHealth
	if err := json.Unmarshal(body, &clusterHealth); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &clusterHealth, nil
}

// getNodesInfo fetches nodes information
func (s *ElasticsearchService) getNodesInfo(connReq *models.TestConnectionRequest) (*models.NodesInfo, error) {
	url := s.buildURL(connReq, "/_nodes")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := s.addAuthentication(req, connReq); err != nil {
		return nil, fmt.Errorf("failed to add authentication: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var nodesInfo models.NodesInfo
	if err := json.Unmarshal(body, &nodesInfo); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &nodesInfo, nil
}

// getIndicesStats fetches indices statistics
func (s *ElasticsearchService) getIndicesStats(connReq *models.TestConnectionRequest) (*models.IndicesStats, error) {
	url := s.buildURL(connReq, "/_stats")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := s.addAuthentication(req, connReq); err != nil {
		return nil, fmt.Errorf("failed to add authentication: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var indicesStats models.IndicesStats
	if err := json.Unmarshal(body, &indicesStats); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &indicesStats, nil
}

// processClusterData processes raw cluster data into frontend-friendly format
func (s *ElasticsearchService) processClusterData(clusterInfo *models.ClusterInfo, clusterHealth *models.ClusterHealth, nodesInfo *models.NodesInfo, indicesStats *models.IndicesStats) *models.ProcessedDashboardData {
	// Process node counts
	nodeCounts := &models.NodeCounts{
		Total: len(nodesInfo.Nodes),
	}

	for _, node := range nodesInfo.Nodes {
		for _, role := range node.Roles {
			switch role {
			case "master":
				nodeCounts.Master++
			case "data", "data_content", "data_hot", "data_warm", "data_cold", "data_frozen":
				nodeCounts.Data++
			case "ingest":
				nodeCounts.Ingest++
			}
		}
	}

	// Process shard counts
	shardCounts := &models.ShardCounts{
		Primary: clusterHealth.ActivePrimaryShards,
		Total:   clusterHealth.ActiveShards,
	}
	shardCounts.Replica = shardCounts.Total - shardCounts.Primary

	// Process index metrics
	indexMetrics := &models.IndexMetrics{
		DocumentCount:  indicesStats.All.Total.Docs.Count,
		DiskUsageBytes: indicesStats.All.Total.Store.SizeInBytes,
		DiskUsage:      formatBytes(indicesStats.All.Total.Store.SizeInBytes),
	}

	return &models.ProcessedDashboardData{
		ClusterInfo:   clusterInfo,
		ClusterHealth: clusterHealth,
		NodeCounts:    nodeCounts,
		ShardCounts:   shardCounts,
		IndexMetrics:  indexMetrics,
	}
}

// GetClusterHealthByConfig fetches cluster health for a specific config
func (s *ElasticsearchService) GetClusterHealthByConfig(connReq *models.TestConnectionRequest) (*models.ClusterHealth, error) {
	return s.getClusterHealth(connReq)
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ExecuteRestRequest executes a generic REST request to the default Elasticsearch cluster
func (s *ElasticsearchService) ExecuteRestRequest(config *models.Config, req *models.ElasticsearchRestRequest) (*models.ElasticsearchRestResponse, error) {
	log.Printf("🔍 Executing ES REST request: %s %s", req.Method, req.Endpoint)

	// Validate the request
	if err := req.Validate(); err != nil {
		log.Printf("❌ REST request validation failed: %v", err)
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   400,
			ErrorDetails: err.Error(),
			ErrorCode:    "VALIDATION_ERROR",
		}, nil
	}

	// Convert config to connection request for URL building
	connReq := &models.TestConnectionRequest{
		Host:                 config.Host,
		Port:                 config.Port,
		SSLOrHTTPS:           config.SSLOrHTTPS,
		AuthenticationMethod: config.AuthenticationMethod,
		Username:             config.Username,
		Password:             config.Password,
	}

	// Clean and normalize the endpoint
	endpoint := strings.TrimSpace(req.Endpoint)
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}

	// Build the URL
	url := s.buildURL(connReq, endpoint)
	log.Printf("🌐 Request URL: %s", url)

	// Prepare request body
	var body io.Reader
	if req.Body != nil && strings.TrimSpace(*req.Body) != "" {
		body = bytes.NewBufferString(*req.Body)
		log.Printf("📄 Request body: %s", *req.Body)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest(strings.ToUpper(req.Method), url, body)
	if err != nil {
		log.Printf("❌ Failed to create HTTP request: %v", err)
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   500,
			ErrorDetails: fmt.Sprintf("HTTP request creation failed: %v", err),
			ErrorCode:    "REQUEST_CREATION_ERROR",
		}, nil
	}

	// Add authentication
	if err := s.addAuthentication(httpReq, connReq); err != nil {
		log.Printf("❌ Authentication setup failed: %v", err)
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   401,
			ErrorDetails: fmt.Sprintf("Authentication error: %v", err),
			ErrorCode:    "AUTH_ERROR",
		}, nil
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "ElasticGaze/1.0")

	// Make the request
	log.Printf("🚀 Making HTTP request...")
	start := time.Now()
	resp, err := s.client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ HTTP request failed after %v: %v", duration, err)
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   500,
			ErrorDetails: fmt.Sprintf("Connection failed after %v: %v", duration, err),
			ErrorCode:    "CONNECTION_ERROR",
		}, nil
	}
	defer resp.Body.Close()

	log.Printf("📊 Response status: %d, Duration: %v", resp.StatusCode, duration)

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response body: %v", err)
		return &models.ElasticsearchRestResponse{
			Success:      false,
			StatusCode:   resp.StatusCode,
			ErrorDetails: fmt.Sprintf("Failed to read response: %v", err),
			ErrorCode:    "RESPONSE_READ_ERROR",
		}, nil
	}

	// Check if the response is successful (2xx status codes)
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !success {
		log.Printf("⚠️ Elasticsearch returned error status %d", resp.StatusCode)
	} else {
		log.Printf("✅ Request completed successfully")
	}

	return &models.ElasticsearchRestResponse{
		Success:    success,
		StatusCode: resp.StatusCode,
		Response:   string(responseBody),
	}, nil
}

// buildURL constructs the full URL for Elasticsearch API calls
func (s *ElasticsearchService) buildURL(connReq *models.TestConnectionRequest, endpoint string) string {
	scheme := "http"
	if connReq.SSLOrHTTPS {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s:%s%s", scheme, connReq.Host, connReq.Port, endpoint)
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
