<script>
	import { createEventDispatcher } from 'svelte';
	
	const dispatch = createEventDispatcher();
	
	// Use $props() for Svelte 5 runes mode
	let {
		method = 'GET'
	} = $props();
	
	// Postman-style method colors
	function getMethodColor(methodName) {
		const colors = {
			'GET': 'text-green-500',
			'POST': 'text-orange-500',
			'PUT': 'text-blue-500',
			'DELETE': 'text-red-500',
			'PATCH': 'text-purple-500',
			'HEAD': 'text-gray-500',
			'OPTIONS': 'text-indigo-500'
		};
		return colors[methodName] || 'text-gray-400';
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
	<option value="GET" class="text-green-500 font-semibold">GET</option>
	<option value="POST" class="text-orange-500 font-semibold">POST</option>
	<option value="PUT" class="text-blue-500 font-semibold">PUT</option>
	<option value="DELETE" class="text-red-500 font-semibold">DELETE</option>
	<option value="PATCH" class="text-purple-500 font-semibold">PATCH</option>
	<option value="HEAD" class="text-gray-500 font-semibold">HEAD</option>
	<option value="OPTIONS" class="text-indigo-500 font-semibold">OPTIONS</option>
</select>