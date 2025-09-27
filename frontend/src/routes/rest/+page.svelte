<script>
	// REST page component
	let method = 'GET';
	let endpoint = '';
	let requestBody = '{}';
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
			"acknowledged": true,
			"shards_acknowledged": true,
			"index": "test-index"
		}, null, 2);
	}
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-4 theme-text-primary">REST API</h1>
	<div class="theme-bg-secondary p-6 shadow-fluent rounded-lg">
		<div class="mb-4">
			<div class="flex gap-2 mb-2">
				<select 
					bind:value={method}
					class={`border theme-border p-2 theme-bg-tertiary rounded font-semibold ${getMethodColor(method)}`}
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
					class="flex-1 border theme-border p-2 theme-bg-tertiary theme-text-primary rounded"
				/>
				<button 
					onclick={handleRequest}
					class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4 rounded"
				>
					Send
				</button>
			</div>
			
			<div class="mb-2">
				<label class="block mb-1 theme-text-primary font-medium">Request Body</label>
				<textarea 
					bind:value={requestBody}
					rows="5"
					class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary font-mono rounded"
				></textarea>
			</div>
			
			{#if responseData}
				<div>
					<label class="block mb-1 theme-text-primary font-medium">Response</label>
					<pre class="border theme-border p-2 theme-bg-tertiary theme-text-primary font-mono overflow-auto rounded">{responseData}</pre>
				</div>
			{/if}
		</div>
	</div>
</div>