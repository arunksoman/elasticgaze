<script>
	import { createEventDispatcher } from 'svelte';
	import RequestMethodSelector from './RequestMethodSelector.svelte';
	import RequestUrlInput from './RequestUrlInput.svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		method = 'GET',
		endpoint = '',
		isLoading = false
	} = $props();
	
	function handleSendRequest() {
		dispatch('send', {
			method,
			endpoint
		});
	}
	
	function handleMethodChange(newMethod) {
		method = newMethod;
		dispatch('methodChange', newMethod);
	}
	
	function handleEndpointChange(newEndpoint) {
		endpoint = newEndpoint;
		dispatch('endpointChange', newEndpoint);
	}
</script>

<!-- URL and Method Section -->
<div class="flex gap-3 mb-6">
	<RequestMethodSelector {method} on:change={(e) => handleMethodChange(e.detail)} />
	<RequestUrlInput {endpoint} on:change={(e) => handleEndpointChange(e.detail)} />
	<button 
		onclick={handleSendRequest}
		disabled={isLoading}
		class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-2 rounded font-medium transition-colors"
	>
		{isLoading ? 'Sending...' : 'Send'}
	</button>
</div>