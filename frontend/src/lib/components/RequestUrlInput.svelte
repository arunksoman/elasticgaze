<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	export let endpoint = '';
	
	function handleInput(event) {
		endpoint = event.target.value;
		dispatch('change', endpoint);
	}
</script>

<input 
	type="text" 
	value={endpoint}
	on:input={handleInput}
	class="flex-1 border theme-border p-2 theme-bg-tertiary theme-text-primary rounded"
/>