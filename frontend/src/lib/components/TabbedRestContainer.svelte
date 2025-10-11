<script>
	import TabBar from './TabBar.svelte';
	import RestRequestForm from './RestRequestForm.svelte';
	import RequestTabs from './RequestTabs.svelte';
	import ResponseViewer from './ResponseViewer.svelte';
	import ResizableSplitter from './ResizableSplitter.svelte';
	import { ExecuteElasticsearchRequest } from '$lib/wailsjs/go/main/App.js';
	import { tabStore } from '$lib/stores/tabStore.js';
	
	// Subscribe to tab store
	let tabState = $state({});
	let activeTab = $state(null);
	
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
	
	// Data update helpers
	function updateTabData(data) {
		if (activeTab) {
			tabStore.updateTabData(activeTab.id, data);
		}
	}
	
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
</script>

<div class="h-screen flex flex-col">
	<!-- Tab bar -->
	<TabBar 
		tabs={tabState.tabs || []}
		activeTabId={tabState.activeTabId}
		on:tabSelect={handleTabSelect}
		on:tabClose={handleTabClose}
		on:tabAdd={handleTabAdd}
	/>
	
	<!-- Active tab content -->
	{#if activeTab}
		<div class="flex-1 flex flex-col min-h-0">
			<div class="p-6 pb-0 flex-shrink-0">
				<h1 class="text-2xl font-medium mb-6 theme-text-primary">REST API</h1>
				
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
</div>