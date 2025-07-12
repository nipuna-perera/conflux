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
	<title>Full-Stack App - Home</title>
	<meta name="description" content="Welcome to our full-stack application" />
</svelte:head>

<!-- Home page content -->
<div class="container mx-auto px-4 py-8">
	<section class="text-center py-16">
		<h1 class="text-4xl font-bold mb-4">
			Welcome to Full-Stack App
		</h1>
		
		<p class="text-xl text-gray-600 mb-8">
			Modern web application with Go backend and SvelteKit frontend
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
			<h3 class="text-lg font-semibold mb-2">Go Backend</h3>
			<p class="text-gray-600">Robust API built with Go, featuring clean architecture and database flexibility.</p>
		</div>
		
		<div class="text-center p-6 bg-white rounded-lg shadow-sm">
			<h3 class="text-lg font-semibold mb-2">SvelteKit Frontend</h3>
			<p class="text-gray-600">Modern, reactive frontend with TypeScript and Tailwind CSS.</p>
		</div>
		
		<div class="text-center p-6 bg-white rounded-lg shadow-sm">
			<h3 class="text-lg font-semibold mb-2">Authentication</h3>
			<p class="text-gray-600">Secure JWT-based authentication with session management.</p>
		</div>
	</section>
</div>
