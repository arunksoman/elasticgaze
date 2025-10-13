import { GetDefaultConfig } from '$lib/wailsjs/go/main/App.js';

/**
 * Service for building connection URLs and managing base URL display
 */
export class ConnectionUrlService {
	/**
	 * Gets the default connection configuration
	 * @returns {Promise<Object|null>} Default connection config or null
	 */
	static async getDefaultConnection() {
		try {
			if (typeof window === 'undefined' || !window.go?.main?.App?.GetDefaultConfig) {
				return null;
			}
			return await GetDefaultConfig();
		} catch (error) {
			console.error('Failed to get default connection:', error);
			return null;
		}
	}

	/**
	 * Builds the base URL from connection configuration
	 * @param {Object} config - Connection configuration object
	 * @returns {string} Base URL (e.g., "https://localhost:9200")
	 */
	static buildBaseUrl(config) {
		if (!config || !config.host) {
			return '';
		}

		const protocol = config.ssl_or_https ? 'https' : 'http';
		const port = config.port && config.port !== '80' && config.port !== '443' 
			? `:${config.port}` 
			: '';
		
		return `${protocol}://${config.host}${port}`;
	}

	/**
	 * Creates a shortened display version of the base URL for UI
	 * @param {string} baseUrl - Full base URL
	 * @returns {string} Shortened display URL
	 */
	static createDisplayUrl(baseUrl) {
		if (!baseUrl) {
			return 'No Connection';
		}

		try {
			const url = new URL(baseUrl);
			const host = url.hostname;
			
			// If it's a localhost or IP, show it fully
			if (host === 'localhost' || host === '127.0.0.1' || /^\d+\.\d+\.\d+\.\d+$/.test(host)) {
				return baseUrl;
			}
			
			// For domain names, show protocol + truncated host
			const parts = host.split('.');
			if (parts.length > 2) {
				const displayHost = `${parts[0]}...${parts[parts.length - 1]}`;
				return `${url.protocol}//${displayHost}${url.port ? `:${url.port}` : ''}`;
			}
			
			return baseUrl;
		} catch (error) {
			return baseUrl.length > 25 ? baseUrl.substring(0, 22) + '...' : baseUrl;
		}
	}

	/**
	 * Gets the current default connection's base URL
	 * @returns {Promise<{baseUrl: string, displayUrl: string, connectionName: string}>}
	 */
	static async getCurrentBaseUrl() {
		const config = await this.getDefaultConnection();
		if (!config) {
			return {
				baseUrl: '',
				displayUrl: 'No Connection',
				connectionName: ''
			};
		}

		const baseUrl = this.buildBaseUrl(config);
		const displayUrl = this.createDisplayUrl(baseUrl);
		
		return {
			baseUrl,
			displayUrl,
			connectionName: config.connection_name || 'Unknown Connection'
		};
	}

	/**
	 * Builds the complete URL from base URL and endpoint
	 * @param {string} baseUrl - Base URL
	 * @param {string} endpoint - API endpoint path
	 * @returns {string} Complete URL
	 */
	static buildCompleteUrl(baseUrl, endpoint) {
		if (!baseUrl) {
			return endpoint;
		}

		// Clean the endpoint
		const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
		
		return baseUrl + cleanEndpoint;
	}
}