<script>
	import { onMount } from 'svelte';
	import TabbedRestContainer from '$lib/components/TabbedRestContainer.svelte';
	import { warmMonacoEditor } from '$lib/services/monacoPreloader.js';
	
	// Aggressively warm Monaco Editor as soon as REST page loads
	onMount(() => {
		warmMonacoEditor();
	});
</script>

<TabbedRestContainer />