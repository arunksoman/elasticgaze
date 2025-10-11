<script>
	import { onMount, onDestroy, createEventDispatcher } from 'svelte';
	import * as monaco from 'monaco-editor';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		value = '',
		language = 'json',
		readOnly = false,
		height = '200px',
		theme = 'vs-dark',
		fontSize = 14,
		tabSize = 2,
		wordWrap = 'off',
		lineNumbers = 'on',
		minimap = false,
		scrollBeyondLastLine = false,
		formatOnPaste = true,
		formatOnType = true,
		autoClosingBrackets = 'always',
		autoClosingQuotes = 'always',
		placeholder = '',
		folding = true,
		showFoldingControls = 'always'
	} = $props();
	
	// Zoom functionality
	let currentZoom = $state(100); // Zoom percentage
	let baseFontSize = fontSize;
	let isFocused = $state(false); // Track if this editor is focused
	
	// Global reference to track active editor
	let editorId = Math.random().toString(36).substr(2, 9); // Unique ID for this editor instance
	
	let editor;
	let container;
	let resizeObserver;
	let handleWindowResize;
	let isUpdatingFromExternal = false;
	
	onMount(() => {
		// Create inline workers to avoid path issues
		if (typeof self !== 'undefined') {
			self.MonacoEnvironment = {
				getWorker: function (workerId, label) {
					// Create a simple inline worker
					const workerScript = `
						// Simple worker that handles Monaco messages
						self.onmessage = function(e) {
							try {
								const { id, method, args } = e.data;
								
								// Simple response for common Monaco worker methods
								let result = null;
								
								switch (method) {
									case 'getSemanticDiagnostics':
									case 'getSyntacticDiagnostics':
									case 'getSuggestionDiagnostics':
									case 'getCompilerOptionsDiagnostics':
										result = []; // Return empty array for diagnostics
										break;
									case 'getCompletionsAtPosition':
										result = { entries: [] }; // Return empty completions
										break;
									case 'getQuickInfoAtPosition':
									case 'getDefinitionAtPosition':
									case 'getTypeDefinitionAtPosition':
									case 'getImplementationAtPosition':
									case 'getReferencesAtPosition':
										result = undefined; // Return undefined for position-based queries
										break;
									case 'getRenameInfo':
										result = { canRename: false }; // Disable rename
										break;
									case 'findRenameLocations':
										result = []; // Return empty rename locations
										break;
									case 'getNavigateToItems':
										result = []; // Return empty navigation items
										break;
									case 'getFormattingEditsForDocument':
									case 'getFormattingEditsForRange':
										result = []; // Return empty formatting edits
										break;
									default:
										result = null;
								}
								
								// Send response back
								self.postMessage({
									seq: e.data.seq,
									type: 'response',
									success: true,
									request_seq: id,
									command: method,
									body: result
								});
							} catch (error) {
								// Send error response
								self.postMessage({
									seq: e.data.seq,
									type: 'response',
									success: false,
									request_seq: e.data.id,
									command: e.data.method,
									message: error.message
								});
							}
						};
						
						// Handle termination
						self.addEventListener('error', function(e) {
							console.warn('Monaco worker error:', e);
						});
					`;
					
					const blob = new Blob([workerScript], { type: 'application/javascript' });
					return new Worker(URL.createObjectURL(blob));
				}
			};
		}
		
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
		
		// Configure JSON mode with working inline worker
		monaco.languages.json.jsonDefaults.setModeConfiguration({
			documentFormattingEdits: true,
			documentRangeFormattingEdits: true,
			completionItems: true,
			hovers: true,
			documentSymbols: true,
			tokens: true,
			colors: true,
			foldingRanges: true,
			diagnostics: true, // Re-enable with our inline worker
			selectionRanges: true
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
			if (!isUpdatingFromExternal) {
				const newValue = editor.getValue();
				value = newValue;
				dispatch('change', newValue);
			}
		});
		
		// Add placeholder functionality
		if (placeholder) {
			updatePlaceholder();
		}
		
		// Add keyboard shortcut for formatting (Shift+Alt+F)
		editor.addCommand(monaco.KeyMod.Shift | monaco.KeyMod.Alt | monaco.KeyCode.KeyF, () => {
			formatDocument();
		});
		
		// Add focus and blur event listeners with debugging
		editor.onDidFocusEditorText(() => {
			isFocused = true;
			// Set this editor as the global active editor
			window.activeMonacoEditorId = editorId;
			window.activeMonacoEditor = {
				editor: editor,
				zoomIn: zoomIn,
				zoomOut: zoomOut,
				resetZoom: resetZoom,
				type: readOnly ? 'Response' : 'Request'
			};
			console.log(`Editor focused: ${readOnly ? 'Response' : 'Request'} Editor (ID: ${editorId})`);
		});
		
		editor.onDidBlurEditorText(() => {
			isFocused = false;
			// Clear global active editor if it's this one
			if (window.activeMonacoEditorId === editorId) {
				window.activeMonacoEditorId = null;
				window.activeMonacoEditor = null;
			}
			console.log(`Editor blurred: ${readOnly ? 'Response' : 'Request'} Editor (ID: ${editorId})`);
		});
		
		// Additional focus detection for widget focus
		editor.onDidFocusEditorWidget(() => {
			isFocused = true;
			// Set this editor as the global active editor
			window.activeMonacoEditorId = editorId;
			window.activeMonacoEditor = {
				editor: editor,
				zoomIn: zoomIn,
				zoomOut: zoomOut,
				resetZoom: resetZoom,
				type: readOnly ? 'Response' : 'Request'
			};
			console.log(`Editor widget focused: ${readOnly ? 'Response' : 'Request'} Editor (ID: ${editorId})`);
		});
		
		editor.onDidBlurEditorWidget(() => {
			// Only blur if neither text nor widget has focus
			if (!editor.hasTextFocus() && !editor.hasWidgetFocus()) {
				isFocused = false;
				// Clear global active editor if it's this one
				if (window.activeMonacoEditorId === editorId) {
					window.activeMonacoEditorId = null;
					window.activeMonacoEditor = null;
				}
				console.log(`Editor widget blurred: ${readOnly ? 'Response' : 'Request'} Editor (ID: ${editorId})`);
			}
		});
		
		// Set up global keyboard shortcuts if not already set up
		if (!window.monacoZoomListenerSetup) {
			window.monacoZoomListenerSetup = true;
			
			const handleGlobalKeydown = (e) => {
				// Check for Ctrl+Plus (Ctrl+=)
				if (e.ctrlKey && (e.key === '=' || e.key === '+')) {
					e.preventDefault();
					if (window.activeMonacoEditor) {
						console.log(`Global Zoom In: ${window.activeMonacoEditor.type} Editor`);
						window.activeMonacoEditor.zoomIn();
					}
				}
				// Check for Ctrl+Minus
				else if (e.ctrlKey && e.key === '-') {
					e.preventDefault();
					if (window.activeMonacoEditor) {
						console.log(`Global Zoom Out: ${window.activeMonacoEditor.type} Editor`);
						window.activeMonacoEditor.zoomOut();
					}
				}
				// Check for Ctrl+0
				else if (e.ctrlKey && e.key === '0') {
					e.preventDefault();
					if (window.activeMonacoEditor) {
						console.log(`Global Zoom Reset: ${window.activeMonacoEditor.type} Editor`);
						window.activeMonacoEditor.resetZoom();
					}
				}
			};
			
			window.addEventListener('keydown', handleGlobalKeydown);
			
			// Store reference for cleanup
			window.monacoZoomKeydownHandler = handleGlobalKeydown;
		}
		
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
	
	// Zoom functions
	function zoomIn() {
		if (currentZoom < 300) { // Max zoom 300%
			currentZoom += 10;
			updateEditorFontSize();
		}
	}
	
	function zoomOut() {
		if (currentZoom > 50) { // Min zoom 50%
			currentZoom -= 10;
			updateEditorFontSize();
		}
	}
	
	function resetZoom() {
		currentZoom = 100;
		updateEditorFontSize();
	}
	
	function updateEditorFontSize() {
		if (editor) {
			const newFontSize = Math.round((baseFontSize * currentZoom) / 100);
			editor.updateOptions({ fontSize: newFontSize });
		}
	}
	
	// Update base font size when fontSize prop changes
	$effect(() => {
		baseFontSize = fontSize;
		updateEditorFontSize();
	});
	
	// Export zoom functions for external use
	export { formatDocument, zoomIn, zoomOut, resetZoom };
	
	onDestroy(() => {
		// Clear this editor from global references if it's active
		if (window.activeMonacoEditorId === editorId) {
			window.activeMonacoEditorId = null;
			window.activeMonacoEditor = null;
		}
		
		if (resizeObserver) {
			resizeObserver.disconnect();
		}
		if (editor) {
			editor.dispose();
		}
		// Clean up window resize listener
		window.removeEventListener('resize', handleWindowResize);
	});
	
	// Update theme when it changes
	$effect(() => {
		if (editor) {
			monaco.editor.setTheme(theme === 'dark' ? 'custom-dark' : 'custom-light');
		}
	});
	
	// Handle external value changes (avoid binding issues)
	$effect(() => {
		if (editor && editor.getValue() !== value) {
			isUpdatingFromExternal = true;
			editor.setValue(value || '');
			updatePlaceholder();
			isUpdatingFromExternal = false;
		}
	});
</script>

<style>
	:global(.placeholder-text) {
		color: #6b7280 !important;
		font-style: italic;
		opacity: 0.7;
	}
	
	.editor-container {
		position: relative;
	}
	
	.zoom-indicator {
		position: absolute;
		top: 8px;
		right: 8px;
		background: rgba(0, 0, 0, 0.7);
		color: white;
		padding: 4px 8px;
		border-radius: 4px;
		font-size: 11px;
		font-family: monospace;
		z-index: 10;
		user-select: none;
		pointer-events: none;
		transition: opacity 0.2s ease;
	}
	
	.zoom-indicator.light {
		background: rgba(255, 255, 255, 0.9);
		color: #333;
		border: 1px solid rgba(0, 0, 0, 0.1);
	}
</style>

<div class="editor-container" style="height: {height}; width: 100%;">
	<div bind:this={container} style="height: 100%; width: 100%;" class="border theme-border rounded min-h-0"></div>
	
	<!-- Zoom level indicator -->
	{#if currentZoom !== 100}
		<div class="zoom-indicator {theme === 'light' ? 'light' : ''}" 
			 title="Current zoom level. Use Ctrl+Plus/Minus to zoom (when focused), Ctrl+0 to reset">
			{currentZoom}%
		</div>
	{/if}
</div>