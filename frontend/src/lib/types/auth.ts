// TypeScript type definitions for authentication
// Defines interfaces for login, registration, and auth responses
// Ensures type safety in authentication flows
import type { User } from './user';

export interface LoginRequest {
	email: string;
	password: string;
}

export interface RegisterRequest {
	email: string;
	password: string;
	first_name: string;
	last_name: string;
}

export interface AuthResponse {
	token: string;
	expires_in: number;
	user: User;
}

export interface LoginFormData {
	email: string;
	password: string;
	remember?: boolean;
}

export interface RegisterFormData {
	email: string;
	password: string;
	confirmPassword: string;
	firstName: string;
	lastName: string;
	terms: boolean;
}
