<script>
	import { onMount } from 'svelte';
	
	// Props using $props()
	let {
		direction = 'horizontal',
		defaultSplit = 50,
		minSize = 20,
		maxSize = 80,
		splitterSize = 8,
		className = '',
		panel1,
		panel2
	} = $props();
	
	// State
	let container = $state();
	let splitter = $state();
	let isDragging = $state(false);
	let splitPercentage = $state(defaultSplit);
	
	// Calculate sizes based on direction using $derived
	const isHorizontal = $derived(direction === 'horizontal');
	const panel1Size = $derived(splitPercentage);
	const panel2Size = $derived(100 - splitPercentage);
	
	// Mouse/touch event handlers
	function handleMouseDown(event) {
		isDragging = true;
		document.addEventListener('mousemove', handleMouseMove);
		document.addEventListener('mouseup', handleMouseUp);
		event.preventDefault();
	}
	
	function handleMouseMove(event) {
		if (!isDragging || !container) return;
		
		const rect = container.getBoundingClientRect();
		let percentage;
		
		if (isHorizontal) {
			const relativeY = event.clientY - rect.top;
			percentage = (relativeY / rect.height) * 100;
		} else {
			const relativeX = event.clientX - rect.left;
			percentage = (relativeX / rect.width) * 100;
		}
		
		// Clamp percentage within bounds
		splitPercentage = Math.max(minSize, Math.min(maxSize, percentage));
	}
	
	function handleMouseUp() {
		isDragging = false;
		document.removeEventListener('mousemove', handleMouseMove);
		document.removeEventListener('mouseup', handleMouseUp);
	}
	
	// Touch event handlers for mobile support
	function handleTouchStart(event) {
		isDragging = true;
		document.addEventListener('touchmove', handleTouchMove);
		document.addEventListener('touchend', handleTouchEnd);
		event.preventDefault();
	}
	
	function handleTouchMove(event) {
		if (!isDragging || !container) return;
		
		const touch = event.touches[0];
		const rect = container.getBoundingClientRect();
		let percentage;
		
		if (isHorizontal) {
			const relativeY = touch.clientY - rect.top;
			percentage = (relativeY / rect.height) * 100;
		} else {
			const relativeX = touch.clientX - rect.left;
			percentage = (relativeX / rect.width) * 100;
		}
		
		splitPercentage = Math.max(minSize, Math.min(maxSize, percentage));
	}
	
	function handleTouchEnd() {
		isDragging = false;
		document.removeEventListener('touchmove', handleTouchMove);
		document.removeEventListener('touchend', handleTouchEnd);
	}
	
	// Cleanup on component destroy
	onMount(() => {
		return () => {
			document.removeEventListener('mousemove', handleMouseMove);
			document.removeEventListener('mouseup', handleMouseUp);
			document.removeEventListener('touchmove', handleTouchMove);
			document.removeEventListener('touchend', handleTouchEnd);
		};
	});
</script>

<div 
	bind:this={container}
	class="resizable-splitter-container {className} {isHorizontal ? 'flex flex-col' : 'flex flex-row'} h-full w-full"
	class:dragging={isDragging}
	class:vertical={!isHorizontal}
>
	<!-- Panel 1 -->
	<div 
		class="panel panel-1 overflow-hidden"
		style="{isHorizontal ? 'height' : 'width'}: {panel1Size}%; {isHorizontal ? 'width: 100%' : 'height: 100%'}"
	>
		{#if panel1}
			{@render panel1()}
		{/if}
	</div>
	
	<!-- Splitter -->
	<button
		bind:this={splitter}
		class="splitter theme-bg-secondary hover:theme-bg-accent transition-colors duration-200 select-none flex items-center justify-center border-0 outline-none"
		class:cursor-ns-resize={isHorizontal}
		class:cursor-ew-resize={!isHorizontal}
		style="{isHorizontal ? 'height' : 'width'}: {splitterSize}px; {isHorizontal ? 'width: 100%' : 'height: 100%'}"
		onmousedown={handleMouseDown}
		ontouchstart={handleTouchStart}
		aria-label="Resize panels - drag to adjust panel sizes"
		title="Drag to resize panels"
	>
		<!-- Visual indicator - always visible -->
		<div class="splitter-handle">
			{#if isHorizontal}
				<div class="flex items-center justify-center w-full h-full">
					<div class="flex flex-col gap-0.5">
						<div class="w-6 h-0.5 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
						<div class="w-6 h-0.5 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
						<div class="w-6 h-0.5 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
					</div>
				</div>
			{:else}
				<div class="flex items-center justify-center w-full h-full">
					<div class="flex gap-0.5">
						<div class="w-0.5 h-6 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
						<div class="w-0.5 h-6 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
						<div class="w-0.5 h-6 bg-gray-400 dark:bg-gray-500 rounded-full"></div>
					</div>
				</div>
			{/if}
		</div>
	</button>
	
	<!-- Panel 2 -->
	<div 
		class="panel panel-2 overflow-hidden"
		style="{isHorizontal ? 'height' : 'width'}: {panel2Size}%; {isHorizontal ? 'width: 100%' : 'height: 100%'}"
	>
		{#if panel2}
			{@render panel2()}
		{/if}
	</div>
</div>

<style>
	.resizable-splitter-container {
		user-select: none;
	}
	
	.resizable-splitter-container.dragging {
		user-select: none;
	}
	
	.resizable-splitter-container.dragging * {
		pointer-events: none;
	}
	
	/* Override cursor during drag */
	.resizable-splitter-container.dragging {
		cursor: ns-resize !important;
	}
	
	.resizable-splitter-container.dragging.vertical {
		cursor: ew-resize !important;
	}
	
	.splitter {
		flex-shrink: 0;
		position: relative;
		z-index: 1;
		border-radius: 2px;
		background: rgba(156, 163, 175, 0.1);
		padding: 0;
		margin: 0;
		border: 1px solid rgba(156, 163, 175, 0.2);
	}
	
	.splitter:focus {
		outline: 2px solid rgb(59, 130, 246);
		outline-offset: 1px;
	}
	
	.splitter:hover {
		background-color: rgba(59, 130, 246, 0.15) !important;
		border-color: rgba(59, 130, 246, 0.3);
	}
	
	.splitter:hover .splitter-handle > div > div > div {
		background-color: rgb(59, 130, 246) !important;
	}
	
	.splitter:active {
		background-color: rgba(37, 99, 235, 0.2) !important;
		border-color: rgba(37, 99, 235, 0.4);
	}
	
	.splitter:active .splitter-handle > div > div > div {
		background-color: rgb(37, 99, 235) !important;
	}
	
	.panel {
		flex-shrink: 0;
	}
	
	/* Prevent text selection during drag */
	.dragging {
		-webkit-user-select: none;
		-moz-user-select: none;
		-ms-user-select: none;
		user-select: none;
	}
	
	/* Enhanced cursor styles */
	.cursor-ns-resize {
		cursor: ns-resize;
	}
	
	.cursor-ew-resize {
		cursor: ew-resize;
	}
	
	/* Additional visual feedback */
	.splitter-handle {
		transition: all 0.2s ease-in-out;
	}
	
	.splitter:hover .splitter-handle {
		transform: scale(1.1);
	}
</style>