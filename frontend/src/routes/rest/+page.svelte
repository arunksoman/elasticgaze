<script>
	import MonacoEditor from '$lib/MonacoEditor.svelte';
	import { theme } from '$lib/theme.js';
	
	// REST page component
	let method = 'GET';
	let endpoint = '';
	let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	let responseData = '';
	
	// Postman-style method colors
	function getMethodColor(methodName) {
		const colors = {
			'GET': 'text-green-500',
			'POST': 'text-orange-500',
			'PUT': 'text-blue-500',
			'DELETE': 'text-red-500',
			// 'PATCH': 'text-purple-500',
			// 'HEAD': 'text-gray-500',
			// 'OPTIONS': 'text-indigo-500'
		};
		return colors[methodName] || 'text-gray-400';
	}
	
	function handleRequest() {
		// Placeholder for REST API request functionality
		console.log('REST Request:', method, endpoint, requestBody);
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
	}
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-6 theme-text-primary">REST API</h1>
	
	<!-- URL and Method Section -->
	<div class="flex gap-3 mb-6">
		<select 
			bind:value={method}
			class={`border theme-border p-3 theme-bg-tertiary rounded font-semibold ${getMethodColor(method)} min-w-[100px]`}
		>
			<option value="GET" class="text-green-500 font-semibold">GET</option>
			<option value="POST" class="text-orange-500 font-semibold">POST</option>
			<option value="PUT" class="text-blue-500 font-semibold">PUT</option>
			<option value="DELETE" class="text-red-500 font-semibold">DELETE</option>
			<!-- <option value="PATCH" class="text-purple-500 font-semibold">PATCH</option> -->
			<!-- <option value="HEAD" class="text-gray-500 font-semibold">HEAD</option> -->
			<!-- <option value="OPTIONS" class="text-indigo-500 font-semibold">OPTIONS</option> -->
		</select>
		<input 
			type="text" 
			bind:value={endpoint}
			placeholder="/_cluster/health"
			class="flex-1 border theme-border p-3 theme-bg-tertiary theme-text-primary rounded"
		/>
		<button 
			onclick={handleRequest}
			class="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded font-medium transition-colors"
		>
			Send
		</button>
	</div>
	
	<!-- Request Body Section -->
	<div class="mb-6">
		<span class="block mb-3 theme-text-primary font-medium text-lg">Request Body</span>
		<MonacoEditor 
			bind:value={requestBody} 
			language="json" 
			height="200px"
			theme={$theme}
		/>
	</div>
	
	<!-- Response Section -->
	{#if responseData}
		<div>
			<span class="block mb-3 theme-text-primary font-medium text-lg">Response</span>
			<MonacoEditor 
				bind:value={responseData} 
				language="json" 
				height="300px"
				readOnly={true}
				theme={$theme}
			/>
		</div>
	{/if}
</div>