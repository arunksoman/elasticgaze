<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { theme } from '$lib/theme.js';
	import WindowControls from '$lib/components/WindowControls.svelte';
	import SidebarMenu from '$lib/components/SidebarMenu.svelte';
	import ConnectionWarning from '$lib/components/ConnectionWarning.svelte';
	import { 
		connectionWarningStatus, 
		updateConnectionWarningStatus, 
		triggerConnectionCheck 
	} from '$lib/stores/connectionWarningStore.js';
	import { connectionUpdateTrigger } from '$lib/stores/connectionUpdateStore.js';
	import { collectionsOpen } from '$lib/stores/collectionsStore.js';
	import { warmMonacoEditor } from '$lib/services/monacoPreloader.js';
	// Wails runtime controls
	import {
		WindowIsMaximised,
		LogInfo
	} from '$lib/wailsjs/runtime/runtime';
	import { GetDefaultConfig } from '$lib/wailsjs/go/main/App';
	
	let { children } = $props();
	
	let isMax = $state(false);
	let currentTheme = $state('light');
	let lastCheckedPath = $state('');
	
	// Connection status strip state
	let defaultConnectionColor = $state('');
	let defaultConnectionName = $state('');
	let showConnectionStrip = $state(false);
	
	// Pages where the connection warning should not be shown
	const excludedPages = ['/about', '/connections'];
	
	// Pages where the connection status strip should not be shown
	const stripExcludedPages = ['/about', '/connections'];
	
	// Subscribe to connection status from store
	const connectionStatus = $derived($connectionWarningStatus);
	const updateTrigger = $derived($connectionUpdateTrigger);
	
	// Check if current page should show the connection warning
	const shouldShowWarning = $derived(() => {
		// Don't hide warning immediately if we're checking the connection
		// to prevent flicker during retry operations
		if (connectionStatus.isChecking && connectionStatus.connectionState !== 'working') {
			return !excludedPages.includes(page.url.pathname);
		}
		
		if (connectionStatus.connectionState === 'working') return false;
		return !excludedPages.includes(page.url.pathname);
	});
	
	// Check if current page should show the connection status strip
	const shouldShowStatusStrip = $derived(() => {
		return !stripExcludedPages.includes(page.url.pathname) && 
		       connectionStatus.connectionState === 'working' && 
		       connectionStatus.hasDefault && 
		       defaultConnectionColor && 
		       defaultConnectionColor.trim() !== '';
	});
	
	// Reactive check when page changes
	$effect(() => {
		const currentPath = page.url.pathname;
		
		// Check connection when:
		// 1. Coming from connections page to any other page
		// 2. Initial load
		// 3. Navigating to non-excluded pages
		if (
			lastCheckedPath === '/connections' && currentPath !== '/connections' ||
			lastCheckedPath === '' ||
			(!excludedPages.includes(currentPath) && hasWails())
		) {
			checkDefaultConnection();
		}
		
		lastCheckedPath = currentPath;
	});
	
	// Reactive effect to monitor connection status changes
	$effect(() => {
		// When connection status changes to working, load the default connection data
		if (connectionStatus.connectionState === 'working' && connectionStatus.hasDefault) {
			loadDefaultConnectionData();
		} else if (connectionStatus.connectionState !== 'working') {
			clearConnectionStripData();
		}
	});
	
	// Listen for connection update triggers from other components
	$effect(() => {
		// When trigger changes, reload connection data if we have a working connection
		if (updateTrigger > 0 && connectionStatus.connectionState === 'working' && hasWails()) {
			loadDefaultConnectionData();
		}
	});
	
	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
	}
	
	async function checkDefaultConnection() {
		if (!hasWails()) return;
		
		try {
			triggerConnectionCheck();
			// Try to use the new TestDefaultConnection method when available
			// @ts-ignore - Wails runtime methods
			if (window.go?.main?.App?.TestDefaultConnection) {
				// @ts-ignore - Wails runtime methods
				const result = await window.go.main.App.TestDefaultConnection();
				if (result.success) {
					updateConnectionWarningStatus(true, true, false, '', 'working');
					// Load default connection data for the strip
					await loadDefaultConnectionData();
				} else if (result.error_code === 'NO_DEFAULT_CONNECTION') {
					updateConnectionWarningStatus(false, false, false, result.message, 'missing');
					clearConnectionStripData();
				} else {
					updateConnectionWarningStatus(true, false, false, result.message, 'failed');
					clearConnectionStripData();
				}
			} else {
				// Fallback to the existing method
				const { GetDefaultConfig } = await import('$lib/wailsjs/go/main/App');
				const defaultConfig = await GetDefaultConfig();
				updateConnectionWarningStatus(true, true, false, '', 'working');
				// Update strip data with the loaded config
				if (defaultConfig) {
					defaultConnectionColor = defaultConfig.env_indicator_color || '';
					defaultConnectionName = defaultConfig.connection_name || '';
				}
			}
		} catch (error) {
			updateConnectionWarningStatus(false, false, false, 'No default connection found', 'missing');
			clearConnectionStripData();
		}
	}
	
	// Function to load default connection data for the status strip
	async function loadDefaultConnectionData() {
		if (!hasWails()) return;
		
		try {
			const defaultConfig = await GetDefaultConfig();
			if (defaultConfig) {
				defaultConnectionColor = defaultConfig.env_indicator_color || '';
				defaultConnectionName = defaultConfig.connection_name || '';
			} else {
				clearConnectionStripData();
			}
		} catch (error) {
			clearConnectionStripData();
		}
	}
	
	// Function to clear connection strip data
	function clearConnectionStripData() {
		defaultConnectionColor = '';
		defaultConnectionName = '';
	}
	
	// Function to convert environment color names to actual color values
	function getEnvironmentColorValue(colorName) {
		const environmentColors = [
			{ name: 'Red', value: 'red', color: '#ef4444' },
			{ name: 'Orange', value: 'orange', color: '#f97316' },
			{ name: 'Yellow', value: 'yellow', color: '#eab308' },
			{ name: 'Green', value: 'green', color: '#22c55e' },
			{ name: 'Dodger Blue', value: 'dodgerblue', color: '#3b82f6' },
			{ name: 'Purple', value: 'purple', color: '#a855f7' },
			{ name: 'Pink', value: 'pink', color: '#ec4899' }
		];
		
		const colorObj = environmentColors.find(c => c.value === colorName);
		return colorObj ? colorObj.color : '#3b82f6'; // Default to dodger blue
	}
	
	// Function to handle when user sets up a connection
	function handleConnectionConfigured() {
		// Re-check the connection status after user might have configured one
		checkDefaultConnection();
	}
	
	onMount(() => {
		// Start warming Monaco Editor in the background for better performance
		warmMonacoEditor();
		
		// Window initialization for Wails
		if (hasWails()) {
			// Assume the window starts in a normal (not maximized) state.
			isMax = false;
			LogInfo('UI assumes initial state is not maximized.');
	
			// For debugging, we can still check the actual state.
			try {
				const actualState = !!WindowIsMaximised();
				if (actualState) {
					LogInfo('Warning: Window actually started in a maximized state.');
				}
			} catch (e) {
				LogInfo('Could not check initial window state: ' + e);
			}
			
			// Check for default connection
			checkDefaultConnection();
		}
		
		// Theme initialization
		if (typeof window !== 'undefined') {
			// Get theme from localStorage or default to 'light'
			const savedTheme = localStorage.getItem('theme') || 'light';
			console.log('Initializing theme:', savedTheme);
			console.log('HTML classes before:', document.documentElement.className);
			currentTheme = savedTheme;
			theme.set(savedTheme);
			
			// Apply theme to DOM by adding/removing dark class to html element
			if (savedTheme === 'dark') {
				document.documentElement.classList.add('dark');
				console.log('Applied dark theme');
			} else {
				document.documentElement.classList.remove('dark');
				console.log('Applied light theme');
			}
			console.log('HTML classes after:', document.documentElement.className);
		}
		
		// Subscribe to theme changes
		const unsubscribe = theme.subscribe(value => {
			console.log('Theme changed to:', value);
			console.log('HTML classes before change:', document.documentElement.className);
			currentTheme = value;
			if (typeof window !== 'undefined') {
				if (value === 'dark') {
					document.documentElement.classList.add('dark');
					console.log('DOM updated to dark theme');
				} else {
					document.documentElement.classList.remove('dark');
					console.log('DOM updated to light theme');
				}
				console.log('HTML classes after change:', document.documentElement.className);
			}
		});
		
		// Cleanup subscription on destroy
		return unsubscribe;
	});
</script>
<div class="flex h-screen theme-bg-primary transition-colors duration-300">
	<!-- Sidebar Menu Component -->
	<SidebarMenu bind:currentTheme />
	
	<!-- Main Content -->
	<main class="flex-1 overflow-auto p-6 pt-12 transition-all duration-300 relative">
		{#if shouldShowWarning()}
			<div class="transition-all duration-500 ease-in-out">
				<ConnectionWarning 
					onRedirect={handleConnectionConfigured}
					connectionState={connectionStatus.connectionState}
					errorMessage={connectionStatus.errorMessage}
				/>
			</div>
		{:else}
			<div class="transition-all duration-500 ease-in-out">
				{@render children?.()}
			</div>
		{/if}
	</main>
	
	<!-- Connection Status Strip - fixed to bottom of viewport -->
	{#if shouldShowStatusStrip()}
		<div 
			class="fixed bottom-0 left-0 right-0 h-[10px] md:h-[10px] sm:h-[8px] transition-all duration-300 ease-in-out hover:h-[12px] sm:hover:h-[10px] z-50"
			style="background-color: {getEnvironmentColorValue(defaultConnectionColor)};"
			role="status"
			aria-label="Connected to {defaultConnectionName} ({defaultConnectionColor} environment)"
			title="Connected to {defaultConnectionName} ({defaultConnectionColor} environment)"
		>
			<!-- Subtle glow effect -->
			<div 
				class="absolute -top-[2px] left-0 right-0 h-[2px] opacity-30 transition-all duration-300 hover:opacity-50 hover:h-[3px]"
				style="background-color: {getEnvironmentColorValue(defaultConnectionColor)};"
			></div>
		</div>
	{/if}
</div>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Window Controls Component -->
<WindowControls bind:isMax />
