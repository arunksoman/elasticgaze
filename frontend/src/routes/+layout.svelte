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
		theme.set(t);
		document.documentElement.classList.toggle('dark', t === 'dark');
		// Save to localStorage
		if (typeof window !== 'undefined') {
			localStorage.setItem('theme', t);
		}
		closeSettings();
	}

	onMount(() => {
		// Load theme from localStorage or default to 'light'
		if (typeof window !== 'undefined') {
			const savedTheme = localStorage.getItem('theme') || 'light';
			theme.set(savedTheme);
			document.documentElement.classList.toggle('dark', savedTheme === 'dark');
		}
	});
</script>
<div class={`flex h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300`}>
	<!-- Sidebar -->
	<nav class={`flex flex-col justify-between h-full py-4 px-2 bg-white dark:bg-gray-800 shadow-lg transition-all duration-300 ${expanded ? 'w-56' : 'w-16'} relative z-10`}>
		<div>
			<!-- Hamburger -->
					<button
						class={`flex items-center justify-center w-12 h-12 mb-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition`}
						onclick={toggleSidebar}
						aria-label="Toggle menu"
					>
						<img src={HamburgerIcon} alt="Menu" class="w-6 h-6 dark:invert transition-all duration-300" />
					</button>
			<!-- Menu Items -->
			<ul class="mt-2">
				{#each menu as item, idx}
					<li>
									<button
										class={`relative flex items-center w-full h-12 mb-1 rounded-lg transition-colors group
											${activeIndex === idx ? 'text-purple-600 dark:text-purple-400' : 'text-gray-600 dark:text-gray-300'}
											hover:bg-gray-100 dark:hover:bg-gray-700`}
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
											class={`w-6 h-6 mx-4 transition-all duration-300 ${activeIndex === idx ? '' : 'dark:invert'}`}
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
						class="flex items-center w-full h-12 rounded-lg transition group hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-300"
						onclick={openSettings}
						aria-label="Settings"
					>
						<img src={SettingsIcon} alt="Settings" class="w-6 h-6 mx-4 dark:invert transition-all duration-300" />
						{#if expanded}
							<span class="ml-2 text-base font-medium whitespace-nowrap">Settings</span>
						{/if}
					</button>
		</div>
		<!-- Settings Popup -->
		{#if showSettings}
			<div class="absolute left-0 bottom-16 w-48 bg-white dark:bg-gray-800 shadow-xl rounded-lg p-4 z-20">
				<div class="flex flex-col gap-2">
					<span class="font-semibold mb-2 text-gray-700 dark:text-gray-200">Choose Theme</span>
					  <button class="py-2 px-4 rounded bg-gray-100 dark:bg-gray-700 hover:bg-purple-100 dark:hover:bg-purple-900 transition" onclick={() => setTheme('light')}>Light</button>
					  <button class="py-2 px-4 rounded bg-gray-100 dark:bg-gray-700 hover:bg-purple-100 dark:hover:bg-purple-900 transition" onclick={() => setTheme('dark')}>Dark</button>
					  <button class="mt-2 text-xs text-gray-500 hover:text-gray-700 dark:hover:text-gray-300" onclick={closeSettings}>Close</button>
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
		<span class="w-4 h-4 inline-block bg-[#222] dark:bg-[#ddd] hover:bg-[#444] dark:hover:bg-[#bbb] transition-colors duration-300" style="mask-image: url('/icons/minimize.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/minimize.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;"></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title={isMax ? 'Restore' : 'Maximize'} onclick={handleToggleMaximise} aria-label={isMax ? 'Restore' : 'Maximize'}>
		<span class="w-4 h-4 inline-block bg-[#222] dark:bg-[#ddd] hover:bg-[#444] dark:hover:bg-[#bbb] transition-colors duration-300" style={`mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;`}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-red-500/35 hover:bg-red-500/[0.06] active:translate-y-[0.5px]" title="Close" onclick={handleClose} aria-label="Close">
		<span class="w-4 h-4 inline-block bg-[#222] dark:bg-[#ddd] hover:bg-red-500 transition-colors duration-300" style="mask-image: url('/icons/close.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/close.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;"></span>
	</button>
	<div class="absolute -top-2 -right-2 -bottom-2 -left-2 pointer-events-none" aria-hidden="true"></div>
</div>
