<script>
	/**
	 * @fileoverview Elasticsearch Connections Management Page
	 * 
	 * Main page component for managing Elasticsearch connections. Provides functionality to:
	 * - View all saved connections in a grid layout
	 * - Add new connections via modal form
	 * - Edit existing connections
	 * - Delete connections with confirmation
	 * - Test connection validity
	 * - Set default connection for the application
	 * - Display toast notifications for user feedback
	 * 
	 * This page orchestrates multiple child components and manages the overall state
	 * for connection CRUD operations.
	 */
	
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Toast from '$lib/Toast.svelte';
	import ConnectionsList from '$lib/components/ConnectionsList.svelte';
	import ConnectionForm from '$lib/components/ConnectionForm.svelte';
	import { connections, connectionService } from '$lib/stores/connectionStore.js';
	import { refreshConnectionStatus } from '$lib/stores/connectionWarningStore.js';
	import { triggerConnectionUpdate } from '$lib/stores/connectionUpdateStore.js';
	import { CreateConfig, UpdateConfig, GetAllConfigs, DeleteConfig } from '$lib/wailsjs/go/main/App.js';
	
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
	 * @typedef {Object} ToastData
	 * @property {string} message - Toast message text
	 * @property {('success'|'error'|'info'|'warning')} type - Toast type for styling
	 * @property {number} [duration] - Auto-hide duration in milliseconds
	 * @property {('fade'|'slide')} [animation] - Animation type
	 * @property {string} [errorCode] - Optional error code for error toasts
	 * @property {string} [errorDetails] - Detailed error information
	 */
	
	// ===== COMPONENT STATE =====
	/** @type {boolean} Whether the connection form modal is visible */
	let showForm = false;
	
	/** @type {Connection|null} Connection being edited (null for new connection) */
	let editingConnection = null;
	
	/** @type {Connection} Form data for connection creation/editing */
	let formData = connectionService.getDefaultFormData();
	
	// ===== TOAST STATE =====
	/** @type {boolean} Whether toast notification is visible */
	let toastShow = false;
	
	/** @type {string} Current toast message */
	let toastMessage = '';
	
	/** @type {('success'|'error'|'info'|'warning')} Current toast type */
	let toastType = 'success';
	
	/** @type {number} Toast auto-hide duration in milliseconds */
	let toastDuration = 1500;
	
	/** @type {('fade'|'slide')} Toast animation type */
	let toastAnimation = 'fade';
	
	/** @type {string} Optional error code for error toasts */
	let toastErrorCode = '';
	
	/** @type {string} Detailed error information for error toasts */
	let toastErrorDetails = '';
	
	// ===== TESTING STATE =====
	/** @type {string|null} ID of connection currently being tested */
	let testingConnectionId = null;
	
	// ===== STORE SUBSCRIPTION =====
	/** @type {Connection[]} Array of all connections from the store */
	let connectionsArray = [];
	
	/** Store unsubscribe function for cleanup */
	const unsubscribe = connections.subscribe(value => {
		connectionsArray = value;
	});
	
	/**
	 * Loads connections from Go backend API and updates the store
	 */
	async function loadConnectionsFromBackend() {
		try {
			const configs = await GetAllConfigs();
			const mappedConnections = configs.map(mapConfigToConnection);
			connections.set(mappedConnections);
		} catch (error) {
			console.error('Error loading connections from backend:', error);
			showToast('Failed to load connections from database', 'error', 3000);
			// Fallback to localStorage
			connectionService.load();
		}
	}
	
	/**
	 * Initialize component by loading saved connections from backend
	 */
	onMount(() => {
		loadConnectionsFromBackend();
	});

	// Cleanup subscription on destroy
	import { onDestroy } from 'svelte';
	
	/**
	 * Cleanup store subscription when component is destroyed
	 */
	onDestroy(() => {
		unsubscribe();
	});

	// ===== EVENT HANDLERS =====
	
	/**
	 * Converts frontend Connection object to backend CreateConfigRequest format
	 * @param {Connection} connection - Frontend connection object
	 * @returns {Object} Backend CreateConfigRequest object
	 */
	function mapConnectionToCreateRequest(connection) {
		return {
			connection_name: connection.name,
			env_indicator_color: connection.environmentColor,
			host: connection.host,
			port: connection.port.toString(),
			ssl_or_https: connection.useSSL,
			authentication_method: connection.authType,
			username: connection.authType === 'basic' ? connection.username || null : null,
			password: connection.authType === 'basic' ? connection.password || null : null,
			set_as_default: connection.isDefault
		};
	}
	
	/**
	 * Converts frontend Connection object to backend UpdateConfigRequest format
	 * @param {Connection} connection - Frontend connection object
	 * @returns {Object} Backend UpdateConfigRequest object
	 */
	function mapConnectionToUpdateRequest(connection) {
		return {
			connection_name: connection.name,
			env_indicator_color: connection.environmentColor,
			host: connection.host,
			port: connection.port.toString(),
			ssl_or_https: connection.useSSL,
			authentication_method: connection.authType,
			username: connection.authType === 'basic' ? connection.username || null : null,
			password: connection.authType === 'basic' ? connection.password || null : null,
			set_as_default: connection.isDefault
		};
	}
	
	/**
	 * Converts backend Config object to frontend Connection format
	 * @param {Object} config - Backend config object
	 * @returns {Connection} Frontend connection object
	 */
	function mapConfigToConnection(config) {
		return {
			id: config.id.toString(),
			name: config.connection_name,
			host: config.host,
			port: parseInt(config.port),
			useSSL: config.ssl_or_https,
			authType: config.authentication_method,
			username: config.username || '',
			password: config.password || '',
			apiKey: '', // API key not supported in backend yet
			isDefault: config.set_as_default,
			environmentColor: config.env_indicator_color
		};
	}
	
	/**
	 * Opens the connection form modal for creating or editing a connection
	 * @param {Connection|null} connection - Connection to edit, or null for new connection
	 */
	function openForm(connection = null) {
		editingConnection = connection;
		if (connection) {
			formData = { ...connection };
		} else {
			formData = connectionService.getDefaultFormData();
		}
		showForm = true;
	}

	/**
	 * Closes the connection form modal and resets editing state
	 */
	function closeForm() {
		showForm = false;
		editingConnection = null;
	}

	/**
	 * Saves a connection (creates new or updates existing) using Go backend API
	 * @param {Connection} connectionData - Connection data to save
	 */
	async function saveConnection(connectionData) {
		try {
			let result;
			
			if (editingConnection) {
				// Update existing connection
				const updateRequest = mapConnectionToUpdateRequest(connectionData);
				result = await UpdateConfig(parseInt(editingConnection.id), updateRequest);
			} else {
				// Create new connection
				const createRequest = mapConnectionToCreateRequest(connectionData);
				result = await CreateConfig(createRequest);
			}
			
			// Success - close form and reload connections
			closeForm();
			await loadConnectionsFromBackend();
			
			// Refresh connection warning status
			await refreshConnectionStatus();
			
			// Trigger connection update for layout strip
			triggerConnectionUpdate();
			
			// Show success toast
			const action = editingConnection ? 'updated' : 'created';
			showToast(`Connection successfully ${action}`, 'success', 2000);
			
		} catch (error) {
			console.error('Error saving connection:', error);
			
			// Handle validation errors - extract meaningful error messages
			let errorMessage = 'Failed to save connection. Please try again.';
			
			if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
				errorMessage = error.message;
			} else if (typeof error === 'string') {
				errorMessage = error;
			} else if (error && typeof error.toString === 'function') {
				errorMessage = error.toString();
			}
			
			// Check for specific validation errors and provide user-friendly messages
			if (errorMessage.includes('only one default connection is allowed')) {
				showToast('Only one default connection is allowed. Please uncheck "Set as default" or edit the existing default connection.', 'error', 5000, 'slide');
			} else if (errorMessage.includes('validation error')) {
				showToast('Please check your connection details and try again.', 'error', 3000);
			} else if (errorMessage.includes('connection name is required')) {
				showToast('Connection name is required. Please enter a name for your connection.', 'error', 3000);
			} else if (errorMessage.includes('host is required')) {
				showToast('Host is required. Please enter the Elasticsearch host address.', 'error', 3000);
			} else {
				// Show the actual error message from the backend
				showToast(`Failed to save connection: ${errorMessage}`, 'error', 4000);
			}
		}
	}

	/**
	 * Deletes a connection by ID using Go backend API
	 * @param {string} connectionId - ID of connection to delete
	 */
	async function deleteConnection(connectionId) {
		try {
			await DeleteConfig(parseInt(connectionId));
			await loadConnectionsFromBackend();
			
			// Refresh connection warning status
			await refreshConnectionStatus();
			
			// Trigger connection update for layout strip
			triggerConnectionUpdate();
			
			showToast('Connection deleted successfully', 'success', 2000);
		} catch (error) {
			console.error('Error deleting connection:', error);
			showToast('Failed to delete connection', 'error', 3000);
		}
	}

	/**
	 * Sets a connection as the default connection using Go backend API
	 * @param {string} connectionId - ID of connection to set as default
	 */
	async function setAsDefault(connectionId) {
		try {
			// Get the current connection
			const connection = connectionsArray.find(c => c.id === connectionId);
			if (!connection) {
				throw new Error('Connection not found');
			}
			
			// Update the connection to set as default
			const updateRequest = {
				...mapConnectionToUpdateRequest(connection),
				set_as_default: true
			};
			
			await UpdateConfig(parseInt(connectionId), updateRequest);
			await loadConnectionsFromBackend();
			
			// Refresh connection warning status
			await refreshConnectionStatus();
			
			// Trigger connection update for layout strip
			triggerConnectionUpdate();
			
			showToast('Default connection updated successfully', 'success', 2000);
		} catch (error) {
			console.error('Error setting default connection:', error);
			
			// Extract meaningful error message
			let errorMessage = 'Failed to set default connection';
			if (error && typeof error === 'object' && 'message' in error && typeof error.message === 'string') {
				errorMessage = error.message;
			} else if (typeof error === 'string') {
				errorMessage = error;
			} else if (error && typeof error.toString === 'function') {
				errorMessage = error.toString();
			}
			
			// Handle the case where trying to set multiple defaults
			if (errorMessage.includes('only one default connection is allowed')) {
				showToast('Only one default connection is allowed. Another connection is already set as default.', 'error', 4000, 'slide');
			} else {
				showToast(`Failed to set default connection: ${errorMessage}`, 'error', 3000);
			}
		}
	}

	/**
	 * Handles test start event by tracking which connection is being tested
	 * @param {string} testingId - ID of connection being tested
	 */
	function handleTestStart(testingId) {
		testingConnectionId = testingId;
	}

	/**
	 * Handles test end event by clearing the testing state
	 */
	function handleTestEnd() {
		testingConnectionId = null;
	}

	/**
	 * Navigates back to the home page
	 */
	function goBack() {
		goto('/');
	}
	
	// ===== TOAST UTILITY FUNCTIONS =====
	
	/**
	 * Shows a toast notification with specified parameters
	 * @param {string} message - Toast message text
	 * @param {('success'|'error'|'info'|'warning')} type - Toast type for styling
	 * @param {number} duration - Auto-hide duration in milliseconds
	 * @param {('fade'|'slide')} animation - Animation type
	 * @param {string} errorCode - Optional error code for error toasts
	 * @param {string} errorDetails - Detailed error information
	 */
	function showToast(message, type = 'success', duration = 1500, animation = 'fade', errorCode = '', errorDetails = '') {
		toastMessage = message;
		toastType = type;
		toastDuration = duration;
		toastAnimation = animation;
		toastErrorCode = errorCode;
		toastErrorDetails = errorDetails;
		toastShow = true;
	}

	/**
	 * Hides the current toast and clears error details
	 */
	function hideToast() {
		toastShow = false;
		// Clear error details when hiding
		toastErrorCode = '';
		toastErrorDetails = '';
	}

	/**
	 * Handles toast events from child components
	 * @param {ToastData} toastData - Toast configuration object
	 */
	function handleToast(toastData) {
		const { message, type, duration, animation, errorCode, errorDetails } = toastData;
		showToast(message, type, duration, animation, errorCode, errorDetails);
	}

	/**
	 * Handles add connection button click by opening the form
	 */
	function handleAddConnection() {
		openForm();
	}
</script>

<div class="w-full max-w-6xl mx-auto mt-2">
	<div class="flex items-center gap-4 mb-6">
		<button 
			onclick={goBack}
			class="p-2 rounded-md theme-bg-secondary theme-text-primary hover:theme-bg-tertiary transition-colors"
			title="Back to Home"
			aria-label="Back to Home"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
		</button>
		<h1 class="text-2xl font-medium theme-text-primary">Elasticsearch Connections</h1>
		{#if connectionsArray.length > 0}
			<button 
				onclick={handleAddConnection}
				class="ml-auto bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors flex items-center gap-2"
			>
				<img src="/icons/create.svg" alt="" class="w-4 h-4" style="filter: brightness(0) invert(1);" />
				Add Connection
			</button>
		{/if}
	</div>
	
	<!-- Connections List Component -->
	<ConnectionsList
		connections={connectionsArray}
		{testingConnectionId}
		onadd={handleAddConnection}
		onteststart={handleTestStart}
		ontestend={handleTestEnd}
		onedit={(connection) => openForm(connection)}
		ondelete={deleteConnection}
		onsetdefault={setAsDefault}
		ontoast={handleToast}
	/>
</div>

<!-- Connection Form Modal Component -->
<ConnectionForm
	bind:show={showForm}
	{editingConnection}
	bind:formData
	onclose={closeForm}
	onsave={saveConnection}
	ontoast={handleToast}
/>

<!-- Toast Component -->
<Toast 
	bind:show={toastShow}
	message={toastMessage}
	type={toastType}
	duration={toastDuration}
	animation={toastAnimation}
	errorCode={toastErrorCode}
	errorDetails={toastErrorDetails}
	on:hide={hideToast}
/>