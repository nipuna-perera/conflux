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
		
		// Page data interface
		interface PageData {}
		
		// Error interface for custom error pages
		interface Error {}
		
		// Platform interface for deployment-specific features
		interface Platform {}
	}
}

export {};
