<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
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
	import HomeIcon from '../../build/icons/home.svg';
	import HamburgerIcon from '../../build/icons/hamburger.svg';
	import NodesIcon from '../../build/icons/nodes.svg';
	import ShardsIcon from '../../build/icons/shards.svg';
	import IndicesIcon from '../../build/icons/index.svg';
	import SearchIcon from '../../build/icons/search.svg';
	import RestIcon from '../../build/icons/rest.svg';
	import SnapshotIcon from '../../build/icons/snapshot.svg';
	import SettingsIcon from '../../build/icons/settings.svg';

	let expanded = $state(false);
	let hoverIndex = null;
	let showSettings = $state(false);
	let theme = 'light';

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
		menu.findIndex((item) => item.route === $page.url.pathname)
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
	  theme = t;
	  document.documentElement.classList.toggle('dark', t === 'dark');
	  closeSettings();
	}

	onMount(() => {
	  if (theme === 'dark') {
	    document.documentElement.classList.add('dark');
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
						<img src={HamburgerIcon} alt="Menu" class="w-6 h-6" />
					</button>
			<!-- Menu Items -->
			<ul class="mt-2">
				{#each menu as item, idx}
					<li>
									<button
										class={`flex items-center w-full h-12 mb-1 rounded-lg transition group
											${activeIndex === idx ? 'bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300 shadow-md' : 'text-gray-600 dark:text-gray-300'}
											hover:bg-gray-100 dark:hover:bg-gray-700`}
										onclick={() => handleMenuClick(idx)}
										onmouseenter={() => handleHover(idx)}
										onmouseleave={handleMouseLeave}
									>
										<img
											src={item.icon}
											alt={item.name}
											class={`w-6 h-6 mx-4 transition-all duration-300 ${
												activeIndex === idx
													? 'filter drop-shadow-lg'
													: 'filter-none'
											}`}
											style={activeIndex === idx ? 'filter: drop-shadow(0 0 5px #A020F0);' : ''}
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
						<img src={SettingsIcon} alt="Settings" class="w-6 h-6 mx-4" />
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
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] active:translate-y-[0.5px]" title="Minimize" onclick={handleMinimise} aria-label="Minimize">
		<span class="w-4 h-4 inline-block bg-[#222] hover:bg-[#444]" style="mask-image: url('/icons/minimize.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/minimize.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;"></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] active:translate-y-[0.5px]" title={isMax ? 'Restore' : 'Maximize'} onclick={handleToggleMaximise} aria-label={isMax ? 'Restore' : 'Maximize'}>
		<span class="w-4 h-4 inline-block bg-[#222] hover:bg-[#444]" style={`mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;`}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-red-500/35 hover:bg-red-500/[0.06] active:translate-y-[0.5px]" title="Close" onclick={handleClose} aria-label="Close">
		<span class="w-4 h-4 inline-block bg-[#222] hover:bg-red-500" style="mask-image: url('/icons/close.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/close.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;"></span>
	</button>
	<div class="absolute -top-2 -right-2 -bottom-2 -left-2 pointer-events-none" aria-hidden="true"></div>
</div>
