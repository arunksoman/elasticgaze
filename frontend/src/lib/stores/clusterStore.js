import { writable } from 'svelte/store';

// Store for the currently selected cluster configuration
export const selectedCluster = writable(null);

// Store for cluster dashboard data
export const clusterDashboardData = writable(null);

// Store for loading state
export const dashboardLoading = writable(false);

// Store for error state
export const dashboardError = writable(null);

// Function to update selected cluster and refresh dashboard data
export function updateSelectedCluster(config) {
	selectedCluster.set(config);
	// The dashboard components will react to this change
}

// Function to clear dashboard data
export function clearDashboardData() {
	clusterDashboardData.set(null);
	dashboardError.set(null);
}