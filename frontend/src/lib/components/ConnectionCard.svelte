<script>
	import { TestConnection } from '$lib/wailsjs/go/main/App.js';
	
	// Props
	export let connection;
	export let testingConnectionId = null;
	
	// Event props (Svelte 5 way)
	export let onteststart;
	export let ontestend;
	export let onedit;
	export let ondelete;
	export let onsetdefault;
	export let ontoast;

	function getEnvironmentColorValue(colorName) {
		const environmentColors = [
			{ name: 'Red', value: 'red', color: '#ef4444' },
			{ name: 'Orange', value: 'orange', color: '#f97316' },
			{ name: 'Yellow', value: 'yellow', color: '#eab308' },
			{ name: 'Green', value: 'green', color: '#22c55e' },
			{ name: 'Dodger Blue', value: 'dodgerblue', color: '#3b82f6' },
			{ name: 'Purple', value: 'purple', color: '#a855f7' },
			{ name: 'Pink', value: 'pink', color: '#ec4899' }
		];
		
		const colorObj = environmentColors.find(c => c.value === colorName);
		return colorObj ? colorObj.color : '#3b82f6'; // Default to dodger blue
	}

	// Test connection function
	async function testConnection() {
		onteststart?.(connection.id);
		
		try {
			// Prepare the test request
			const testRequest = {
				host: connection.host,
				port: connection.port.toString(),
				ssl_or_https: connection.useSSL,
				authentication_method: connection.authType,
				username: connection.authType === 'basic' ? connection.username : undefined,
				password: connection.authType === 'basic' ? connection.password : undefined,
				api_key: connection.authType === 'apikey' ? connection.apiKey : undefined
			};
			
			// Call the backend test function
			const response = await TestConnection(testRequest);
			
			if (response.success) {
				let successMessage = 'Connection successful';
				if (response.cluster_name) {
					successMessage += ` (${response.cluster_name})`;
				}
				if (response.version) {
					successMessage += ` - ES ${response.version}`;
				}
				ontoast?.({ message: successMessage, type: 'success', duration: 2000 });
			} else {
				const errorMessage = response.message || 'Connection failed';
				const errorCode = response.error_code || '';
				const errorDetails = response.error_details || '';
				ontoast?.({ 
					message: errorMessage, 
					type: 'error', 
					duration: 5000,
					animation: 'slide',
					errorCode,
					errorDetails
				});
			}
		} catch (error) {
			console.error('Test connection error:', error);
			ontoast?.({ 
				message: 'Connection test failed', 
				type: 'error', 
				duration: 5000,
				animation: 'slide',
				errorCode: 'UNKNOWN_ERROR',
				errorDetails: (error instanceof Error ? error.message : String(error)) || 'Unknown error occurred'
			});
		} finally {
			ontestend?.();
		}
	}

	function editConnection() {
		onedit?.(connection);
	}

	function deleteConnection() {
		if (confirm('Are you sure you want to delete this connection?')) {
			ondelete?.(connection.id);
		}
	}

	function setAsDefault() {
		onsetdefault?.(connection.id);
	}
</script>

<div class="theme-bg-secondary border theme-border rounded-lg p-4">
	<div class="flex items-center justify-between">
		<div class="flex-1">
			<div class="flex items-center gap-3 mb-2">
				<!-- Environment Color Indicator -->
				<div 
					class="w-3 h-3 rounded-full flex-shrink-0" 
					style="background-color: {getEnvironmentColorValue(connection.environmentColor || 'dodgerblue')}"
					title="Environment: {connection.environmentColor || 'dodgerblue'}"
				></div>
				<h3 class="text-lg font-medium theme-text-primary">{connection.name}</h3>
				{#if connection.isDefault}
					<span class="bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100 text-xs px-2 py-1 rounded-full font-medium">
						Default
					</span>
				{/if}
			</div>
			<div class="text-sm theme-text-secondary space-y-1">
				<p><strong>Host:</strong> {connection.host}:{connection.port}</p>
				<p><strong>SSL:</strong> {connection.useSSL ? 'Enabled' : 'Disabled'}</p>
				<p><strong>Auth:</strong> {connection.authType === 'basic' ? 'Basic Auth' : connection.authType === 'apikey' ? 'API Key' : 'None'}</p>
				{#if connection.authType === 'basic'}
					<p><strong>Username:</strong> {connection.username}</p>
				{/if}
			</div>
		</div>
		<div class="flex items-center gap-2">
			<button 
				onclick={testConnection}
				disabled={testingConnectionId === connection.id}
				class="bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors flex items-center gap-1"
				title="Test Connection"
			>
				{#if testingConnectionId === connection.id}
					<svg class="w-3 h-3 animate-spin" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Testing...
				{:else}
					Test
				{/if}
			</button>
			{#if !connection.isDefault}
				<button 
					onclick={setAsDefault}
					class="bg-yellow-600 hover:bg-yellow-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors"
					title="Set as Default"
				>
					Set Default
				</button>
			{/if}
			<button 
				onclick={editConnection}
				class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors flex items-center gap-1"
				title="Edit Connection"
			>
				<img src="/icons/edit.svg" alt="" class="w-3 h-3" style="filter: brightness(0) invert(1);" />
				Edit
			</button>
			<button 
				onclick={deleteConnection}
				class="bg-red-600 hover:bg-red-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors flex items-center gap-1"
				title="Delete Connection"
			>
				<img src="/icons/delete.svg" alt="" class="w-3 h-3" style="filter: brightness(0) invert(1);" />
				Delete
			</button>
		</div>
	</div>
</div>