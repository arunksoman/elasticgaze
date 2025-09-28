/**
 * @fileoverview Type definitions for the ElasticGaze application
 */

/**
 * @typedef {Object} Connection
 * @property {string} id - Unique connection identifier
 * @property {string} name - Display name for the connection
 * @property {string} host - Elasticsearch host address
 * @property {number} port - Elasticsearch port number
 * @property {boolean} useSSL - Whether to use SSL/HTTPS
 * @property {('basic'|'apikey'|'none')} authType - Authentication type
 * @property {string} username - Username for basic auth
 * @property {string} password - Password for basic auth
 * @property {string} apiKey - API key for API key auth
 * @property {boolean} isDefault - Whether this is the default connection
 * @property {string} environmentColor - Environment color indicator
 */

/**
 * @typedef {Object} FormErrors
 * @property {string} name - Error message for connection name
 * @property {string} host - Error message for host
 * @property {string} port - Error message for port
 * @property {string} username - Error message for username
 * @property {string} password - Error message for password
 * @property {string} apiKey - Error message for API key
 */

/**
 * @typedef {Object} ToastData
 * @property {string} message - Toast message to display
 * @property {('success'|'error'|'warning'|'info')} [type='success'] - Toast type
 * @property {number} [duration=1500] - Toast duration in milliseconds
 * @property {string} [animation='fade'] - Toast animation type
 * @property {string} [errorCode] - Error code for error toasts
 * @property {string} [errorDetails] - Error details for error toasts
 */

/**
 * @typedef {Object} EnvironmentColor
 * @property {string} name - Display name of the color
 * @property {string} value - Color value/key
 * @property {string} color - Hex color code
 */

/**
 * @typedef {Object} TestConnectionRequest
 * @property {string} host - Elasticsearch host
 * @property {string} port - Port as string
 * @property {boolean} ssl_or_https - Whether to use SSL
 * @property {string} authentication_method - Auth method
 * @property {string} [username] - Username for basic auth
 * @property {string} [password] - Password for basic auth
 * @property {string} [api_key] - API key for API key auth
 */

/**
 * @typedef {Object} TestConnectionResponse
 * @property {boolean} success - Whether the test was successful
 * @property {string} [message] - Response message
 * @property {string} [cluster_name] - Elasticsearch cluster name
 * @property {string} [version] - Elasticsearch version
 * @property {string} [error_code] - Error code if failed
 * @property {string} [error_details] - Error details if failed
 */