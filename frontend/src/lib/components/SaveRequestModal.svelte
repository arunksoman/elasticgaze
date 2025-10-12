<script>
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';
	import CollectionsTree from './CollectionsTree.svelte';
	import Toast from '$lib/Toast.svelte';
	import { 
		GetAllCollectionTrees, 
		CreateRestRequest,
		EnsureDefaultCollection
	} from '$lib/wailsjs/go/main/App.js';

	let { 
		isOpen = $bindable(false),
		requestData = {},
		onSave = () => {},
		onCancel = () => {}
	} = $props();

	let collectionsData = $state([]);
	let expandedNodes = $state(new Set());
	let selectedNode = $state(null);
	let requestName = $state('');
	let isLoading = $state(false);
	let isSaving = $state(false);
	
	// Toast state
	let toastShow = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	let toastDuration = $state(1500);
	let toastAnimation = $state('fade');

	// Toast function
	function showToast(message, type = 'success', duration = 1500, animation = 'fade') {
		toastMessage = message;
		toastType = type;
		toastDuration = duration;
		toastAnimation = animation;
		toastShow = true;
	}

	// Load collections data
	async function loadCollections() {
		isLoading = true;
		try {
			await EnsureDefaultCollection();
			const trees = await GetAllCollectionTrees();
			collectionsData = trees || [];
			
			// Auto-expand all collections for better UX
			const autoExpanded = new Set();
			trees?.forEach(collection => {
				autoExpanded.add(`collection-${collection.id}`);
			});
			expandedNodes = autoExpanded;
		} catch (error) {
			console.error('Failed to load collections:', error);
			collectionsData = [];
		} finally {
			isLoading = false;
		}
	}

	// Load collections when modal opens
	$effect(() => {
		if (isOpen) {
			loadCollections();
			// Generate a default name from request data
			if (requestData.method && requestData.endpoint) {
				const endpoint = requestData.endpoint || requestData.baseEndpoint || '';
				const path = endpoint.split('/').pop() || 'request';
				requestName = `${requestData.method?.toUpperCase()} ${path}`;
			} else {
				requestName = 'New Request';
			}
			selectedNode = null;
		}
	});

	// Handle modal close
	function closeModal() {
		isOpen = false;
		requestName = '';
		selectedNode = null;
		onCancel();
	}

	// Handle save request
	async function saveRequest() {
		if (!requestName.trim()) {
			showToast('Please enter a request name', 'error');
			return;
		}

		if (!selectedNode) {
			showToast('Please select a collection or folder', 'error');
			return;
		}

		isSaving = true;
		try {
			const createRequest = {
				name: requestName.trim(),
				method: requestData.method || 'GET',
				url: requestData.endpoint || requestData.baseEndpoint || '',
				body: requestData.requestBody?.trim() ? requestData.requestBody : null,
				description: requestData.description?.trim() ? requestData.description : null,
				collection_id: selectedNode.type === 'collection' ? selectedNode.id : getCollectionId(selectedNode),
				folder_id: selectedNode.type === 'folder' ? selectedNode.id : null
			};

			console.log('Creating request:', createRequest);
			
			const createdRequest = await CreateRestRequest(createRequest);
			
			showToast('Request saved successfully', 'success');
			
			// Call the onSave callback with the created request
			onSave(createdRequest);
			
			// Close modal
			closeModal();
			
		} catch (error) {
			console.error('Failed to save request:', error);
			showToast('Failed to save request', 'error');
		} finally {
			isSaving = false;
		}
	}

	// Helper to get collection ID from a folder node
	function getCollectionId(node) {
		return collectionsData.find(collection => 
			collection.id === node.id || hasChild(collection, node.id)
		)?.id || collectionsData[0]?.id;
	}

	function hasChild(parent, childId) {
		if (!parent.children) return false;
		return parent.children.some(child => 
			child.id === childId || hasChild(child, childId)
		);
	}

	// Handle keyboard shortcuts
	function handleKeydown(event) {
		if (event.key === 'Escape') {
			closeModal();
		} else if (event.key === 'Enter' && event.ctrlKey) {
			saveRequest();
		}
	}

	// Handle backdrop click
	function handleBackdropClick(event) {
		if (event.target === event.currentTarget) {
			closeModal();
		}
	}

	// Focus name input when modal opens
	function focusNameInput(element) {
		element.focus();
		element.select();
	}
</script>

<!-- Modal Backdrop -->
{#if isOpen}
	<div 
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
		transition:fly={{ duration: 200, opacity: 0 }}
		onclick={handleBackdropClick}
		onkeydown={handleKeydown}
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
		tabindex="-1"
	>
		<!-- Modal Content -->
		<div 
			class="theme-bg-secondary rounded-lg shadow-xl w-full max-w-md mx-4 theme-border border"
			transition:fly={{ duration: 300, y: 50, easing: quintOut }}
			role="document"
		>
			<!-- Modal Header -->
			<div class="flex items-center justify-between p-6 border-b theme-border">
				<h2 id="modal-title" class="text-lg font-semibold theme-text-primary">Save Request</h2>
				<button
					class="theme-text-secondary hover:theme-text-primary transition-colors"
					onclick={closeModal}
					aria-label="Close modal"
				>
					<img src="/icons/x.svg" alt="Close" class="w-5 h-5 theme-icon" />
				</button>
			</div>

			<!-- Modal Body -->
			<div class="p-6 space-y-4">
				<!-- Request Name Input -->
				<div>
					<label for="request-name" class="block text-sm font-medium theme-text-primary mb-2">
						Request Name
					</label>
					<input
						id="request-name"
						type="text"
						bind:value={requestName}
						class="w-full px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
						placeholder="Enter request name"
						use:focusNameInput
					/>
				</div>

				<!-- Collection/Folder Selection -->
				<div>
					<div id="collections-label" class="block text-sm font-medium theme-text-primary mb-2">
						Save to Collection/Folder
					</div>
					<div 
						id="collections-tree"
						role="tree"
						aria-labelledby="collections-label"
						class="border theme-border rounded-md max-h-64 overflow-y-auto theme-bg-primary"
					>
						{#if isLoading}
							<div class="p-4 text-center theme-text-secondary">
								Loading collections...
							</div>
						{:else if collectionsData.length === 0}
							<div class="p-4 text-center theme-text-secondary">
								No collections found
							</div>
						{:else}
							<div class="p-2">
								<CollectionsTree
									{collectionsData}
									bind:expandedNodes
									bind:selectedNode
									mode="select"
									showRequests={false}
									allowDragDrop={false}
									showActions={false}
								/>
							</div>
						{/if}
					</div>
					{#if selectedNode}
						<div class="mt-2 text-sm theme-text-secondary">
							Selected: <span class="font-medium theme-text-primary">{selectedNode.name}</span>
							{#if selectedNode.type === 'folder'}
								<span class="text-xs">in {collectionsData.find(c => c.id === getCollectionId(selectedNode))?.name || ''}</span>
							{/if}
						</div>
					{/if}
				</div>

				<!-- Request Preview -->
				<div class="theme-bg-tertiary rounded-md p-3">
					<div class="text-xs font-medium theme-text-secondary mb-2">Request Preview</div>
					<div class="space-y-1">
						<div class="flex items-center gap-2">
							<span class="text-xs font-mono font-bold px-1.5 py-0.5 rounded bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300">
								{requestData.method?.toUpperCase() || 'GET'}
							</span>
							<span class="text-xs theme-text-primary truncate">
								{requestData.endpoint || requestData.baseEndpoint || 'No URL specified'}
							</span>
						</div>
						{#if requestData.requestBody?.trim()}
							<div class="text-xs theme-text-secondary">Has request body</div>
						{/if}
						{#if requestData.description?.trim()}
							<div class="text-xs theme-text-secondary">Has description</div>
						{/if}
					</div>
				</div>
			</div>

			<!-- Modal Footer -->
			<div class="flex items-center justify-end gap-3 p-6 border-t theme-border">
				<button
					class="px-4 py-2 text-sm font-medium theme-text-secondary hover:theme-text-primary transition-colors"
					onclick={closeModal}
					disabled={isSaving}
				>
					Cancel
				</button>
				<button
					class="px-4 py-2 text-sm font-medium bg-blue-600 hover:bg-blue-700 text-white rounded-md disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
					onclick={saveRequest}
					disabled={isSaving || !requestName.trim() || !selectedNode}
				>
					{#if isSaving}
						Saving...
					{:else}
						Save Request
					{/if}
				</button>
			</div>

			<!-- Keyboard shortcuts hint -->
			<div class="px-6 pb-4 text-xs theme-text-secondary">
				<kbd class="px-1 theme-bg-tertiary theme-text-primary rounded">Ctrl+Enter</kbd> to save, 
				<kbd class="px-1 theme-bg-tertiary theme-text-primary rounded">Esc</kbd> to cancel
			</div>
		</div>
	</div>
{/if}

<!-- Toast Component -->
<Toast 
	bind:show={toastShow} 
	message={toastMessage} 
	type={toastType} 
	duration={toastDuration} 
	animation={toastAnimation} />