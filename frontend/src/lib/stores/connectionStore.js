import { writable } from 'svelte/store';

/**
 * @typedef {Object} Connection
 * @property {string|null} id - Unique connection identifier (null for new connections)
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
 * Svelte writable store containing array of connection objects
 * @type {import('svelte/store').Writable<Array<Connection>>}
 */
export const connections = writable([]);

/**
 * Connection management service providing CRUD operations and localStorage persistence
 * @namespace connectionService
 */
export const connectionService = {
	/**
	 * Loads connections from localStorage and updates the store
	 * @memberof connectionService
	 * @function
	 * @returns {Array<Connection>} Array of loaded connections
	 */
	load() {
		const stored = localStorage.getItem('elasticsearch-connections');
		if (stored) {
			const loadedConnections = JSON.parse(stored);
			connections.set(loadedConnections);
			return loadedConnections;
		}
		return [];
	},

	/**
	 * Saves connections array to localStorage and updates the store
	 * @memberof connectionService
	 * @function
	 * @param {Array<Connection>} connectionsArray - Array of connection objects to save
	 * @returns {void}
	 */
	save(connectionsArray) {
		localStorage.setItem('elasticsearch-connections', JSON.stringify(connectionsArray));
		connections.set(connectionsArray);
	},

	/**
	 * Adds a new connection to the store and localStorage
	 * @memberof connectionService
	 * @function
	 * @param {Connection} newConnection - New connection object to add
	 * @returns {Promise<Connection>} Promise resolving to the added connection with generated ID
	 */
	add(newConnection) {
		return new Promise((resolve) => {
			connections.update(currentConnections => {
				const connectionData = {
					...newConnection,
					id: Date.now().toString()
				};

				// If this connection is set as default, remove default from others
				let updatedConnections = currentConnections;
				if (connectionData.isDefault) {
					updatedConnections = currentConnections.map(c => ({ ...c, isDefault: false }));
				}

				const finalConnections = [...updatedConnections, connectionData];
				this.save(finalConnections);
				resolve(connectionData);
				return finalConnections;
			});
		});
	},

	/**
	 * Updates an existing connection in the store and localStorage
	 * @memberof connectionService
	 * @function
	 * @param {Connection} updatedConnection - Updated connection object
	 * @returns {Promise<Connection|null>} Promise resolving to updated connection or null if not found
	 */
	update(updatedConnection) {
		return new Promise((resolve) => {
			connections.update(currentConnections => {
				const index = currentConnections.findIndex(c => c.id === updatedConnection.id);
				if (index !== -1) {
					let updatedConnections = [...currentConnections];
					updatedConnections[index] = updatedConnection;

					// If this connection is set as default, remove default from others
					if (updatedConnection.isDefault) {
						updatedConnections = updatedConnections.map(c => 
							c.id === updatedConnection.id ? c : { ...c, isDefault: false }
						);
					}

					this.save(updatedConnections);
					resolve(updatedConnection);
					return updatedConnections;
				}
				resolve(null);
				return currentConnections;
			});
		});
	},

	/**
	 * Deletes a connection from the store and localStorage
	 * @memberof connectionService
	 * @function
	 * @param {string} connectionId - ID of the connection to delete
	 * @returns {Promise<boolean>} Promise resolving to true when deletion is complete
	 */
	delete(connectionId) {
		return new Promise((resolve) => {
			connections.update(currentConnections => {
				const filteredConnections = currentConnections.filter(c => c.id !== connectionId);
				this.save(filteredConnections);
				resolve(true);
				return filteredConnections;
			});
		});
	},

	/**
	 * Sets a connection as the default connection
	 * @memberof connectionService
	 * @function
	 * @param {string} connectionId - ID of the connection to set as default
	 * @returns {Promise<boolean>} Promise resolving to true when update is complete
	 */
	setAsDefault(connectionId) {
		return new Promise((resolve) => {
			connections.update(currentConnections => {
				const updatedConnections = currentConnections.map(c => ({
					...c,
					isDefault: c.id === connectionId
				}));
				this.save(updatedConnections);
				resolve(true);
				return updatedConnections;
			});
		});
	},

	/**
	 * Gets the default connection from the store
	 * @memberof connectionService
	 * @function
	 * @returns {Promise<Connection|null>} Promise resolving to default connection or null if none set
	 */
	getDefault() {
		return new Promise((resolve) => {
			connections.subscribe(currentConnections => {
				const defaultConnection = currentConnections.find(c => c.isDefault);
				resolve(defaultConnection || null);
			})();
		});
	},

	/**
	 * Gets a connection by its ID
	 * @memberof connectionService
	 * @function
	 * @param {string} connectionId - ID of the connection to retrieve
	 * @returns {Promise<Connection|null>} Promise resolving to connection object or null if not found
	 */
	getById(connectionId) {
		return new Promise((resolve) => {
			connections.subscribe(currentConnections => {
				const connection = currentConnections.find(c => c.id === connectionId);
				resolve(connection || null);
			})();
		});
	},

	/**
	 * Returns default form data object for new connections
	 * @memberof connectionService
	 * @function
	 * @returns {Connection} Default form data object
	 */
	getDefaultFormData() {
		return {
			id: null,
			name: '',
			host: '',
			port: 9200,
			username: '',
			password: '',
			useSSL: false,
			apiKey: '',
			authType: 'basic',
			isDefault: false,
			environmentColor: 'dodgerblue'
		};
	}
};