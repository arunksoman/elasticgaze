export namespace models {
	
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

