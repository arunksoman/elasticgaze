<script>
	import { TestConnection } from '$lib/wailsjs/go/main/App.js';
	
	// Props
	export let show = false;
	export let editingConnection = null;
	export let formData = {
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
	
	// Event props (Svelte 5 way)
	export let onclose;
	export let onsave;
	export let ontoast;

	// Environment color options
	const environmentColors = [
		{ name: 'Red', value: 'red', color: '#ef4444' },
		{ name: 'Orange', value: 'orange', color: '#f97316' },
		{ name: 'Yellow', value: 'yellow', color: '#eab308' },
		{ name: 'Green', value: 'green', color: '#22c55e' },
		{ name: 'Dodger Blue', value: 'dodgerblue', color: '#3b82f6' },
		{ name: 'Purple', value: 'purple', color: '#a855f7' },
		{ name: 'Pink', value: 'pink', color: '#ec4899' }
	];

	// Form validation
	let errors = {
		name: '',
		host: '',
		port: '',
		username: '',
		password: '',
		apiKey: ''
	};

	function validateForm() {
		errors = {
			name: '',
			host: '',
			port: '',
			username: '',
			password: '',
			apiKey: ''
		};
		
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
		
		return Object.values(errors).every(error => error === '');
	}

	function getEnvironmentColorValue(colorName) {
		const colorObj = environmentColors.find(c => c.value === colorName);
		return colorObj ? colorObj.color : '#3b82f6'; // Default to dodger blue
	}

	function closeForm() {
		errors = {
			name: '',
			host: '',
			port: '',
			username: '',
			password: '',
			apiKey: ''
		};
		onclose?.();
	}

	function saveConnection() {
		if (!validateForm()) {
			return;
		}
		onsave?.(formData);
	}

	// Test connection from form data
	async function testFormConnection() {
		if (!formData.host.trim()) {
			ontoast?.({ message: 'Host is required for testing', type: 'error' });
			return;
		}
		
		try {
			const testRequest = {
				host: formData.host,
				port: formData.port.toString(),
				ssl_or_https: formData.useSSL,
				authentication_method: formData.authType,
				username: formData.authType === 'basic' ? formData.username : undefined,
				password: formData.authType === 'basic' ? formData.password : undefined,
				api_key: formData.authType === 'apikey' ? formData.apiKey : undefined
			};
			
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
		}
	}

	// Clear errors when form data changes
	$: if (formData) {
		// Only clear errors that are not currently showing (to avoid flicker)
		if (Object.values(errors).some(error => error !== '')) {
			errors = {
				name: '',
				host: '',
				port: '',
				username: '',
				password: '',
				apiKey: ''
			};
		}
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4 pt-16">
		<div class="theme-bg-primary rounded-lg shadow-xl max-w-2xl w-full max-h-[550px] overflow-y-auto" style="border-top: 5px solid {getEnvironmentColorValue(formData.environmentColor)};">
			<div class="p-4">
				<div class="flex items-center justify-between mb-4">
					<h2 class="text-lg font-medium theme-text-primary">
						{editingConnection ? 'Edit Connection' : 'Add New Connection'}
					</h2>
					<button 
						onclick={closeForm}
						class="p-1.5 rounded-md theme-text-secondary hover:theme-text-primary transition-colors"
						aria-label="Close form"
					>
						<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
						</svg>
					</button>
				</div>
				
				<form onsubmit={(e) => { e.preventDefault(); saveConnection(); }} class="space-y-2">
					<!-- Environment Color Selector -->
					<fieldset>
						<legend class="block text-sm font-medium theme-text-primary mb-2">Environment Color</legend>
						<div class="flex gap-2 flex-wrap">
							{#each environmentColors as colorOption}
								<button
									type="button"
									onclick={() => formData.environmentColor = colorOption.value}
									class="w-4 h-4 rounded-full border transition-all duration-200 hover:scale-125 {formData.environmentColor === colorOption.value ? 'border-gray-400 shadow-lg' : 'border-gray-200 dark:border-gray-600'}"
									style="background-color: {colorOption.color};"
									title={colorOption.name}
									aria-label="Select {colorOption.name} color"
								></button>
							{/each}
						</div>
					</fieldset>
					
					<!-- Connection Name -->
					<div>
						<label for="name" class="block text-sm font-medium theme-text-primary mb-1">Connection Name</label>
						<input 
							type="text" 
							id="name"
							bind:value={formData.name}
							class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
							placeholder="Production ES Cluster"
						/>
						{#if errors.name}
							<p class="text-red-500 text-xs mt-1">{errors.name}</p>
						{/if}
					</div>
					
					<!-- Host and Port -->
					<div class="grid grid-cols-3 gap-2">
						<div class="col-span-2">
							<label for="host" class="block text-sm font-medium theme-text-primary mb-1">Host</label>
							<input 
								type="text" 
								id="host"
								bind:value={formData.host}
								class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
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
								class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
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
							class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
						>
							<option value="none">No Authentication</option>
							<option value="basic">Basic Auth (Username/Password)</option>
							<option value="apikey">API Key</option>
						</select>
					</div>
					
					<!-- Basic Auth Fields -->
					{#if formData.authType === 'basic'}
						<div class="grid grid-cols-2 gap-2">
							<div>
								<label for="username" class="block text-sm font-medium theme-text-primary mb-1">Username</label>
								<input 
									type="text" 
									id="username"
									bind:value={formData.username}
									class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
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
									class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
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
								class="w-full border theme-border p-1.5 theme-bg-tertiary theme-text-primary rounded-md text-sm"
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
					<div class="flex justify-between gap-2 pt-3">
						<button 
							type="button"
							onclick={testFormConnection}
							class="bg-green-600 hover:bg-green-700 text-white px-3 py-1.5 rounded-md font-medium transition-colors text-sm flex items-center gap-1"
						>
							<img src="/icons/test.svg" alt="" class="w-4 h-4" style="filter: brightness(0) invert(1);" />
							Test Connection
						</button>
						<div class="flex gap-2">
							<button 
								type="button"
								onclick={closeForm}
								class="px-3 py-1.5 border theme-border theme-text-primary rounded-md hover:theme-bg-secondary transition-colors text-sm"
							>
								Cancel
							</button>
							<button 
								type="submit"
								class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded-md font-medium transition-colors text-sm flex items-center gap-1"
							>
								<img src="/icons/{editingConnection ? 'edit' : 'create'}.svg" alt="" class="w-4 h-4" style="filter: brightness(0) invert(1);" />
								{editingConnection ? 'Update' : 'Save'} Connection
							</button>
						</div>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}