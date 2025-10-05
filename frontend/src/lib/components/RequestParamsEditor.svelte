<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	export let params = [];
	
	// Ensure there's always at least one empty row for new parameters
	$: if (params.length === 0 || params[params.length - 1].key || params[params.length - 1].value) {
		// Only add if the last row has content or if there are no rows
		if (params.length === 0 || (params[params.length - 1].key && params[params.length - 1].value)) {
			params = [...params, { key: '', value: '', enabled: true }];
		}
	}
	
	function addParam() {
		params = [...params, { key: '', value: '', enabled: true }];
		dispatch('paramsChange', params);
	}
	
	function removeParam(index) {
		params = params.filter((_, i) => i !== index);
		dispatch('paramsChange', params);
	}
	
	function updateParam(index, field, value) {
		params = params.map((param, i) => 
			i === index ? { ...param, [field]: value } : param
		);
		
		// Auto-add new parameter when user starts typing in the last row's value field
		if (field === 'value' && value && index === params.length - 1) {
			params = [...params, { key: '', value: '', enabled: true }];
		}
		
		// Clean up empty rows (except the last one for new entries)
		if (field === 'key' || field === 'value') {
			const nonEmptyParams = params.filter((param, i) => {
				// Keep the last row even if empty (for new entries)
				if (i === params.length - 1) return true;
				// Keep rows that have either key or value
				return param.key || param.value;
			});
			
			// Ensure we always have at least one empty row at the end
			if (nonEmptyParams.length === 0 || 
				(nonEmptyParams[nonEmptyParams.length - 1].key || nonEmptyParams[nonEmptyParams.length - 1].value)) {
				nonEmptyParams.push({ key: '', value: '', enabled: true });
			}
			
			params = nonEmptyParams;
		}
		
		dispatch('paramsChange', params);
	}
	
	function toggleParam(index) {
		params = params.map((param, i) => 
			i === index ? { ...param, enabled: !param.enabled } : param
		);
		dispatch('paramsChange', params);
	}
</script>

<div class="space-y-3">
	<div class="flex justify-between items-center">
		<span class="block theme-text-primary font-medium text-lg">Parameters</span>
		<button 
			on:click={addParam}
			class="bg-gray-500 hover:bg-gray-600 text-white px-3 py-1 rounded text-sm font-medium transition-colors"
			title="Add additional parameter row"
		>
			+ Add Row
		</button>
	</div>
	
	{#if params.length === 0}
		<div class="text-center py-8 theme-text-secondary">
			<p>No parameters added yet.</p>
			<p class="text-sm mt-1">Start typing in the fields below to add parameters.</p>
		</div>
	{:else}
		<div class="space-y-2">
			<!-- Header -->
			<div class="grid grid-cols-12 gap-2 text-sm font-medium theme-text-secondary px-2">
				<div class="col-span-1"></div>
				<div class="col-span-5">Key</div>
				<div class="col-span-5">Value</div>
				<div class="col-span-1"></div>
			</div>
			
			<!-- Parameter rows -->
			{#each params as param, index}
				<div class="grid grid-cols-12 gap-2 items-center">
					<!-- Enabled checkbox -->
					<div class="col-span-1 flex justify-center">
						<input 
							type="checkbox" 
							checked={param.enabled}
							on:change={() => toggleParam(index)}
							class="w-4 h-4 text-blue-600 rounded"
						/>
					</div>
					
					<!-- Key input -->
					<div class="col-span-5">
						<input 
							type="text" 
							value={param.key}
							on:input={(e) => updateParam(index, 'key', e.target.value)}
							placeholder="Parameter key"
							class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded text-sm"
							class:opacity-50={!param.enabled}
						/>
					</div>
					
					<!-- Value input -->
					<div class="col-span-5">
						<input 
							type="text" 
							value={param.value}
							on:input={(e) => updateParam(index, 'value', e.target.value)}
							placeholder="Parameter value"
							class="w-full border theme-border p-2 theme-bg-tertiary theme-text-primary rounded text-sm"
							class:opacity-50={!param.enabled}
						/>
					</div>
					
					<!-- Remove button -->
					<div class="col-span-1 flex justify-center">
						<button 
							on:click={() => removeParam(index)}
							class="text-red-500 hover:text-red-700 p-1"
							title="Remove parameter"
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>