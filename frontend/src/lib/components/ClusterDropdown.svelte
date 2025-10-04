<script>
	import { onMount } from 'svelte';
	import { GetAllConfigs, GetClusterHealthForAllConfigs, GetDefaultConfig, GetClusterDashboardDataByConfig } from '$lib/wailsjs/go/main/App';
	import { selectedCluster } from '$lib/stores/clusterStore.js';

	let showDropdown = $state(false);
	let configs = $state([]);
	let healthMap = $state({});
	let clusterNamesMap = $state({});
	let selectedConfig = $state(null);
	let loading = $state(true);

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

	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
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
		<div class="flex items-center px-3 py-2 theme-bg-secondary theme-text-secondary border theme-border">
			<div class="w-4 h-4 mr-2 animate-spin">
				<div class="w-full h-full border-2 border-current border-r-transparent rounded-full"></div>
			</div>
			<span class="text-sm">Loading...</span>
		</div>
	{:else if selectedConfig}
		<button
			class="flex items-center px-3 py-2 theme-bg-secondary theme-text-primary theme-hover border theme-border transition-colors duration-200 min-w-[200px]"
			onclick={toggleDropdown}
		>
			<div 
				class="w-3 h-3 mr-2 rounded-full {getHealthColor(selectedConfig.connection_name)}"
				title="Health status"
			></div>
			<span class="text-sm font-medium truncate flex-1 text-left">
				{getClusterName(selectedConfig)}
			</span>
			<svg 
				class="w-4 h-4 ml-2 theme-icon transition-transform duration-200 {showDropdown ? 'rotate-180' : ''}" 
				fill="none" 
				stroke="currentColor" 
				viewBox="0 0 24 24"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
			</svg>
		</button>

		{#if showDropdown}
			<div class="absolute top-full left-0 mt-1 w-full theme-bg-secondary border theme-border shadow-lg z-50 max-h-60 overflow-y-auto">
				{#each configs as config (config.id)}
					<button
						class="w-full flex items-center px-3 py-2 text-left theme-hover theme-text-primary transition-colors duration-200 {selectedConfig?.id === config.id ? 'theme-bg-tertiary' : ''}"
						onclick={() => selectCluster(config)}
					>
						<div 
							class="w-3 h-3 mr-2 rounded-full {getHealthColor(config.connection_name)}"
							title="Health status"
						></div>
						<div class="flex-1 min-w-0">
							<div class="text-sm font-medium truncate">{getClusterName(config)}</div>
							<div class="text-xs theme-text-secondary truncate">{config.connection_name} ({config.host}:{config.port})</div>
						</div>
						{#if config.set_as_default}
							<span class="ml-2 text-xs theme-text-secondary">(default)</span>
						{/if}
					</button>
				{/each}

				{#if configs.length === 0}
					<div class="px-3 py-2 text-sm theme-text-secondary">
						No cluster configurations found
					</div>
				{/if}
			</div>
		{/if}
	{:else}
		<div class="flex items-center px-3 py-2 theme-bg-secondary theme-text-secondary border theme-border">
			<span class="text-sm">No clusters configured</span>
		</div>
	{/if}
</div>