import { writable } from 'svelte/store';

// Initialize theme from localStorage if available, otherwise use 'light'
function createTheme() {
	const { subscribe, set, update } = writable('light');

	return {
		subscribe,
		set: (value) => {
			set(value);
			if (typeof window !== 'undefined') {
				localStorage.setItem('theme', value);
			}
		},
		toggle: () => update(n => n === 'light' ? 'dark' : 'light'),
		init: () => {
			if (typeof window !== 'undefined') {
				const stored = localStorage.getItem('theme');
				if (stored) {
					set(stored);
				}
			}
		}
	};
}

export const theme = createTheme();
