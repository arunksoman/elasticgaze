import { writable } from 'svelte/store';

// Create a writable store for REST page data
const createRestStore = () => {
	const initialState = {
		method: 'GET',
		endpoint: '',
		baseEndpoint: '',
		params: [],
		requestBody: '',  // Empty by default - no default request body
		responseData: '',
		isLoading: false
	};

	const { subscribe, set, update } = writable(initialState);

	return {
		subscribe,
		// Method to update method
		setMethod: (method) => update(store => ({ ...store, method })),
		
		// Method to update endpoint
		setEndpoint: (endpoint) => update(store => ({ ...store, endpoint })),
		
		// Method to update base endpoint
		setBaseEndpoint: (baseEndpoint) => update(store => ({ ...store, baseEndpoint })),
		
		// Method to update params
		setParams: (params) => update(store => ({ ...store, params })),
		
		// Method to update request body
		setRequestBody: (requestBody) => update(store => ({ ...store, requestBody })),
		
		// Method to update response data
		setResponseData: (responseData) => update(store => ({ ...store, responseData })),
		
		// Method to update loading state
		setLoading: (isLoading) => update(store => ({ ...store, isLoading })),
		
		// Method to reset only response data (useful for clearing response)
		clearResponse: () => update(store => ({ ...store, responseData: '' })),
		
		// Method to reset everything to initial state
		reset: () => set(initialState)
	};
};

export const restStore = createRestStore();