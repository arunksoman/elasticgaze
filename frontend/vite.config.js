import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [tailwindcss(), sveltekit()],
	optimizeDeps: {
		include: ['monaco-editor']
	},
	build: {
		rollupOptions: {
			output: {
				manualChunks: (id) => {
					// Split vendor dependencies into separate chunks
					if (id.includes('node_modules')) {
						if (id.includes('monaco-editor')) {
							return 'monaco-editor';
						}
						if (id.includes('svelte') || id.includes('@sveltejs')) {
							return 'svelte-vendor';
						}
						return 'vendor';
					}
				}
			}
		},
		// Increase chunk size warning limit for Monaco Editor
		chunkSizeWarningLimit: 1000
	},
	define: {
		// Global Monaco Environment
		global: 'globalThis',
	},
	worker: {
		format: 'es'
	}
});
