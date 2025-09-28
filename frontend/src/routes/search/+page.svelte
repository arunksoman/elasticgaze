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
	<h1 class="text-2xl font-medium mb-4 theme-text-primary">Search</h1>
	<div class="mb-4">
		<form onsubmit={(e) => { e.preventDefault(); handleSearch(); }} class="flex gap-2">
			<input 
				type="text" 
				bind:value={searchQuery}
				placeholder="Enter search query..."
				class="flex-1 border theme-border p-2 theme-bg-tertiary theme-text-primary rounded"
			/>
			<button 
				type="submit"
				class="fluent-btn bg-fluent-purple hover:bg-fluent-purple-light text-white px-4 rounded"
			>
				Search
			</button>
		</form>
	</div>
	
	{#if searchResults.length > 0}
		<div class="border theme-border p-4 rounded">
			<h2 class="text-lg font-medium mb-2 theme-text-primary">Results</h2>
			<ul class="divide-y theme-border">
				{#each searchResults as result}
					<li class="py-2">
						<div class="text-fluent-purple dark:text-fluent-purple-light">
							Score: {result.score} | Index: {result.index}
						</div>
						<div class="theme-text-secondary">{@html result.highlight}</div>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>