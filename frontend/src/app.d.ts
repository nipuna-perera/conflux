// Global TypeScript declarations for SvelteKit
// Defines app-wide types, interfaces, and ambient declarations
// Provides type safety for SvelteKit-specific features
declare global {
	namespace App {
		// User interface for authenticated sessions
		interface Locals {
			user?: {
				id: number;
				email: string;
				first_name: string;
				last_name: string;
			};
		}
		
		// Page data interface - extend as needed for typed page data
		// eslint-disable-next-line @typescript-eslint/no-empty-object-type
		interface PageData {}
		
		// Error interface for custom error pages - extend as needed
		// eslint-disable-next-line @typescript-eslint/no-empty-object-type
		interface Error {}
		
		// Platform interface for deployment-specific features - extend as needed
		// eslint-disable-next-line @typescript-eslint/no-empty-object-type
		interface Platform {}
	}
}

export {};
