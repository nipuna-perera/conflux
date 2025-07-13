<!-- Templates listing page -->
<!-- Browse and search configuration templates -->
<!-- Create new configurations from templates -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { templatesStore, configsStore } from '$lib/stores/config';
	
	// Redirect if not authenticated
	$: if (!$auth.isLoading && !$auth.isAuthenticated) {
		goto('/auth/login');
	}
	
	let searchQuery = '';
	let selectedCategory = '';
	let currentPage = 1;
	const pageSize = 12;
	
	// Categories for filtering
	const categories = [
		'',
		'torrenting',
		'media',
		'networking',
		'monitoring',
		'automation',
		'other'
	];
	
	// Load templates on mount and when filters change
	onMount(() => loadTemplates());
	
	$: {
		// Reactive loading when filters change
		if (searchQuery !== undefined || selectedCategory !== undefined) {
			currentPage = 1;
			loadTemplates();
		}
	}
	
	async function loadTemplates() {
		await templatesStore.loadTemplates({
			search: searchQuery || undefined,
			category: selectedCategory || undefined,
			page: currentPage,
			limit: pageSize
		});
	}
	
	async function createFromTemplate(templateId: number, templateName: string) {
		const configName = prompt(`Enter a name for your ${templateName} configuration:`);
		if (configName && configName.trim()) {
			try {
				const newConfig = await configsStore.createConfig(templateId, configName.trim());
				goto(`/configs/${newConfig.id}`);
			} catch (error) {
				alert(`Failed to create configuration: ${error instanceof Error ? error.message : 'Unknown error'}`);
			}
		}
	}
	
	function handleSearch(event: Event) {
		const target = event.target as HTMLInputElement;
		searchQuery = target.value;
	}
	
	function handleCategoryChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		selectedCategory = target.value;
	}
	
	function nextPage() {
		if (currentPage * pageSize < $templatesStore.total) {
			currentPage++;
			loadTemplates();
		}
	}
	
	function prevPage() {
		if (currentPage > 1) {
			currentPage--;
			loadTemplates();
		}
	}
</script>

<svelte:head>
	<title>Configuration Templates - Conflux</title>
	<meta name="description" content="Browse and use configuration templates for self-hosted applications" />
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
					<h1 class="text-xl font-semibold">Configuration Templates</h1>
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
					<label for="search" class="block text-sm font-medium text-gray-700 mb-1">Search Templates</label>
					<input
						id="search"
						type="text"
						placeholder="Search by name or description..."
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						value={searchQuery}
						on:input={handleSearch}
					/>
				</div>
				<div>
					<label for="category" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
					<select
						id="category"
						class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
						value={selectedCategory}
						on:change={handleCategoryChange}
					>
						<option value="">All Categories</option>
						{#each categories.slice(1) as category}
							<option value={category}>{category.charAt(0).toUpperCase() + category.slice(1)}</option>
						{/each}
					</select>
				</div>
			</div>
		</div>
		
		<!-- Templates Grid -->
		{#if $templatesStore.loading}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each Array(6) as _}
					<div class="bg-white shadow rounded-lg p-6 animate-pulse">
						<div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
						<div class="h-3 bg-gray-200 rounded w-1/2 mb-4"></div>
						<div class="h-3 bg-gray-200 rounded w-full mb-2"></div>
						<div class="h-3 bg-gray-200 rounded w-2/3"></div>
					</div>
				{/each}
			</div>
		{:else if $templatesStore.error}
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error loading templates</h3>
						<p class="mt-1 text-sm text-red-700">{$templatesStore.error}</p>
					</div>
				</div>
			</div>
		{:else if $templatesStore.templates.length > 0}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each $templatesStore.templates as template}
					<div class="bg-white shadow rounded-lg p-6 hover:shadow-lg transition-shadow">
						<div class="flex items-start justify-between mb-4">
							<div class="flex-1">
								<h3 class="text-lg font-medium text-gray-900 mb-1">{template.display_name}</h3>
								<p class="text-sm text-gray-600 mb-2">{template.description}</p>
								<div class="flex items-center space-x-2">
									<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
										{template.category}
									</span>
									<span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
										{template.format.toUpperCase()}
									</span>
								</div>
							</div>
						</div>
						
						<div class="border-t pt-4">
							<div class="flex justify-between items-center">
								<div class="text-sm text-gray-500">
									v{template.version}
								</div>
								<div class="space-x-2">
									<button
										class="inline-flex items-center px-3 py-1 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
										on:click={() => goto(`/templates/${template.id}`)}
									>
										Preview
									</button>
									<button
										class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
										on:click={() => createFromTemplate(template.id, template.display_name)}
									>
										Use Template
									</button>
								</div>
							</div>
						</div>
					</div>
				{/each}
			</div>
			
			<!-- Pagination -->
			{#if $templatesStore.total > pageSize}
				<div class="mt-8 flex items-center justify-between">
					<div class="text-sm text-gray-700">
						Showing {((currentPage - 1) * pageSize) + 1} to {Math.min(currentPage * pageSize, $templatesStore.total)} of {$templatesStore.total} templates
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
							disabled={currentPage * pageSize >= $templatesStore.total}
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
				<h3 class="mt-2 text-sm font-medium text-gray-900">No templates found</h3>
				<p class="mt-1 text-sm text-gray-500">
					{searchQuery || selectedCategory ? 'Try adjusting your search criteria' : 'No templates are currently available'}
				</p>
			</div>
		{/if}
	</div>
</div>
