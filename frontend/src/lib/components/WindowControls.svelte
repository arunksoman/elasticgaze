<script>
	import { goto } from '$app/navigation';
	import {
		WindowMinimise,
		WindowToggleMaximise,
		WindowIsMaximised,
		Quit,
		LogInfo
	} from '$lib/wailsjs/runtime/runtime';

	let { isMax = $bindable(false) } = $props();

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

<!-- Window Controls (top-right) -->
<div class="fixed top-2 right-2 flex gap-2 items-center z-[1000]" style="-webkit-app-region: no-drag;" aria-label="Window controls">
	<!-- Connection Button -->
	<div class="mr-5">
		<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title="Elasticsearch Connections" onclick={handleConnection} aria-label="Elasticsearch Connections">
			<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-icon); mask-image: url('/icons/connect_elastic.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/connect_elastic.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label="Connection icon"></span>
		</button>
	</div>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title="Minimize" onclick={handleMinimise} aria-label="Minimize">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-icon); mask-image: url('/icons/minimize.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/minimize.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label="Minimize icon"></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-sm hover:bg-black/[0.02] dark:hover:bg-white/[0.05] active:translate-y-[0.5px]" title={isMax ? 'Restore' : 'Maximize'} onclick={handleToggleMaximise} aria-label={isMax ? 'Restore' : 'Maximize'}>
		<span class="w-4 h-4 inline-block transition-colors duration-300" style={`background-color: var(--window-control-icon); mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/${isMax ? 'restore' : 'maximize'}.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;`} onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-icon)'} role="img" aria-label={isMax ? 'Restore icon' : 'Maximize icon'}></span>
	</button>
	<button class="appearance-none border-none outline-none p-1.5 rounded-md bg-transparent cursor-pointer flex items-center justify-center transition-all duration-120 hover:shadow-red-500/35 hover:bg-red-500/[0.06] active:translate-y-[0.5px]" title="Close" onclick={handleClose} aria-label="Close">
		<span class="w-4 h-4 inline-block transition-colors duration-300" style="background-color: var(--window-control-close-icon); mask-image: url('/icons/close.svg'); mask-repeat: no-repeat; mask-position: center; mask-size: contain; -webkit-mask-image: url('/icons/close.svg'); -webkit-mask-repeat: no-repeat; -webkit-mask-position: center; -webkit-mask-size: contain;" onmouseenter={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon-hover)'} onmouseleave={(e) => e.target.style.backgroundColor = 'var(--window-control-close-icon)'} role="img" aria-label="Close icon"></span>
	</button>
	<div class="absolute -top-2 -right-2 -bottom-2 -left-2 pointer-events-none" aria-hidden="true"></div>
</div>