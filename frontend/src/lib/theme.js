import { writable } from 'svelte/store';

// Initialize theme with 'light' as default
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
				const stored = localStorage.getItem('theme') || 'light';
				set(stored);
				return stored;
			}
			return 'light';
		}
	};
}

export const theme = createTheme();
