<script>
	import { createEventDispatcher } from 'svelte';
	import { UpdateRestRequest } from '$lib/wailsjs/go/main/App';
	import { tabStore } from '$lib/stores/tabStore.js';
	import Toast from '$lib/Toast.svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		tab,
		isActive = false,
		showCloseButton = true
	} = $props();
	
	let isEditing = $state(false);
	let editingName = $state('');
	let inputElement = $state(null);
	
	// Reset editing state when tab becomes inactive
	$effect(() => {
		if (!isActive && isEditing) {
			cancelNameEdit();
		}
	});
	
	// Toast state
	let toastVisible = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	
	function showToast(message, type = 'success', duration = 1500, animation = 'fade') {
		toastMessage = message;
		toastType = type;
		toastVisible = true;
		
		setTimeout(() => {
			toastVisible = false;
		}, duration);
	}
	
	function handleTabClick() {
		if (!isEditing) {
			dispatch('select', tab.id);
		}
	}
	
	function handleTabKeydown(event) {
		if (!isEditing && (event.key === 'Enter' || event.key === ' ')) {
			event.preventDefault();
			dispatch('select', tab.id);
		}
	}
	
	function handleCloseClick(event) {
		event.stopPropagation(); // Prevent tab selection when closing
		dispatch('close', tab.id);
	}
	
	function handleTitleDoubleClick(event) {
		event.stopPropagation();
		if (tab.data.requestId) { // Only allow editing if it's a saved request
			isEditing = true;
			editingName = tab.title; // Use current tab title
			setTimeout(() => {
				if (inputElement) {
					inputElement.focus();
					inputElement.select();
				}
			}, 0);
		}
	}
	
	function handleNameKeydown(event) {
		if (event.key === 'Enter') {
			saveNameChange();
		} else if (event.key === 'Escape') {
			cancelNameEdit();
		}
		event.stopPropagation();
	}
	
	function handleNameBlur() {
		saveNameChange();
	}
	
	async function saveNameChange() {
		if (!isEditing) return;
		
		const trimmedName = editingName.trim();
		if (trimmedName && trimmedName !== tab.title && tab.data.requestId) {
			try {
				// First update the backend
				await UpdateRestRequest(tab.data.requestId, {
					name: trimmedName
				});
				
				// Then update the tab store - use the current tab's ID explicitly
				tabStore.updateTabTitle(tab.id, trimmedName);
				
				// Dispatch the event to refresh collections sidebar
				dispatch('nameChanged', { 
					tabId: tab.id, 
					requestId: tab.data.requestId, 
					newName: trimmedName 
				});
				
				showToast('Request name updated successfully', 'success');
			} catch (error) {
				console.error('Failed to update request name:', error);
				showToast(`Failed to update request name: ${error.message || error}`, 'error');
			}
		}
		
		isEditing = false;
		editingName = '';
	}
	
	function cancelNameEdit() {
		isEditing = false;
		editingName = '';
	}
	
	function handleSaveClick(event) {
		event.stopPropagation();
		dispatch('save', tab.id);
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
	<span class="tab-title text-sm font-medium flex items-center" ondblclick={handleTitleDoubleClick} role="button" tabindex="-1">
		{#if isEditing}
			<input
				bind:this={inputElement}
				bind:value={editingName}
				onkeydown={handleNameKeydown}
				onblur={handleNameBlur}
				onclick={(e) => e.stopPropagation()}
				class="bg-transparent border-none outline-none text-sm font-medium w-full min-w-0"
				placeholder="Request name"
			/>
		{:else}
			<span class="truncate" title={tab.data.requestId ? 'Double-click to edit name' : tab.title}>
				{tab.title}
			</span>
		{/if}
		{#if tab.isModified}
			<span class="modified-indicator ml-1 w-2 h-2 bg-orange-500 rounded-full" title="Unsaved changes"></span>
		{/if}
	</span>
	
	<!-- Save button for unsaved new requests -->
	{#if tab.isModified && !tab.data.requestId}
		<button
			class="save-button ml-2 w-5 h-5 flex items-center justify-center rounded transition-colors"
			onclick={handleSaveClick}
			title="Save request (Ctrl+S)"
			aria-label="Save request"
		>
			<svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="m19 21-7-4-7 4V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"></path>
			</svg>
		</button>
	{/if}
	
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
	
	.save-button {
		opacity: 0;
		transition: opacity 0.2s ease;
		color: var(--text-secondary);
	}
	
	.save-button:hover {
		background-color: #dcfce7; /* green-100 */
		color: #16a34a; /* green-600 */
	}
	
	:root.dark .save-button:hover {
		background-color: #14532d; /* green-900 */
		color: #4ade80; /* green-400 */
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
	
	.tab-item:hover .save-button {
		opacity: 1;
	}
	
	.modified-indicator {
		flex-shrink: 0;
	}
</style>

<!-- Toast notification -->
{#if toastVisible}
	<Toast 
		message={toastMessage} 
		type={toastType}
		visible={toastVisible}
		onHide={() => toastVisible = false}
	/>
{/if}