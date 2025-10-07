<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		endpoint = ''
	} = $props();
	
	function handleInput(event) {
		endpoint = event.target.value;
		dispatch('change', endpoint);
	}
</script>

<input 
	type="text" 
	value={endpoint}
	oninput={handleInput}
	class="flex-1 border theme-border p-2 theme-bg-tertiary theme-text-primary rounded"
/>