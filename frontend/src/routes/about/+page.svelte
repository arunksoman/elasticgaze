<script>
	// About page component
	import { onMount, onDestroy } from 'svelte';
	
	onMount(() => {
		// Add padding to main content when about page is mounted
		const main = document.querySelector('main');
		if (main) {
			main.style.paddingBottom = '60px';
		}
	});
	
	onDestroy(() => {
		// Remove padding when leaving about page
		const main = document.querySelector('main');
		if (main) {
			main.style.paddingBottom = '1.5rem'; // Reset to original padding
		}
	});
</script>

<div class="p-6 relative min-h-full">
	<h1 class="text-2xl font-medium mb-6 theme-text-primary">About ElasticGaze</h1>
	
	<!-- Project Inspiration Section -->
	<div class="mb-8">
		<h2 class="text-xl font-medium mb-4 theme-text-primary">Project Inspiration</h2>
		<div class="border theme-border p-6 rounded-lg">
			<p class="mb-4 theme-text-secondary leading-relaxed">
				ElasticGaze is inspired by ElasticVue, the popular Elasticsearch management interface. 
				While ElasticVue provided a solid foundation for Elasticsearch administration, 
				ElasticGaze aims to deliver significant improvements and enhancements over its predecessor.
			</p>
			<p class="mb-4 theme-text-secondary leading-relaxed">
				Our vision is to create a more intuitive, feature-rich, and performant tool that makes 
				Elasticsearch management effortless. ElasticGaze incorporates modern web technologies, 
				improved user experience design, and advanced features that weren't available in the original.
			</p>
			<p class="theme-text-secondary leading-relaxed">
				Whether you're a developer, DevOps engineer, or data analyst working with Elasticsearch, 
				ElasticGaze provides the tools you need to monitor, manage, and optimize your clusters 
				with unprecedented ease and efficiency.
			</p>
		</div>
	</div>

	<!-- License Section -->
	<div class="mb-8">
		<h2 class="text-xl font-medium mb-4 theme-text-primary">License</h2>
		<div class="border theme-border p-6 rounded-lg">
			<h3 class="text-lg font-medium mb-3 theme-text-primary">MIT License</h3>
			<div class="theme-text-secondary leading-relaxed text-sm font-mono theme-bg-tertiary border theme-border p-4 rounded">
				<p class="mb-2">Copyright (c) 2025 ElasticGaze</p>
				<p class="mb-4">
					Permission is hereby granted, free of charge, to any person obtaining a copy
					of this software and associated documentation files (the "Software"), to deal
					in the Software without restriction, including without limitation the rights
					to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
					copies of the Software, and to permit persons to whom the Software is
					furnished to do so, subject to the following conditions:
				</p>
				<p class="mb-4">
					The above copyright notice and this permission notice shall be included in all
					copies or substantial portions of the Software.
				</p>
				<p>
					THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
					IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
					FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
					AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
					LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
					OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
					SOFTWARE.
				</p>
			</div>
		</div>
	</div>

	<!-- Sticky Footer -->
	<div class="fixed bottom-0 left-0 right-0 py-4 text-center bg-transparent">
		<p class="theme-text-primary text-sm font-medium">
			Made with ❤️ in <img src="/icons/india.svg" alt="India Flag" class="inline w-4 h-4 ml-1 align-middle" />
		</p>
	</div>
</div>