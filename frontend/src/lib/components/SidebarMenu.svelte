<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { theme } from '$lib/theme.js';
	import { onDestroy } from 'svelte';

	let { currentTheme = $bindable('light') } = $props();

	const HomeIcon = '/icons/home.svg';
	const AppIcon = '/icons/eg_icon1.svg';
	const NodesIcon = '/icons/nodes.svg';
	const ShardsIcon = '/icons/shards.svg';
	const IndicesIcon = '/icons/index.svg';
	const SearchIcon = '/icons/search.svg';
	const RestIcon = '/icons/rest.svg';
	const SnapshotIcon = '/icons/snapshot.svg';
	const AboutIcon = '/icons/about.svg';
	const SettingsIcon = '/icons/settings.svg';

	let expanded = $state(false);
	let permanentlyExpanded = $state(false);
	let hoverIndex = null;
	let showSettings = $state(false);
	let hoverTimeout = null;

	const menu = [
		{ name: 'Home', icon: HomeIcon, route: '/' },
		{ name: 'Nodes', icon: NodesIcon, route: '/nodes' },
		{ name: 'Shards', icon: ShardsIcon, route: '/shards' },
		{ name: 'Indices', icon: IndicesIcon, route: '/indices' },
		{ name: 'Search', icon: SearchIcon, route: '/search' },
		{ name: 'REST', icon: RestIcon, route: '/rest' },
		{ name: 'Snapshot', icon: SnapshotIcon, route: '/snapshot' },
		{ name: 'About', icon: AboutIcon, route: '/about' }
	];

	const activeIndex = $derived(
		menu.findIndex((item) => item.route === page.url.pathname)
	);

	function handleMenuClick(idx) {
		goto(menu[idx].route);
	}

	function toggleSidebar() {
		permanentlyExpanded = !permanentlyExpanded;
		expanded = permanentlyExpanded;
	}

	function handleHover(idx) {
		if (permanentlyExpanded) return; // Don't auto-expand if permanently expanded
		
		if (hoverTimeout) {
			clearTimeout(hoverTimeout);
			hoverTimeout = null;
		}
		hoverIndex = idx;
		expanded = true;
	}

	function handleMouseLeave() {
		if (permanentlyExpanded) return; // Don't auto-collapse if permanently expanded
		
		// Don't immediately collapse, use a delay to prevent flickering
		if (hoverTimeout) {
			clearTimeout(hoverTimeout);
		}
		
		hoverTimeout = setTimeout(() => {
			hoverIndex = null;
			expanded = false;
			hoverTimeout = null;
		}, 300); // 300ms delay before collapsing
	}

	function handleSidebarEnter() {
		if (permanentlyExpanded) return; // Don't handle hover if permanently expanded
		
		// Clear any pending collapse when entering sidebar area
		if (hoverTimeout) {
			clearTimeout(hoverTimeout);
			hoverTimeout = null;
		}
	}

	function handleSidebarLeave() {
		if (permanentlyExpanded) return; // Don't handle hover if permanently expanded
		
		// Only collapse when leaving the entire sidebar
		if (hoverTimeout) {
			clearTimeout(hoverTimeout);
		}
		
		hoverTimeout = setTimeout(() => {
			hoverIndex = null;
			expanded = false;
			hoverTimeout = null;
		}, 200); // Shorter delay when leaving entire sidebar
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
		// Don't close settings popup to allow multiple theme switches
	}

	// Cleanup timeout on component destroy
	onDestroy(() => {
		if (hoverTimeout) {
			clearTimeout(hoverTimeout);
		}
	});
</script>

<!-- Sidebar -->
<nav 
	class={`flex flex-col h-full py-4 px-2 theme-bg-secondary shadow-lg transition-all duration-300 ${expanded ? 'w-56' : 'w-16'} relative z-10`}
	onmouseenter={handleSidebarEnter}
	onmouseleave={handleSidebarLeave}
>
	<div class="flex-1">
		<!-- App Icon -->
		<button
			class="flex items-center w-full h-12 mb-2 rounded-lg transition theme-text-primary theme-hover"
			onclick={toggleSidebar}
			aria-label="Toggle menu"
		>
			<img src={AppIcon} alt="ElasticGaze" class="w-6 h-6 mx-4 transition-all duration-300" />
			{#if expanded}
				<span class="ml-1 text-lg font-semibold whitespace-nowrap theme-text-primary">ElasticGaze</span>
			{/if}
		</button>
		<!-- Menu Items -->
		<ul class="mt-2">
			{#each menu as item, idx}
				<li>
					<button
						class={`relative flex items-center w-full h-12 mb-1 rounded-lg transition-colors group theme-hover ${activeIndex === idx ? 'text-purple-600' : 'theme-text-primary'}`}
						onclick={() => handleMenuClick(idx)}
						onmouseenter={() => handleHover(idx)}
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
	<!-- Bottom section with Settings and Version -->
	<div class="flex flex-col">
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
		<!-- Version Number -->
		<div class="mb-4 px-4">
			<span class="text-xs theme-text-secondary opacity-70">V0.0.1</span>
		</div>
	</div>
	<!-- Settings Popup -->
	{#if showSettings}
		<div class="absolute left-0 bottom-16 w-48 shadow-xl rounded-lg p-4 z-20 theme-bg-secondary">
			<div class="flex flex-col gap-4">
				<span class="font-semibold mb-2 theme-text-primary">Choose Theme</span>
				
				<!-- Theme Toggle Switch -->
				<div class="flex items-center justify-center">
					<button
						class="relative inline-flex items-center w-16 h-6 bg-gray-300 dark:bg-gray-600 rounded-full p-1 transition-colors duration-300 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
						onclick={() => setTheme(currentTheme === 'light' ? 'dark' : 'light')}
						aria-label="Toggle theme"
					>
						<!-- Toggle slider -->
						<div 
							class={`absolute w-6 h-6 bg-white rounded-full shadow-lg transform transition-transform duration-300 flex items-center justify-center ${
								currentTheme === 'dark' ? 'translate-x-8' : 'translate-x-0'
							}`}
						>
							<!-- Icon inside the toggle -->
							<img 
								src={currentTheme === 'light' ? '/icons/bright.svg' : '/icons/dark.svg'} 
								alt={currentTheme === 'light' ? 'Light mode' : 'Dark mode'}
								class="w-4 h-4"
								style={currentTheme === 'light' 
									? 'filter: brightness(0) saturate(100%) invert(69%) sepia(100%) saturate(1769%) hue-rotate(16deg) brightness(108%) contrast(98%);' 
									: 'filter: brightness(0) saturate(100%) invert(20%) sepia(10%) saturate(500%) hue-rotate(200deg) brightness(95%) contrast(90%);'}
							/>
						</div>
						
						<!-- Background icons -->
						<div class="flex items-center justify-between w-full px-1.5">
							<img 
								src="/icons/bright.svg" 
								alt="Light mode"
								class={`w-3 h-3 transition-opacity duration-300 ${
									currentTheme === 'light' ? 'opacity-0' : 'opacity-60'
								}`}
								style="filter: brightness(0) saturate(100%) invert(100%);"
							/>
							<img 
								src="/icons/dark.svg" 
								alt="Dark mode"
								class={`w-3 h-3 transition-opacity duration-300 ${
									currentTheme === 'dark' ? 'opacity-0' : 'opacity-60'
								}`}
								style="filter: brightness(0) saturate(100%) invert(100%);"
							/>
						</div>
					</button>
				</div>
				
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