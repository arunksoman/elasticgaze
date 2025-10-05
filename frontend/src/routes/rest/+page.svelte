<script>
	import RestRequestForm from '$lib/components/RestRequestForm.svelte';
	import RequestBodyEditor from '$lib/components/RequestBodyEditor.svelte';
	import ResponseViewer from '$lib/components/ResponseViewer.svelte';
	
	// REST page component
	let method = 'GET';
	let endpoint = '';
	let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	let responseData = '';
	let isLoading = false;
	
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
	
	<RequestBodyEditor bind:requestBody />
	
	<ResponseViewer {responseData} />
</div>