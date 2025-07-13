<!-- Home page component -->
<!-- Landing page with application overview and navigation -->
<!-- Entry point for unauthenticated users -->
<script lang="ts">
	import { auth } from '$lib/stores/auth';
	
	// Reactive statements for authentication state
	$: isAuthenticated = $auth.isAuthenticated;
	$: user = $auth.user;
</script>

<!-- SEO meta tags for home page -->
<svelte:head>
	<title>Conflux - Configuration Management for Self-Hosted Apps</title>
	<meta name="description" content="Manage configuration files for your self-hosted applications with version control, format conversion, and template support" />
</svelte:head>

<!-- Home page content -->
<div class="container mx-auto px-4 py-8">
	<section class="text-center py-16">
		<h1 class="text-4xl font-bold mb-4">
			Welcome to Conflux
		</h1>
		
		<p class="text-xl text-gray-600 mb-8">
			Configuration management for self-hosted applications with version control and format flexibility
		</p>
		
		{#if isAuthenticated}
			<!-- Authenticated user content -->
			<div class="space-y-4">
				<p class="text-lg">Welcome back, {user?.first_name}!</p>
				<a href="/dashboard" class="inline-block bg-blue-600 text-white px-6 py-3 rounded-md hover:bg-blue-700">
					Go to Dashboard
				</a>
			</div>
		{:else}
			<!-- Guest user content -->
			<div class="space-x-4">
				<a href="/auth/login" class="inline-block bg-blue-600 text-white px-6 py-3 rounded-md hover:bg-blue-700">
					Login
				</a>
				<a href="/auth/register" class="inline-block border border-blue-600 text-blue-600 px-6 py-3 rounded-md hover:bg-blue-50">
					Sign Up
				</a>
			</div>
		{/if}
	</section>
	
	<section class="grid md:grid-cols-3 gap-8 mt-16">
		<div class="text-center p-6 bg-white rounded-lg shadow-sm">
			<h3 class="text-lg font-semibold mb-2">Configuration Templates</h3>
			<p class="text-gray-600">Pre-built templates for popular self-hosted applications like cross-seed, with automatic format detection.</p>
		</div>
		
		<div class="text-center p-6 bg-white rounded-lg shadow-sm">
			<h3 class="text-lg font-semibold mb-2">Version Control</h3>
			<p class="text-gray-600">Track changes to your configurations with full version history and easy rollback capabilities.</p>
		</div>
		
		<div class="text-center p-6 bg-white rounded-lg shadow-sm">
			<h3 class="text-lg font-semibold mb-2">Format Conversion</h3>
			<p class="text-gray-600">Convert between YAML, JSON, TOML, and ENV formats with real-time validation and syntax highlighting.</p>
		</div>
	</section>
</div>
