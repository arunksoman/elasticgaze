<script>
	import { onMount } from 'svelte';
	import { fly, slide } from 'svelte/transition';
	import { quintOut } from 'svelte/easing';
	import { sidebarExpanded } from '$lib/stores/sidebarStore.js';
	import { tabStore } from '$lib/stores/tabStore.js';
	import Toast from '$lib/Toast.svelte';
	import { 
		GetAllCollectionTrees, 
		CreateCollection, 
		CreateFolder, 
		CreateRestRequest,
		UpdateCollection,
		UpdateFolder,
		UpdateRestRequest,
		DeleteCollection,
		DeleteFolder,
		DeleteRestRequest,
		EnsureDefaultCollection
	} from '$lib/wailsjs/go/main/App.js';

	let { 
		isOpen = $bindable(false)
	} = $props();

	let collectionsData = $state([]);
	let expandedNodes = $state(new Set());
	let editingNode = $state(null);
	let editingName = $state('');
	let contextMenu = $state(null);
	let isLoading = $state(false);
	let draggedItem = $state(null);
	let dropTarget = $state(null);
	
	// Toast state
	let toastShow = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	let toastDuration = $state(1500);
	let toastAnimation = $state('fade');

	// Get sidebar width based on expanded state
	const sidebarWidth = $derived($sidebarExpanded ? 224 : 64); // w-56 = 224px, w-16 = 64px

	// Toast function
	function showToast(message, type = 'success', duration = 1500, animation = 'fade') {
		toastMessage = message;
		toastType = type;
		toastDuration = duration;
		toastAnimation = animation;
		toastShow = true;
	}

	// Custom transition that slides from the sidebar edge
	function slideFromSidebar(node, { duration = 500, easing = quintOut }) {
		return {
			duration,
			easing,
			css: (t) => {
				const width = t * 320; // 320px = w-80 in Tailwind (20rem)
				return `
					width: ${width}px;
					max-width: 320px;
					overflow: hidden;
				`;
			}
		};
	}

	// Focus action for inline editing  
	function focus(element) {
		element.focus();
		element.select();
	}

	// Load collections data
	async function loadCollections() {
		isLoading = true;
		try {
			// Ensure default collection exists
			await EnsureDefaultCollection();
			// Load all collection trees
			const trees = await GetAllCollectionTrees();
			collectionsData = trees || [];
		} catch (error) {
			console.error('Failed to load collections:', error);
			collectionsData = [];
		} finally {
			isLoading = false;
		}
	}

	// Load collections on mount
	onMount(() => {
		loadCollections();
	});

	// Watch for isOpen changes to reload data
	$effect(() => {
		if (isOpen) {
			loadCollections();
		}
	});

	// Toggle node expansion
	function toggleNode(nodeId, nodeType) {
		const key = `${nodeType}-${nodeId}`;
		if (expandedNodes.has(key)) {
			expandedNodes.delete(key);
		} else {
			expandedNodes.add(key);
		}
		expandedNodes = new Set(expandedNodes); // Trigger reactivity
	}

	// Check if node is expanded
	function isNodeExpanded(nodeId, nodeType) {
		return expandedNodes.has(`${nodeType}-${nodeId}`);
	}

	// Handle request selection
	function selectRequest(request) {
		console.log('Request selected:', request);
		
		// Load request data into a new tab
		const newTabData = {
			method: request.method || 'GET',
			endpoint: request.url || '',
			baseEndpoint: request.url || '',
			params: [],
			requestBody: request.body || '',
			description: request.description || '',
			responseData: null,
			isLoading: false
		};
		
		// Add a new tab with the loaded request data
		tabStore.addTab(request.name || 'Loaded Request', newTabData);
	}

	// Start editing
	function startEditing(node) {
		editingNode = node;
		editingName = node.name;
		contextMenu = null;
	}

	// Cancel editing
	function cancelEditing() {
		editingNode = null;
		editingName = '';
	}

	// Save edit
	async function saveEdit() {
		if (!editingNode || !editingName.trim()) return;
		
		try {
			if (editingNode.type === 'collection') {
				await UpdateCollection(editingNode.id, { name: editingName.trim() });
			} else if (editingNode.type === 'folder') {
				await UpdateFolder(editingNode.id, { name: editingName.trim() });
			} else if (editingNode.type === 'request') {
				await UpdateRestRequest(editingNode.id, { name: editingName.trim() });
			}
			await loadCollections();
		} catch (error) {
			console.error('Failed to update:', error);
		}
		
		editingNode = null;
		editingName = '';
	}

	// Handle key events for editing
	function handleEditKeydown(event) {
		if (event.key === 'Enter') {
			event.preventDefault();
			saveEdit();
		} else if (event.key === 'Escape') {
			event.preventDefault();
			cancelEditing();
		}
	}

	// Context menu actions
	function showContextMenu(event, node) {
		event.preventDefault();
		event.stopPropagation();
		
		contextMenu = {
			x: event.clientX,
			y: event.clientY,
			node: node
		};
	}

	function hideContextMenu() {
		contextMenu = null;
	}

	// Create new items
	async function createNewCollection() {
		try {
			await CreateCollection({ 
				name: 'New Collection',
				description: 'A new REST request collection'
			});
			await loadCollections();
		} catch (error) {
			console.error('Failed to create collection:', error);
		}
		hideContextMenu();
	}

	async function createNewFolder(parentNode) {
		try {
			const collectionId = parentNode.type === 'collection' ? parentNode.id : 
				getCollectionId(parentNode);
			const parentFolderId = parentNode.type === 'folder' ? parentNode.id : null;
			
			await CreateFolder({
				name: 'New Folder',
				collection_id: collectionId,
				parent_folder_id: parentFolderId
			});
			await loadCollections();
		} catch (error) {
			console.error('Failed to create folder:', error);
		}
		hideContextMenu();
	}

	async function createNewRequest(parentNode) {
		try {
			const collectionId = parentNode.type === 'collection' ? parentNode.id : 
				getCollectionId(parentNode);
			const folderId = parentNode.type === 'folder' ? parentNode.id : null;
			
			await CreateRestRequest({
				name: 'New Request',
				method: 'GET',
				url: '/',
				collection_id: collectionId,
				folder_id: folderId
			});
			await loadCollections();
		} catch (error) {
			console.error('Failed to create request:', error);
		}
		hideContextMenu();
	}

	// Delete items
	async function deleteNode(node) {
		if (!confirm(`Are you sure you want to delete "${node.name}"?`)) return;
		
		try {
			if (node.type === 'collection') {
				await DeleteCollection(node.id);
			} else if (node.type === 'folder') {
				await DeleteFolder(node.id);
			} else if (node.type === 'request') {
				await DeleteRestRequest(node.id);
			}
			await loadCollections();
		} catch (error) {
			console.error('Failed to delete:', error);
		}
		hideContextMenu();
	}

	// Helper to get collection ID for nested nodes
	function getCollectionId(node) {
		// For this implementation, we assume each tree starts with a collection
		// In a more complex implementation, you might need to traverse up the tree
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

	// Handle clicks outside to close context menu and editing
	function handleDocumentClick(event) {
		if (contextMenu && !event.target.closest('.context-menu')) {
			hideContextMenu();
		}
	}

	onMount(() => {
		document.addEventListener('click', handleDocumentClick);
		return () => {
			document.removeEventListener('click', handleDocumentClick);
		};
	});

	// Get icon for request method
	function getMethodIcon(method) {
		const icons = {
			'GET': '/icons/get.svg',
			'POST': '/icons/post.svg',
			'PUT': '/icons/put.svg',
			'DELETE': '/icons/delete.svg',
			'PATCH': '/icons/patch.svg',
			'HEAD': '/icons/head.svg',
			'OPTIONS': '/icons/options.svg'
		};
		return icons[method?.toUpperCase()] || '/icons/get.svg';
	}

	// Get method color class
	function getMethodColorClass(method) {
		const colors = {
			'GET': 'text-blue-600 dark:text-blue-400',
			'POST': 'text-green-600 dark:text-green-400',
			'PUT': 'text-orange-600 dark:text-orange-400',
			'DELETE': 'text-red-600 dark:text-red-400',
			'PATCH': 'text-purple-600 dark:text-purple-400',
			'HEAD': 'text-gray-600 dark:text-gray-400',
			'OPTIONS': 'text-yellow-600 dark:text-yellow-400'
		};
		return colors[method?.toUpperCase()] || 'text-gray-600 dark:text-gray-400';
	}

	// Drag and drop functions
	/**
	 * Handle drag start event for a request
	 */
	function handleDragStart(event, node) {
		if (node.type !== 'request') return; // Only allow dragging requests
		
		draggedItem = node;
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/plain', ''); // Required for Firefox
		
		// Add visual feedback
		event.target.style.opacity = '0.5';
	}

	function handleDragEnd(event) {
		draggedItem = null;
		dropTarget = null;
		event.target.style.opacity = '1';
	}

	function handleDragOver(event, node) {
		// Only allow dropping on collections and folders
		if (node.type === 'request') return;
		
		event.preventDefault();
		event.dataTransfer.dropEffect = 'move';
		dropTarget = node;
	}

	function handleDragLeave(event, node) {
		if (dropTarget?.id === node.id && dropTarget?.type === node.type) {
			dropTarget = null;
		}
	}

	async function handleDrop(event, targetNode) {
		event.preventDefault();
		
		if (!draggedItem || !targetNode || draggedItem.type !== 'request') {
			return;
		}

		// Don't allow dropping on self or if already in the same location
		if (draggedItem.id === targetNode.id && draggedItem.type === targetNode.type) {
			return;
		}

		try {
			// Capture names before potential nullification
			const draggedItemName = draggedItem.name;
			const targetNodeName = targetNode.name;
			
			// Determine the new folder ID and collection ID
			let newFolderId = null;
			let newCollectionId = null;
			
			if (targetNode.type === 'folder') {
				newFolderId = targetNode.id;
				// When moving to a folder, keep current collection (don't change collection_id)
			} else if (targetNode.type === 'collection') {
				// Moving to root of collection (no folder, but specific collection)
				newFolderId = null;
				newCollectionId = targetNode.id;
			}

			// Build update request object
			const updateRequest = {};
			
			// Always update folder_id, using special handling for null values
			if (targetNode.type === 'folder') {
				updateRequest.folder_id = newFolderId;
			} else if (targetNode.type === 'collection') {
				// For collection root, we need to explicitly set folder_id to null
				// Use a special marker that the backend can recognize
				updateRequest.folder_id = -1; // Use -1 to indicate null
			}
			
			// Only update collection_id when moving to a different collection
			if (newCollectionId !== null) {
				updateRequest.collection_id = newCollectionId;
			}
			
			console.log('🔄 Updating request:', draggedItem.id, 'with:', updateRequest);
			
			const result = await UpdateRestRequest(draggedItem.id, updateRequest);
			console.log('✅ Update result:', result);

			// Reload collections to reflect the change
			console.log('🔄 Reloading collections...');
			await loadCollections();
			console.log('✅ Collections reloaded:', JSON.stringify($state.snapshot(collectionsData), null, 2));
			
			// Show success message using captured names
			showToast(`Request "${draggedItemName}" moved to ${targetNodeName}`, 'success');
			
			// Expand the target if it's a folder to show the moved request
			if (targetNode.type === 'folder') {
				const key = `${targetNode.type}-${targetNode.id}`;
				expandedNodes.add(key);
				expandedNodes = new Set(expandedNodes);
			}

		} catch (error) {
			console.error('Failed to move request:', error);
			showToast('Failed to move request', 'error');
		} finally {
			draggedItem = null;
			dropTarget = null;
		}
	}
</script>

<!-- Collections Sidebar Overlay -->
{#if isOpen}
	<div 
		class="fixed top-0 h-full w-80 max-w-80 z-10 theme-bg-primary border-r theme-border shadow-xl"
		style="left: {sidebarWidth}px;"
		transition:slideFromSidebar={{ duration: 500, easing: quintOut }}
	>
		<!-- Header -->
		<div class="p-4 border-b theme-border">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-2">
					<button
						class="p-1.5 rounded theme-hover theme-text-secondary hover:theme-text-primary"
						onclick={() => isOpen = false}
						title="Close Collections"
					>
						<img src="/icons/chevrons-left.svg" alt="Close Collections" class="w-4 h-4 theme-icon" />
					</button>
					<h2 class="text-lg font-semibold theme-text-primary">Collections</h2>
				</div>
				<div class="flex items-center gap-2">
					<button
						class="p-1.5 rounded theme-hover theme-text-secondary hover:theme-text-primary"
						onclick={createNewCollection}
						title="New Collection"
					>
						<img src="/icons/plus.svg" alt="Add" class="w-4 h-4 theme-icon" />
					</button>
				</div>
			</div>
		</div>

		<!-- Content -->
		<div class="flex-1 overflow-y-auto p-2" style="height: calc(100vh - 80px);">
			{#if isLoading}
				<div class="flex items-center justify-center py-8">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600"></div>
				</div>
			{:else if collectionsData.length === 0}
				<div class="text-center py-8 theme-text-secondary">
					<p>No collections yet</p>
					<button
						class="mt-2 px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700"
						onclick={createNewCollection}
					>
						Create First Collection
					</button>
				</div>
			{:else}
				<div class="space-y-1">
					{#each collectionsData as collection (collection.id)}
						{@render TreeNode(collection, 0)}
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}

<!-- Context Menu -->
{#if contextMenu}
	<div 
		class="fixed z-50 min-w-48 py-2 theme-bg-secondary border theme-border rounded-lg shadow-lg"
		style="left: {contextMenu.x}px; top: {contextMenu.y}px;"
	>
		<div class="context-menu">
			<button
				class="w-full px-4 py-2 text-left text-sm theme-text-primary hover:theme-bg-tertiary"
				onclick={() => startEditing(contextMenu.node)}
			>
				Rename
			</button>
			
			{#if contextMenu.node.type === 'collection' || contextMenu.node.type === 'folder'}
				<button
					class="w-full px-4 py-2 text-left text-sm theme-text-primary hover:theme-bg-tertiary"
					onclick={() => createNewFolder(contextMenu.node)}
				>
					New Folder
				</button>
				<button
					class="w-full px-4 py-2 text-left text-sm theme-text-primary hover:theme-bg-tertiary"
					onclick={() => createNewRequest(contextMenu.node)}
				>
					New Request
				</button>
			{/if}
			
			<hr class="my-1 theme-border" />
			
			<button
				class="w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20"
				onclick={() => deleteNode(contextMenu.node)}
			>
				Delete
			</button>
		</div>
	</div>
{/if}

<!-- Tree Node Component -->
{#snippet TreeNode(node, level)}
	<div class="select-none">
		<div 
			class="flex items-center gap-1 py-1 px-2 rounded cursor-pointer theme-hover group {
				node.type === 'request' ? 'cursor-grab active:cursor-grabbing' : ''
			} {
				dropTarget?.id === node.id && dropTarget?.type === node.type ? 'bg-purple-100 dark:bg-purple-900/30 border-2 border-purple-300 dark:border-purple-600' : ''
			}"
			style="padding-left: {(level * 16) + 8}px"
			role="button"
			tabindex="0"
			draggable={node.type === 'request'}
			ondragstart={(e) => handleDragStart(e, node)}
			ondragend={handleDragEnd}
			ondragover={(e) => handleDragOver(e, node)}
			ondragleave={(e) => handleDragLeave(e, node)}
			ondrop={(e) => handleDrop(e, node)}
			onclick={() => {
				if (node.type === 'request') {
					selectRequest(node);
				} else {
					toggleNode(node.id, node.type);
				}
			}}
			onkeydown={(e) => {
				if (e.key === 'Enter' || e.key === ' ') {
					e.preventDefault();
					if (node.type === 'request') {
						selectRequest(node);
					} else {
						toggleNode(node.id, node.type);
					}
				}
			}}
			oncontextmenu={(e) => showContextMenu(e, node)}
		>
			<!-- Expand/Collapse Icon -->
			{#if node.children && node.children.length > 0}
				<button
					class="p-0.5 rounded theme-text-secondary hover:theme-text-primary"
					onclick={(e) => {
						e.stopPropagation();
						toggleNode(node.id, node.type);
					}}
				>
					<img 
						src="/icons/triangle.svg" 
						alt="Expand" 
						class="w-2 h-2 theme-icon transition-transform {isNodeExpanded(node.id, node.type) ? 'rotate-90' : ''}" 
					/>
				</button>
			{:else}
				<div class="w-4"></div>
			{/if}

			<!-- Node Icon -->
			{#if node.type === 'collection'}
				<img src="/icons/folder.svg" alt="Collection" class="w-4 h-4 theme-icon" />
			{:else if node.type === 'folder'}
				<img src="/icons/folder.svg" alt="Folder" class="w-4 h-4 theme-icon" />
			{:else if node.type === 'request'}
				<span class="text-xs font-mono font-bold px-1 rounded {getMethodColorClass(node.method)}">
					{node.method?.toUpperCase() || 'GET'}
				</span>
			{/if}

			<!-- Node Name/Edit Input -->
			<div class="flex-1 min-w-0">
				{#if editingNode && editingNode.id === node.id && editingNode.type === node.type}
					<input
						type="text"
						bind:value={editingName}
						class="w-full px-1 py-0 text-sm bg-transparent border-b theme-border theme-text-primary focus:outline-none focus:border-purple-500 focus-input"
						onkeydown={handleEditKeydown}
						onblur={saveEdit}
						use:focus
					/>
				{:else}
					<span class="text-sm theme-text-primary truncate block">
						{node.name}
					</span>
				{/if}
			</div>

			<!-- Action Buttons (visible on hover) -->
			<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
				<button
					class="p-0.5 rounded theme-text-secondary hover:theme-text-primary"
					onclick={(e) => {
						e.stopPropagation();
						showContextMenu(e, node);
					}}
					title="More actions"
				>
					<img src="/icons/more-horizontal.svg" alt="More" class="w-3 h-3 theme-icon" />
				</button>
			</div>
		</div>

		<!-- Children -->
		{#if node.children && node.children.length > 0 && isNodeExpanded(node.id, node.type)}
			{#each node.children as child (child.id + child.type)}
				{@render TreeNode(child, level + 1)}
			{/each}
		{/if}
	</div>
{/snippet}

<!-- Toast Component -->
<Toast 
	bind:show={toastShow} 
	message={toastMessage} 
	type={toastType} 
	duration={toastDuration} 
	animation={toastAnimation} 
/>

<style>
	/* Custom styles for focused input - handled by focus action */
</style>