<script>
	import MonacoEditor from '$lib/MonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	export let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	export let height = '250px';
	export let title = 'Request Body';
</script>

<div class="mb-6">
	<span class="block mb-3 theme-text-primary font-medium text-lg">{title}</span>
	<MonacoEditor 
		bind:value={requestBody} 
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
	/>
</div>