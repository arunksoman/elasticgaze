<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import ClusterDropdown from './ClusterDropdown.svelte';
	import { connectionWarningStatus } from '$lib/stores/connectionWarningStore.js';
	import { sidebarExpanded } from '$lib/stores/sidebarStore.js';
	import { collectionsOpen } from '$lib/stores/collectionsStore.js';
	import {
		WindowMinimise,
		WindowToggleMaximise,
		WindowIsMaximised,
		Quit,
		LogInfo
	} from '$lib/wailsjs/runtime/runtime';

	let { isMax = $bindable(false) } = $props();

	// Pages where the cluster dropdown should not be shown
	const excludedPages = ['/about', '/connections'];
	
	// Subscribe to connection status from store
	const connectionStatus = $derived($connectionWarningStatus);
	
	// Subscribe to sidebar states with default values
	const isSidebarExpanded = $derived($sidebarExpanded ?? false);
	const isCollectionsOpen = $derived($collectionsOpen ?? false);
	
	// Calculate dynamic left position for draggable area
	let draggableLeftPosition = $state(72); // Default fallback value
	
	$effect(() => {
		let leftPos = 0;
		
		// Add main sidebar width
		leftPos += isSidebarExpanded ? 224 : 64; // w-56 = 224px, w-16 = 64px
		
		// Add collections sidebar width if open and on REST page
		if (isCollectionsOpen && page.url.pathname === '/rest') {
			leftPos += 320; // w-80 = 320px
		}
		
		const finalPos = leftPos + 8; // Add small buffer
		draggableLeftPosition = finalPos;
	});
	
	// Check if cluster dropdown should be shown
	const shouldShowClusterDropdown = $derived(() => {
		const currentPath = page.url.pathname;
		return !excludedPages.includes(currentPath) && 
		       connectionStatus.connectionState === 'working';
	});

	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
	}

	function handleMinimise() {
		if (hasWails()) {
			try { WindowMinimise(); } catch {}
		}
	}

	function handleToggleMaximise() {
		LogInfo('Toggle button clicked!');
		if (hasWails()) {
			try {
				WindowToggleMaximise();
				const maximized = !!WindowIsMaximised();
				LogInfo('Window is now maximized: ' + maximized);
				isMax = maximized;
			} catch (e) {
				LogInfo('Error during toggle: ' + e);
			}
		} else {
			LogInfo('Wails runtime not found.');
		}
	}

	function handleClose() {
		if (hasWails()) {
			try { Quit(); } catch {}
		} else {
			// Fallback when running purely in browser dev
			window.close();
		}
	}

	function handleConnection() {
		// Navigate to connections management page
		goto('/connections');
	}
</script>

<!-- Window Controls and Draggable Title Bar -->
<!-- Dynamic draggable area that adjusts to sidebar states -->
<div class="fixed top-0 h-12 z-[999]" style="left: {draggableLeftPosition}px; right: 150px; --wails-draggable:drag" aria-hidden="true">
	<!-- Draggable area that dynamically adapts to sidebar widths -->
</div>

<!-- Window Controls (top-right) -->
<div class="fixed top-2 right-2 flex gap-2 items-center z-[1000]" aria-label="Window controls" style="--wails-draggable:no-drag">
	<!-- Cluster Dropdown (only shown when connection is working and not on excluded pages) -->
	{#if shouldShowClusterDropdown()}
		<div class="mr-4" style="--wails-draggable:no-drag">
			<ClusterDropdown />
		</div>
	{/if}

	<!-- Connection Button -->
	<div class="mr-5" style="--wails-draggable:no-drag">
		<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title="Elasticsearch Connections" onclick={handleConnection} aria-label="Elasticsearch Connections">
			<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-icon); mask-image: url('/icons/connect_elastic.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/connect_elastic.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label="Connection icon"></span>
		</button>
	</div>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title="Minimize" onclick={handleMinimise} aria-label="Minimize" style="--wails-draggable:no-drag">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-icon); mask-image: url('/icons/minimize.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/minimize.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label="Minimize icon"></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title={isMax ? 'Restore' : 'Maximize'} onclick={handleToggleMaximise} aria-label={isMax ? 'Restore' : 'Maximize'} style="--wails-draggable:no-drag">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style={`background-color: var(--window-control-icon); mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;`} onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label={isMax ? 'Restore icon' : 'Maximize icon'}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-red-500/35 hover:bg-red-500/[0.06] active:translate-y-[0.5px]" title="Close" onclick={handleClose} aria-label="Close" style="--wails-draggable:no-drag">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-close-icon); mask-image: url('/icons/close.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/close.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon)'} role="img" aria-label="Close icon"></span>
	</button>
	<div class="absolute -top-2 -right-2 -bottom-2 -left-2 pointer-events-none" aria-hidden="true"></div>
</div>