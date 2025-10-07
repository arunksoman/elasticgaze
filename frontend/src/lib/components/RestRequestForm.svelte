<script>
	import { createEventDispatcher } from 'svelte';
	import RequestMethodSelector from './RequestMethodSelector.svelte';
	import RequestUrlInput from './RequestUrlInput.svelte';
	
	const dispatch = createEventDispatcher();
	
	export let method = 'GET';
	export let endpoint = '';
	export let isLoading = false;
	
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
		on:click={handleSendRequest}
		disabled={isLoading}
		class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-2 rounded font-medium transition-colors"
	>
		{isLoading ? 'Sending...' : 'Send'}
	</button>
</div>