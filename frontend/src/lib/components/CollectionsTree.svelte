<script>
	import { onMount } from 'svelte';
	
	let { 
		collectionsData = [],
		expandedNodes = $bindable(new Set()),
		mode = 'display', // 'display' or 'select'
		selectedNode = $bindable(null),
		showRequests = true,
		allowDragDrop = true,
		showActions = true,
		editingNode = $bindable(null),
		editingName = $bindable(''),
		onNodeSelect = () => {},
		onNodeEdit = () => {},
		onNodeAction = () => {},
		onDragDrop = () => {}
	} = $props();

	let draggedItem = $state(null);
	let dropTarget = $state(null);

	// Focus action for inline editing  
	function focus(element) {
		element.focus();
		element.select();
	}

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

	// Handle node selection
	function handleNodeClick(node) {
		if (mode === 'select') {
			if (node.type === 'collection' || node.type === 'folder') {
				selectedNode = node;
			}
		} else if (node.type === 'request') {
			onNodeSelect(node);
		} else {
			toggleNode(node.id, node.type);
		}
	}

	// Start editing
	function startEditing(node) {
		editingNode = node;
		editingName = node.name;
	}

	// Cancel editing
	function cancelEditing() {
		editingNode = null;
		editingName = '';
	}

	// Save edit
	function saveEdit() {
		if (!editingNode || !editingName.trim()) {
			cancelEditing();
			return;
		}
		
		onNodeEdit(editingNode, editingName.trim());
		editingNode = null;
		editingName = '';
	}

	// Handle key events for editing
	function handleEditKeydown(event) {
		event.stopPropagation();
		
		if (event.key === 'Enter') {
			event.preventDefault();
			saveEdit();
		} else if (event.key === 'Escape') {
			event.preventDefault();
			cancelEditing();
		}
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
	function handleDragStart(event, node) {
		if (!allowDragDrop || node.type !== 'request') return;
		
		draggedItem = node;
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/plain', '');
		event.target.style.opacity = '0.5';
	}

	function handleDragEnd(event) {
		draggedItem = null;
		dropTarget = null;
		event.target.style.opacity = '1';
	}

	function handleDragOver(event, node) {
		if (!allowDragDrop || node.type === 'request') return;
		
		event.preventDefault();
		event.dataTransfer.dropEffect = 'move';
		dropTarget = node;
	}

	function handleDragLeave(event, node) {
		if (dropTarget?.id === node.id && dropTarget?.type === node.type) {
			dropTarget = null;
		}
	}

	function handleDrop(event, targetNode) {
		if (!allowDragDrop) return;
		
		event.preventDefault();
		
		if (!draggedItem || !targetNode || draggedItem.type !== 'request') {
			return;
		}

		onDragDrop(draggedItem, targetNode);
		
		draggedItem = null;
		dropTarget = null;
	}
</script>

<!-- Tree Node Component -->
{#snippet TreeNode(node, level)}
	<div class="select-none">
		<div 
			class="flex items-center gap-1 py-1 px-2 rounded cursor-pointer group {
				mode === 'select' && (node.type === 'collection' || node.type === 'folder') 
					? 'hover:bg-blue-100 dark:hover:bg-blue-900/30' 
					: 'theme-hover'
			} {
				mode === 'select' && selectedNode?.id === node.id && selectedNode?.type === node.type
					? 'bg-blue-200 dark:bg-blue-800/50 border border-blue-300 dark:border-blue-600'
					: ''
			} {
				node.type === 'request' && allowDragDrop ? 'cursor-grab active:cursor-grabbing' : ''
			} {
				dropTarget?.id === node.id && dropTarget?.type === node.type 
					? 'bg-purple-100 dark:bg-purple-900/30 border-2 border-purple-300 dark:border-purple-600' 
					: ''
			}"
			style="padding-left: {(level * 16) + 8}px"
			role="button"
			tabindex="0"
			draggable={allowDragDrop && node.type === 'request'}
			ondragstart={(e) => handleDragStart(e, node)}
			ondragend={handleDragEnd}
			ondragover={(e) => handleDragOver(e, node)}
			ondragleave={(e) => handleDragLeave(e, node)}
			ondrop={(e) => handleDrop(e, node)}
			onclick={() => handleNodeClick(node)}
			ondblclick={() => {
				if (showActions) {
					startEditing(node);
				}
			}}
			onkeydown={(e) => {
				if (editingNode && editingNode.id === node.id && editingNode.type === node.type) {
					return;
				}
				
				if (e.key === 'Enter' || e.key === ' ') {
					e.preventDefault();
					handleNodeClick(node);
				}
			}}
			oncontextmenu={(e) => {
				if (showActions) {
					e.preventDefault();
					onNodeAction(e, node);
				}
			}}
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
			{:else if node.type === 'request' && showRequests}
				<span class="text-xs font-mono font-bold px-1 rounded {getMethodColorClass(node.method)}">
					{node.method?.toUpperCase() || 'GET'}
				</span>
			{/if}

			<!-- Node Name/Edit Input -->
			<div class="flex-1 min-w-0">
				{#if editingNode && editingNode.id === node.id && editingNode.type === node.type && showActions}
					<input
						type="text"
						bind:value={editingName}
						class="text-sm bg-transparent border border-purple-500 theme-text-primary focus:outline-none focus:border-purple-600 rounded px-1 py-0.5 min-w-0"
						style="width: {Math.max(60, editingName.length * 8 + 16)}px;"
						onclick={(e) => e.stopPropagation()}
						onkeydown={handleEditKeydown}
						onblur={(e) => {
							e.stopPropagation();
							saveEdit();
						}}
						onfocus={(e) => e.stopPropagation()}
						use:focus
					/>
				{:else}
					<div class="inline-block">
						<span 
							class="text-sm theme-text-primary cursor-pointer theme-hover px-1 py-0.5 rounded inline-block"
							role="button"
							tabindex="0"
							onclick={(e) => {
								if (showActions) {
									e.stopPropagation();
									startEditing(node);
								}
							}}
							onkeydown={(e) => {
								if (showActions && (e.key === 'Enter' || e.key === ' ')) {
									e.preventDefault();
									e.stopPropagation();
									startEditing(node);
								}
							}}
							title={showActions ? "Click to edit" : ""}
							aria-label={showActions ? `Edit ${node.name}` : node.name}
						>
							{node.name}
						</span>
					</div>
				{/if}
			</div>

			<!-- Action Buttons (visible on hover) -->
			{#if showActions}
				<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
					<button
						class="p-0.5 rounded theme-text-secondary hover:theme-text-primary"
						onclick={(e) => {
							e.stopPropagation();
							onNodeAction(e, node);
						}}
						title="More actions"
					>
						<img src="/icons/more-horizontal.svg" alt="More" class="w-3 h-3 theme-icon" />
					</button>
				</div>
			{/if}

			<!-- Selection indicator for select mode -->
			{#if mode === 'select' && (node.type === 'collection' || node.type === 'folder')}
				<div class="flex items-center">
					{#if selectedNode?.id === node.id && selectedNode?.type === node.type}
						<img src="/icons/check.svg" alt="Selected" class="w-4 h-4 text-blue-600 dark:text-blue-400" />
					{/if}
				</div>
			{/if}
		</div>

		<!-- Children -->
		{#if node.children && node.children.length > 0 && isNodeExpanded(node.id, node.type)}
			{#each node.children as child (child.id + child.type)}
				{#if showRequests || child.type !== 'request'}
					{@render TreeNode(child, level + 1)}
				{/if}
			{/each}
		{/if}
	</div>
{/snippet}

<!-- Tree Container -->
<div class="collections-tree">
	{#each collectionsData as collection (collection.id + collection.type)}
		{@render TreeNode(collection, 0)}
	{/each}
</div>