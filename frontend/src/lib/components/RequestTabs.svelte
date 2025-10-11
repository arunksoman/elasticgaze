<script>
	import { createEventDispatcher } from 'svelte';
	import RequestParamsEditor from './RequestParamsEditor.svelte';
	import RequestBodyEditor from './RequestBodyEditor.svelte';
	import DescriptionEditor from './DescriptionEditor.svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		activeTab = 'body',  // Default to body tab
		params = [],
		requestBody = '',  // No default request body
		description = ''   // No default description
	} = $props();
	
	function setActiveTab(tab) {
		activeTab = tab;
	}
	
	function handleParamsChange(event) {
		params = event.detail;
		dispatch('paramsChange', params);
	}
	
	function handleRequestBodyChange(event) {
		requestBody = event.detail;
		dispatch('requestBodyChange', requestBody);
	}
	
	function handleDescriptionChange(event) {
		description = event.detail;
		dispatch('descriptionChange', description);
	}
</script>

<div class="h-full flex flex-col">
	<!-- Tab Headers -->
	<div class="flex border-b theme-border mb-4 flex-shrink-0">
		<button 
			onclick={() => setActiveTab('body')}
			class="px-4 py-2 font-medium text-sm transition-colors border-b-2 {activeTab === 'body' ? 'border-blue-500 theme-text-primary' : 'border-transparent theme-text-secondary hover:theme-text-primary'}"
		>
			Request Body
		</button>
		<button 
			onclick={() => setActiveTab('params')}
			class="px-4 py-2 font-medium text-sm transition-colors border-b-2 {activeTab === 'params' ? 'border-blue-500 theme-text-primary' : 'border-transparent theme-text-secondary hover:theme-text-primary'}"
		>
			Params
		</button>
		<button 
			onclick={() => setActiveTab('description')}
			class="px-4 py-2 font-medium text-sm transition-colors border-b-2 {activeTab === 'description' ? 'border-blue-500 theme-text-primary' : 'border-transparent theme-text-secondary hover:theme-text-primary'}"
		>
			Description
		</button>
	</div>
	
	<!-- Tab Content -->
	<div class="flex-1 min-h-0">
		{#if activeTab === 'description'}
			<DescriptionEditor 
				{description}
				on:change={handleDescriptionChange}
			/>
		{:else if activeTab === 'params'}
			<RequestParamsEditor 
				{params}
				on:paramsChange={handleParamsChange}
			/>
		{:else if activeTab === 'body'}
			<RequestBodyEditor 
				{requestBody}
				title=""
				on:change={handleRequestBodyChange}
			/>
		{/if}
	</div>
</div>