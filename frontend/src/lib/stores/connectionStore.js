import { writable } from 'svelte/store';

// Create a writable store for connections
export const connections = writable([]);

// Connection management functions
export const connectionService = {
	// Load connections from localStorage
	load() {
		const stored = localStorage.getItem('elasticsearch-connections');
		if (stored) {
			const loadedConnections = JSON.parse(stored);
			connections.set(loadedConnections);
			return loadedConnections;
		}
		return [];
	},

	// Save connections to localStorage
	save(connectionsArray) {
		localStorage.setItem('elasticsearch-connections', JSON.stringify(connectionsArray));
		connections.set(connectionsArray);
	},

	// Add a new connection
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

	// Update an existing connection
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

	// Delete a connection
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

	// Set a connection as default
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

	// Get default connection
	getDefault() {
		return new Promise((resolve) => {
			connections.subscribe(currentConnections => {
				const defaultConnection = currentConnections.find(c => c.isDefault);
				resolve(defaultConnection || null);
			})();
		});
	},

	// Get connection by ID
	getById(connectionId) {
		return new Promise((resolve) => {
			connections.subscribe(currentConnections => {
				const connection = currentConnections.find(c => c.id === connectionId);
				resolve(connection || null);
			})();
		});
	},

	// Reset form data to default values
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