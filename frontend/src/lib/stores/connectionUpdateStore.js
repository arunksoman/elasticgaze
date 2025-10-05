import { writable } from 'svelte/store';

/**
 * Store to trigger connection data refreshes across components
 * This is used to communicate between the connections page and layout
 * when connections are modified
 */
export const connectionUpdateTrigger = writable(0);

/**
 * Triggers a connection update refresh across all listening components
 * Call this after modifying connections to ensure all UI elements update
 */
export function triggerConnectionUpdate() {
	connectionUpdateTrigger.update(n => n + 1);
}