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
	
	// Pages where the connection warning should not be shown
	const excludedPages = ['/about', '/connections'];
	
	// Subscribe to connection status from store
	const connectionStatus = $derived($connectionWarningStatus);
	
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
				} else if (result.error_code === 'NO_DEFAULT_CONNECTION') {
					updateConnectionWarningStatus(false, false, false, result.message, 'missing');
				} else {
					updateConnectionWarningStatus(true, false, false, result.message, 'failed');
				}
			} else {
				// Fallback to the existing method
				const { GetDefaultConfig } = await import('$lib/wailsjs/go/main/App');
				await GetDefaultConfig();
				updateConnectionWarningStatus(true, true, false, '', 'working');
			}
		} catch (error) {
			updateConnectionWarningStatus(false, false, false, 'No default connection found', 'missing');
		}
	}
	
	// Function to handle when user sets up a connection
	function handleConnectionConfigured() {
		// Re-check the connection status after user might have configured one
		checkDefaultConnection();
	}
	
	onMount(() => {
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
	<main class="flex-1 overflow-auto p-6 transition-all duration-300">
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
</div>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Window Controls Component -->
<WindowControls bind:isMax />
