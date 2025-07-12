// Server-side hooks for SvelteKit
// Handles authentication, session management, and request preprocessing
// Runs on the server for every request before page rendering
import type { Handle } from '@sveltejs/kit';

// Handle function runs on every server-side request
// Performs authentication checks and populates locals with user data
export const handle: Handle = async ({ event, resolve }) => {
	// Server-side authentication logic:
	// - Extract JWT token from cookies
	// - Validate token with backend API
	// - Populate event.locals.user with authenticated user data
	// - Handle authentication errors gracefully
	
	const response = await resolve(event);
	return response;
};
