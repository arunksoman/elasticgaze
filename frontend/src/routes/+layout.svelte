<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { theme } from '$lib/theme.js';
	import WindowControls from '$lib/components/WindowControls.svelte';
	import SidebarMenu from '$lib/components/SidebarMenu.svelte';
	// Wails runtime controls
	import {
		WindowIsMaximised,
		LogInfo
	} from '$lib/wailsjs/runtime/runtime';
	
	let { children } = $props();
	
	let isMax = $state(false);
	let currentTheme = $state('light');
	
	function hasWails() {
		return typeof window !== 'undefined' && !!window.runtime;
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
		{@render children?.()}
	</main>
</div>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Window Controls Component -->
<WindowControls bind:isMax />
