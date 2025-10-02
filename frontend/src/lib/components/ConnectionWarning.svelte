<script>
	import { goto } from '$app/navigation';
	
	let { onRedirect = () => {} } = $props();
	
	function handleConnectionRedirect() {
		onRedirect();
		goto('/connections');
	}
</script>

<div class="flex flex-col items-center justify-center min-h-[50vh] p-8 text-center theme-bg-primary transition-colors duration-300">
	<!-- Connection Lost Icon -->
	<div class="mb-6">
		<img 
			src="/icons/connection_lost.svg" 
			alt="Connection Lost" 
			class="w-24 h-24 opacity-60 dark:opacity-80"
		/>
	</div>
	
	<!-- Warning Message -->
	<div class="mb-8 max-w-md">
		<h2 class="text-2xl font-semibold theme-text-primary mb-4">
			No Default Connection Found
		</h2>
		<p class="theme-text-secondary text-lg leading-relaxed">
			A default Elasticsearch connection is required for this application to work properly. 
			Please configure a connection and set it as default to continue.
		</p>
	</div>
	
	<!-- Redirect Button -->
	<button 
		onclick={handleConnectionRedirect}
		class="px-6 py-3 bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 text-white font-medium rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
	>
		Configure Connections
	</button>
</div>