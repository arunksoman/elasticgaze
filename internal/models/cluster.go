package models

// ClusterInfo represents cluster information response
type ClusterInfo struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUUID string `json:"cluster_uuid"`
	Version     struct {
		Number                           string `json:"number"`
		BuildFlavor                      string `json:"build_flavor"`
		BuildType                        string `json:"build_type"`
		BuildHash                        string `json:"build_hash"`
		BuildDate                        string `json:"build_date"`
		LuceneVersion                    string `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

// ClusterHealth represents cluster health response
type ClusterHealth struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"` // green, yellow, red
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

// NodesInfo represents nodes information response
type NodesInfo struct {
	ClusterName string `json:"cluster_name"`
	Nodes       map[string]struct {
		Name             string                 `json:"name"`
		TransportAddress string                 `json:"transport_address"`
		Host             string                 `json:"host"`
		IP               string                 `json:"ip"`
		Version          string                 `json:"version"`
		BuildFlavor      string                 `json:"build_flavor"`
		BuildType        string                 `json:"build_type"`
		BuildHash        string                 `json:"build_hash"`
		Roles            []string               `json:"roles"`
		Attributes       map[string]string      `json:"attributes"`
		Settings         map[string]interface{} `json:"settings"`
	} `json:"nodes"`
}

// IndicesStats represents indices statistics response
type IndicesStats struct {
	Shards struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	All struct {
		Primaries struct {
			Docs struct {
				Count   int64 `json:"count"`
				Deleted int64 `json:"deleted"`
			} `json:"docs"`
			Store struct {
				SizeInBytes             int64 `json:"size_in_bytes"`
				TotalDataSetSizeInBytes int64 `json:"total_data_set_size_in_bytes,omitempty"`
			} `json:"store"`
		} `json:"primaries"`
		Total struct {
			Docs struct {
				Count   int64 `json:"count"`
				Deleted int64 `json:"deleted"`
			} `json:"docs"`
			Store struct {
				SizeInBytes             int64 `json:"size_in_bytes"`
				TotalDataSetSizeInBytes int64 `json:"total_data_set_size_in_bytes,omitempty"`
			} `json:"store"`
		} `json:"total"`
	} `json:"_all"`
}

// ClusterDashboardData represents the data for the home dashboard
type ClusterDashboardData struct {
	ClusterInfo   *ClusterInfo   `json:"cluster_info"`
	ClusterHealth *ClusterHealth `json:"cluster_health"`
	NodesInfo     *NodesInfo     `json:"nodes_info"`
	IndicesStats  *IndicesStats  `json:"indices_stats"`
}

// NodeCounts represents count of different node types
type NodeCounts struct {
	Master int `json:"master"`
	Data   int `json:"data"`
	Ingest int `json:"ingest"`
	Total  int `json:"total"`
}

// ShardCounts represents count of different shard types
type ShardCounts struct {
	Primary int `json:"primary"`
	Replica int `json:"replica"`
	Total   int `json:"total"`
}

// IndexMetrics represents index-related metrics
type IndexMetrics struct {
	DocumentCount  int64  `json:"document_count"`
	DiskUsage      string `json:"disk_usage"`
	DiskUsageBytes int64  `json:"disk_usage_bytes"`
}

// ProcessedDashboardData represents processed data for the frontend
type ProcessedDashboardData struct {
	ClusterInfo   *ClusterInfo   `json:"cluster_info"`
	ClusterHealth *ClusterHealth `json:"cluster_health"`
	NodeCounts    *NodeCounts    `json:"node_counts"`
	ShardCounts   *ShardCounts   `json:"shard_counts"`
	IndexMetrics  *IndexMetrics  `json:"index_metrics"`
}
