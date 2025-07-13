<!-- User configurations listing page -->
<!-- Manage user configurations with search and filtering -->
<!-- Access to edit, delete, and export configurations -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { configsStore, templatesStore } from '$lib/stores/config';
	import { configAPI } from '$lib/utils/configApi';
	import type { ConfigFormat } from '$lib/types/config';
	
	// Redirect if not authenticated
	$: {
		if (!$auth.isLoading && !$auth.isAuthenticated) {
			goto('/auth/login');
		}
	}
	
	let searchQuery = '';
	let selectedTemplate = '';
	let currentPage = 1;
	const pageSize = 20;
	
	// Load data on mount
	onMount(() => {
		loadConfigs();
		templatesStore.loadTemplates(); // For filter dropdown
	});
	
	// Reactive loading when filters change
	$: {
		if (searchQuery !== undefined || selectedTemplate !== undefined) {
			currentPage = 1;
			loadConfigs();
		}
	}
	
	async function loadConfigs() {
		await configsStore.loadConfigs({
			template_id: selectedTemplate ? parseInt(selectedTemplate) : undefined,
			page: currentPage,
			limit: pageSize
		});
	}
	
	async function deleteConfig(configId: number, configName: string) {
		if (confirm(`Are you sure you want to delete "${configName}"? This action cannot be undone.`)) {
			try {
				await configsStore.deleteConfig(configId);
			} catch (error) {
				alert(`Failed to delete configuration: ${error instanceof Error ? error.message : 'Unknown error'}`);
			}
		}
	}
	
	async function exportConfig(configId: number, format: ConfigFormat, configName: string) {
		try {
			const blob = await configAPI.exportConfig(configId, format);
			const url = window.URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = `${configName}.${format}`;
			document.body.appendChild(a);
			a.click();
			window.URL.revokeObjectURL(url);
			document.body.removeChild(a);
		} catch (error) {
			alert(`Failed to export configuration: ${error instanceof Error ? error.message : 'Unknown error'}`);
		}
	}
	
	function handleSearch(event: Event) {
		const target = event.target as HTMLInputElement;
		searchQuery = target.value;
	}
	
	function handleTemplateChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		selectedTemplate = target.value;
	}
	
	function nextPage() {
		if (currentPage * pageSize < $configsStore.total) {
			currentPage++;
			loadConfigs();
		}
	}
	
	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			loadConfigs();
		}
	}
	
	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<svelte:head>
	<title>My Configurations - Conflux</title>
	<meta name="description" content="Manage your configuration files with version control" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Header -->
	<div class="bg-white shadow">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<a href="/dashboard" class="text-blue-600 hover:text-blue-500 mr-4">
						← Dashboard
					</a>
					<h1 class="text-xl font-semibold">My Configurations</h1>
				</div>
				<div class="flex items-center">
					<a
						href="/templates"
						class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
					>
						Create New
					</a>
				</div>
			</div>
		</div>
	</div>
	
	<!-- Main Content -->
	<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
		<!-- Search and Filters -->
		<div class="bg-white shadow rounded-lg p-6 mb-6">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div>
					<label for="search" class="block text-sm font-medium text-gray-700 mb-1">Search Configurations</label>
					<input
						id="search"
						type="text"
						placeholder="Search by name..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						value={searchQuery}
						on:input={handleSearch}
					/>
				</div>
				<div>
					<label for="template" class="block text-sm font-medium text-gray-700 mb-1">Filter by Template</label>
					<select
						id="template"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						value={selectedTemplate}
						on:change={handleTemplateChange}
					>
						<option value="">All Templates</option>
						{#each $templatesStore.templates as template (template.id)}
							<option value={template.id}>{template.display_name}</option>
						{/each}
					</select>
				</div>
			</div>
		</div>
		
		<!-- Configurations List -->
		{#if $configsStore.loading}
			<div class="bg-white shadow rounded-lg">
				<div class="px-6 py-4 border-b border-gray-200">
					<div class="h-4 bg-gray-200 rounded w-1/4 animate-pulse"></div>
				</div>
				<div class="divide-y divide-gray-200">
					{#each Array(5) as _, i (i)}
						<div class="px-6 py-4 animate-pulse">
							<div class="flex items-center justify-between">
								<div class="flex-1">
									<div class="h-4 bg-gray-200 rounded w-1/3 mb-2"></div>
									<div class="h-3 bg-gray-200 rounded w-1/4"></div>
								</div>
								<div class="flex space-x-2">
									<div class="h-8 bg-gray-200 rounded w-16"></div>
									<div class="h-8 bg-gray-200 rounded w-16"></div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else if $configsStore.error}
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error loading configurations</h3>
						<p class="mt-1 text-sm text-red-700">{$configsStore.error}</p>
					</div>
				</div>
			</div>
		{:else if $configsStore.configs.length > 0}
			<div class="bg-white shadow rounded-lg">
				<div class="px-6 py-4 border-b border-gray-200">
					<h3 class="text-lg font-medium text-gray-900">
						{$configsStore.total} Configuration{$configsStore.total !== 1 ? 's' : ''}
					</h3>
				</div>
				<div class="divide-y divide-gray-200">
					{#each $configsStore.configs as config (config.id)}
						<div class="px-6 py-4 hover:bg-gray-50">
							<div class="flex items-center justify-between">
								<div class="flex-1 min-w-0">
									<div class="flex items-center space-x-3">
										<h4 class="text-sm font-medium text-gray-900 truncate">
											{config.name}
										</h4>
										<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
											{config.format.toUpperCase()}
										</span>
										{#if config.template}
											<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
												{config.template.display_name}
											</span>
										{/if}
									</div>
									<div class="mt-1 flex items-center space-x-4 text-sm text-gray-500">
										<span>Updated {formatDate(config.updated_at)}</span>
										{#if config.description}
											<span class="truncate max-w-md">{config.description}</span>
										{/if}
									</div>
								</div>
								<div class="flex items-center space-x-2">
									<!-- Export dropdown -->
									<div class="relative inline-block text-left">
										<button
											class="inline-flex items-center px-3 py-1 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
											type="button"
											on:click={(e) => {
												const menu = e.currentTarget?.nextElementSibling;
												if (menu) menu.classList.toggle('hidden');
											}}
										>
											Export
											<svg class="ml-1 -mr-1 h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
												<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
											</svg>
										</button>
										<div class="hidden absolute right-0 z-10 mt-2 w-32 origin-top-right bg-white shadow-lg ring-1 ring-black ring-opacity-5">
											<div class="py-1">
												<button
													class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => exportConfig(config.id, 'yaml', config.name)}
												>
													YAML
												</button>
												<button
													class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => exportConfig(config.id, 'json', config.name)}
												>
													JSON
												</button>
												<button
													class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => exportConfig(config.id, 'toml', config.name)}
												>
													TOML
												</button>
												<button
													class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
													on:click={() => exportConfig(config.id, 'env', config.name)}
												>
													ENV
												</button>
											</div>
										</div>
									</div>
									
									<a
										href="/configs/{config.id}"
										class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
									>
										Edit
									</a>
									
									<button
										class="inline-flex items-center px-3 py-1 border border-red-300 text-sm leading-4 font-medium rounded-md text-red-700 bg-white hover:bg-red-50"
										on:click={() => deleteConfig(config.id, config.name)}
									>
										Delete
									</button>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
			
			<!-- Pagination -->
			{#if $configsStore.total > pageSize}
				<div class="mt-6 flex items-center justify-between">
					<div class="text-sm text-gray-700">
						Showing {((currentPage - 1) * pageSize) + 1} to {Math.min(currentPage * pageSize, $configsStore.total)} of {$configsStore.total} configurations
					</div>
					<div class="flex space-x-2">
						<button
							class="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
							disabled={currentPage === 1}
							on:click={prevPage}
						>
							Previous
						</button>
						<button
							class="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
							disabled={currentPage * pageSize >= $configsStore.total}
							on:click={nextPage}
						>
							Next
						</button>
					</div>
				</div>
			{/if}
		{:else}
			<div class="text-center py-12">
				<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
				</svg>
				<h3 class="mt-2 text-sm font-medium text-gray-900">No configurations found</h3>
				<p class="mt-1 text-sm text-gray-500">
					{searchQuery || selectedTemplate ? 'Try adjusting your search criteria' : 'Get started by creating your first configuration'}
				</p>
				<div class="mt-6">
					<a
						href="/templates"
						class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
					>
						Create from Template
					</a>
				</div>
			</div>
		{/if}
	</div>
</div>
