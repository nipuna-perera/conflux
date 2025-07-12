// Layout load function for root layout
// Provides global data and configuration for all pages
// Handles universal application setup and data fetching
import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ fetch, url }) => {
	// Global layout data loading:
	// - Load application configuration
	// - Fetch user preferences
	// - Handle global error states
	// - Provide data to all child routes
	
	return {
		// Global data available to all pages
		config: {
			apiUrl: '/api',
			appName: 'Full-Stack App'
		}
	};
};
