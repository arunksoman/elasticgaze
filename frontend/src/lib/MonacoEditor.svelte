<script>
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';
	
	export let value = '';
	export let language = 'json';
	export let readOnly = false;
	export let height = '200px';
	export let theme = 'vs-dark';
	
	let editor;
	let container;
	
	onMount(() => {
		// Configure Monaco Editor
		monaco.editor.defineTheme('custom-dark', {
			base: 'vs-dark',
			inherit: true,
			rules: [],
			colors: {
				'editor.background': '#1f2937',
				'editor.foreground': '#f9fafb'
			}
		});
		
		monaco.editor.defineTheme('custom-light', {
			base: 'vs',
			inherit: true,
			rules: [],
			colors: {
				'editor.background': '#ffffff',
				'editor.foreground': '#111827'
			}
		});
		
		// Create editor instance
		editor = monaco.editor.create(container, {
			value: value,
			language: language,
			theme: theme === 'dark' ? 'custom-dark' : 'custom-light',
			readOnly: readOnly,
			minimap: { enabled: false },
			scrollBeyondLastLine: false,
			fontSize: 14,
			lineNumbers: 'on',
			roundedSelection: false,
			scrollbar: {
				verticalScrollbarSize: 8,
				horizontalScrollbarSize: 8
			},
			automaticLayout: true
		});
		
		// Listen for content changes
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});
	});
	
	onDestroy(() => {
		if (editor) {
			editor.dispose();
		}
	});
	
	// Update editor when value changes externally
	$: if (editor && editor.getValue() !== value) {
		editor.setValue(value);
	}
	
	// Update theme when it changes
	$: if (editor) {
		monaco.editor.setTheme(theme === 'dark' ? 'custom-dark' : 'custom-light');
	}
</script>

<div bind:this={container} style="height: {height};" class="border theme-border rounded"></div>