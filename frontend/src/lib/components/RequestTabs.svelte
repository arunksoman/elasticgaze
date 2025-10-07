<script>
	import { createEventDispatcher } from 'svelte';
	import RequestParamsEditor from './RequestParamsEditor.svelte';
	import RequestBodyEditor from './RequestBodyEditor.svelte';
	
	const dispatch = createEventDispatcher();
	
	export let activeTab = 'body';
	export let params = [];
	export let requestBody = '{\n  "query": {\n    "match_all": {}\n  }\n}';
	
	function setActiveTab(tab) {
		activeTab = tab;
	}
	
	function handleParamsChange(event) {
		params = event.detail;
		dispatch('paramsChange', params);
	}
</script>

<div class="h-full flex flex-col">
	<!-- Tab Headers -->
	<div class="flex border-b theme-border mb-4 flex-shrink-0">
		<button 
			on:click={() => setActiveTab('body')}
			class="px-4 py-2 font-medium text-sm transition-colors border-b-2 {activeTab === 'body' ? 'border-blue-500 theme-text-primary' : 'border-transparent theme-text-secondary hover:theme-text-primary'}"
		>
			Request Body
		</button>
		<button 
			on:click={() => setActiveTab('params')}
			class="px-4 py-2 font-medium text-sm transition-colors border-b-2 {activeTab === 'params' ? 'border-blue-500 theme-text-primary' : 'border-transparent theme-text-secondary hover:theme-text-primary'}"
		>
			Params
		</button>
	</div>
	
	<!-- Tab Content -->
	<div class="flex-1 min-h-0">
		{#if activeTab === 'params'}
			<RequestParamsEditor 
				bind:params 
				on:paramsChange={handleParamsChange}
			/>
		{:else if activeTab === 'body'}
			<RequestBodyEditor 
				bind:requestBody 
				title=""
			/>
		{/if}
	</div>
</div>