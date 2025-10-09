<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		tab,
		isActive = false,
		showCloseButton = true
	} = $props();
	
	function handleTabClick() {
		dispatch('select', tab.id);
	}
	
	function handleTabKeydown(event) {
		if (event.key === 'Enter' || event.key === ' ') {
			event.preventDefault();
			dispatch('select', tab.id);
		}
	}
	
	function handleCloseClick(event) {
		event.stopPropagation(); // Prevent tab selection when closing
		dispatch('close', tab.id);
	}
</script>

<div 
	class="tab-item flex items-center px-4 py-2 cursor-pointer border-b-2 transition-colors duration-200 {isActive ? 'active-tab' : 'inactive-tab'}"
	onclick={handleTabClick}
	onkeydown={handleTabKeydown}
	role="tab"
	tabindex="0"
	aria-selected={isActive}
>
	<!-- Tab title with modified indicator -->
	<span class="tab-title text-sm font-medium flex items-center">
		{tab.title}
		{#if tab.isModified}
			<span class="modified-indicator ml-1 w-2 h-2 bg-orange-500 rounded-full" title="Unsaved changes"></span>
		{/if}
	</span>
	
	<!-- Close button -->
	{#if showCloseButton}
		<button
			class="close-button ml-2 w-5 h-5 flex items-center justify-center rounded transition-colors"
			onclick={handleCloseClick}
			title="Close tab"
			aria-label="Close tab"
		>
			<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<line x1="18" y1="6" x2="6" y2="18"></line>
				<line x1="6" y1="6" x2="18" y2="18"></line>
			</svg>
		</button>
	{/if}
</div>

<style>
	.tab-item {
		min-width: 120px;
		max-width: 200px;
		position: relative;
		user-select: none;
	}
	
	.active-tab {
		border-color: #3b82f6; /* blue-500 */
		background-color: var(--bg-secondary);
		color: var(--text-primary);
	}
	
	.inactive-tab {
		border-color: transparent;
		color: var(--text-secondary);
	}
	
	.inactive-tab:hover {
		background-color: var(--hover-bg);
		color: var(--text-primary);
	}
	
	.tab-title {
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		flex: 1;
	}
	
	.close-button {
		opacity: 0;
		transition: opacity 0.2s ease;
		color: var(--text-secondary);
	}
	
	.close-button:hover {
		background-color: #fecaca; /* red-200 */
		color: #dc2626; /* red-600 */
	}
	
	:root.dark .close-button:hover {
		background-color: #7f1d1d; /* red-900 */
		color: #f87171; /* red-400 */
	}
	
	.tab-item:hover .close-button {
		opacity: 1;
	}
	
	.modified-indicator {
		flex-shrink: 0;
	}
</style>