<script>
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';
	
	export let value = '';
	export let language = 'json';
	export let readOnly = false;
	export let height = '200px';
	export let theme = 'vs-dark';
	export let fontSize = 14;
	export let tabSize = 2;
	export let wordWrap = 'off';
	export let lineNumbers = 'on';
	export let minimap = false;
	export let scrollBeyondLastLine = false;
	export let formatOnPaste = true;
	export let formatOnType = true;
	export let autoClosingBrackets = 'always';
	export let autoClosingQuotes = 'always';
	export let placeholder = '';
	export let folding = true;
	export let showFoldingControls = 'always';
	
	let editor;
	let container;
	let resizeObserver;
	let handleWindowResize;
	
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
		
		// Configure JSON language services for better formatting support
		monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
			validate: true,
			allowComments: false,
			schemas: [],
			enableSchemaRequest: false
		});
		
		// Ensure JSON formatter is available by registering document formatting provider
		monaco.languages.registerDocumentFormattingEditProvider('json', {
			provideDocumentFormattingEdits: (model, options, token) => {
				try {
					const text = model.getValue();
					const formatted = JSON.stringify(JSON.parse(text), null, options.tabSize);
					return [{
						range: model.getFullModelRange(),
						text: formatted
					}];
				} catch (e) {
					// If JSON is invalid, return no edits
					return [];
				}
			}
		});
		
		// Register folding range provider for JSON to ensure folding works
		monaco.languages.registerFoldingRangeProvider('json', {
			provideFoldingRanges: function(model, context, token) {
				const foldingRanges = [];
				const text = model.getValue();
				const lines = text.split('\n');
				
				const stack = [];
				for (let i = 0; i < lines.length; i++) {
					const line = lines[i];
					const trimmed = line.trim();
					
					// Opening braces/brackets
					if (trimmed.includes('{') || trimmed.includes('[')) {
						stack.push({ start: i + 1, type: trimmed.includes('{') ? 'object' : 'array' });
					}
					
					// Closing braces/brackets
					if (trimmed.includes('}') || trimmed.includes(']')) {
						const lastOpen = stack.pop();
						if (lastOpen && i > lastOpen.start) {
							foldingRanges.push({
								start: lastOpen.start,
								end: i + 1,
								kind: monaco.languages.FoldingRangeKind.Region
							});
						}
					}
				}
				
				return foldingRanges;
			}
		});
		
		// Create editor instance
		editor = monaco.editor.create(container, {
			value: value,
			language: language,
			theme: theme === 'dark' ? 'custom-dark' : 'custom-light',
			readOnly: readOnly,
			minimap: { enabled: minimap },
			scrollBeyondLastLine: scrollBeyondLastLine,
			fontSize: fontSize,
			lineNumbers: lineNumbers,
			roundedSelection: false,
			scrollbar: {
				verticalScrollbarSize: 8,
				horizontalScrollbarSize: 8
			},
			automaticLayout: true,
			folding: folding,
			foldingStrategy: 'indentation', // Use indentation strategy for more reliable folding
			showFoldingControls: showFoldingControls,
			foldingHighlight: true,
			glyphMargin: true, // Ensure glyph margin is enabled for folding controls
			// JSON-specific enhancements
			formatOnPaste: formatOnPaste,
			formatOnType: formatOnType,
			// Enable bracket matching and auto-closing
			matchBrackets: 'always',
			autoClosingBrackets: autoClosingBrackets,
			autoClosingQuotes: autoClosingQuotes,
			// Better selection and indentation
			selectOnLineNumbers: true,
			detectIndentation: true,
			insertSpaces: true,
			tabSize: tabSize,
			wordWrap: wordWrap,
			// Enable context menu with formatting options even in read-only mode
			contextmenu: true,
			readOnlyMessage: { value: 'Cannot edit in read-only mode' }
		});
		
		// Override context menu to enable formatting in read-only mode
		if (readOnly) {
			editor.addAction({
				id: 'format-document-readonly',
				label: 'Format Document',
				contextMenuGroupId: '1_modification',
				contextMenuOrder: 1,
				run: function(ed) {
					// Temporarily make editor writable for formatting
					const model = ed.getModel();
					if (model) {
						ed.updateOptions({ readOnly: false });
						ed.getAction('editor.action.formatDocument').run().then(() => {
							ed.updateOptions({ readOnly: true });
						});
					}
				}
			});
		}
		
		// Listen for content changes
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});
		
		// Add placeholder functionality
		if (placeholder) {
			updatePlaceholder();
		}
		
		// Add keyboard shortcut for formatting (Shift+Alt+F)
		editor.addCommand(monaco.KeyMod.Shift | monaco.KeyMod.Alt | monaco.KeyCode.KeyF, () => {
			formatDocument();
		});
		
		// Set up ResizeObserver to handle container size changes
		if (typeof ResizeObserver !== 'undefined') {
			resizeObserver = new ResizeObserver((entries) => {
				if (editor) {
					// Multiple approaches to ensure layout updates
					requestAnimationFrame(() => {
						editor.layout();
						// Force a second layout after a short delay for stubborn cases
						setTimeout(() => {
							editor.layout();
						}, 100);
					});
				}
			});
			
			// Observe the container for size changes
			resizeObserver.observe(container);
		}
		
		// Additional fallback: listen for window resize events
		const windowResizeHandler = () => {
			if (editor) {
				setTimeout(() => editor.layout(), 50);
			}
		};
		window.addEventListener('resize', windowResizeHandler);
		
		// Store reference for cleanup
		handleWindowResize = windowResizeHandler;
	});
	
	// Function to handle placeholder display
	function updatePlaceholder() {
		if (!editor || !placeholder) return;
		
		const model = editor.getModel();
		if (!model) return;
		
		// Show placeholder when editor is empty
		if (model.getValue() === '') {
			editor.deltaDecorations([], [{
				range: new monaco.Range(1, 1, 1, 1),
				options: {
					after: {
						content: placeholder,
						inlineClassName: 'placeholder-text'
					}
				}
			}]);
		} else {
			editor.deltaDecorations(editor.deltaDecorations([], []), []);
		}
	}
	
	// Function to format the document
	function formatDocument() {
		if (editor) {
			editor.getAction('editor.action.formatDocument').run();
		}
	}
	
	// Export format function for external use
	export { formatDocument };
	
	onDestroy(() => {
		if (resizeObserver) {
			resizeObserver.disconnect();
		}
		if (editor) {
			editor.dispose();
		}
		// Clean up window resize listener
		window.removeEventListener('resize', handleWindowResize);
	});
	
	// Update editor when value changes externally
	$: if (editor && editor.getValue() !== value) {
		editor.setValue(value);
		if (placeholder) {
			updatePlaceholder();
		}
	}
	
	// Update theme when it changes
	$: if (editor) {
		monaco.editor.setTheme(theme === 'dark' ? 'custom-dark' : 'custom-light');
	}
</script>

<style>
	:global(.placeholder-text) {
		color: #6b7280 !important;
		font-style: italic;
		opacity: 0.7;
	}
</style>

<div bind:this={container} style="height: {height}; width: 100%;" class="border theme-border rounded min-h-0"></div>