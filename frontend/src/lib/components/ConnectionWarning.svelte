<script>
	import { goto } from '$app/navigation';
	import { refreshConnectionStatus } from '$lib/stores/connectionWarningStore.js';
	
	let { 
		onRedirect = () => {},
		connectionState = 'missing', // 'missing', 'failed'
		errorMessage = ''
	} = $props();
	
	let isRetrying = $state(false);
	let retrySuccessful = $state(false);
	
	function handleConnectionRedirect() {
		onRedirect();
		goto('/connections');
	}
	
	async function handleRetryConnection() {
		if (isRetrying) return;
		
		try {
			isRetrying = true;
			retrySuccessful = false;
			
			// Add a small delay to show the retrying state
			await new Promise(resolve => setTimeout(resolve, 800));
			await refreshConnectionStatus();
			
			// Mark as successful and add a delay to ensure state stabilizes
			retrySuccessful = true;
			await new Promise(resolve => setTimeout(resolve, 1000));
			
			// If we get here, the retry was truly successful
		} catch (error) {
			console.error('Retry failed:', error);
			retrySuccessful = false;
			// Add a small delay before resetting to prevent flicker
			await new Promise(resolve => setTimeout(resolve, 300));
		} finally {
			isRetrying = false;
		}
	}
	
	// Determine content based on connection state
	const content = $derived.by(() => {
		if (connectionState === 'failed') {
			return {
				title: 'Default Connection Failed',
				message: `The default Elasticsearch connection failed to connect. ${errorMessage || 'Please check your connection settings and try again.'}`,
				buttonText: 'Fix Connection'
			};
		} else if (connectionState === 'missing') {
			return {
				title: 'No Default Connection Found',
				message: 'A default Elasticsearch connection is required for this application to work properly. Please configure a connection and set it as default to continue.',
				buttonText: 'Configure Connections'
			};
		} else {
			// Handle 'unknown' or any other state
			return {
				title: 'Connection Check Required',
				message: 'Checking connection status... If this persists, please configure a default Elasticsearch connection.',
				buttonText: 'Configure Connections'
			};
		}
	});
</script>

<div class="flex flex-col items-center justify-center min-h-[50vh] p-8 text-center theme-bg-primary transition-all duration-500 ease-in-out">
	<!-- Connection Lost Icon -->
	<div class="mb-6 transition-all duration-300">
		<img 
			src="/icons/connection_lost.svg" 
			alt="Connection Lost" 
			class="w-24 h-24 opacity-60 dark:opacity-80"
		/>
	</div>
	
	<!-- Warning Message -->
	<div class="mb-8 max-w-md transition-all duration-300">
		<h2 class="text-2xl font-semibold text-gray-900 dark:text-gray-100 theme-text-primary mb-4">
			{content.title}
		</h2>
		<p class="text-gray-700 dark:text-gray-300 theme-text-secondary text-lg leading-relaxed">
			{content.message}
		</p>
	</div>
	
	<!-- Action Buttons -->
	<div class="flex gap-3 flex-wrap justify-center transition-all duration-300">
		{#if connectionState === 'failed'}
			<!-- Retry Button for Failed Connections -->
			<button 
				onclick={handleRetryConnection}
				disabled={isRetrying}
				class="px-6 py-3 w-44 bg-green-600 hover:bg-green-700 disabled:bg-green-400 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
			>
				<span class="inline-block transition-all duration-200">
					{isRetrying ? 'Retrying...' : 'Retry Connection'}
				</span>
			</button>
		{/if}
		
		<!-- Configure/Fix Button -->
		<button 
			onclick={handleConnectionRedirect}
			disabled={isRetrying}
			class="px-6 py-3 w-44 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 disabled:cursor-not-allowed dark:bg-blue-500 dark:hover:bg-blue-600 text-white font-medium rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
		>
			<span class="inline-block transition-all duration-200">
				{content.buttonText}
			</span>
		</button>
	</div>
</div>