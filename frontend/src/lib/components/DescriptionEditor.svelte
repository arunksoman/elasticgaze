<script>
	import { marked } from 'marked';
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	// Props
	let { description = '' } = $props();
	
	// Local state
	let isPreviewMode = $state(false);
	let renderedHtml = $state('');
	let textareaRef = $state(null);
	
	// Configure marked for safety and simplicity
	marked.setOptions({
		breaks: true,
		gfm: true,
		headerIds: false,
		mangle: false
	});
	
	// Render markdown to HTML
	function renderMarkdown(text) {
		if (!text.trim()) {
			return '<p class="theme-text-secondary">Start typing your description...</p>';
		}
		try {
			return marked(text);
		} catch (error) {
			console.warn('Markdown rendering error:', error);
			return `<p class="text-red-500">Error rendering markdown: ${error.message}</p>`;
		}
	}
	
	// Update rendered HTML when description changes
	$effect(() => {
		renderedHtml = renderMarkdown(description);
	});
	
	// Handle input changes
	function handleInput(event) {
		description = event.target.value;
		dispatch('change', description);
	}
	
	// Handle tab key for better UX
	function handleKeydown(event) {
		if (event.key === 'Tab') {
			event.preventDefault();
			
			const start = event.target.selectionStart;
			const end = event.target.selectionEnd;
			const value = event.target.value;
			
			// Insert tab character
			const newValue = value.substring(0, start) + '  ' + value.substring(end);
			event.target.value = newValue;
			description = newValue;
			
			// Restore cursor position
			event.target.selectionStart = event.target.selectionEnd = start + 2;
			
			dispatch('change', description);
		}
	}
	
	// Toggle between edit and preview mode
	function togglePreview() {
		isPreviewMode = !isPreviewMode;
		
		// Focus textarea when switching back to edit mode
		if (!isPreviewMode && textareaRef) {
			setTimeout(() => textareaRef.focus(), 0);
		}
	}
</script>

<div class="h-full flex flex-col">
	<!-- Header with toggle button -->
	<div class="flex items-center justify-between px-3 py-2 border-b theme-border flex-shrink-0">
		<h3 class="text-sm font-medium theme-text-secondary">
			{isPreviewMode ? 'Markdown Preview' : 'Markdown Editor'}
		</h3>
		<button
			onclick={togglePreview}
			class="px-3 py-1 text-xs font-medium rounded-md transition-colors theme-bg-tertiary theme-text-primary hover:theme-hover"
		>
			{isPreviewMode ? 'Edit' : 'Preview'}
		</button>
	</div>
	
	<!-- Content area -->
	<div class="flex-1 min-h-0">
		{#if isPreviewMode}
			<!-- Preview Mode -->
			<div class="h-full px-6 overflow-auto theme-bg-secondary">
				<div class="prose prose-sm max-w-none theme-text-primary">
					{@html renderedHtml}
				</div>
			</div>
		{:else}
			<!-- Edit Mode -->
			<div class="h-full p-3 theme-bg-secondary">
				<textarea
					bind:this={textareaRef}
					bind:value={description}
					oninput={handleInput}
					onkeydown={handleKeydown}
					class="w-full h-full resize-none border-none outline-none 
					       theme-bg-secondary theme-text-primary
					       font-mono text-sm leading-relaxed"
					spellcheck="false"
				></textarea>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Custom scrollbar styling */
	.overflow-auto::-webkit-scrollbar {
		width: 6px;
	}
	
	.overflow-auto::-webkit-scrollbar-track {
		background: transparent;
	}
	
	.overflow-auto::-webkit-scrollbar-thumb {
		background: rgba(156, 163, 175, 0.3);
		border-radius: 3px;
	}
	
	.overflow-auto::-webkit-scrollbar-thumb:hover {
		background: rgba(156, 163, 175, 0.5);
	}
	
	/* Basic prose styling that inherits theme colors */
	:global(.prose) {
		color: var(--text-primary) !important;
	}
	
	/* Remove top margin from first element in prose */
	:global(.prose > *:first-child) {
		margin-top: 1em !important;
	}
	
	:global(.prose *) {
		color: inherit;
	}
	
	:global(.prose h1, .prose h2, .prose h3, .prose h4, .prose h5, .prose h6) {
		color: var(--text-primary) !important;
		font-weight: 600;
		margin-top: 1.5em;
		margin-bottom: 0.5em;
	}
	
	:global(.prose h1) { font-size: 1.5em; }
	:global(.prose h2) { font-size: 1.3em; }
	:global(.prose h3) { font-size: 1.1em; }
	
	:global(.prose p) {
		margin: 1em 0;
		color: var(--text-primary) !important;
	}
	
	:global(.prose ul, .prose ol) {
		padding-left: 1.5em;
		margin: 1em 0;
		color: var(--text-primary) !important;
	}
	
	:global(.prose li) {
		margin: 0.25em 0;
		color: var(--text-primary) !important;
	}
	
	:global(.prose strong) {
		font-weight: 600;
		color: var(--text-primary) !important;
	}
	
	:global(.prose em) {
		font-style: italic;
		color: var(--text-primary) !important;
	}
	
	:global(.prose code) {
		background: rgba(156, 163, 175, 0.1);
		padding: 2px 4px;
		border-radius: 3px;
		font-size: 0.875em;
		font-family: ui-monospace, SFMono-Regular, "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace;
		color: var(--text-primary) !important;
	}
	
	:global(.prose pre) {
		background: rgba(156, 163, 175, 0.1);
		padding: 12px;
		border-radius: 6px;
		font-size: 0.875em;
		line-height: 1.5;
		overflow-x: auto;
		margin: 1em 0;
		color: var(--text-primary) !important;
	}
	
	:global(.prose pre code) {
		background: transparent;
		padding: 0;
		color: var(--text-primary) !important;
	}
	
	:global(.prose blockquote) {
		border-left: 4px solid rgba(59, 130, 246, 0.5);
		padding-left: 1em;
		margin: 1em 0;
		font-style: italic;
		opacity: 0.8;
		color: var(--text-primary) !important;
	}
	
	:global(.prose table) {
		border-collapse: collapse;
		margin: 1em 0;
		width: 100%;
	}
	
	:global(.prose th, .prose td) {
		border: 1px solid rgba(156, 163, 175, 0.3);
		padding: 8px 12px;
		text-align: left;
	}
	
	:global(.prose th) {
		font-weight: 600;
		background: rgba(156, 163, 175, 0.1);
	}
	
	:global(.prose a) {
		color: #3b82f6;
		text-decoration: underline;
	}
</style>