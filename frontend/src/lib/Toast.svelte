<script>
	import { createEventDispatcher, onDestroy } from 'svelte';
	
	export let message = '';
	export let type = 'success'; // 'success', 'error', 'info', 'warning'
	export let duration = 1500; // milliseconds
	export let show = false;
	export let animation = 'fade'; // 'fade', 'slide'
	export let errorCode = ''; // Optional error code
	export let errorDetails = ''; // Detailed error information
	
	const dispatch = createEventDispatcher();
	
	let toastElement;
	let timer;
	let isVisible = false;
	
	// Watch show prop and handle animations
	$: if (show) {
		showToast();
	} else {
		isVisible = false;
	}
	
	function showToast() {
		isVisible = true;
		
		// Auto-hide after duration if duration > 0
		if (duration > 0) {
			clearTimeout(timer);
			timer = setTimeout(() => {
				hideToast();
			}, duration);
		}
	}
	
	function hideToast() {
		isVisible = false;
		// Wait for animation to complete before setting show to false
		setTimeout(() => {
			show = false;
			dispatch('hide');
		}, 300); // Match CSS transition duration
	}
	
	// Cleanup timer when component is destroyed
	onDestroy(() => {
		clearTimeout(timer);
	});
	
	// Get the appropriate styling based on type
	function getToastClass(type, animation, isVisible) {
		const baseClasses = 'fixed bottom-4 left-4 px-4 py-3 rounded-lg shadow-lg z-50 max-w-sm transform transition-all duration-300 ease-in-out';
		
		const typeClasses = {
			success: 'bg-green-600 text-white border-l-4 border-green-400',
			error: 'bg-red-600 text-white border-l-4 border-red-400',
			info: 'bg-blue-600 text-white border-l-4 border-blue-400',
			warning: 'bg-yellow-600 text-black border-l-4 border-yellow-400'
		};
		
		let animationClasses = '';
		if (animation === 'fade') {
			animationClasses = isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95';
		} else if (animation === 'slide') {
			animationClasses = isVisible ? 'translate-x-0 opacity-100' : '-translate-x-full opacity-0';
		}
		
		return `${baseClasses} ${typeClasses[type] || typeClasses.success} ${animationClasses}`;
	}
	
	// Get the appropriate icon based on type
	function getIcon(type) {
		const icons = {
			success: '✓',
			error: '✕',
			info: 'ℹ',
			warning: '⚠'
		};
		return icons[type] || icons.success;
	}
</script>

{#if show}
	<div 
		bind:this={toastElement}
		class={getToastClass(type, animation, isVisible)}
		role="alert"
		aria-live="polite"
	>
		<div class="flex items-start gap-2">
			<span class="text-lg font-semibold flex-shrink-0 mt-0.5">{getIcon(type)}</span>
			<div class="flex-1 min-w-0">
				<div class="text-sm font-medium break-words">{message}</div>
				{#if errorCode}
					<div class="text-xs opacity-90 mt-1">
						<span class="font-mono bg-black bg-opacity-20 px-1.5 py-0.5 rounded">
							Code: {errorCode}
						</span>
					</div>
				{/if}
				{#if errorDetails && type === 'error'}
					<details class="mt-2 text-xs">
						<summary class="cursor-pointer hover:opacity-80 select-none">
							Show details
						</summary>
						<div class="mt-1 p-2 bg-black bg-opacity-20 rounded text-xs font-mono whitespace-pre-wrap break-all">
							{errorDetails}
						</div>
					</details>
				{/if}
			</div>
			<button 
				on:click={hideToast}
				class="text-lg hover:opacity-75 transition-opacity flex-shrink-0 ml-2"
				aria-label="Close toast"
			>
				×
			</button>
		</div>
	</div>
{/if}

<style>
	/* Animation classes for smooth enter/exit */
	:global(.toast-enter) {
		transform: translateX(-100%);
		opacity: 0;
	}
	
	:global(.toast-enter-active) {
		transform: translateX(0);
		opacity: 1;
		transition: all 0.3s ease-out;
	}
	
	:global(.toast-exit) {
		transform: translateX(0);
		opacity: 1;
	}
	
	:global(.toast-exit-active) {
		transform: translateX(-100%);
		opacity: 0;
		transition: all 0.3s ease-in;
	}
</style>