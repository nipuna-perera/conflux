// Authentication state management using Svelte stores
// Manages user session, login state, and authentication persistence
// Provides reactive authentication state across the application
import { writable } from 'svelte/store';
import type { User } from '$lib/types/user';

// Authentication store interface
interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	isLoading: boolean;
	token: string | null;
}

// Initial authentication state
const initialState: AuthState = {
	user: null,
	isAuthenticated: false,
	isLoading: false,
	token: null
};

// Create authentication store
function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initialState);
	
	return {
		subscribe,
		
		// Login action - updates store with user data and token
		login: (user: User, token: string) => {
			// Store authentication data and persist to localStorage
			localStorage.setItem('auth_token', token);
			localStorage.setItem('auth_user', JSON.stringify(user));
			
			set({
				user,
				token,
				isAuthenticated: true,
				isLoading: false
			});
		},
		
		// Logout action - clears authentication state
		logout: () => {
			// Clear authentication data and remove from localStorage
			localStorage.removeItem('auth_token');
			localStorage.removeItem('auth_user');
			
			set(initialState);
		},
		
		// Initialize auth state from stored token
		init: async () => {
			// Load authentication state from localStorage
			// Validate stored token with backend
			update(state => ({ ...state, isLoading: true }));
			
			const token = localStorage.getItem('auth_token');
			const userStr = localStorage.getItem('auth_user');
			
			if (token && userStr) {
				try {
					const user = JSON.parse(userStr);
					set({
						user,
						token,
						isAuthenticated: true,
						isLoading: false
					});
				} catch (error) {
					console.error('Failed to parse stored user data:', error);
					localStorage.removeItem('auth_token');
					localStorage.removeItem('auth_user');
					set({ ...initialState, isLoading: false });
				}
			} else {
				set({ ...initialState, isLoading: false });
			}
		},
		
		// Update user profile data
		updateUser: (user: User) => {
			// Update user data in store
			localStorage.setItem('auth_user', JSON.stringify(user));
			update(state => ({ ...state, user }));
		}
	};
}

export const auth = createAuthStore();
