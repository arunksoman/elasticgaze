<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		method = 'GET'
	} = $props();
	
	// Postman-style method colors - match CollectionsSidebar colors
	function getMethodColor(methodName) {
		const colors = {
			'GET': 'text-blue-600 dark:text-blue-400',
			'POST': 'text-green-600 dark:text-green-400',
			'PUT': 'text-orange-600 dark:text-orange-400',
			'DELETE': 'text-red-600 dark:text-red-400',
			'PATCH': 'text-purple-600 dark:text-purple-400',
			'HEAD': 'text-gray-600 dark:text-gray-400',
			'OPTIONS': 'text-yellow-600 dark:text-yellow-400'
		};
		return colors[methodName] || 'text-gray-600 dark:text-gray-400';
	}
	
	function handleChange(event) {
		method = event.target.value;
		dispatch('change', method);
	}
</script>

<select 
	value={method}
	onchange={handleChange}
	class={`border theme-border p-2 theme-bg-tertiary rounded font-semibold ${getMethodColor(method)} min-w-[100px]`}
>
	<option value="GET" class="text-blue-600 dark:text-blue-400 font-semibold">GET</option>
	<option value="POST" class="text-green-600 dark:text-green-400 font-semibold">POST</option>
	<option value="PUT" class="text-orange-600 dark:text-orange-400 font-semibold">PUT</option>
	<option value="DELETE" class="text-red-600 dark:text-red-400 font-semibold">DELETE</option>
	<option value="PATCH" class="text-purple-600 dark:text-purple-400 font-semibold">PATCH</option>
	<option value="HEAD" class="text-gray-600 dark:text-gray-400 font-semibold">HEAD</option>
	<option value="OPTIONS" class="text-yellow-600 dark:text-yellow-400 font-semibold">OPTIONS</option>
</select>