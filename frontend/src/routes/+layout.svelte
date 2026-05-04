<script lang="ts">
	import Nav from './Nav.svelte';
	import '../app.css';
    import MenuBar from './MenuBar.svelte';

	let { children } = $props();
</script>

<div class="flex flex-col bg-gray-100">
	<MenuBar></MenuBar>
	<div class="h-screen overflow-y-auto pt-26 sm:pt-16 pb-10">
		<div class="flex flex-row px-8">
			<div class="hidden sm:block w-40 flex-shrink-0">
				<Nav />
			</div>
			<div class="flex flex-grow-1 flex-shrink-1 flex-basis-0 min-w-0 flex-col">
				{@render children()}
			</div>
		</div>
	</div>
</div>