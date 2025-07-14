<!-- Dashboard page component -->
<!-- Protected area for authenticated users -->
<!-- Shows configuration overview and management features -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { configsStore, templatesStore } from '$lib/stores/config';
	import { apiClient } from '$lib/utils/api';
	
	// Redirect if not authenticated
	$: if (!$auth.isLoading && !$auth.isAuthenticated) {
		goto('/auth/login');
	}
	
	$: user = $auth.user;
	
	// Load dashboard data
	onMount(async () => {
		if ($auth.isAuthenticated) {
			configsStore.loadConfigs({ limit: 5 }); // Load recent configs
			templatesStore.loadTemplates({ limit: 6 }); // Load popular templates
		}
	});
	
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
	<title>Dashboard - Conflux</title>
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
			<!-- Welcome Section -->
			<div class="bg-white overflow-hidden shadow rounded-lg mb-6">
				<div class="px-4 py-5 sm:p-6">
					<h3 class="text-lg leading-6 font-medium text-gray-900">Welcome to Conflux</h3>
					<p class="mt-1 text-sm text-gray-600">
						Manage your configuration files with version control and format flexibility.
					</p>
					<div class="mt-4">
						<a href="/configs" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700">
							Manage Configurations
						</a>
						<a href="/templates" class="ml-3 inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50">
							Browse Templates
						</a>
					</div>
				</div>
			</div>
			
			<!-- Quick Stats -->
			<div class="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4 mb-6">
				<div class="bg-white overflow-hidden shadow rounded-lg">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-blue-500 rounded-md flex items-center justify-center">
									<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
									</svg>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Configurations</dt>
									<dd class="text-lg font-medium text-gray-900">{$configsStore.total}</dd>
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
									<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
									</svg>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Templates</dt>
									<dd class="text-lg font-medium text-gray-900">{$templatesStore.total}</dd>
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
									<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h8a1 1 0 011 1v2m3 0H4a1 1 0 00-1 1v16a1 1 0 001 1h16a1 1 0 001-1V5a1 1 0 00-1-1z" />
									</svg>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Formats</dt>
									<dd class="text-lg font-medium text-gray-900">4</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
				
				<div class="bg-white overflow-hidden shadow rounded-lg">
					<div class="p-5">
						<div class="flex items-center">
							<div class="flex-shrink-0">
								<div class="w-8 h-8 bg-orange-500 rounded-md flex items-center justify-center">
									<svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>
								</div>
							</div>
							<div class="ml-5 w-0 flex-1">
								<dl>
									<dt class="text-sm font-medium text-gray-500 truncate">Recent Activity</dt>
									<dd class="text-lg font-medium text-gray-900">Today</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
			</div>

			<!-- Recent Configurations -->
			<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
				<div class="bg-white shadow rounded-lg">
					<div class="px-4 py-5 sm:p-6">
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg leading-6 font-medium text-gray-900">Recent Configurations</h3>
							<a href="/configs" class="text-sm text-blue-600 hover:text-blue-500">View all</a>
						</div>
						
						{#if $configsStore.loading}
							<div class="animate-pulse space-y-3">
								{#each Array(3) as _, i (i)}
									<div class="h-4 bg-gray-200 rounded w-3/4"></div>
								{/each}
							</div>
						{:else if $configsStore.configs.length > 0}
							<div class="space-y-3">
								{#each $configsStore.configs.slice(0, 5) as config (config.id)}
									<div class="flex items-center justify-between p-3 border border-gray-200 rounded-md">
										<div>
											<p class="text-sm font-medium text-gray-900">{config.name}</p>
											<p class="text-xs text-gray-500">{config.format.toUpperCase()} • {new Date(config.updated_at).toLocaleDateString()}</p>
										</div>
										<a href="/configs/{config.id}" class="text-blue-600 hover:text-blue-500 text-sm">Edit</a>
									</div>
								{/each}
							</div>
						{:else}
							<div class="text-center py-6">
								<p class="text-gray-500">No configurations yet</p>
								<a href="/templates" class="mt-2 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-blue-700 bg-blue-100 hover:bg-blue-200">
									Create from template
								</a>
							</div>
						{/if}
					</div>
				</div>

				<!-- Popular Templates -->
				<div class="bg-white shadow rounded-lg">
					<div class="px-4 py-5 sm:p-6">
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg leading-6 font-medium text-gray-900">Popular Templates</h3>
							<a href="/templates" class="text-sm text-blue-600 hover:text-blue-500">View all</a>
						</div>
						
						{#if $templatesStore.loading}
							<div class="animate-pulse space-y-3">
								{#each Array(3) as _, i (i)}
									<div class="h-4 bg-gray-200 rounded w-3/4"></div>
								{/each}
							</div>
						{:else if $templatesStore.templates.length > 0}
							<div class="space-y-3">
								{#each $templatesStore.templates.slice(0, 5) as template (template.id)}
									<div class="flex items-center justify-between p-3 border border-gray-200 rounded-md">
										<div>
											<p class="text-sm font-medium text-gray-900">{template.display_name}</p>
											<p class="text-xs text-gray-500">{template.category} • {template.format.toUpperCase()}</p>
										</div>
										<button class="text-blue-600 hover:text-blue-500 text-sm">Use</button>
									</div>
								{/each}
							</div>
						{:else}
							<div class="text-center py-6">
								<p class="text-gray-500">No templates available</p>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
