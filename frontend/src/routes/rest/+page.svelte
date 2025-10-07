<script>
	import RestRequestForm from '$lib/components/RestRequestForm.svelte';
	import RequestTabs from '$lib/components/RequestTabs.svelte';
	import ResponseViewer from '$lib/components/ResponseViewer.svelte';
	import ResizableSplitter from '$lib/components/ResizableSplitter.svelte';
	import { ExecuteElasticsearchRequest } from '$lib/wailsjs/go/main/App.js';
	
	// REST page component
	let method = 'GET';
	let endpoint = '';
	let baseEndpoint = '';
	let params = [];
	let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	let responseData = '';
	let isLoading = false;
	
	// Function to build URL with parameters
	function buildUrlWithParams(base, parameters) {
		if (!base) return '';
		
		const enabledParams = parameters.filter(p => p.enabled && p.key && p.value);
		if (enabledParams.length === 0) return base;
		
		const queryString = enabledParams
			.map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`)
			.join('&');
		
		return base + (base.includes('?') ? '&' : '?') + queryString;
	}
	
	// Function to parse URL and extract base endpoint and parameters
	function parseUrl(url) {
		if (!url) return { base: '', params: [] };
		
		const [base, queryString] = url.split('?');
		if (!queryString) return { base: url, params: [] };
		
		const params = queryString.split('&').map(param => {
			const [key, value] = param.split('=');
			return {
				key: decodeURIComponent(key || ''),
				value: decodeURIComponent(value || ''),
				enabled: true
			};
		}).filter(p => p.key);
		
		return { base, params };
	}
	
	// Watch for endpoint changes to update base and params
	let lastBuiltUrl = '';
	$: {
		const currentBuiltUrl = buildUrlWithParams(baseEndpoint, params);
		if (endpoint !== currentBuiltUrl && endpoint !== lastBuiltUrl) {
			const parsed = parseUrl(endpoint);
			baseEndpoint = parsed.base;
			params = parsed.params;
		}
		lastBuiltUrl = currentBuiltUrl;
	}
	
	// Watch for params changes to update endpoint
	function handleParamsChange(event) {
		params = event.detail;
		const newUrl = buildUrlWithParams(baseEndpoint, params);
		if (newUrl !== endpoint) {
			endpoint = newUrl;
			lastBuiltUrl = newUrl;
		}
	}
	
	// Watch for base endpoint changes in the input
	$: {
		// Extract base endpoint when user types directly in URL input
		if (endpoint && !endpoint.includes('?')) {
			baseEndpoint = endpoint;
		}
	}
	
	function handleRequest(event) {
		const { method: requestMethod, endpoint: requestEndpoint } = event.detail;
		
		// Set loading state
		isLoading = true;
		
		console.log('REST Request:', requestMethod, requestEndpoint, requestBody);
		
		// Prepare the request object for the Go backend
		const esRequest = {
			method: requestMethod,
			endpoint: requestEndpoint,
			body: requestBody && requestBody.trim() ? requestBody : null
		};
		
		// Call the Go function
		ExecuteElasticsearchRequest(esRequest)
			.then(response => {
				console.log('ES Response:', response);
				
				if (response.success) {
					// Parse and pretty-print the response
					try {
						const parsedResponse = JSON.parse(response.response);
						responseData = JSON.stringify(parsedResponse, null, 2);
					} catch (e) {
						// If it's not JSON, just display as-is
						responseData = response.response;
					}
				} else {
					// Show error response
					const errorResponse = {
						error: true,
						status_code: response.status_code,
						error_code: response.error_code,
						error_details: response.error_details
					};
					responseData = JSON.stringify(errorResponse, null, 2);
				}
			})
			.catch(error => {
				console.error('Request failed:', error);
				const errorResponse = {
					error: true,
					message: 'Request failed',
					details: error.toString()
				};
				responseData = JSON.stringify(errorResponse, null, 2);
			})
			.finally(() => {
				isLoading = false;
			});
	}
	
	// Handle splitter resize - dispatch window resize event as backup
	function handleSplitterResize(percentage) {
		// Dispatch window resize event that Monaco's automaticLayout can detect
		setTimeout(() => {
			window.dispatchEvent(new Event('resize'));
		}, 50);
	}
</script>

<div class="h-screen flex flex-col">
	<div class="p-6 pb-0 flex-shrink-0">
		<h1 class="text-2xl font-medium mb-6 theme-text-primary">REST API</h1>
		
		<RestRequestForm 
			bind:method 
			bind:endpoint 
			{isLoading}
			on:send={handleRequest}
		/>
	</div>
	
	<div class="flex-1 p-6 pt-0 min-h-0">
		<ResizableSplitter 
			defaultSplit={45} 
			minSize={25} 
			maxSize={75}
			className="h-full"
			onResize={handleSplitterResize}
		>
			{#snippet panel1()}
				<!-- Request Section -->
				<div class="h-full flex flex-col">
					<RequestTabs 
						bind:params 
						bind:requestBody 
						on:paramsChange={handleParamsChange}
					/>
				</div>
			{/snippet}
			
			{#snippet panel2()}
				<!-- Response Section -->
				<div class="h-full flex flex-col">
					<ResponseViewer {responseData} />
				</div>
			{/snippet}
		</ResizableSplitter>
	</div>
</div>