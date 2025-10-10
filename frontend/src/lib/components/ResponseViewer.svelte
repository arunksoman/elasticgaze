<script>
	import LazyMonacoEditor from './LazyMonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	// Use $props() for Svelte 5 runes mode
	let {
		responseData = '',
		height = '100%',
		title = 'Response'
	} = $props();
</script>

<div class="h-full flex flex-col">
	<span class="block mb-3 theme-text-primary font-medium text-lg flex-shrink-0">{title}</span>
	<div class="flex-1 min-h-0">
		{#if responseData}
			<LazyMonacoEditor 
				value={responseData} 
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
				loadingHeight={height}
				loadingText="Loading response viewer..."
			/>
		{:else}
			<div class="h-full border theme-border rounded flex items-center justify-center theme-bg-secondary">
				<p class="theme-text-secondary text-center">
					<span class="block text-sm mb-2">No response yet</span>
					<span class="text-xs opacity-75">Send a request to see the response here</span>
				</p>
			</div>
		{/if}
	</div>
</div>