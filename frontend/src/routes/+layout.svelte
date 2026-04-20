<script lang="ts">
	import Header from './Header.svelte';
	import '../app.css';

	let { children } = $props();
</script>

<div class="flex flex-col min-h-screen">
	<Header />

	<main class="flex grow-1 flex-col w-full">
		{@render children()}
	</main>

	<footer class="flex flex-col justify-center items-center p-1 sm:px-0">
		<p class="font-normal leading-5">
			visit <a class="font-bold text-orange no-underline hover:underline" href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to learn about SvelteKit
		</p>
	</footer>
</div>