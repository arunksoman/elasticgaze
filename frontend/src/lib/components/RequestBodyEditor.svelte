<script>
	import { createEventDispatcher } from 'svelte';
	import MonacoEditor from '$lib/MonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	const dispatch = createEventDispatcher();
	
	export let requestBody = '';  // No default request body
	export let height = '100%';
	export let title = 'Request Body';
	
	function handleChange(event) {
		requestBody = event.detail;
		dispatch('change', requestBody);
	}
</script>

<div class="mb-6 h-full flex flex-col">
	<span class="block mb-3 theme-text-primary font-medium text-lg flex-shrink-0">{title}</span>
	<div class="flex-1 min-h-0">
		<MonacoEditor 
			value={requestBody}
			language="json" 
			{height}
			theme={$theme}
			fontSize={13}
			tabSize={2}
			wordWrap="on"
			formatOnPaste={true}
			formatOnType={false}
			folding={true}
			showFoldingControls="always"
			placeholder="Enter your JSON request body here..."
			on:change={handleChange}
		/>
	</div>
</div>