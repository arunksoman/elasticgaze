<script>
	import { createEventDispatcher } from 'svelte';
	import Tab from './Tab.svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		tabs = [],
		activeTabId = null,
		showAddButton = true,
		maxTabs = 10
	} = $props();
	
	function handleTabSelect(event) {
		dispatch('tabSelect', event.detail);
	}
	
	function handleTabClose(event) {
		dispatch('tabClose', event.detail);
	}
	
	function handleAddTab() {
		if (tabs.length < maxTabs) {
			dispatch('tabAdd');
		}
	}
</script>

<div class="tab-bar flex items-center border-b theme-border theme-bg-secondary">
	<!-- Tab list container with scroll -->
	<div class="tab-list flex overflow-x-auto flex-1 min-w-0">
		{#each tabs as tab (tab.id)}
			<Tab 
				{tab}
				isActive={tab.id === activeTabId}
				showCloseButton={tabs.length > 1}
				on:select={handleTabSelect}
				on:close={handleTabClose}
			/>
		{/each}
	</div>
	
	<!-- Add new tab button -->
	{#if showAddButton && tabs.length < maxTabs}
		<button
			class="add-tab-button px-3 py-2 theme-text-secondary hover:theme-text-primary transition-colors rounded"
			onclick={handleAddTab}
			title="Add new tab"
			aria-label="Add new tab"
		>
			<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<line x1="12" y1="5" x2="12" y2="19"></line>
				<line x1="5" y1="12" x2="19" y2="12"></line>
			</svg>
		</button>
	{/if}
</div>

<style>
	.tab-bar {
		height: 40px;
		flex-shrink: 0;
	}
	
	.tab-list {
		scrollbar-width: none; /* Firefox */
		-ms-overflow-style: none; /* IE/Edge */
	}
	
	.tab-list::-webkit-scrollbar {
		display: none; /* Chrome/Safari */
	}
	
	.add-tab-button {
		flex-shrink: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		border-left: 1px solid var(--border-color);
	}
	
	.add-tab-button:hover {
		background-color: var(--hover-bg);
	}
</style>