import { writable } from 'svelte/store';
// Note: TestDefaultConnection will be available after Wails bindings are regenerated
// import { TestDefaultConnection } from '$lib/wailsjs/go/main/App';

// Connection status store for UI state management
export const connectionWarningStatus = writable({
	hasDefault: true,
	isWorking: true,
	isChecking: true,
	lastChecked: 0,
	errorMessage: '',
	connectionState: 'unknown', // 'working', 'missing', 'failed', 'unknown'
	isStabilizing: false // Prevents rapid state changes
});

/**
 * Function to update connection status
 * @param {boolean} hasDefault - Whether a default connection exists
 * @param {boolean} isWorking - Whether the connection is working
 * @param {boolean} isChecking - Whether we're currently checking
 * @param {string} errorMessage - Error message if any
 * @param {string} connectionState - Current connection state
 * @param {boolean} isStabilizing - Whether we're stabilizing the state
 */
export function updateConnectionWarningStatus(hasDefault, isWorking = true, isChecking = false, errorMessage = '', connectionState = 'unknown', isStabilizing = false) {
	connectionWarningStatus.set({
		hasDefault: Boolean(hasDefault),
		isWorking: Boolean(isWorking),
		isChecking: Boolean(isChecking),
		lastChecked: Date.now(),
		errorMessage: String(errorMessage),
		connectionState: String(connectionState),
		isStabilizing: Boolean(isStabilizing)
	});
}

// Function to trigger connection check
export function triggerConnectionCheck() {
	connectionWarningStatus.update(status => ({
		...status,
		isChecking: true
	}));
}

// Function to check and update connection status (for use in other components)
export async function refreshConnectionStatus() {
	if (typeof window === 'undefined' || !window.runtime) return;
	
	try {
		triggerConnectionCheck();
		// Note: This will use TestDefaultConnection when Wails bindings are regenerated
		// For now, we'll use a fallback approach
		// @ts-ignore - Wails runtime methods
		if (window.go?.main?.App?.TestDefaultConnection) {
			// @ts-ignore - Wails runtime methods
			const result = await window.go.main.App.TestDefaultConnection();
			if (result.success) {
				updateConnectionWarningStatus(true, true, false, '', 'working');
			} else if (result.error_code === 'NO_DEFAULT_CONNECTION') {
				updateConnectionWarningStatus(false, false, false, result.message, 'missing');
			} else {
				updateConnectionWarningStatus(true, false, false, result.message, 'failed');
			}
		} else {
			// Fallback to old method for backward compatibility
			const { GetDefaultConfig } = await import('$lib/wailsjs/go/main/App');
			await GetDefaultConfig();
			updateConnectionWarningStatus(true, true, false, '', 'working');
		}
	} catch (error) {
		// If GetDefaultConfig fails, it likely means no default connection exists
		updateConnectionWarningStatus(false, false, false, 'No default connection found', 'missing');
	}
}

// Special function for retry operations that prevents bounce effects
export async function retryConnectionStatus() {
	if (typeof window === 'undefined' || !window.runtime) return;
	
	try {
		// Set stabilizing state to prevent UI bounce
		updateConnectionWarningStatus(false, false, true, 'Retrying...', 'failed', true);
		
		// @ts-ignore - Wails runtime methods
		if (window.go?.main?.App?.TestDefaultConnection) {
			// @ts-ignore - Wails runtime methods
			const result = await window.go.main.App.TestDefaultConnection();
			if (result.success) {
				// Success - set working state and maintain stabilizing for a moment
				updateConnectionWarningStatus(true, true, false, '', 'working', true);
				// After a delay, remove stabilizing flag
				setTimeout(() => {
					updateConnectionWarningStatus(true, true, false, '', 'working', false);
				}, 1500);
			} else {
				// Still failed - update with new error
				updateConnectionWarningStatus(true, false, false, result.message, 'failed', false);
			}
		} else {
			// Fallback method
			const { GetDefaultConfig } = await import('$lib/wailsjs/go/main/App');
			await GetDefaultConfig();
			updateConnectionWarningStatus(true, true, false, '', 'working', true);
			setTimeout(() => {
				updateConnectionWarningStatus(true, true, false, '', 'working', false);
			}, 1500);
		}
	} catch (error) {
		// Still failed
		const errorMsg = error && typeof error === 'object' && 'message' in error ? error.message : 'Connection test failed';
		updateConnectionWarningStatus(true, false, false, String(errorMsg), 'failed', false);
	}
}

// Debounced version to prevent excessive calls
/** @type {any} */
let refreshTimeout = null;

/**
 * Debounced refresh to prevent excessive API calls
 * @param {number} delay - Delay in milliseconds (default 500ms)
 */
export function debouncedRefreshConnectionStatus(delay = 500) {
	if (refreshTimeout) {
		clearTimeout(refreshTimeout);
	}
	refreshTimeout = setTimeout(() => {
		refreshConnectionStatus();
		refreshTimeout = null;
	}, delay);
}