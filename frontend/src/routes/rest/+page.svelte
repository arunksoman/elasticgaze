<script>
	import RestRequestForm from '$lib/components/RestRequestForm.svelte';
	import RequestTabs from '$lib/components/RequestTabs.svelte';
	import ResponseViewer from '$lib/components/ResponseViewer.svelte';
	
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
		
		// Placeholder for REST API request functionality
		console.log('REST Request:', requestMethod, requestEndpoint, requestBody);
		
		// Simulate API call delay
		setTimeout(() => {
			responseData = JSON.stringify({
				"took": 5,
				"timed_out": false,
				"_shards": {
					"total": 1,
					"successful": 1,
					"skipped": 0,
					"failed": 0
				},
				"hits": {
					"total": {
						"value": 10000,
						"relation": "gte"
					},
					"max_score": 1.0,
					"hits": [
						{
							"_index": "test-index",
							"_type": "_doc",
							"_id": "1",
							"_score": 1.0,
							"_source": {
								"title": "Sample Document",
								"content": "This is a sample Elasticsearch document"
							}
						}
					]
				}
			}, null, 2);
			isLoading = false;
		}, 1000);
	}
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-6 theme-text-primary">REST API</h1>
	
	<RestRequestForm 
		bind:method 
		bind:endpoint 
		{isLoading}
		on:send={handleRequest}
	/>
	
	<RequestTabs 
		bind:params 
		bind:requestBody 
		on:paramsChange={handleParamsChange}
	/>
	
	<ResponseViewer {responseData} />
</div>