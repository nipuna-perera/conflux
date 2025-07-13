// Client-side hooks for SvelteKit
// Handles client-side error reporting and navigation events
// Runs in the browser for enhanced user experience
import type { HandleClientError } from '@sveltejs/kit';

// HandleClientError manages client-side error reporting
// Provides graceful error handling and user feedback
export const handleError: HandleClientError = ({ error, event: _event }) => {
	// Client-side error handling:
	// - Log errors to monitoring service
	// - Show user-friendly error messages
	// - Handle authentication errors
	// - Prevent sensitive data exposure
	
	console.error('Client error:', error);
	
	return {
		message: 'An unexpected error occurred. Please try again.'
	};
};
