<!-- Individual configuration editor page -->
<!-- Edit configuration content with syntax highlighting and validation -->
<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { editorStore } from '$lib/stores/config';
	import { configAPI } from '$lib/utils/configApi';
	import type { UserConfig, ConfigFormat, ConfigVersion } from '$lib/types/config';
	
	// Redirect if not authenticated
	$: if (!$auth.isLoading && !$auth.isAuthenticated) {
		goto('/auth/login');
	}
	
	const configId = parseInt($page.params.id);
	
	// Configuration edit state
	let editingName = false;
	let editingDescription = false;
	let nameInput = '';
	let descriptionInput = '';
	
	// Version comparison
	let showVersionComparison = false;
	let selectedVersions: number[] = [];
	const formatOptions: ConfigFormat[] = ['yaml', 'json', 'toml', 'env'];
	
	onMount(async () => {
		await editorStore.loadConfig(configId);
	});
	
	// Reactive updates for input fields
	$: if ($editorStore.config && !editingName && !editingDescription) {
		nameInput = $editorStore.config.name;
		descriptionInput = $editorStore.config.description || '';
	}
	
	async function saveConfig() {
		if (!$editorStore.config) return;
		
		try {
			// Update metadata if changed
			if (nameInput !== $editorStore.config.name || descriptionInput !== $editorStore.config.description) {
				await configAPI.updateUserConfig($editorStore.config.id, {
					name: nameInput,
					description: descriptionInput,
					format: $editorStore.format
				});
			}
			
			// Save content changes
			if ($editorStore.isDirty) {
				await editorStore.saveConfig('Manual save');
			}
			
			editingName = false;
			editingDescription = false;
		} catch (err) {
			alert(`Failed to save configuration: ${err instanceof Error ? err.message : 'Unknown error'}`);
		}
	}
	
	async function convertFormat(newFormat: ConfigFormat) {
		if ($editorStore.format === newFormat) return;
		
		try {
			await editorStore.changeFormat(newFormat);
		} catch (err) {
			alert(`Failed to convert format: ${err instanceof Error ? err.message : 'Unknown error'}`);
		}
	}
	
	async function validateConfig() {
		try {
			await editorStore.validateConfig();
			if ($editorStore.isValid) {
				alert('Configuration is valid!');
			} else {
				alert(`Validation failed: ${$editorStore.validationError}`);
			}
		} catch (err) {
			alert(`Validation failed: ${err instanceof Error ? err.message : 'Unknown error'}`);
		}
	}
	
	async function restoreVersion(versionId: number) {
		if (!confirm('Are you sure you want to restore this version? Your current changes will be lost.')) {
			return;
		}
		
		try {
			await editorStore.restoreVersion(versionId);
		} catch (err) {
			alert(`Failed to restore version: ${err instanceof Error ? err.message : 'Unknown error'}`);
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
	
	function toggleVersionSelection(versionId: number) {
		if (selectedVersions.includes(versionId)) {
			selectedVersions = selectedVersions.filter(id => id !== versionId);
		} else if (selectedVersions.length < 2) {
			selectedVersions = [...selectedVersions, versionId];
		}
	}
	
	// Auto-save functionality
	$: if ($editorStore.content !== undefined && $editorStore.config) {
		// Debounced auto-save could be implemented here
	}
</script>

<svelte:head>
	<title>{$editorStore.config?.name || 'Loading...'} - Conflux</title>
	<meta name="description" content="Edit configuration file with syntax highlighting and validation" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
	<!-- Header -->
	<div class="bg-white shadow">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center space-x-4">
					<a href="/configs" class="text-blue-600 hover:text-blue-500">
						← Back to Configurations
					</a>
					{#if $editorStore.config}
						<div class="flex items-center space-x-2">
							{#if editingName}
								<input
									type="text"
									bind:value={nameInput}
									class="text-xl font-semibold border-b-2 border-blue-500 bg-transparent focus:outline-none"
									on:blur={() => editingName = false}
									on:keydown={(e) => e.key === 'Enter' && (editingName = false)}
								/>
							{:else}
								<button
									class="text-xl font-semibold text-left hover:text-blue-600"
									on:click={() => editingName = true}
								>
									{nameInput}
								</button>
							{/if}
							<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
								{$editorStore.config.format.toUpperCase()}
							</span>
						</div>
					{/if}
				</div>
				<div class="flex items-center space-x-2">
					<button
						class="inline-flex items-center px-3 py-1 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
						on:click={validateConfig}
					>
						Validate
					</button>
					<button
						class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
						disabled={$editorStore.saving}
						on:click={saveConfig}
					>
						{$editorStore.saving ? 'Saving...' : 'Save'}
					</button>
				</div>
			</div>
		</div>
	</div>
	
	{#if $editorStore.loading}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<div class="animate-pulse">
				<div class="h-4 bg-gray-200 rounded w-1/4 mb-4"></div>
				<div class="h-96 bg-gray-200 rounded"></div>
			</div>
		</div>
	{:else if $editorStore.error}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<p class="mt-1 text-sm text-red-700">{$editorStore.error}</p>
					</div>
				</div>
			</div>
		</div>
	{:else if $editorStore.config}
		<div class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
			<!-- Configuration Details -->
			<div class="bg-white shadow rounded-lg p-6 mb-6">
				<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
					<div class="lg:col-span-2">
						<div class="mb-4">
							<div class="block text-sm font-medium text-gray-700 mb-1">Description</div>
							{#if editingDescription}
								<textarea
									bind:value={descriptionInput}
									class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500"
									rows="2"
									placeholder="Add a description..."
									on:blur={() => editingDescription = false}
								></textarea>
							{:else}
								<button
									class="text-sm text-gray-600 text-left hover:text-blue-600 min-h-[2rem] p-2 border border-transparent hover:border-gray-300 rounded w-full"
									on:click={() => editingDescription = true}
								>
									{descriptionInput || 'Click to add description...'}
								</button>
							{/if}
						</div>
						
						<div class="flex items-center space-x-4 text-sm text-gray-500">
							<span>Last updated: {formatDate($editorStore.config.updated_at)}</span>
							{#if $editorStore.config.template}
								<span>Template: {$editorStore.config.template.display_name}</span>
							{/if}
						</div>
					</div>
					
					<div>
						<div class="block text-sm font-medium text-gray-700 mb-1">Format Conversion</div>
						<div class="grid grid-cols-2 gap-2">
							{#each formatOptions as format}
								<button
									class="px-3 py-2 text-sm rounded-md border {$editorStore.format === format ? 'bg-blue-100 border-blue-300 text-blue-700' : 'bg-gray-50 border-gray-300 text-gray-700 hover:bg-gray-100'}"
									on:click={() => convertFormat(format)}
									disabled={$editorStore.format === format}
								>
									{format.toUpperCase()}
								</button>
							{/each}
						</div>
					</div>
				</div>
			</div>
			
			<!-- Editor and Sidebar -->
			<div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
				<!-- Main Editor -->
				<div class="lg:col-span-3">
					<div class="bg-white shadow rounded-lg">
						<div class="px-6 py-4 border-b border-gray-200">
							<h3 class="text-lg font-medium text-gray-900">Configuration Content</h3>
						</div>
						<div class="p-6">
							<textarea
								bind:value={$editorStore.content}
								class="w-full h-96 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
								placeholder="Enter your configuration content..."
							></textarea>
						</div>
					</div>
				</div>
				
				<!-- Sidebar with Versions -->
				<div class="lg:col-span-1">
					<div class="bg-white shadow rounded-lg">
						<div class="px-6 py-4 border-b border-gray-200">
							<h3 class="text-lg font-medium text-gray-900">Version History</h3>
						</div>
						<div class="max-h-96 overflow-y-auto">
							<div class="divide-y divide-gray-200">
								{#each $editorStore.versions as version}
									<div class="p-4 hover:bg-gray-50">
										<div class="flex items-center justify-between">
											<div class="flex-1 min-w-0">
												<p class="text-sm font-medium text-gray-900">
													v{version.version}
												</p>
												<p class="text-xs text-gray-500">
													{formatDate(version.created_at)}
												</p>
												{#if version.change_note}
													<p class="text-xs text-gray-600 mt-1 truncate">
														{version.change_note}
													</p>
												{/if}
											</div>
											<div class="flex items-center space-x-1">
												<input
													type="checkbox"
													class="h-4 w-4 text-blue-600 border-gray-300 rounded"
													checked={selectedVersions.includes(version.id)}
													on:change={() => toggleVersionSelection(version.id)}
													disabled={selectedVersions.length >= 2 && !selectedVersions.includes(version.id)}
												/>
												<button
													class="text-xs text-blue-600 hover:text-blue-500"
													on:click={() => restoreVersion(version.id)}
												>
													Restore
												</button>
											</div>
										</div>
									</div>
								{/each}
							</div>
						</div>
						
						{#if selectedVersions.length === 2}
							<div class="px-6 py-4 border-t border-gray-200">
								<button
									class="w-full px-3 py-2 text-sm font-medium text-blue-600 bg-blue-50 rounded-md hover:bg-blue-100"
									on:click={() => showVersionComparison = true}
								>
									Compare Versions
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Version Comparison Modal -->
{#if showVersionComparison}
	<div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg max-w-6xl w-full max-h-[90vh] overflow-hidden">
			<div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
				<h3 class="text-lg font-medium text-gray-900">Version Comparison</h3>
				<button
					class="text-gray-400 hover:text-gray-500"
					on:click={() => showVersionComparison = false}
				>
					<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
					</svg>
				</button>
			</div>
			<div class="p-6 overflow-y-auto">
				<p class="text-sm text-gray-600 mb-4">
					Showing differences between selected versions. This would typically include a side-by-side diff view.
				</p>
				<div class="bg-gray-50 p-4 rounded-md">
					<p class="text-sm text-gray-700">
						Version comparison feature would be implemented here with a proper diff library.
					</p>
				</div>
			</div>
		</div>
	</div>
{/if}
