<!-- Root layout component for the entire application -->
<!-- Provides common UI elements, navigation, and global styles -->
<!-- Wraps all pages with consistent layout structure -->
<script lang="ts">
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth';
	import '../app.css'; // Global styles and Tailwind CSS
	
	// Initialize authentication state on app load
	onMount(() => {
		auth.init();
	});
</script>

<!-- Global layout structure -->
<div class="min-h-screen flex flex-col">
	<header class="bg-white shadow-sm border-b">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between h-16">
				<div class="flex items-center">
					<h1 class="text-xl font-semibold">Full-Stack App</h1>
				</div>
				<nav class="flex items-center space-x-4">
					<a href="/" class="text-gray-600 hover:text-gray-900">Home</a>
					<a href="/auth/login" class="text-gray-600 hover:text-gray-900">Login</a>
				</nav>
			</div>
		</div>
	</header>
	
	<main class="flex-1">
		<!-- Page content is rendered here -->
		<slot />
	</main>
	
	<footer class="bg-gray-100 border-t">
		<div class="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
			<p class="text-center text-gray-600">&copy; 2025 Full-Stack App. All rights reserved.</p>
		</div>
	</footer>
</div>
