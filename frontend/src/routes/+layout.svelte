<script>
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onMount } from 'svelte';
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
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children?.()}

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
