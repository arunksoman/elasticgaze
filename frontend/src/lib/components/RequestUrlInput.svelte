<script lang="ts">
	import { createEventDispatcher, onMount, onDestroy, tick } from 'svelte';
	import { ConnectionUrlService } from '$lib/services/connectionUrlService.js';
	import { connectionUpdateTrigger } from '$lib/stores/connectionUpdateStore.js';
	import { UrlAutocompleteService, type UrlSuggestion } from '$lib/services/urlAutocompleteService.js';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		endpoint = '',
		method = 'GET'
	} = $props();
	
	// State for base URL information
	let baseUrlData = $state({
		baseUrl: '',
		displayUrl: 'Loading...',
		connectionName: ''
	});
	
	let showTooltip = $state(false);
	let isLoading = $state(true);
	
	// Autocomplete state
	let suggestions = $state<UrlSuggestion[]>([]);
	let showSuggestions = $state(false);
	let selectedSuggestionIndex = $state(-1);
	let inputElement: HTMLInputElement;
	let suggestionsContainer = $state<HTMLDivElement>();
	
	// Load base URL on component mount
	onMount(async () => {
		await loadBaseUrl();
	});
	
	// Cleanup on component destroy
	onDestroy(() => {
		clearTimeout(debounceTimer);
	});
	
	// Listen for connection updates
	$effect(() => {
		// React to connection update triggers
		$connectionUpdateTrigger;
		loadBaseUrl();
	});
	
	// Debounce timer for input handling
	let debounceTimer: ReturnType<typeof setTimeout>;
	let previousMethod = method;
	let isMouseOverSuggestions = $state(false);
	
	// Handle method changes without causing reactive loops
	$effect(() => {
		if (method !== previousMethod) {
			previousMethod = method;
			// Clear existing suggestions and update if there's enough input
			if (endpoint.trim() && endpoint.trim().length >= 2) {
				clearTimeout(debounceTimer);
				debounceTimer = setTimeout(() => {
					updateSuggestions();
				}, 100);
			} else {
				suggestions = [];
				showSuggestions = false;
				selectedSuggestionIndex = -1;
			}
		}
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
		
		// Debounce suggestions update to prevent excessive calls
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(() => {
			updateSuggestions();
		}, 150);
		
		// Dispatch only the endpoint, let the parent handle complete URL building
		dispatch('change', endpoint);
	}
	
	function updateSuggestions() {
		try {
			if (!endpoint.trim() || endpoint.trim().length < 2) {
				suggestions = [];
				showSuggestions = false;
				selectedSuggestionIndex = -1;
				return;
			}
			
			const newSuggestions = UrlAutocompleteService.getSuggestions(method, endpoint, 20);
			
			// Only update if suggestions count changed or first suggestion is different
			const shouldUpdate = newSuggestions.length !== suggestions.length ||
				(newSuggestions.length > 0 && suggestions.length > 0 && newSuggestions[0].path !== suggestions[0].path);
			
			if (shouldUpdate) {
				suggestions = newSuggestions;
				showSuggestions = newSuggestions.length > 0;
				selectedSuggestionIndex = -1;
			}
		} catch (error) {
			console.error('Error updating suggestions:', error);
			suggestions = [];
			showSuggestions = false;
			selectedSuggestionIndex = -1;
		}
	}
	
	function handleKeydown(event: KeyboardEvent) {
		if (!showSuggestions || suggestions.length === 0) {
			if (event.key === 'Escape') {
				showSuggestions = false;
				selectedSuggestionIndex = -1;
			}
			return;
		}
		
		switch (event.key) {
			case 'ArrowDown':
				event.preventDefault();
				selectedSuggestionIndex = Math.min(selectedSuggestionIndex + 1, suggestions.length - 1);
				scrollToSelectedSuggestion();
				break;
			case 'ArrowUp':
				event.preventDefault();
				selectedSuggestionIndex = Math.max(selectedSuggestionIndex - 1, -1);
				scrollToSelectedSuggestion();
				break;
			case 'Enter':
				event.preventDefault();
				if (selectedSuggestionIndex >= 0) {
					selectSuggestion(suggestions[selectedSuggestionIndex]);
				}
				break;
			case 'Escape':
				event.preventDefault();
				showSuggestions = false;
				selectedSuggestionIndex = -1;
				break;
		}
	}
	
	async function scrollToSelectedSuggestion() {
		await tick();
		if (selectedSuggestionIndex >= 0 && suggestionsContainer) {
			const selectedElement = suggestionsContainer.children[selectedSuggestionIndex] as HTMLElement;
			if (selectedElement) {
				selectedElement.scrollIntoView({ block: 'nearest' });
			}
		}
	}
	
	function selectSuggestion(suggestion: UrlSuggestion) {
		// Clear any pending debounced updates
		clearTimeout(debounceTimer);
		
		endpoint = suggestion.path;
		showSuggestions = false;
		selectedSuggestionIndex = -1;
		suggestions = []; // Clear suggestions to prevent re-renders
		isMouseOverSuggestions = false; // Reset mouse state
		
		dispatch('change', endpoint);
		
		// Use setTimeout to ensure the focus happens after the blur
		setTimeout(() => {
			inputElement?.focus();
		}, 0);
	}
	
	function handleFocus() {
		if (endpoint.trim() && endpoint.trim().length >= 2 && suggestions.length > 0) {
			showSuggestions = true;
		}
	}
	
	function handleBlur() {
		// Don't hide suggestions if mouse is over them
		if (isMouseOverSuggestions) {
			return;
		}
		
		// Delay to allow mouse click on suggestions
		setTimeout(() => {
			if (!isMouseOverSuggestions) {
				showSuggestions = false;
				selectedSuggestionIndex = -1;
			}
		}, 300);
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

<div class="flex-1 relative">
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
			bind:this={inputElement}
			type="text" 
			value={endpoint}
			oninput={handleInput}
			onkeydown={handleKeydown}
			onfocus={handleFocus}
			onblur={handleBlur}
			class="flex-1 px-3 py-2 theme-bg-tertiary theme-text-primary border-0 outline-none focus:theme-bg-primary h-full"
			autocomplete="off"
			spellcheck="false"
		/>
	</div>
	
	<!-- Suggestions Dropdown -->
	{#if showSuggestions && suggestions.length > 0}
		<div 
			bind:this={suggestionsContainer}
			class="absolute top-full left-0 right-0 z-50 mt-1 theme-bg-primary border theme-border rounded-md shadow-lg max-h-96 overflow-y-auto"
			role="listbox"
			tabindex="-1"
			onmouseenter={() => isMouseOverSuggestions = true}
			onmouseleave={() => isMouseOverSuggestions = false}
			onmousedown={(e) => e.preventDefault()}
		>
			{#each suggestions as suggestion, index}
				<button
					class="w-full px-3 py-2 text-left text-sm theme-text-primary hover:theme-bg-secondary transition-colors border-b theme-border last:border-b-0 font-mono {selectedSuggestionIndex === index ? 'theme-bg-secondary' : ''}"
					role="option"
					aria-selected={selectedSuggestionIndex === index}
					onclick={() => selectSuggestion(suggestion)}
					onmouseup={() => selectSuggestion(suggestion)}
					onmouseenter={() => selectedSuggestionIndex = index}
					onmousedown={(e) => e.preventDefault()}
				>
					<div class="flex items-center justify-between">
						<span class="truncate">{suggestion.path}</span>
						<span class="text-xs theme-text-muted ml-2 flex-shrink-0">
							{method}
						</span>
					</div>
				</button>
			{/each}
		</div>
	{/if}
	
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