<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	// Wails runtime controls
	import {
		WindowMinimise,
		WindowIsNormal,
		WindowToggleMaximise,
		WindowIsMaximised,
		Quit,
		LogInfo
	} from '$lib/wailsjs/runtime/runtime';
	
	let { children } = $props();
	
	let isMax = $state(false);
	
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
	import { goto } from '$app/navigation';
	import { theme } from '$lib/theme.js';
	
	const HomeIcon = '/icons/home.svg';
	const HamburgerIcon = '/icons/hamburger.svg';
	const NodesIcon = '/icons/nodes.svg';
	const ShardsIcon = '/icons/shards.svg';
	const IndicesIcon = '/icons/index.svg';
	const SearchIcon = '/icons/search.svg';
	const RestIcon = '/icons/rest.svg';
	const SnapshotIcon = '/icons/snapshot.svg';
	const SettingsIcon = '/icons/settings.svg';

	let expanded = $state(false);
	let hoverIndex = null;
	let showSettings = $state(false);
	
	// Reactive theme value
	let currentTheme = $state('light');

	const menu = [
	  { name: 'Home', icon: HomeIcon, route: '/' },
	  { name: 'Nodes', icon: NodesIcon, route: '/nodes' },
	  { name: 'Shards', icon: ShardsIcon, route: '/shards' },
	  { name: 'Indices', icon: IndicesIcon, route: '/indices' },
	  { name: 'Search', icon: SearchIcon, route: '/search' },
	  { name: 'Rest', icon: RestIcon, route: '/rest' },
	  { name: 'Snapshot', icon: SnapshotIcon, route: '/snapshot' }
	];

	const activeIndex = $derived(
		menu.findIndex((item) => item.route === page.url.pathname)
	);

	function handleMenuClick(idx) {
	  goto(menu[idx].route);
	}

	function toggleSidebar() {
	  expanded = !expanded;
	}

	function handleHover(idx) {
	  hoverIndex = idx;
	  expanded = true;
	}

	function handleMouseLeave() {
	  hoverIndex = null;
	  expanded = false;
	}

	function openSettings() {
	  showSettings = true;
	}

	function closeSettings() {
	  showSettings = false;
	}

	function setTheme(t) {
		console.log('setTheme called with:', t);
		currentTheme = t;
		theme.set(t);
		
		// Apply theme to DOM by adding/removing dark class to html element
		if (typeof window !== 'undefined') {
			if (t === 'dark') {
				document.documentElement.classList.add('dark');
				console.log('setTheme: Applied dark theme to DOM');
			} else {
				document.documentElement.classList.remove('dark');
				console.log('setTheme: Applied light theme to DOM');
			}
			// Save to localStorage
			localStorage.setItem('theme', t);
			console.log('Theme saved to localStorage:', t);
		}
		closeSettings();
	}
</script>
<div class="flex h-screen theme-bg-primary transition-colors duration-300">
	<!-- Sidebar -->
	<nav class={`flex flex-col justify-between h-full py-4 px-2 theme-bg-secondary shadow-lg transition-all duration-300 ${expanded ? 'w-56' : 'w-16'} relative z-10`}>
		<div>
			<!-- Hamburger -->
					<button
						class="flex items-center justify-center w-12 h-12 mb-2 rounded-lg transition theme-text-primary theme-hover"
						onclick={toggleSidebar}
						aria-label="Toggle menu"
					>
						<img src={HamburgerIcon} alt="Menu" class="w-6 h-6 transition-all duration-300 theme-icon" />
					</button>
			<!-- Menu Items -->
			<ul class="mt-2">
				{#each menu as item, idx}
					<li>
									<button
										class={`relative flex items-center w-full h-12 mb-1 rounded-lg transition-colors group theme-hover ${activeIndex === idx ? 'text-purple-600' : 'theme-text-primary'}`}
										onclick={() => handleMenuClick(idx)}
										onmouseenter={() => handleHover(idx)}
										onmouseleave={handleMouseLeave}
									>
										{#if activeIndex === idx}
											<div class="absolute top-0 left-0 right-0 h-0.5 bg-purple-600"></div>
										{/if}
										<img
											src={item.icon}
											alt={item.name}
											class={`w-6 h-6 mx-4 transition-all duration-300 ${activeIndex === idx ? '' : 'theme-icon'}`}
											style={activeIndex === idx ? 'filter: brightness(0) saturate(100%) invert(34%) sepia(98%) saturate(2546%) hue-rotate(259deg) brightness(99%) contrast(92%);' : ''}
										/>
										{#if expanded}
											<span class="ml-2 text-base font-medium whitespace-nowrap">{item.name}</span>
										{/if}
									</button>
					</li>
				{/each}
			</ul>
		</div>
		<!-- Settings Button -->
		<div class="mb-2">
					<button
						class="flex items-center w-full h-12 rounded-lg transition group theme-text-primary theme-hover"
						onclick={openSettings}
						aria-label="Settings"
					>
						<img src={SettingsIcon} alt="Settings" class="w-6 h-6 mx-4 transition-all duration-300 theme-icon" />
						{#if expanded}
							<span class="ml-2 text-base font-medium whitespace-nowrap">Settings</span>
						{/if}
					</button>
		</div>
		<!-- Settings Popup -->
		{#if showSettings}
			<div class="absolute left-0 bottom-16 w-48 shadow-xl rounded-lg p-4 z-20 theme-bg-secondary">
				<div class="flex flex-col gap-2">
					<span class="font-semibold mb-2 theme-text-primary">Choose Theme</span>
					<button 
						class={`py-2 px-4 rounded transition ${
							currentTheme === 'light' 
								? 'bg-purple-500 text-white hover:bg-purple-600' 
								: 'theme-bg-tertiary theme-text-secondary hover:bg-purple-100'
						}`} 
						onclick={() => setTheme('light')}
					>
						Light
					</button>
					<button 
						class={`py-2 px-4 rounded transition ${
							currentTheme === 'dark' 
								? 'bg-purple-500 text-white hover:bg-purple-600' 
								: 'theme-bg-tertiary theme-text-secondary hover:bg-purple-100'
						}`}
						onclick={() => setTheme('dark')}
					>
						Dark
					</button>
					<button 
						class="mt-2 text-xs transition theme-text-muted hover:theme-text-secondary"
						onclick={closeSettings}
					>
						Close
					</button>
				</div>
			</div>
		{/if}
	</nav>
		<!-- Main Content -->
		<main class="flex-1 overflow-auto p-6 transition-all duration-300">
			{@render children?.()}
		</main>
</div>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>


<!-- Window Controls (top-right) -->
<div class="fixed top-2 right-2 flex gap-2 items-center z-[1000]" style="-webkit-app-region: no-drag;" aria-label="Window controls">
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title="Minimize" onclick={handleMinimise} aria-label="Minimize">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-icon); mask-image: url('/icons/minimize.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/minimize.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title={isMax ? 'Restore' : 'Maximize'} onclick={handleToggleMaximise} aria-label={isMax ? 'Restore' : 'Maximize'}>
		<span class="w-4 h-4 inline-block transition-colors duration-300" style={`background-color: var(--window-control-icon); mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;`} onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-red-500/35 hover:bg-red-500/[0.06] active:translate-y-[0.5px]" title="Close" onclick={handleClose} aria-label="Close">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-close-icon); mask-image: url('/icons/close.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/close.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon)'}></span>
	</button>
	<div class="absolute -top-2 -right-2 -bottom-2 -left-2 pointer-events-none" aria-hidden="true"></div>
</div>
