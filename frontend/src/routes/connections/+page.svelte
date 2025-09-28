<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Toast from '$lib/Toast.svelte';
	import ConnectionsList from '$lib/components/ConnectionsList.svelte';
	import ConnectionForm from '$lib/components/ConnectionForm.svelte';
	import { connections, connectionService } from '$lib/stores/connectionStore.js';
	
	// Component state
	let showForm = false;
	let editingConnection = null;
	let formData = connectionService.getDefaultFormData();
	
	// Toast state
	let toastShow = false;
	let toastMessage = '';
	let toastType = 'success';
	let toastDuration = 1500;
	let toastAnimation = 'fade';
	let toastErrorCode = '';
	let toastErrorDetails = '';
	
	// Testing state
	let testingConnectionId = null;
	
	// Subscribe to connections store
	let connectionsArray = [];
	const unsubscribe = connections.subscribe(value => {
		connectionsArray = value;
	});
	
	onMount(() => {
		connectionService.load();
	});

	// Cleanup subscription on destroy
	import { onDestroy } from 'svelte';
	onDestroy(() => {
		unsubscribe();
	});

	// Event handlers
	function openForm(connection = null) {
		editingConnection = connection;
		if (connection) {
			formData = { ...connection };
		} else {
			formData = connectionService.getDefaultFormData();
		}
		showForm = true;
	}

	function closeForm() {
		showForm = false;
		editingConnection = null;
	}

	async function saveConnection(connectionData) {
		try {
			if (editingConnection) {
				await connectionService.update(connectionData);
			} else {
				await connectionService.add(connectionData);
			}
			closeForm();
		} catch (error) {
			console.error('Error saving connection:', error);
		}
	}

	async function deleteConnection(connectionId) {
		try {
			await connectionService.delete(connectionId);
		} catch (error) {
			console.error('Error deleting connection:', error);
		}
	}

	async function setAsDefault(connectionId) {
		try {
			await connectionService.setAsDefault(connectionId);
		} catch (error) {
			console.error('Error setting default connection:', error);
		}
	}

	function handleTestStart(testingId) {
		testingConnectionId = testingId;
	}

	function handleTestEnd() {
		testingConnectionId = null;
	}

	function goBack() {
		goto('/');
	}
	
	// Toast utility functions
	function showToast(message, type = 'success', duration = 1500, animation = 'fade', errorCode = '', errorDetails = '') {
		toastMessage = message;
		toastType = type;
		toastDuration = duration;
		toastAnimation = animation;
		toastErrorCode = errorCode;
		toastErrorDetails = errorDetails;
		toastShow = true;
	}

	function hideToast() {
		toastShow = false;
		// Clear error details when hiding
		toastErrorCode = '';
		toastErrorDetails = '';
	}

	function handleToast(toastData) {
		const { message, type, duration, animation, errorCode, errorDetails } = toastData;
		showToast(message, type, duration, animation, errorCode, errorDetails);
	}

	function handleAddConnection() {
		openForm();
	}
</script>

<div class="p-6">
	<div class="flex items-center gap-4 mb-6">
		<button 
			onclick={goBack}
			class="p-2 rounded-md theme-bg-secondary theme-text-primary hover:theme-bg-tertiary transition-colors"
			title="Back to Home"
			aria-label="Back to Home"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
			</svg>
		</button>
		<h1 class="text-2xl font-medium theme-text-primary">Elasticsearch Connections</h1>
		{#if connectionsArray.length > 0}
			<button 
				onclick={handleAddConnection}
				class="ml-auto bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors flex items-center gap-2"
			>
				<img src="/icons/create.svg" alt="" class="w-4 h-4" style="filter: brightness(0) invert(1);" />
				Add Connection
			</button>
		{/if}
	</div>
	
	<!-- Connections List Component -->
	<ConnectionsList
		connections={connectionsArray}
		{testingConnectionId}
		onadd={handleAddConnection}
		onteststart={handleTestStart}
		ontestend={handleTestEnd}
		onedit={(connection) => openForm(connection)}
		ondelete={deleteConnection}
		onsetdefault={setAsDefault}
		ontoast={handleToast}
	/>
</div>

<!-- Connection Form Modal Component -->
<ConnectionForm
	bind:show={showForm}
	{editingConnection}
	bind:formData
	onclose={closeForm}
	onsave={saveConnection}
	ontoast={handleToast}
/>

<!-- Toast Component -->
<Toast 
	bind:show={toastShow}
	message={toastMessage}
	type={toastType}
	duration={toastDuration}
	animation={toastAnimation}
	errorCode={toastErrorCode}
	errorDetails={toastErrorDetails}
	on:hide={hideToast}
/>