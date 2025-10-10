<script>
	import { createEventDispatcher } from 'svelte';
	import LazyMonacoEditor from './LazyMonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		requestBody = '',  // No default request body
		height = '100%',
		title = 'Request Body'
	} = $props();
	
	function handleChange(event) {
		requestBody = event.detail;
		dispatch('change', requestBody);
	}
</script>

<div class="mb-6 h-full flex flex-col">
	<span class="block mb-3 theme-text-primary font-medium text-lg flex-shrink-0">{title}</span>
	<div class="flex-1 min-h-0">
		<LazyMonacoEditor 
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
			loadingHeight={height}
			loadingText="Loading request editor..."
			on:change={handleChange}
		/>
	</div>
</div>