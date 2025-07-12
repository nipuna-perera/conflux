// TypeScript type definitions for user entities
// Defines user data structures and profile interfaces
// Ensures type safety across user-related components
export interface User {
	id: number;
	email: string;
	first_name: string;
	last_name: string;
	created_at: string;
	updated_at: string;
}

export interface UpdateProfileData {
	first_name?: string;
	last_name?: string;
	email?: string;
}

export interface UserProfile extends User {
	// Extended profile data
	avatar_url?: string;
	bio?: string;
	location?: string;
}
