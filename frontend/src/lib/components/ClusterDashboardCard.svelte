<script>
	import { onMount } from 'svelte';
	import { GetClusterDashboardData, GetClusterDashboardDataByConfig } from '$lib/wailsjs/go/main/App';
	import { selectedCluster } from '$lib/stores/clusterStore.js';

	let isFlipped = $state(false);
	let dashboardData = $state(null);
	let loading = $state(true);
	let error = $state(null);

	// Subscribe to selected cluster changes
	const currentCluster = $derived($selectedCluster);

	// Health icon mapping for dashboard card
	const healthIcons = {
		green: '/icons/check-circle.svg',
		yellow: '/icons/warning.svg',
		red: '/icons/close-circle.svg'
	};

	// Health status color mapping for dropdown dots
	function getHealthColor(status) {
		switch (status) {
			case 'green': return 'bg-green-500';
			case 'yellow': return 'bg-yellow-500';
			case 'red': return 'bg-red-500';
			default: return 'bg-red-500';
		}
	}

	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
	}

	async function loadDashboardData() {
		if (!hasWails()) return;
		
		try {
			loading = true;
			error = null;
			
			if (currentCluster) {
				dashboardData = await GetClusterDashboardDataByConfig(currentCluster.id);
			} else {
				dashboardData = await GetClusterDashboardData();
			}
		} catch (err) {
			console.error('Failed to load dashboard data:', err);
			error = err.message || 'Failed to load cluster data';
		} finally {
			loading = false;
		}
	}

	function flipCard() {
		isFlipped = !isFlipped;
	}

	function getHealthIcon(status) {
		return healthIcons[status] || healthIcons.red;
	}

	function getHealthIconClass(status) {
		switch (status) {
			case 'green': return 'text-green-500';
			case 'yellow': return 'text-yellow-500';
			case 'red': return 'text-red-500';
			default: return 'text-red-500';
		}
	}

	function getHealthIconStyle(status) {
		switch (status) {
			case 'green': return 'filter: invert(36%) sepia(91%) saturate(1752%) hue-rotate(89deg) brightness(103%) contrast(107%);'; // Matches text-green-500 (#10b981)
			case 'yellow': return 'filter: invert(71%) sepia(87%) saturate(1493%) hue-rotate(2deg) brightness(103%) contrast(105%);'; // Matches text-yellow-500 (#eab308)
			case 'red': return 'filter: invert(27%) sepia(96%) saturate(5314%) hue-rotate(343deg) brightness(101%) contrast(107%);'; // Matches text-red-500 (#ef4444)
			default: return 'filter: invert(27%) sepia(96%) saturate(5314%) hue-rotate(343deg) brightness(101%) contrast(107%);';
		}
	}

	function formatBytes(bytes) {
		if (!bytes) return '0 B';
		const units = ['B', 'KB', 'MB', 'GB', 'TB'];
		let size = bytes;
		let unitIndex = 0;
		
		while (size >= 1024 && unitIndex < units.length - 1) {
			size /= 1024;
			unitIndex++;
		}
		
		return `${size.toFixed(1)} ${units[unitIndex]}`;
	}

	// Reactive effect to reload data when selectedCluster changes
	$effect(() => {
		loadDashboardData();
	});

	onMount(() => {
		loadDashboardData();
	});
</script>

<div class="w-full max-w-6xl mx-auto mt-8">
	<!-- Simple Card with Toggle -->
	<div class="overflow-hidden">
		
		<!-- Card Header -->
		<div class="border-b theme-border px-6 py-4 flex items-center justify-between">
			<div class="flex items-center">
				{#if dashboardData?.cluster_health?.status}
					<img 
						src={getHealthIcon(dashboardData.cluster_health.status)} 
						alt="Cluster health" 
						class="w-5 h-5 mr-3"
						style={getHealthIconStyle(dashboardData.cluster_health.status)}
					/>
				{/if}
				<h2 class="text-xl font-bold theme-text-primary">
					{dashboardData?.cluster_info?.cluster_name || 'Cluster Dashboard'}
				</h2>
				{#if dashboardData?.cluster_health?.status}
					<span class="ml-3 px-2 py-1 text-xs font-medium rounded theme-bg-tertiary {getHealthIconClass(dashboardData.cluster_health.status)} capitalize">
						{dashboardData.cluster_health.status}
					</span>
				{/if}
			</div>
			
			<button 
				onclick={flipCard}
				class="px-3 py-1 text-sm theme-text-secondary hover:theme-text-primary transition-colors duration-200 border theme-border rounded hover:theme-bg-tertiary lg:hidden"
			>
				{isFlipped ? 'Show Summary' : 'Show Details'}
			</button>
		</div>

		<!-- Card Content -->
		<div class="p-6">
			{#if loading}
				<div class="flex items-center justify-center h-64">
					<div class="text-center theme-text-secondary">
						<div class="w-8 h-8 mx-auto mb-4 animate-spin">
							<div class="w-full h-full border-2 border-current border-r-transparent rounded-full"></div>
						</div>
						<p>Loading cluster data...</p>
					</div>
				</div>
			{:else if error}
				<div class="flex items-center justify-center h-64">
					<div class="text-center theme-text-secondary">
						<img src="/icons/warning.svg" alt="Error" class="w-8 h-8 mx-auto mb-4 theme-icon" />
						<p class="mb-2">Failed to load cluster data</p>
						<p class="text-sm theme-text-muted">{error}</p>
					</div>
				</div>
			{:else if dashboardData}
				
				<!-- Summary View - Always visible on large screens, toggle on small screens -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 {isFlipped ? 'hidden lg:grid' : 'grid'}">
						
						<!-- Cluster Info Card -->
						<div class="border theme-border rounded-lg p-4">
							<div class="flex items-center mb-3">
								<img src="/icons/about.svg" alt="Info" class="w-5 h-5 mr-2 theme-icon" />
								<h3 class="font-semibold theme-text-primary">Cluster Info</h3>
							</div>
							<div class="space-y-2">
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide">Name</p>
									<p class="text-sm font-medium theme-text-primary">{dashboardData.cluster_info?.cluster_name || 'Unknown'}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide">Health</p>
									<div class="flex items-center">
										<img 
											src={getHealthIcon(dashboardData.cluster_health?.status)} 
											alt="Health" 
											class="w-4 h-4 mr-2"
											style={getHealthIconStyle(dashboardData.cluster_health?.status)}
										/>
										<p class="text-sm font-medium {getHealthIconClass(dashboardData.cluster_health?.status)} capitalize">{dashboardData.cluster_health?.status || 'Unknown'}</p>
									</div>
								</div>
							</div>
						</div>

						<!-- Nodes Card -->
						<div class="border theme-border rounded-lg p-4">
							<div class="flex items-center mb-3">
								<img src="/icons/nodes.svg" alt="Nodes" class="w-5 h-5 mr-2 theme-icon" />
								<h3 class="font-semibold theme-text-primary">Nodes</h3>
							</div>
							<div class="space-y-2">
								<div class="flex justify-between">
									<span class="text-sm theme-text-secondary">Master:</span>
									<span class="text-sm font-medium theme-text-primary">{dashboardData.node_counts?.master || 0}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-sm theme-text-secondary">Data:</span>
									<span class="text-sm font-medium theme-text-primary">{dashboardData.node_counts?.data || 0}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-sm theme-text-secondary">Ingest:</span>
									<span class="text-sm font-medium theme-text-primary">{dashboardData.node_counts?.ingest || 0}</span>
								</div>
								<div class="flex justify-between border-t theme-border pt-2">
									<span class="text-sm font-medium theme-text-secondary">Total:</span>
									<span class="text-lg font-bold theme-text-primary">{dashboardData.node_counts?.total || 0}</span>
								</div>
							</div>
						</div>

						<!-- Shards Card -->
						<div class="border theme-border rounded-lg p-4">
							<div class="flex items-center mb-3">
								<img src="/icons/shards.svg" alt="Shards" class="w-5 h-5 mr-2 theme-icon" />
								<h3 class="font-semibold theme-text-primary">Shards</h3>
							</div>
							<div class="space-y-2">
								<div class="flex justify-between">
									<span class="text-sm theme-text-secondary">Primary:</span>
									<span class="text-sm font-medium theme-text-primary">{dashboardData.shard_counts?.primary || 0}</span>
								</div>
								<div class="flex justify-between">
									<span class="text-sm theme-text-secondary">Replica:</span>
									<span class="text-sm font-medium theme-text-primary">{dashboardData.shard_counts?.replica || 0}</span>
								</div>
								<div class="flex justify-between border-t theme-border pt-2">
									<span class="text-sm font-medium theme-text-secondary">Total:</span>
									<span class="text-lg font-bold theme-text-primary">{dashboardData.shard_counts?.total || 0}</span>
								</div>
							</div>
						</div>

						<!-- Indices Card -->
						<div class="border theme-border rounded-lg p-4">
							<div class="flex items-center mb-3">
								<img src="/icons/index.svg" alt="Indices" class="w-5 h-5 mr-2 theme-icon" />
								<h3 class="font-semibold theme-text-primary">Indices</h3>
							</div>
							<div class="space-y-2">
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide">Documents</p>
									<p class="text-lg font-bold theme-text-primary">{dashboardData.index_metrics?.document_count?.toLocaleString() || '0'}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide">Storage</p>
									<p class="text-lg font-bold theme-text-primary">{dashboardData.index_metrics?.disk_usage || '0 B'}</p>
								</div>
							</div>
						</div>

					</div>
				
				<!-- Detailed View - Toggle on small screens, always visible on large screens -->
				<div class="space-y-6 max-h-96 overflow-y-auto {isFlipped ? 'block' : 'hidden lg:block'} {!isFlipped ? 'lg:mt-8' : ''}">
						
						<!-- Cluster Information -->
						<div class="border-b-2 theme-border pb-6">
							<h3 class="text-lg font-semibold theme-text-primary mb-4 flex items-center">
								<img src="/icons/about.svg" alt="Cluster" class="w-5 h-5 mr-2 theme-icon" />
								Cluster Information
							</h3>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4 border theme-border rounded-lg p-4">
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Node Name</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_info?.name || 'N/A'}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Cluster Name</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_info?.cluster_name || 'N/A'}</p>
								</div>
								<div class="md:col-span-2">
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Cluster UUID</p>
									<p class="text-xs font-mono theme-text-secondary break-all">{dashboardData.cluster_info?.cluster_uuid || 'N/A'}</p>
								</div>
								<div class="md:col-span-2">
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Tagline</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_info?.tagline || 'N/A'}</p>
								</div>
							</div>
						</div>

						<!-- Version Information -->
						{#if dashboardData.cluster_info?.version}
							<div>
								<h3 class="text-lg font-semibold theme-text-primary mb-4 flex items-center">
									<img src="/icons/settings.svg" alt="Version" class="w-5 h-5 mr-2 theme-icon" />
									Version Information
								</h3>
								<div class="grid grid-cols-1 md:grid-cols-2 gap-4 border-2 theme-border rounded-lg p-4">
									<div>
										<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Version</p>
										<p class="text-sm theme-text-primary">{dashboardData.cluster_info.version.number}</p>
									</div>
									<div>
										<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Build Type</p>
										<p class="text-sm theme-text-primary">{dashboardData.cluster_info.version.build_type}</p>
									</div>
									<div>
										<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Lucene Version</p>
										<p class="text-sm theme-text-primary">{dashboardData.cluster_info.version.lucene_version}</p>
									</div>
									<div>
										<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Build Date</p>
										<p class="text-sm theme-text-primary">{dashboardData.cluster_info.version.build_date}</p>
									</div>
								</div>
							</div>
						{/if}

						<!-- Cluster Health -->
						<div>
							<h3 class="text-lg font-semibold theme-text-primary mb-4 flex items-center">
								<img 
									src={getHealthIcon(dashboardData.cluster_health?.status)} 
									alt="Health" 
									class="w-5 h-5 mr-2"
									style={getHealthIconStyle(dashboardData.cluster_health?.status)}
								/>
								Cluster Health
							</h3>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4 border theme-border rounded-lg p-4">
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Status</p>
									<div class="flex items-center">
										<img 
											src={getHealthIcon(dashboardData.cluster_health?.status)} 
											alt="Health" 
											class="w-4 h-4 mr-2"
											style={getHealthIconStyle(dashboardData.cluster_health?.status)}
										/>
										<p class="text-sm {getHealthIconClass(dashboardData.cluster_health?.status)} font-medium capitalize">{dashboardData.cluster_health?.status || 'unknown'}</p>
									</div>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Active Shards</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_health?.active_shards || 0}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Relocating Shards</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_health?.relocating_shards || 0}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Unassigned Shards</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_health?.unassigned_shards || 0}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Pending Tasks</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_health?.number_of_pending_tasks || 0}</p>
								</div>
								<div>
									<p class="text-xs theme-text-muted uppercase tracking-wide mb-1">Active Shards %</p>
									<p class="text-sm theme-text-primary">{dashboardData.cluster_health?.active_shards_percent_as_number?.toFixed(1) || 0}%</p>
								</div>
							</div>
						</div>

					</div>
			{/if}
		</div>

	</div>
</div>

<style>
	/* Simple responsive grid adjustments */
	@media (max-width: 768px) {
		.grid.lg\\:grid-cols-4 {
			grid-template-columns: repeat(1, minmax(0, 1fr));
		}
	}
</style>