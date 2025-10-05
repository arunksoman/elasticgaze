<script>
	import MonacoEditor from '$lib/MonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	export let responseData = '';
	export let height = '400px';
	export let title = 'Response';
</script>

{#if responseData}
	<div>
		<span class="block mb-3 theme-text-primary font-medium text-lg">{title}</span>
		<MonacoEditor 
			bind:value={responseData} 
			language="json" 
			{height}
			readOnly={true}
			theme={$theme}
			fontSize={12}
			tabSize={2}
			wordWrap="on"
			lineNumbers="on"
			minimap={true}
			formatOnPaste={false}
			formatOnType={false}
			folding={true}
			showFoldingControls="always"
		/>
	</div>
{/if}