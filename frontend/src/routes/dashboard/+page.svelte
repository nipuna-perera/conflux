<!-- Dashboard page component -->
<!-- Protected area for authenticated users -->
<!-- Shows user profile and application features -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { apiClient } from '$lib/utils/api';
	
	// Redirect if not authenticated
	$: if (!$auth.isLoading && !$auth.isAuthenticated) {
		goto('/auth/login');
	}
	
	$: user = $auth.user;
	
	// Handle logout
	async function handleLogout() {
		try {
			await apiClient.logout();
		} catch (error) {
			console.error('Logout error:', error);
		} finally {
			auth.logout();
			goto('/');
		}
	}
</script>

<svelte:head>
	<title>Dashboard - Full-Stack App</title>
</svelte:head>

{#if $auth.isLoading}
	<div class="min-h-screen flex items-center justify-center">
		<div class="text-center">
			<div class="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
			<p class="mt-4 text-gray-600">Loading...</p>
		</div>
	</div>
{:else if user}
	<div class="min-h-screen bg-gray-50">
		<!-- Dashboard Header -->
		<div class="bg-white shadow">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div class="flex justify-between h-16">
					<div class="flex items-center">
						<h1 class="text-xl font-semibold">Dashboard</h1>
					</div>
					<div class="flex items-center space-x-4">
						<span class="text-gray-700">Welcome, {user.first_name}!</span>
						<button
							on:click={handleLogout}
							class="bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700"
						>
							Logout
						</button>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Dashboard Content -->
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<!-- User Profile Card -->
			<div class="bg-white overflow-hidden shadow rounded-lg mb-6">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900">Profile Information</h3>
					<div class="mt-5 grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
						<div>
							<label class="block text-sm font-medium text-gray-700">First Name</label>
							<div class="mt-1 text-sm text-gray-900">{user.first_name}</div>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700">Last Name</label>
							<div class="mt-1 text-sm text-gray-900">{user.last_name}</div>
						</div>
						<div class="sm:col-span-2">
							<label class="block text-sm font-medium text-gray-700">Email</label>
							<div class="mt-1 text-sm text-gray-900">{user.email}</div>
						</div>
					</div>
				</div>
			</div>
			
			<!-- Feature Cards -->
			<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
				<div class="bg-white overflow-hidden shadow rounded-lg">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
									<span class="text-white font-bold">P</span>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Profile</dt>
									<dd class="text-lg font-medium text-gray-900">Manage your account</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
				
				<div class="bg-white overflow-hidden shadow rounded-lg">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center">
									<span class="text-white font-bold">S</span>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Settings</dt>
									<dd class="text-lg font-medium text-gray-900">Configure preferences</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
				
				<div class="bg-white overflow-hidden shadow rounded-lg">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-purple-500 rounded-md flex items-center justify-center">
									<span class="text-white font-bold">A</span>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Analytics</dt>
									<dd class="text-lg font-medium text-gray-900">View insights</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
