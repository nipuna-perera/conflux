// API client utilities for backend communication
// Provides centralized HTTP client with authentication and error handling
// Abstracts API calls and manages request/response transformation
import type { LoginRequest, RegisterRequest, AuthResponse } from '$lib/types/auth';
import type { User, UpdateProfileData } from '$lib/types/user';

export class ApiClient {
	private baseUrl: string;
	
	constructor(baseUrl?: string) {
		// Use environment variable for backend URL, fallback to proxy in development
		this.baseUrl = baseUrl || this.getApiBaseUrl();
	}
	
	private getApiBaseUrl(): string {
		if (import.meta.env.DEV) {
			console.log('API client: Getting base URL, window exists:', typeof window !== 'undefined');
		}
		
		// For browser context, detect environment and use appropriate URL
		if (typeof window !== 'undefined') {
			const hostname = window.location.hostname;
			if (import.meta.env.DEV) {
				console.log('API client: Browser hostname:', hostname);
			}
			
			// If hostname is localhost, we're in local development
			if (hostname === 'localhost' || hostname === '127.0.0.1') {
				if (import.meta.env.DEV) {
					console.log('API client: Local development mode - using localhost:8080');
				}
				return 'http://localhost:8080/api';
			} else {
				// If not localhost, we're likely in Docker container accessing via host
				if (import.meta.env.DEV) {
					console.log('API client: Docker container mode - using localhost:8080');
				}
				return 'http://localhost:8080/api';
			}
		}
		
		// For server-side rendering in Docker, use container service name
		if (import.meta.env.DEV) {
			console.log('API client: Server-side rendering in Docker - using backend:8080');
		}
		return 'http://backend:8080/api';
	}
	
	// Generic HTTP request method with authentication
	// Automatically includes JWT tokens and handles common errors
	async request<T = unknown>(endpoint: string, options: RequestInit = {}): Promise<T> {
		// API request implementation:
		// - Add authentication headers
		// - Handle request/response JSON transformation
		// - Manage error responses and token refresh
		// - Provide consistent error handling
		
		const url = `${this.baseUrl}${endpoint}`;
		if (import.meta.env.DEV) {
			console.log(`API client: Making request to ${url}`);
		}
		const token = localStorage.getItem('auth_token');
		
		const headers: Record<string, string> = {
			'Content-Type': 'application/json',
			...(options.headers as Record<string, string>),
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
		// Ensure this only runs in browser context
		if (typeof window === 'undefined') {
			throw new Error('Login can only be called from browser context');
		}
		
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
