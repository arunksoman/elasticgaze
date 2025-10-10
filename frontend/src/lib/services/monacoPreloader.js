// Monaco Editor preloader service
/** @type {Promise<any> | null} */
let monacoPromise = null;
let isPreloading = false;
let isLoaded = false;

/**
 * Start preloading Monaco Editor in the background
 * This can be called early in the app lifecycle
 */
export function preloadMonacoEditor() {
	if (isPreloading || isLoaded) {
		return monacoPromise;
	}
	
	isPreloading = true;
	console.log('🚀 Starting Monaco Editor preload...');
	
	monacoPromise = import('$lib/MonacoEditor.svelte')
		.then((module) => {
			console.log('✅ Monaco Editor preloaded successfully');
			isLoaded = true;
			isPreloading = false;
			return module;
		})
		.catch((error) => {
			console.error('❌ Monaco Editor preload failed:', error);
			isPreloading = false;
			monacoPromise = null;
			throw error;
		});
	
	return monacoPromise;
}

/**
 * Get the preloaded Monaco Editor or start loading it
 * Returns a promise that resolves to the Monaco Editor module
 */
export function getMonacoEditor() {
	if (monacoPromise) {
		return monacoPromise;
	}
	
	return preloadMonacoEditor();
}

/**
 * Check if Monaco Editor is already loaded
 */
export function isMonacoEditorLoaded() {
	return isLoaded;
}

/**
 * Check if Monaco Editor is currently being preloaded
 */
export function isMonacoEditorPreloading() {
	return isPreloading;
}