// SvelteKit configuration file
// Configures build adapter, preprocessing, and framework settings
// Determines how the application is built and deployed
import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Enables TypeScript and PostCSS preprocessing
	preprocess: vitePreprocess(),

	kit: {
		// Node.js adapter for containerized deployment
		adapter: adapter(),
		
		// API proxy configuration for backend communication
		alias: {
			$lib: 'src/lib'
		}
	}
};

export default config;
