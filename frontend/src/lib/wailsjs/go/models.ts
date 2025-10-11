export namespace models {
	
	export class ClusterHealth {
	    cluster_name: string;
	    status: string;
	    timed_out: boolean;
	    number_of_nodes: number;
	    number_of_data_nodes: number;
	    active_primary_shards: number;
	    active_shards: number;
	    relocating_shards: number;
	    initializing_shards: number;
	    unassigned_shards: number;
	    delayed_unassigned_shards: number;
	    number_of_pending_tasks: number;
	    number_of_in_flight_fetch: number;
	    task_max_waiting_in_queue_millis: number;
	    active_shards_percent_as_number: number;
	
	    static createFrom(source: any = {}) {
	        return new ClusterHealth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cluster_name = source["cluster_name"];
	        this.status = source["status"];
	        this.timed_out = source["timed_out"];
	        this.number_of_nodes = source["number_of_nodes"];
	        this.number_of_data_nodes = source["number_of_data_nodes"];
	        this.active_primary_shards = source["active_primary_shards"];
	        this.active_shards = source["active_shards"];
	        this.relocating_shards = source["relocating_shards"];
	        this.initializing_shards = source["initializing_shards"];
	        this.unassigned_shards = source["unassigned_shards"];
	        this.delayed_unassigned_shards = source["delayed_unassigned_shards"];
	        this.number_of_pending_tasks = source["number_of_pending_tasks"];
	        this.number_of_in_flight_fetch = source["number_of_in_flight_fetch"];
	        this.task_max_waiting_in_queue_millis = source["task_max_waiting_in_queue_millis"];
	        this.active_shards_percent_as_number = source["active_shards_percent_as_number"];
	    }
	}
	export class ClusterInfo {
	    name: string;
	    cluster_name: string;
	    cluster_uuid: string;
	    // Go type: struct { Number string "json:\"number\""; BuildFlavor string "json:\"build_flavor\""; BuildType string "json:\"build_type\""; BuildHash string "json:\"build_hash\""; BuildDate string "json:\"build_date\""; LuceneVersion string "json:\"lucene_version\""; MinimumWireCompatibilityVersion string "json:\"minimum_wire_compatibility_version\""; MinimumIndexCompatibilityVersion string "json:\"minimum_index_compatibility_version\"" }
	    version: any;
	    tagline: string;
	
	    static createFrom(source: any = {}) {
	        return new ClusterInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cluster_name = source["cluster_name"];
	        this.cluster_uuid = source["cluster_uuid"];
	        this.version = this.convertValues(source["version"], Object);
	        this.tagline = source["tagline"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Config {
	    id: number;
	    connection_name: string;
	    env_indicator_color: string;
	    host: string;
	    port: string;
	    ssl_or_https: boolean;
	    authentication_method: string;
	    username?: string;
	    password?: string;
	    set_as_default: boolean;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.connection_name = source["connection_name"];
	        this.env_indicator_color = source["env_indicator_color"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.ssl_or_https = source["ssl_or_https"];
	        this.authentication_method = source["authentication_method"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.set_as_default = source["set_as_default"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class CreateConfigRequest {
	    connection_name: string;
	    env_indicator_color: string;
	    host: string;
	    port: string;
	    ssl_or_https: boolean;
	    authentication_method: string;
	    username?: string;
	    password?: string;
	    set_as_default: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CreateConfigRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_name = source["connection_name"];
	        this.env_indicator_color = source["env_indicator_color"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.ssl_or_https = source["ssl_or_https"];
	        this.authentication_method = source["authentication_method"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.set_as_default = source["set_as_default"];
	    }
	}
	export class ElasticsearchRestRequest {
	    method: string;
	    endpoint: string;
	    body?: string;
	
	    static createFrom(source: any = {}) {
	        return new ElasticsearchRestRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.endpoint = source["endpoint"];
	        this.body = source["body"];
	    }
	}
	export class ElasticsearchRestResponse {
	    success: boolean;
	    status_code: number;
	    response: string;
	    error_details?: string;
	    error_code?: string;
	
	    static createFrom(source: any = {}) {
	        return new ElasticsearchRestResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.status_code = source["status_code"];
	        this.response = source["response"];
	        this.error_details = source["error_details"];
	        this.error_code = source["error_code"];
	    }
	}
	export class IndexMetrics {
	    document_count: number;
	    disk_usage: string;
	    disk_usage_bytes: number;
	
	    static createFrom(source: any = {}) {
	        return new IndexMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.document_count = source["document_count"];
	        this.disk_usage = source["disk_usage"];
	        this.disk_usage_bytes = source["disk_usage_bytes"];
	    }
	}
	export class NodeCounts {
	    master: number;
	    data: number;
	    ingest: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new NodeCounts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.master = source["master"];
	        this.data = source["data"];
	        this.ingest = source["ingest"];
	        this.total = source["total"];
	    }
	}
	export class ShardCounts {
	    primary: number;
	    replica: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new ShardCounts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.primary = source["primary"];
	        this.replica = source["replica"];
	        this.total = source["total"];
	    }
	}
	export class ProcessedDashboardData {
	    cluster_info?: ClusterInfo;
	    cluster_health?: ClusterHealth;
	    node_counts?: NodeCounts;
	    shard_counts?: ShardCounts;
	    index_metrics?: IndexMetrics;
	
	    static createFrom(source: any = {}) {
	        return new ProcessedDashboardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cluster_info = this.convertValues(source["cluster_info"], ClusterInfo);
	        this.cluster_health = this.convertValues(source["cluster_health"], ClusterHealth);
	        this.node_counts = this.convertValues(source["node_counts"], NodeCounts);
	        this.shard_counts = this.convertValues(source["shard_counts"], ShardCounts);
	        this.index_metrics = this.convertValues(source["index_metrics"], IndexMetrics);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class TestConnectionRequest {
	    host: string;
	    port: string;
	    ssl_or_https: boolean;
	    authentication_method: string;
	    username?: string;
	    password?: string;
	    api_key?: string;
	
	    static createFrom(source: any = {}) {
	        return new TestConnectionRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.ssl_or_https = source["ssl_or_https"];
	        this.authentication_method = source["authentication_method"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.api_key = source["api_key"];
	    }
	}
	export class TestConnectionResponse {
	    success: boolean;
	    message: string;
	    cluster_name?: string;
	    version?: string;
	    error_details?: string;
	    error_code?: string;
	
	    static createFrom(source: any = {}) {
	        return new TestConnectionResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.cluster_name = source["cluster_name"];
	        this.version = source["version"];
	        this.error_details = source["error_details"];
	        this.error_code = source["error_code"];
	    }
	}
	export class UpdateConfigRequest {
	    connection_name?: string;
	    env_indicator_color?: string;
	    host?: string;
	    port?: string;
	    ssl_or_https?: boolean;
	    authentication_method?: string;
	    username?: string;
	    password?: string;
	    set_as_default?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UpdateConfigRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connection_name = source["connection_name"];
	        this.env_indicator_color = source["env_indicator_color"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.ssl_or_https = source["ssl_or_https"];
	        this.authentication_method = source["authentication_method"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.set_as_default = source["set_as_default"];
	    }
	}

}

export namespace service {
	
	export class CacheInfo {
	    exists: boolean;
	    size: number;
	    modTime: string;
	    cacheKey: string;
	    cachePath: string;
	    isExpired: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CacheInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.exists = source["exists"];
	        this.size = source["size"];
	        this.modTime = source["modTime"];
	        this.cacheKey = source["cacheKey"];
	        this.cachePath = source["cachePath"];
	        this.isExpired = source["isExpired"];
	    }
	}

}

