<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	
	// Connection management state
	let connections = [];
	let showForm = false;
	let editingConnection = null;
	let formData = {
		id: null,
		name: '',
		host: '',
		port: 9200,
		username: '',
		password: '',
		useSSL: false,
		apiKey: '',
		authType: 'basic', // 'basic', 'apikey', 'none'
		isDefault: false
	};
	
	// Form validation
	let errors = {};
	
	onMount(() => {
		loadConnections();
	});
	
	function loadConnections() {
		// Load from localStorage for now - in production, this could be from a backend/config file
		const stored = localStorage.getItem('elasticsearch-connections');
		if (stored) {
			connections = JSON.parse(stored);
		}
	}
	
	function saveConnections() {
		localStorage.setItem('elasticsearch-connections', JSON.stringify(connections));
	}
	
	function validateForm() {
		errors = {};
		
		if (!formData.name.trim()) {
			errors.name = 'Connection name is required';
		}
		
		if (!formData.host.trim()) {
			errors.host = 'Host is required';
		}
		
		if (formData.port < 1 || formData.port > 65535) {
			errors.port = 'Port must be between 1 and 65535';
		}
		
		if (formData.authType === 'basic') {
			if (!formData.username.trim()) {
				errors.username = 'Username is required for basic auth';
			}
			if (!formData.password.trim()) {
				errors.password = 'Password is required for basic auth';
			}
		}
		
		if (formData.authType === 'apikey') {
			if (!formData.apiKey.trim()) {
				errors.apiKey = 'API Key is required for API key auth';
			}
		}
		
		return Object.keys(errors).length === 0;
	}
	
	function openForm(connection = null) {
		editingConnection = connection;
		if (connection) {
			formData = { ...connection };
		} else {
			formData = {
				id: null,
				name: '',
				host: '',
				port: 9200,
				username: '',
				password: '',
				useSSL: false,
				apiKey: '',
				authType: 'basic',
				isDefault: false
			};
		}
		showForm = true;
		errors = {};
	}
	
	function closeForm() {
		showForm = false;
		editingConnection = null;
		errors = {};
	}
	
	function saveConnection() {
		if (!validateForm()) {
			return;
		}
		
		const connectionData = { ...formData };
		
		if (editingConnection) {
			// Update existing connection
			const index = connections.findIndex(c => c.id === editingConnection.id);
			if (index !== -1) {
				connections[index] = connectionData;
			}
		} else {
			// Add new connection
			connectionData.id = Date.now().toString();
			connections = [...connections, connectionData];
		}
		
		// If this connection is set as default, remove default from others
		if (connectionData.isDefault) {
			connections = connections.map(c => 
				c.id === connectionData.id ? c : { ...c, isDefault: false }
			);
		}
		
		saveConnections();
		closeForm();
	}
	
	function deleteConnection(connectionId) {
		if (confirm('Are you sure you want to delete this connection?')) {
			connections = connections.filter(c => c.id !== connectionId);
			saveConnections();
		}
	}
	
	function setAsDefault(connectionId) {
		connections = connections.map(c => ({
			...c,
			isDefault: c.id === connectionId
		}));
		saveConnections();
	}
	
	function testConnection(connection) {
		// Placeholder for connection testing
		alert(`Testing connection to ${connection.name}...\nThis would attempt to connect to: ${connection.host}:${connection.port}`);
	}
	
	function goBack() {
		goto('/');
	}
</script>

<div class="p-6">
	<div class="flex items-center gap-4 mb-6">
		<button 
			onclick={goBack}
			class="p-2 rounded-md theme-bg-secondary theme-text-primary hover:theme-bg-tertiary transition-colors"
			title="Back to Home"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
		</button>
		<h1 class="text-2xl font-medium theme-text-primary">Elasticsearch Connections</h1>
	</div>
	
	<!-- Connections List -->
	{#if connections.length === 0}
		<div class="text-center py-12">
			<div class="w-16 h-16 mx-auto mb-4 theme-text-secondary">
				<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
			</div>
			<h3 class="text-lg font-medium theme-text-primary mb-2">No connections configured</h3>
			<p class="theme-text-secondary mb-4">Add your first Elasticsearch connection to get started.</p>
			<button 
				onclick={() => openForm()}
				class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors"
			>
				Add Connection
			</button>
		</div>
	{:else}
		<div class="grid gap-4">
			{#each connections as connection (connection.id)}
				<div class="theme-bg-secondary border theme-border rounded-lg p-4">
					<div class="flex items-center justify-between">
						<div class="flex-1">
							<div class="flex items-center gap-3 mb-2">
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
								onclick={() => testConnection(connection)}
								class="bg-green-600 hover:bg-green-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors"
								title="Test Connection"
							>
								Test
							</button>
							{#if !connection.isDefault}
								<button 
									onclick={() => setAsDefault(connection.id)}
									class="bg-yellow-600 hover:bg-yellow-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors"
									title="Set as Default"
								>
									Set Default
								</button>
							{/if}
							<button 
								onclick={() => openForm(connection)}
								class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors"
								title="Edit Connection"
							>
								Edit
							</button>
							<button 
								onclick={() => deleteConnection(connection.id)}
								class="bg-red-600 hover:bg-red-700 text-white px-3 py-1.5 rounded text-sm font-medium transition-colors"
								title="Delete Connection"
							>
								Delete
							</button>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Connection Form Modal -->
{#if showForm}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
		<div class="theme-bg-primary rounded-lg shadow-xl max-w-2xl w-full max-h-[580px] overflow-y-auto">
			<div class="p-4">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium theme-text-primary">
						{editingConnection ? 'Edit Connection' : 'Add New Connection'}
					</h2>
					<button 
						onclick={closeForm}
						class="p-1.5 rounded-md theme-text-secondary hover:theme-text-primary transition-colors"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
				
				<form onsubmit={(e) => { e.preventDefault(); saveConnection(); }} class="space-y-3">
					<!-- Connection Name -->
					<div>
						<label for="name" class="block text-sm font-medium theme-text-primary mb-1">Connection Name</label>
						<input 
							type="text" 
							id="name"
							bind:value={formData.name}
							class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
							placeholder="Production ES Cluster"
						/>
						{#if errors.name}
							<p class="text-red-500 text-xs mt-1">{errors.name}</p>
						{/if}
					</div>
					
					<!-- Host and Port -->
					<div class="grid grid-cols-3 gap-3">
						<div class="col-span-2">
							<label for="host" class="block text-sm font-medium theme-text-primary mb-1">Host</label>
							<input 
								type="text" 
								id="host"
								bind:value={formData.host}
								class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
								placeholder="localhost or https://example.com"
							/>
							{#if errors.host}
								<p class="text-red-500 text-xs mt-1">{errors.host}</p>
							{/if}
						</div>
						<div>
							<label for="port" class="block text-sm font-medium theme-text-primary mb-1">Port</label>
							<input 
								type="number" 
								id="port"
								bind:value={formData.port}
								class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
								placeholder="9200"
								min="1"
								max="65535"
							/>
							{#if errors.port}
								<p class="text-red-500 text-xs mt-1">{errors.port}</p>
							{/if}
						</div>
					</div>
					
					<!-- SSL -->
					<div class="flex items-center py-1">
						<input 
							type="checkbox" 
							id="useSSL"
							bind:checked={formData.useSSL}
							class="mr-2"
						/>
						<label for="useSSL" class="text-sm font-medium theme-text-primary">Use SSL/HTTPS</label>
					</div>
					
					<!-- Authentication Type -->
					<div>
						<label for="authType" class="block text-sm font-medium theme-text-primary mb-1">Authentication</label>
						<select 
							id="authType"
							bind:value={formData.authType}
							class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
						>
							<option value="none">No Authentication</option>
							<option value="basic">Basic Auth (Username/Password)</option>
							<option value="apikey">API Key</option>
						</select>
					</div>
					
					<!-- Basic Auth Fields -->
					{#if formData.authType === 'basic'}
						<div class="grid grid-cols-2 gap-3">
							<div>
								<label for="username" class="block text-sm font-medium theme-text-primary mb-1">Username</label>
								<input 
									type="text" 
									id="username"
									bind:value={formData.username}
									class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
									placeholder="elastic"
								/>
								{#if errors.username}
									<p class="text-red-500 text-xs mt-1">{errors.username}</p>
								{/if}
							</div>
							<div>
								<label for="password" class="block text-sm font-medium theme-text-primary mb-1">Password</label>
								<input 
									type="password" 
									id="password"
									bind:value={formData.password}
									class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
									placeholder="••••••••"
								/>
								{#if errors.password}
									<p class="text-red-500 text-xs mt-1">{errors.password}</p>
								{/if}
							</div>
						</div>
					{/if}
					
					<!-- API Key Field -->
					{#if formData.authType === 'apikey'}
						<div>
							<label for="apiKey" class="block text-sm font-medium theme-text-primary mb-1">API Key</label>
							<input 
								type="password" 
								id="apiKey"
								bind:value={formData.apiKey}
								class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded-md text-sm"
								placeholder="Enter your API key"
							/>
							{#if errors.apiKey}
								<p class="text-red-500 text-xs mt-1">{errors.apiKey}</p>
							{/if}
						</div>
					{/if}
					
					<!-- Set as Default -->
					<div class="flex items-center py-1">
						<input 
							type="checkbox" 
							id="isDefault"
							bind:checked={formData.isDefault}
							class="mr-2"
						/>
						<label for="isDefault" class="text-sm font-medium theme-text-primary">Set as default connection</label>
					</div>
					
					<!-- Form Actions -->
					<div class="flex justify-end gap-3 pt-3">
						<button 
							type="button"
							onclick={closeForm}
							class="px-3 py-2 border theme-border theme-text-primary rounded-md hover:theme-bg-secondary transition-colors text-sm"
						>
							Cancel
						</button>
						<button 
							type="submit"
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors text-sm"
						>
							{editingConnection ? 'Update' : 'Save'} Connection
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}