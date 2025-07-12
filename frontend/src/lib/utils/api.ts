// API client utilities for backend communication
// Provides centralized HTTP client with authentication and error handling
// Abstracts API calls and manages request/response transformation
import type { LoginRequest, RegisterRequest, AuthResponse } from '$lib/types/auth';
import type { User, UpdateProfileData } from '$lib/types/user';

export class ApiClient {
	private baseUrl: string;
	
	constructor(baseUrl: string = '/api') {
		this.baseUrl = baseUrl;
	}
	
	// Generic HTTP request method with authentication
	// Automatically includes JWT tokens and handles common errors
	async request(endpoint: string, options: RequestInit = {}): Promise<any> {
		// API request implementation:
		// - Add authentication headers
		// - Handle request/response JSON transformation
		// - Manage error responses and token refresh
		// - Provide consistent error handling
		
		const url = `${this.baseUrl}${endpoint}`;
		const token = localStorage.getItem('auth_token');
		
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
			...options.headers,
		};
		
		if (token) {
			headers.Authorization = `Bearer ${token}`;
		}
		
		const response = await fetch(url, {
			...options,
			headers,
		});
		
		if (!response.ok) {
			const error = await response.json().catch(() => ({ message: 'An error occurred' }));
			throw new Error(error.message || `HTTP ${response.status}`);
		}
		
		return response.json();
	}
	
	// Authentication API methods
	async login(email: string, password: string): Promise<AuthResponse> {
		const loginRequest: LoginRequest = { email, password };
		return this.request('/auth/login', {
			method: 'POST',
			body: JSON.stringify(loginRequest),
		});
	}
	
	async register(userData: RegisterRequest): Promise<AuthResponse> {
		return this.request('/auth/register', {
			method: 'POST',
			body: JSON.stringify(userData),
		});
	}
	
	async logout(): Promise<void> {
		return this.request('/auth/logout', {
			method: 'POST',
		});
	}
	
	// User API methods
	async getProfile(): Promise<User> {
		return this.request('/users/profile');
	}
	
	async updateProfile(userData: UpdateProfileData): Promise<User> {
		return this.request('/users/profile', {
			method: 'PUT',
			body: JSON.stringify(userData),
		});
	}
	
	async getUser(id: number): Promise<User> {
		return this.request(`/users/${id}`);
	}
}

// API client instance for use throughout the application
export const apiClient = new ApiClient();
