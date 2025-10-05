<script>
	import { createEventDispatcher } from 'svelte';
	import RequestMethodSelector from './RequestMethodSelector.svelte';
	import RequestUrlInput from './RequestUrlInput.svelte';
	
	const dispatch = createEventDispatcher();
	
	export let method = 'GET';
	export let endpoint = '';
	export let placeholder = '/_cluster/health';
	export let isLoading = false;
	
	function handleSendRequest() {
		dispatch('send', {
			method,
			endpoint
		});
	}
</script>

<!-- URL and Method Section -->
<div class="flex gap-3 mb-6">
	<RequestMethodSelector bind:method />
	<RequestUrlInput bind:endpoint {placeholder} />
	<button 
		on:click={handleSendRequest}
		disabled={isLoading}
		class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white px-6 py-3 rounded font-medium transition-colors"
	>
		{isLoading ? 'Sending...' : 'Send'}
	</button>
</div>