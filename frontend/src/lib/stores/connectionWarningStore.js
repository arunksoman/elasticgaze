import { writable } from 'svelte/store';
import { GetDefaultConfig } from '$lib/wailsjs/go/main/App';

// Connection status store for UI state management
export const connectionWarningStatus = writable({
	hasDefault: true,
	isChecking: true,
	lastChecked: 0
});

/**
 * Function to update connection status
 * @param {boolean} hasDefault - Whether a default connection exists
 * @param {boolean} isChecking - Whether we're currently checking
 */
export function updateConnectionWarningStatus(hasDefault, isChecking = false) {
	connectionWarningStatus.set({
		hasDefault: Boolean(hasDefault),
		isChecking: Boolean(isChecking),
		lastChecked: Date.now()
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
		await GetDefaultConfig();
		updateConnectionWarningStatus(true, false);
	} catch (error) {
		updateConnectionWarningStatus(false, false);
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