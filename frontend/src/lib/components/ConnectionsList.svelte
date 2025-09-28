<script>
	import ConnectionCard from './ConnectionCard.svelte';
	
	// Props
	export let connections = [];
	export let testingConnectionId = null;
	
	// Event props (Svelte 5 way)
	export let onadd;
	export let onteststart;
	export let ontestend;
	export let onedit;
	export let ondelete;
	export let onsetdefault;
	export let ontoast;

	function addConnection() {
		onadd?.();
	}
</script>

{#if connections.length === 0}
	<div class="text-center py-12">
		<div class="w-16 h-16 mx-auto mb-4 theme-text-secondary">
			<svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
			</svg>
		</div>
		<h3 class="text-lg font-medium theme-text-primary mb-2">No connections configured</h3>
		<p class="theme-text-secondary mb-4">Add your first Elasticsearch connection to get started.</p>
		<button 
			onclick={addConnection}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors flex items-center gap-2"
		>
			<img src="/icons/create.svg" alt="" class="w-4 h-4" style="filter: brightness(0) invert(1);" />
			Add Connection
		</button>
	</div>
{:else}
	<div class="grid gap-4">
		{#each connections as connection (connection.id)}
			<ConnectionCard
				{connection}
				{testingConnectionId}
				onteststart={onteststart}
				ontestend={ontestend}
				onedit={onedit}
				ondelete={ondelete}
				onsetdefault={onsetdefault}
				ontoast={ontoast}
			/>
		{/each}
	</div>
{/if}