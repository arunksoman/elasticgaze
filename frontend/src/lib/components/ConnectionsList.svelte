<script lang="js">
	import ConnectionCard from './ConnectionCard.svelte';
	
	/**
	 * @typedef {Object} Connection
	 * @property {string|null} id - Unique connection identifier (null for new connections)
	 * @property {string} name - Display name for the connection
	 * @property {string} host - Elasticsearch host address
	 * @property {number} port - Elasticsearch port number
	 * @property {boolean} useSSL - Whether to use SSL/HTTPS
	 * @property {('basic'|'apikey'|'none')} authType - Authentication type
	 * @property {string} username - Username for basic auth
	 * @property {string} password - Password for basic auth
	 * @property {string} apiKey - API key for API key auth
	 * @property {boolean} isDefault - Whether this is the default connection
	 * @property {string} environmentColor - Environment color indicator
	 */
	
	/**
	 * @typedef {Object} ToastData
	 * @property {string} message - Toast message to display
	 * @property {('success'|'error'|'warning'|'info')} [type] - Toast type
	 * @property {number} [duration] - Toast duration in milliseconds
	 * @property {string} [animation] - Toast animation type
	 * @property {string} [errorCode] - Error code for error toasts
	 * @property {string} [errorDetails] - Error details for error toasts
	 */
	
	/**
	 * Array of connection objects to display
	 * @type {Array<Connection>}
	 */
	export let connections = [];
	
	/**
	 * ID of the connection currently being tested (for loading state)
	 * @type {string|null}
	 */
	export let testingConnectionId = null;
	
	// Event props (Svelte 5 way)
	/**
	 * Callback function called when add connection button is clicked
	 * @type {function(): void}
	 */
	export let onadd;
	
	/**
	 * Callback function called when connection test starts
	 * @type {function(string): void}
	 */
	export let onteststart;
	
	/**
	 * Callback function called when connection test ends
	 * @type {function(): void}
	 */
	export let ontestend;
	
	/**
	 * Callback function called when edit button is clicked
	 * @type {function(Connection): void}
	 */
	export let onedit;
	
	/**
	 * Callback function called when delete button is clicked
	 * @type {function(string): void}
	 */
	export let ondelete;
	
	/**
	 * Callback function called when set as default button is clicked
	 * @type {function(string): void}
	 */
	export let onsetdefault;
	
	/**
	 * Callback function called to show toast notifications
	 * @type {function(ToastData): void}
	 */
	export let ontoast;

	/**
	 * Triggers the add connection callback
	 * @function
	 * @returns {void}
	 */
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