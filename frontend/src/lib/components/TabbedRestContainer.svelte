<script>
	import TabBar from './TabBar.svelte';
	import RestRequestForm from './RestRequestForm.svelte';
	import RequestTabs from './RequestTabs.svelte';
	import ResponseViewer from './ResponseViewer.svelte';
	import ResizableSplitter from './ResizableSplitter.svelte';
	import CollectionsSidebar from './CollectionsSidebar.svelte';
	import Toast from '$lib/Toast.svelte';
	import { ExecuteElasticsearchRequest, UpdateRestRequest } from '$lib/wailsjs/go/main/App.js';
	import { tabStore } from '$lib/stores/tabStore.js';
	import { collectionsOpen } from '$lib/stores/collectionsStore.js';
	import { onMount } from 'svelte';
	
	// Subscribe to tab store
	let tabState = $state({});
	let activeTab = $state(null);
	
	// Toast state
	let toastShow = $state(false);
	let toastMessage = $state('');
	let toastType = $state('success');
	let toastDuration = $state(1500);
	let toastAnimation = $state('fade');

	// Toast function
	function showToast(message, type = 'success', duration = 1500, animation = 'fade') {
		toastMessage = message;
		toastType = type;
		toastDuration = duration;
		toastAnimation = animation;
		toastShow = true;
	}
	
	$effect(() => {
		const unsubscribe = tabStore.subscribe(store => {
			tabState = store;
			activeTab = store.tabs.find(tab => tab.id === store.activeTabId);
		});
		
		return unsubscribe;
	});
	
	// Function to build URL with parameters
	function buildUrlWithParams(base, parameters) {
		if (!base) return '';
		
		const enabledParams = parameters.filter(p => p.enabled && p.key && p.value);
		if (enabledParams.length === 0) return base;
		
		const queryString = enabledParams
			.map(p => `${encodeURIComponent(p.key)}=${encodeURIComponent(p.value)}`)
			.join('&');
		
		return base + (base.includes('?') ? '&' : '?') + queryString;
	}
	
	// Function to parse URL and extract base endpoint and parameters
	function parseUrl(url) {
		if (!url) return { base: '', params: [] };
		
		const [base, queryString] = url.split('?');
		if (!queryString) return { base: url, params: [] };
		
		const params = queryString.split('&').map(param => {
			const [key, value] = param.split('=');
			return {
				key: decodeURIComponent(key || ''),
				value: decodeURIComponent(value || ''),
				enabled: true
			};
		}).filter(p => p.key);
		
		return { base, params };
	}
	
	// Tab management handlers
	function handleTabSelect(event) {
		tabStore.switchTab(event.detail);
	}
	
	function handleTabClose(event) {
		tabStore.closeTab(event.detail);
	}
	
	function handleTabAdd() {
		tabStore.addTab();
	}
	
	function handleTabNameChanged(event) {
		// Refresh collections sidebar when a tab name is changed
		// Add a small delay to ensure the backend update is complete
		setTimeout(() => {
			if (window.refreshCollections) {
				window.refreshCollections();
			}
		}, 100);
	}
	
	// Data update helpers
	function updateTabData(data) {
		if (activeTab) {
			tabStore.updateTabData(activeTab.id, data);
		}
	}

	// Save current tab data to database
	async function saveTabData() {
		if (!activeTab) {
			console.log('Cannot save: No active tab');
			return;
		}

		if (!activeTab.data.requestId) {
			showToast('Cannot save: This is not a saved request. Create a request from collections to save changes.', 'error', 3000);
			return;
		}

		try {
			const updateData = {
				name: activeTab.title,
				method: activeTab.data.method,
				url: activeTab.data.endpoint,
				body: activeTab.data.requestBody || null,
				description: activeTab.data.description || null
			};

			await UpdateRestRequest(activeTab.data.requestId, updateData);
			
			// Mark tab as saved (remove modified indicator)
			tabStore.markTabSaved(activeTab.id);
			
			// Refresh collections sidebar to reflect changes
			if (window.refreshCollections) {
				window.refreshCollections();
			}
			
			showToast('Request saved successfully', 'success');
			console.log('Request saved successfully');
		} catch (error) {
			console.error('Failed to save request:', error);
			showToast('Failed to save request', 'error');
		}
	}

	// Handle keyboard shortcuts
	function handleKeydown(event) {
		if (event.ctrlKey && event.key === 's') {
			event.preventDefault();
			saveTabData();
		}
	}

	// Add keyboard event listener
	onMount(() => {
		document.addEventListener('keydown', handleKeydown);
		return () => {
			document.removeEventListener('keydown', handleKeydown);
		};
	});
	
	// REST API handlers
	function handleMethodChange(event) {
		updateTabData({ method: event.detail });
	}
	
	function handleEndpointChange(newEndpoint) {
		const parsed = parseUrl(newEndpoint);
		updateTabData({ 
			endpoint: newEndpoint,
			baseEndpoint: parsed.base,
			params: parsed.params
		});
	}
	
	function handleParamsChange(event) {
		const newParams = event.detail;
		const newUrl = buildUrlWithParams(activeTab?.data.baseEndpoint || '', newParams);
		updateTabData({
			params: newParams,
			endpoint: newUrl
		});
	}
	
	function handleRequestBodyChange(event) {
		updateTabData({ requestBody: event.detail });
	}
	
	function handleDescriptionChange(event) {
		updateTabData({ description: event.detail });
	}
	
	function handleRequest(event) {
		if (!activeTab) return;
		
		const { method: requestMethod, endpoint: requestEndpoint } = event.detail;
		
		// Set loading state
		updateTabData({ isLoading: true });
		
		console.log('REST Request:', requestMethod, requestEndpoint, activeTab.data.requestBody);
		
		// Prepare the request object for the Go backend
		const esRequest = {
			method: requestMethod,
			endpoint: requestEndpoint,
			body: activeTab.data.requestBody && activeTab.data.requestBody.trim() ? activeTab.data.requestBody : null
		};
		
		// Call the Go function
		ExecuteElasticsearchRequest(esRequest)
			.then(response => {
				console.log('ES Response:', response);
				
				let finalResponseData;
				if (response.success) {
					// Parse and pretty-print the response
					try {
						const parsedResponse = JSON.parse(response.response);
						finalResponseData = JSON.stringify(parsedResponse, null, 2);
					} catch (e) {
						// If it's not JSON, just display as-is
						finalResponseData = response.response;
					}
				} else {
					// Show error response
					const errorResponse = {
						error: true,
						status_code: response.status_code,
						error_code: response.error_code,
						error_details: response.error_details
					};
					finalResponseData = JSON.stringify(errorResponse, null, 2);
				}
				
				// Update response data
				updateTabData({ responseData: finalResponseData });
			})
			.catch(error => {
				console.error('Request failed:', error);
				const errorResponse = {
					error: true,
					message: 'Request failed',
					details: error.toString()
				};
				const finalResponseData = JSON.stringify(errorResponse, null, 2);
				updateTabData({ responseData: finalResponseData });
			})
			.finally(() => {
				updateTabData({ isLoading: false });
			});
	}
	
	// Handle splitter resize
	function handleSplitterResize(percentage) {
		setTimeout(() => {
			window.dispatchEvent(new Event('resize'));
		}, 50);
	}

	// Collections sidebar handlers
	function handleRequestSelect(request) {
		console.log('Request selected:', request);
	}

	function handleRequestLoad(requestData) {
		// Load request data into a new tab
		const newTabData = {
			method: requestData.method || 'GET',
			endpoint: requestData.url || '',
			baseEndpoint: requestData.url || '',
			params: [],
			requestBody: requestData.body || '',
			description: requestData.description || '',
			responseData: null,
			isLoading: false
		};
		
		// Add a new tab with the loaded request data
		tabStore.addTab(requestData.name || 'Loaded Request', newTabData);
	}
</script>

<div class="h-screen flex flex-col">
	<!-- Tab bar -->
	<TabBar 
		tabs={tabState.tabs || []}
		activeTabId={tabState.activeTabId}
		on:tabSelect={handleTabSelect}
		on:tabClose={handleTabClose}
		on:tabAdd={handleTabAdd}
		on:tabNameChanged={handleTabNameChanged}
	/>
	
	<!-- Active tab content -->
	{#if activeTab}
		<div class="flex-1 flex flex-col min-h-0">
			<div class="p-6 pb-0 flex-shrink-0">
				<div class="flex items-center justify-between mb-6">
					<h1 class="text-2xl font-medium theme-text-primary">REST API</h1>
					{#if activeTab?.isModified}
						<div class="flex items-center gap-2 text-sm theme-text-secondary">
							<span class="w-2 h-2 bg-orange-500 rounded-full"></span>
							{#if activeTab.data.requestId}
								<span>Unsaved changes</span>
								<kbd class="px-2 py-1 bg-gray-200 dark:bg-gray-700 rounded text-xs">Ctrl+S</kbd>
								<span>to save</span>
							{:else}
								<span>Unsaved changes (Create request in collections to save)</span>
							{/if}
						</div>
					{/if}
				</div>
				
				<RestRequestForm 
					method={activeTab.data.method}
					endpoint={activeTab.data.endpoint}
					isLoading={activeTab.data.isLoading}
					on:methodChange={handleMethodChange}
					on:endpointChange={(e) => handleEndpointChange(e.detail)}
					on:send={handleRequest}
				/>
			</div>
			
			<div class="flex-1 p-6 pt-0 min-h-0">
				<ResizableSplitter 
					defaultSplit={45} 
					minSize={25} 
					maxSize={75}
					className="h-full"
					onResize={handleSplitterResize}
				>
					{#snippet panel1()}
						<!-- Request Section -->
						<div class="h-full flex flex-col">
							<RequestTabs 
								params={activeTab.data.params}
								requestBody={activeTab.data.requestBody}
								description={activeTab.data.description}
								on:paramsChange={handleParamsChange}
								on:requestBodyChange={handleRequestBodyChange}
								on:descriptionChange={handleDescriptionChange}
							/>
						</div>
					{/snippet}
					
					{#snippet panel2()}
						<!-- Response Section -->
						<div class="h-full flex flex-col">
							<ResponseViewer responseData={activeTab.data.responseData} />
						</div>
					{/snippet}
				</ResizableSplitter>
			</div>
		</div>
	{:else}
		<!-- Fallback if no active tab -->
		<div class="flex-1 flex items-center justify-center theme-text-secondary">
			<p>No active tab</p>
		</div>
	{/if}

	<!-- Collections Sidebar - Only on REST page -->
	<CollectionsSidebar 
		bind:isOpen={$collectionsOpen}
	/>
	
	<!-- Toast Component -->
	<Toast 
		bind:show={toastShow} 
		message={toastMessage} 
		type={toastType} 
		duration={toastDuration} 
		animation={toastAnimation} 
	/>
</div>