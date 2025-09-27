<script>
	// Search page component
	let searchQuery = '';
	let searchResults = [];
	
	function handleSearch() {
		// Placeholder for search functionality
		console.log('Search query:', searchQuery);
		searchResults = [
			{ id: 1, index: 'logs-2023', score: 0.9, highlight: 'Example search <em>result</em> text' }
		];
	}
</script>

<div class="p-6">
	<h1 class="text-2xl font-medium mb-4">Search</h1>
	<div class="bg-fluent-surface-light dark:bg-fluent-surface-dark p-6 shadow-fluent">
		<div class="mb-4">
			<form onsubmit={(e) => { e.preventDefault(); handleSearch(); }} class="flex gap-2">
				<input 
					type="text" 
					bind:value={searchQuery}
					placeholder="Enter search query..."
					class="flex-1 border border-fluent-border-light dark:border-fluent-border-dark p-2 bg-fluent-bg-light dark:bg-fluent-bg-dark"
				/>
				<button 
					type="submit"
					class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4"
				>
					Search
				</button>
			</form>
		</div>
		
		{#if searchResults.length > 0}
			<div class="border border-fluent-border-light dark:border-fluent-border-dark p-4">
				<h2 class="text-lg font-medium mb-2">Results</h2>
				<ul class="divide-y divide-fluent-border-light dark:divide-fluent-border-dark">
					{#each searchResults as result}
						<li class="py-2">
							<div class="text-fluent-purple dark:text-fluent-purple-light">
								Score: {result.score} | Index: {result.index}
							</div>
							<div>{@html result.highlight}</div>
						</li>
					{/each}
				</ul>
			</div>
		{/if}
	</div>
</div>