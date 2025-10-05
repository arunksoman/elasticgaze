<script>
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { GetAllConfigs, GetClusterHealthForAllConfigs, GetDefaultConfig, GetClusterDashboardDataByConfig } from '$lib/wailsjs/go/main/App';
	import { selectedCluster } from '$lib/stores/clusterStore.js';

	let showDropdown = $state(false);
	let configs = $state([]);
	let healthMap = $state({});
	let clusterNamesMap = $state({});
	let selectedConfig = $state(null);
	let loading = $state(true);
	let lastUpdatePath = $state('');
	let isRefreshing = $state(false);

	// Health status color mapping
	function getHealthColor(configName) {
		const health = healthMap[configName] || 'red';
		switch (health) {
			case 'green': return 'bg-green-500';
			case 'yellow': return 'bg-yellow-500';
			case 'red': return 'bg-red-500';
			default: return 'bg-red-500';
		}
	}

	// Reactive effect to refresh health data when navigating between pages
	$effect(() => {
		const currentPath = page.url.pathname;
		// Only refresh health data when navigating away from and back to home page
		// and avoid refreshing on initial load or same page
		if (lastUpdatePath !== '' && 
			lastUpdatePath !== currentPath && 
			currentPath === '/' && 
			hasWails() && 
			configs.length > 0) {
			// Add a small delay to prevent rapid successive calls
			setTimeout(() => {
				refreshHealthData();
			}, 100);
		}
		lastUpdatePath = currentPath;
	});

	// Reactive effect to sync with selectedCluster store changes
	$effect(() => {
		const storeCluster = $selectedCluster;
		if (storeCluster && storeCluster.id !== selectedConfig?.id) {
			selectedConfig = storeCluster;
		}
	});

	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
	}

	async function refreshHealthData() {
		if (!hasWails() || configs.length === 0 || isRefreshing) return;
		
		try {
			isRefreshing = true;
			// Refresh health status for all configs
			const health = await GetClusterHealthForAllConfigs();
			healthMap = health || {};
		} catch (err) {
			console.warn('Failed to refresh cluster health:', err);
		} finally {
			isRefreshing = false;
		}
	}

	async function loadClusterNames() {
		// Get actual cluster names from Elasticsearch for each config
		for (const config of configs) {
			try {
				const dashboardData = await GetClusterDashboardDataByConfig(config.id);
				if (dashboardData?.cluster_info?.cluster_name) {
					clusterNamesMap[config.id] = dashboardData.cluster_info.cluster_name;
				} else {
					clusterNamesMap[config.id] = config.connection_name;
				}
			} catch (err) {
				console.warn(`Failed to get cluster name for ${config.connection_name}:`, err);
				clusterNamesMap[config.id] = config.connection_name;
			}
		}
		clusterNamesMap = { ...clusterNamesMap }; // Trigger reactivity
	}

	async function loadConfigs() {
		if (!hasWails()) return;
		
		try {
			loading = true;
			// Load all configurations
			const allConfigs = await GetAllConfigs();
			configs = allConfigs || [];

			// Load health status for all configs
			try {
				const health = await GetClusterHealthForAllConfigs();
				healthMap = health || {};
			} catch (err) {
				console.warn('Failed to load cluster health:', err);
				healthMap = {};
			}

			// Load cluster names from Elasticsearch
			await loadClusterNames();

			// Get default config
			try {
				const defaultConfig = await GetDefaultConfig();
				selectedConfig = defaultConfig;
				selectedCluster.set(defaultConfig);
			} catch (err) {
				// No default config found, select first one if available
				if (configs.length > 0) {
					selectedConfig = configs[0];
					selectedCluster.set(configs[0]);
				}
			}
		} catch (error) {
			console.error('Failed to load cluster configurations:', error);
		} finally {
			loading = false;
		}
	}

	function selectCluster(config) {
		selectedConfig = config;
		selectedCluster.set(config);
		showDropdown = false;
	}

	function toggleDropdown() {
		showDropdown = !showDropdown;
	}



	function getClusterName(config) {
		return clusterNamesMap[config.id] || config.connection_name;
	}

	// Close dropdown when clicking outside
	function handleClickOutside(event) {
		if (!event.target.closest('.cluster-dropdown')) {
			showDropdown = false;
		}
	}

	onMount(() => {
		loadConfigs();
		document.addEventListener('click', handleClickOutside);
		return () => document.removeEventListener('click', handleClickOutside);
	});
</script>

<div class="cluster-dropdown relative">
	{#if loading}
		<div class="flex items-center px-2 py-1.5 theme-bg-secondary theme-text-secondary border-0 rounded">
			<div class="w-3 h-3 mr-2 animate-spin">
				<div class="w-full h-full border border-current border-r-transparent rounded-full"></div>
			</div>
			<span class="text-xs">Loading...</span>
		</div>
	{:else if selectedConfig}
		<button
			class="flex items-center px-2 py-1.5 theme-bg-secondary theme-text-primary hover:theme-bg-tertiary rounded transition-colors duration-150 min-w-[160px] text-xs"
			onclick={toggleDropdown}
		>
			<div 
				class="w-2 h-2 mr-1.5 rounded-full {getHealthColor(selectedConfig.connection_name)}"
				title="Health status"
			></div>
			<span class="font-medium truncate flex-1 text-left">
				{getClusterName(selectedConfig)}
			</span>
			<svg 
				class="w-3 h-3 ml-1 theme-icon transition-transform duration-150 {showDropdown ? 'rotate-180' : ''}" 
				fill="none" 
				stroke="currentColor" 
				viewBox="0 0 24 24"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
			</svg>
		</button>

		{#if showDropdown}
			<div class="absolute top-full left-0 mt-0.5 w-full theme-bg-secondary border-0 rounded shadow-sm z-50 max-h-48 overflow-y-auto">
				{#each configs as config (config.id)}
					<button
						class="w-full flex items-center px-2 py-1.5 text-left hover:theme-bg-tertiary theme-text-primary transition-colors duration-150 text-xs {selectedConfig?.id === config.id ? 'theme-bg-tertiary' : ''}"
						onclick={() => selectCluster(config)}
					>
						<div 
							class="w-2 h-2 mr-1.5 rounded-full {getHealthColor(config.connection_name)}"
							title="Health status"
						></div>
						<div class="flex-1 min-w-0">
							<div class="font-medium truncate">{getClusterName(config)}</div>
							<div class="text-[10px] theme-text-secondary truncate">{config.connection_name} • {config.host}:{config.port}</div>
						</div>
						{#if config.set_as_default}
							<span class="ml-1 text-[10px] theme-text-secondary">(default)</span>
						{/if}
					</button>
				{/each}

				{#if configs.length === 0}
					<div class="px-2 py-1.5 text-xs theme-text-secondary">
						No clusters found
					</div>
				{/if}
			</div>
		{/if}
	{:else}
		<div class="flex items-center px-2 py-1.5 theme-bg-secondary theme-text-secondary border-0 rounded">
			<span class="text-xs">No clusters</span>
		</div>
	{/if}
</div>