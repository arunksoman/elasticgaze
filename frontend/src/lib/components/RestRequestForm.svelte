<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import RequestMethodSelector from './RequestMethodSelector.svelte';
	import RequestUrlInput from './RequestUrlInput.svelte';
	import { ConnectionUrlService } from '$lib/services/connectionUrlService.js';
	
	const dispatch = createEventDispatcher();
	
	let baseUrlData = $state({
		baseUrl: '',
		displayUrl: 'Loading...',
		connectionName: ''
	});
	
	// Use $props() for Svelte 5 runes mode
	let {
		method = 'GET',
		endpoint = '',
		isLoading = false
	} = $props();
	
	async function handleSendRequest() {
		// Build complete URL for sending to backend
		const currentBaseUrlData = await ConnectionUrlService.getCurrentBaseUrl();
		const completeUrl = ConnectionUrlService.buildCompleteUrl(currentBaseUrlData.baseUrl, endpoint);
		
		dispatch('send', {
			method,
			endpoint: completeUrl // Send complete URL to backend
		});
	}
	
	function handleMethodChange(newMethod: string) {
		method = newMethod;
		dispatch('methodChange', newMethod);
	}
	
	function handleEndpointChange(newEndpoint: string) {
		endpoint = newEndpoint;
		dispatch('endpointChange', newEndpoint);
	}
</script>

<!-- URL and Method Section -->
<div class="flex items-start gap-3 mb-6">
	<div class="flex-shrink-0">
		<RequestMethodSelector {method} on:change={(e) => handleMethodChange(e.detail)} />
	</div>
	<div class="flex-1 min-w-0">
		<RequestUrlInput {endpoint} on:change={(e) => handleEndpointChange(e.detail)} />
	</div>
	<div class="flex-shrink-0">
		<button 
			onclick={handleSendRequest}
			disabled={isLoading}
			class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-2 rounded font-medium transition-colors whitespace-nowrap h-10"
		>
			{isLoading ? 'Sending...' : 'Send'}
		</button>
	</div>
</div>