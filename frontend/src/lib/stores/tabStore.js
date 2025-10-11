import { writable } from 'svelte/store';

// Tab data structure
const createTab = (id) => ({
	id,
	title: `Tab ${id}`,
	isModified: false,
	data: {
		method: 'GET',
		endpoint: '',
		baseEndpoint: '',
		params: [],
		requestBody: '',
		description: '',
		responseData: '',
		isLoading: false
	}
});

// Create the tab manager store
const createTabStore = () => {
	let tabCounter = 1;
	
	const initialState = {
		tabs: [createTab(tabCounter)], // Start with one tab
		activeTabId: tabCounter,
		nextId: tabCounter + 1
	};

	const { subscribe, set, update } = writable(initialState);

	return {
		subscribe,
		
		// Add a new tab
		addTab: () => update(store => {
			const newTab = createTab(store.nextId);
			return {
				...store,
				tabs: [...store.tabs, newTab],
				activeTabId: newTab.id,
				nextId: store.nextId + 1
			};
		}),
		
		// Close a tab
		closeTab: (tabId) => update(store => {
			// Don't close if it's the only tab
			if (store.tabs.length <= 1) return store;
			
			const newTabs = store.tabs.filter(tab => tab.id !== tabId);
			let newActiveTabId = store.activeTabId;
			
			// If we're closing the active tab, switch to the next available tab
			if (store.activeTabId === tabId) {
				const closedIndex = store.tabs.findIndex(tab => tab.id === tabId);
				// Try to switch to the tab to the right, or left if that's the last tab
				newActiveTabId = newTabs[Math.min(closedIndex, newTabs.length - 1)]?.id || newTabs[0]?.id;
			}
			
			return {
				...store,
				tabs: newTabs,
				activeTabId: newActiveTabId
			};
		}),
		
		// Switch to a different tab
		switchTab: (tabId) => update(store => ({
			...store,
			activeTabId: tabId
		})),
		
		// Update tab data
		updateTabData: (tabId, newData) => update(store => ({
			...store,
			tabs: store.tabs.map(tab => 
				tab.id === tabId 
					? { 
						...tab, 
						data: { ...tab.data, ...newData },
						isModified: true 
					}
					: tab
			)
		})),
		
		// Update tab title
		updateTabTitle: (tabId, title) => update(store => ({
			...store,
			tabs: store.tabs.map(tab => 
				tab.id === tabId ? { ...tab, title } : tab
			)
		})),
		
		// Mark tab as saved (remove modified flag)
		markTabSaved: (tabId) => update(store => ({
			...store,
			tabs: store.tabs.map(tab => 
				tab.id === tabId ? { ...tab, isModified: false } : tab
			)
		})),
		
		// Get active tab data
		getActiveTab: (store) => {
			return store.tabs.find(tab => tab.id === store.activeTabId);
		},
		
		// Get tab by ID
		getTab: (store, tabId) => {
			return store.tabs.find(tab => tab.id === tabId);
		}
	};
};

export const tabStore = createTabStore();