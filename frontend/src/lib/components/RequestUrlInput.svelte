<script>
	export let endpoint = '';
	export let placeholder = '/_cluster/health';
</script>

<input 
	type="text" 
	bind:value={endpoint}
	{placeholder}
	class="flex-1 border theme-border p-3 theme-bg-tertiary theme-text-primary rounded"
/>