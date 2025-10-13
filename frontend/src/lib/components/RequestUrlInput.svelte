<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import { ConnectionUrlService } from '$lib/services/connectionUrlService.js';
	import { connectionUpdateTrigger } from '$lib/stores/connectionUpdateStore.js';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		endpoint = ''
	} = $props();
	
	// State for base URL information
	let baseUrlData = $state({
		baseUrl: '',
		displayUrl: 'Loading...',
		connectionName: ''
	});
	
	let showTooltip = $state(false);
	let isLoading = $state(true);
	
	// Load base URL on component mount
	onMount(async () => {
		await loadBaseUrl();
	});
	
	// Listen for connection updates
	$effect(() => {
		// React to connection update triggers
		$connectionUpdateTrigger;
		loadBaseUrl();
	});
	
	async function loadBaseUrl() {
		isLoading = true;
		try {
			baseUrlData = await ConnectionUrlService.getCurrentBaseUrl();
		} catch (error) {
			console.error('Failed to load base URL:', error);
			baseUrlData = {
				baseUrl: '',
				displayUrl: 'Connection Error',
				connectionName: ''
			};
		} finally {
			isLoading = false;
		}
	}
	
	function handleInput(event: Event) {
		const target = event.target as HTMLInputElement;
		endpoint = target.value;
		// Dispatch only the endpoint, let the parent handle complete URL building
		dispatch('change', endpoint);
	}
	
	function handleMouseEnter() {
		if (baseUrlData.baseUrl) {
			showTooltip = true;
		}
	}
	
	function handleMouseLeave() {
		showTooltip = false;
	}
	
	// Get the complete URL for display
	function getCompleteUrl() {
		if (!baseUrlData.baseUrl) {
			// If no base URL, return the endpoint as-is (might be empty or a complete URL)
			return endpoint.trim();
		}
		return ConnectionUrlService.buildCompleteUrl(baseUrlData.baseUrl, endpoint);
	}
</script>

<div class="flex-1">
	<div class="flex items-stretch border theme-border rounded overflow-hidden h-10">
		<!-- Base URL Display Box -->
		<div 
			class="relative flex items-center px-3 py-2 theme-bg-secondary border-r theme-border cursor-help select-none min-w-0"
			onmouseenter={handleMouseEnter}
			onmouseleave={handleMouseLeave}
			title={baseUrlData.baseUrl ? `${baseUrlData.connectionName}: ${baseUrlData.baseUrl}` : 'No connection configured - Click to refresh'}
			role="button"
			tabindex="0"
		>
			<div class="flex items-center gap-2 min-w-0">
				<!-- Connection indicator -->
				<div class="w-2 h-2 rounded-full flex-shrink-0 {baseUrlData.baseUrl ? 'bg-green-500' : 'bg-red-400'}"></div>
				
			<!-- Base URL text -->
			<span class="text-sm theme-text-secondary truncate min-w-0">
				{isLoading ? 'Loading...' : (baseUrlData.baseUrl ? '{{base_url}}' : 'No Connection')}
			</span>				<!-- Refresh button for connection errors -->
				{#if !isLoading && !baseUrlData.baseUrl}
					<button 
						class="ml-1 p-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
						onclick={loadBaseUrl}
						title="Refresh connection"
						aria-label="Refresh connection"
					>
						<svg class="w-3 h-3 theme-text-secondary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
						</svg>
					</button>
				{/if}
			</div>
			
			<!-- Tooltip -->
			{#if showTooltip && baseUrlData.baseUrl}
				<div class="absolute bottom-full left-0 mb-2 px-3 py-2 theme-bg-primary border theme-border rounded shadow-lg z-50 max-w-sm">
					<div class="text-xs theme-text-primary whitespace-nowrap">
						<div class="font-medium mb-1">{baseUrlData.connectionName}</div>
						<div class="font-mono">{baseUrlData.baseUrl}</div>
						<div class="text-xs opacity-75 mt-1">Click to refresh • Complete URL shown below</div>
					</div>
				</div>
			{:else if showTooltip && !baseUrlData.baseUrl}
				<div class="absolute bottom-full left-0 mb-2 px-3 py-2 theme-bg-primary border theme-border rounded shadow-lg z-50 max-w-sm">
					<div class="text-xs theme-text-primary">
						<div class="font-medium text-red-400 mb-1">No Connection</div>
						<div>Please configure a default connection in Settings</div>
						<div class="text-xs opacity-75 mt-1">Click the refresh button to retry</div>
					</div>
				</div>
			{/if}
		</div>
		
		<!-- Endpoint Input -->
		<input 
			type="text" 
			value={endpoint}
			oninput={handleInput}
			class="flex-1 px-3 py-2 theme-bg-tertiary theme-text-primary border-0 outline-none focus:theme-bg-primary h-full"
		/>
	</div>
	
	<!-- Complete URL Display - Reserve space to prevent layout shifts -->
	<div class="mt-1 h-4 text-xs theme-text-muted font-mono break-all">
		{#if baseUrlData.baseUrl && endpoint.trim()}
			Complete URL: {getCompleteUrl()}
		{:else}
			<!-- Empty space to maintain consistent height -->
			&nbsp;
		{/if}
	</div>
</div>