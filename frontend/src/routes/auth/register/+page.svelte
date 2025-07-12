<!-- Registration page component -->
<!-- Handles user registration with form validation -->
<!-- Creates new user accounts -->
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
	let confirmPassword = '';
	let firstName = '';
	let lastName = '';
	let loading = false;
	let error = '';
	
	// Handle form submission
	async function handleRegister() {
		if (!email || !password || !firstName || !lastName) {
			error = 'Please fill in all fields';
			return;
		}
		
		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}
		
		if (password.length < 8) {
			error = 'Password must be at least 8 characters long';
			return;
		}
		
		loading = true;
		error = '';
		
		try {
			const response = await apiClient.register({
				email,
				password,
				first_name: firstName,
				last_name: lastName
			});
			auth.login(response.user, response.token);
			goto('/dashboard');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Register - Full-Stack App</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
	<div class="max-w-md w-full space-y-8">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Create your account
			</h2>
		</div>
		
		<form class="mt-8 space-y-6" on:submit|preventDefault={handleRegister}>
			{#if error}
				<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
					{error}
				</div>
			{/if}
			
			<div class="grid grid-cols-2 gap-4">
				<div>
					<label for="firstName" class="block text-sm font-medium text-gray-700">First Name</label>
					<input
						id="firstName"
						type="text"
						bind:value={firstName}
						required
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						placeholder="First name"
					/>
				</div>
				
				<div>
					<label for="lastName" class="block text-sm font-medium text-gray-700">Last Name</label>
					<input
						id="lastName"
						type="text"
						bind:value={lastName}
						required
						class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						placeholder="Last name"
					/>
				</div>
			</div>
			
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
			
			<div>
				<label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirm Password</label>
				<input
					id="confirmPassword"
					type="password"
					bind:value={confirmPassword}
					required
					class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
					placeholder="Confirm your password"
				/>
			</div>
			
			<button
				type="submit"
				disabled={loading}
				class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
			>
				{loading ? 'Creating account...' : 'Create account'}
			</button>
		</form>
		
		<div class="text-center">
			<a href="/auth/login" class="text-blue-600 hover:text-blue-500">
				Already have an account? Sign in
			</a>
		</div>
	</div>
</div>
