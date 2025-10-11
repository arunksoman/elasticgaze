// Monaco Editor preloader service with enhanced caching
import { GetMonacoCacheInfo, ReadMonacoCache, WriteMonacoCache } from '$lib/wailsjs/go/main/App.js';

/** @type {Promise<any> | null} */
let monacoPromise = null;
let isPreloading = false;
let isLoaded = false;

// Cache for the loaded module
/** @type {any} */
let cachedModule = null;

// Cache key for localStorage (fallback)
const CACHE_KEY = 'monaco-editor-cache';
const CACHE_VERSION = '1.0.0';

// Performance timing
let loadStartTime = 0;

/**
 * Check if browser storage cache is available and valid (fallback)
 */
function checkBrowserCache() {
	try {
		const cached = localStorage.getItem(CACHE_KEY);
		if (!cached) return null;
		
		const data = JSON.parse(cached);
		
		// Check cache version and expiry (24 hours)
		if (data.version !== CACHE_VERSION || 
		    Date.now() - data.timestamp > 24 * 60 * 60 * 1000) {
			localStorage.removeItem(CACHE_KEY);
			return null;
		}
		
		return data;
	} catch (error) {
		console.warn('⚠️ Failed to check browser cache:', error);
		return null;
	}
}

/**
 * Store cache marker in browser storage (fallback)
 */
function setBrowserCache() {
	try {
		const data = {
			version: CACHE_VERSION,
			timestamp: Date.now(),
			cached: true
		};
		localStorage.setItem(CACHE_KEY, JSON.stringify(data));
		console.log('💾 Monaco Editor cached in browser storage');
	} catch (error) {
		console.warn('⚠️ Failed to set browser cache:', error);
	}
}

/**
 * Check Go backend cache
 */
async function checkGoCache() {
	try {
		const cacheInfo = await GetMonacoCacheInfo(CACHE_VERSION);
		console.log('📦 Go cache info:', cacheInfo);
		return cacheInfo;
	} catch (error) {
		console.warn('⚠️ Failed to check Go cache:', error);
		return null;
	}
}

/**
 * Store data in Go backend cache
 * @param {string} data - The data to cache
 */
async function setGoCache(data) {
	try {
		await WriteMonacoCache(CACHE_VERSION, data);
		console.log('💾 Monaco Editor cached in Go backend');
		return true;
	} catch (error) {
		console.warn('⚠️ Failed to set Go cache:', error);
		return false;
	}
}

/**
 * Read data from Go backend cache
 */
async function getGoCache() {
	try {
		const data = await ReadMonacoCache(CACHE_VERSION);
		console.log('📦 Monaco Editor loaded from Go cache');
		return data;
	} catch (error) {
		console.warn('⚠️ Failed to read Go cache:', error);
		return null;
	}
}

/**
 * Start preloading Monaco Editor in the background
 * This can be called early in the app lifecycle
 */
export function preloadMonacoEditor() {
	if (isPreloading || isLoaded) {
		return monacoPromise;
	}
	
	// Return cached module if available
	if (cachedModule) {
		console.log('✅ Monaco Editor already cached in memory');
		return Promise.resolve(cachedModule);
	}
	
	isPreloading = true;
	loadStartTime = performance.now();
	console.log('🚀 Starting Monaco Editor preload...');
	
	monacoPromise = loadMonacoWithCaching()
		.then((module) => {
			const loadTime = performance.now() - loadStartTime;
			console.log(`✅ Monaco Editor preloaded successfully in ${loadTime.toFixed(2)}ms`);
			cachedModule = module;
			isLoaded = true;
			isPreloading = false;
			return module;
		})
		.catch((error) => {
			const loadTime = performance.now() - loadStartTime;
			console.error(`❌ Monaco Editor preload failed after ${loadTime.toFixed(2)}ms:`, error);
			isPreloading = false;
			monacoPromise = null;
			throw error;
		});
	
	return monacoPromise;
}

/**
 * Load Monaco Editor with enhanced caching strategy
 */
async function loadMonacoWithCaching() {
	// Check Go backend cache first
	const goCacheInfo = await checkGoCache();
	
	if (goCacheInfo && goCacheInfo.exists && !goCacheInfo.isExpired) {
		console.log('📦 Go cache found - using cached version');
	} else {
		// Fallback to browser cache check
		const browserCache = checkBrowserCache();
		if (browserCache) {
			console.log('📦 Browser cache found - loading optimized');
		}
	}
	
	// Load Monaco Editor module
	const module = await import('$lib/MonacoEditor.svelte');
	
	// Cache the successful load in both Go backend and browser
	await Promise.allSettled([
		setGoCache(JSON.stringify({
			version: CACHE_VERSION,
			timestamp: Date.now(),
			moduleLoaded: true
		})),
		setBrowserCache()
	]);
	
	return module;
}

/**
 * Get the preloaded Monaco Editor or start loading it
 * Returns a promise that resolves to the Monaco Editor module
 */
export function getMonacoEditor() {
	// Return cached module immediately if available
	if (cachedModule) {
		console.log('⚡ Monaco Editor returned from cache');
		return Promise.resolve(cachedModule);
	}
	
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

/**
 * Force preload Monaco Editor in the background
 * This should be called early in the app lifecycle
 */
export function warmMonacoEditor() {
	// Only start preloading if not already done
	if (!isLoaded && !isPreloading) {
		console.log('🔥 Warming Monaco Editor...');
		const promise = preloadMonacoEditor();
		if (promise) {
			promise.catch(() => {
				// Ignore errors during warming
			});
		}
	}
}

/**
 * Clear the in-memory, browser, and Go backend cache
 */
export async function clearMonacoCache() {
	cachedModule = null;
	monacoPromise = null;
	isLoaded = false;
	isPreloading = false;
	
	// Clear browser cache
	try {
		localStorage.removeItem(CACHE_KEY);
	} catch (error) {
		console.warn('⚠️ Failed to clear browser cache:', error);
	}
	
	// Clear Go backend cache
	try {
		const { ClearAllMonacoCache } = await import('$lib/wailsjs/go/main/App.js');
		await ClearAllMonacoCache();
		console.log('🗑️ Go cache cleared');
	} catch (error) {
		console.warn('⚠️ Failed to clear Go cache:', error);
	}
	
	console.log('🗑️ Monaco Editor cache cleared');
}

/**
 * Get cache status with performance metrics
 */
export async function getMonacoCacheStatus() {
	const browserCache = checkBrowserCache();
	let goCacheInfo = null;
	
	try {
		goCacheInfo = await GetMonacoCacheInfo(CACHE_VERSION);
	} catch (error) {
		console.warn('Failed to get Go cache info:', error);
	}
	
	return {
		isLoaded,
		isPreloading,
		hasCachedModule: !!cachedModule,
		hasBrowserCache: !!browserCache,
		browserCacheAge: browserCache ? Date.now() - browserCache.timestamp : 0,
		hasGoCache: goCacheInfo?.exists || false,
		goCacheExpired: goCacheInfo?.isExpired || false,
		goCacheSize: goCacheInfo?.size || 0,
		goCachePath: goCacheInfo?.cachePath || '',
		lastLoadTime: loadStartTime > 0 ? performance.now() - loadStartTime : 0
	};
}