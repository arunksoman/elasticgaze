<script>
	import MonacoEditor from '$lib/MonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	export let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	export let height = '100%';
	export let title = 'Request Body';
</script>

<div class="mb-6 h-full flex flex-col">
	<span class="block mb-3 theme-text-primary font-medium text-lg flex-shrink-0">{title}</span>
	<div class="flex-1 min-h-0">
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
</div>