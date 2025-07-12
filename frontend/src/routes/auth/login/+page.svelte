<!-- Login page component -->
<!-- Handles user authentication with form validation -->
<!-- Redirects authenticated users to dashboard -->
<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { apiClient } from '$lib/utils/api';
	
	// Redirect if already authenticated
	$: if ($auth.isAuthenticated) {
		goto('/dashboard');
	}
	
	// Form state
	let email = '';
	let password = '';
	let loading = false;
	let error = '';
	
	// Handle form submission
	async function handleLogin() {
		if (!email || !password) {
			error = 'Please fill in all fields';
			return;
		}
		
		loading = true;
		error = '';
		
		try {
			const response = await apiClient.login(email, password);
			auth.login(response.user, response.token);
			goto('/dashboard');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login - Full-Stack App</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Sign in to your account
			</h2>
		</div>
		
		<form class="mt-8 space-y-6" on:submit|preventDefault={handleLogin}>
			{#if error}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
					{error}
				</div>
			{/if}
			
			<div>
				<label for="email" class="block text-sm font-medium text-gray-700">Email</label>
				<input
					id="email"
					type="email"
					bind:value={email}
					required
					class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
					placeholder="Enter your email"
				/>
			</div>
			
			<div>
				<label for="password" class="block text-sm font-medium text-gray-700">Password</label>
				<input
					id="password"
					type="password"
					bind:value={password}
					required
					class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
					placeholder="Enter your password"
				/>
			</div>
			
			<button
				type="submit"
				disabled={loading}
				class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
			>
				{loading ? 'Signing in...' : 'Sign in'}
			</button>
		</form>
		
		<div class="text-center">
			<a href="/auth/register" class="text-blue-600 hover:text-blue-500">
				Don't have an account? Sign up
			</a>
		</div>
	</div>
</div>
