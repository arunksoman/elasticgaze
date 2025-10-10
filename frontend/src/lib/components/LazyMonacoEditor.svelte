<script>
	import { onMount } from 'svelte';
	import { getMonacoEditor, isMonacoEditorLoaded } from '$lib/services/monacoPreloader.js';
	
	// Props
	let {
		loadingHeight = '200px',
		loadingText = 'Loading editor...',
		...editorProps
	} = $props();
	
	let isLoaded = $state(isMonacoEditorLoaded());
	let MonacoEditor = $state(null);
	let loadError = $state(null);
	
	onMount(async () => {
		if (isLoaded && MonacoEditor) {
			return; // Already loaded
		}
		
		try {
			// Use the preloader service to get Monaco Editor
			const module = await getMonacoEditor();
			MonacoEditor = module.default;
			isLoaded = true;
		} catch (error) {
			console.error('Failed to load Monaco Editor:', error);
			loadError = error;
		}
	});
</script>

{#if loadError}
	<div class="flex items-center justify-center border border-red-300 bg-red-50 text-red-700 rounded-lg p-4" style="height: {loadingHeight}">
		<div class="text-center">
			<p class="font-medium">Failed to load editor</p>
			<p class="text-sm mt-1">{loadError.message}</p>
		</div>
	</div>
{:else if !isLoaded || !MonacoEditor}
	<!-- Loading skeleton -->
	<div class="flex items-center justify-center theme-bg-secondary border border-gray-300 dark:border-gray-600 rounded-lg animate-pulse" style="height: {loadingHeight}">
		<div class="text-center">
			<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500 mx-auto mb-2"></div>
			<p class="theme-text-secondary text-sm">{loadingText}</p>
		</div>
	</div>
{:else}
	<!-- Loaded Monaco Editor -->
	<svelte:component this={MonacoEditor} {...editorProps} />
{/if}