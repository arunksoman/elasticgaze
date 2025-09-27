<script>
	// REST page component
	let method = 'GET';
	let endpoint = '';
	let requestBody = '{}';
	let responseData = '';
	
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
	<h1 class="text-2xl font-medium mb-4 text-gray-900 dark:text-gray-100">REST API</h1>
	<div class="bg-fluent-surface-light dark:bg-fluent-surface-dark p-6 shadow-fluent rounded-lg">
		<div class="mb-4">
			<div class="flex gap-2 mb-2">
				<select 
					bind:value={method}
					class="border border-fluent-border-light dark:border-fluent-border-dark p-2 bg-fluent-bg-light dark:bg-fluent-bg-dark text-gray-900 dark:text-gray-100 rounded"
				>
					<option value="GET">GET</option>
					<option value="POST">POST</option>
					<option value="PUT">PUT</option>
					<option value="DELETE">DELETE</option>
				</select>
				<input 
					type="text" 
					bind:value={endpoint}
					placeholder="/_cluster/health"
					class="flex-1 border border-fluent-border-light dark:border-fluent-border-dark p-2 bg-fluent-bg-light dark:bg-fluent-bg-dark text-gray-900 dark:text-gray-100 rounded"
				/>
				<button 
					onclick={handleRequest}
					class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4 rounded"
				>
					Send
				</button>
			</div>
			
			<div class="mb-2">
				<label class="block mb-1 text-gray-900 dark:text-gray-100 font-medium">Request Body</label>
				<textarea 
					bind:value={requestBody}
					rows="5"
					class="w-full border border-fluent-border-light dark:border-fluent-border-dark p-2 bg-fluent-bg-light dark:bg-fluent-bg-dark text-gray-900 dark:text-gray-100 font-mono rounded"
				></textarea>
			</div>
			
			{#if responseData}
				<div>
					<label class="block mb-1 text-gray-900 dark:text-gray-100 font-medium">Response</label>
					<pre class="border border-fluent-border-light dark:border-fluent-border-dark p-2 bg-fluent-bg-light dark:bg-fluent-bg-dark text-gray-900 dark:text-gray-100 font-mono overflow-auto rounded">{responseData}</pre>
				</div>
			{/if}
		</div>
	</div>
</div>