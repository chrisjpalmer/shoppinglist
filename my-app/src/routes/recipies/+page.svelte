<script lang="ts">
	import { ShoppingListService } from '../../gen/shopping_list_service_pb';
	import { createClient } from "@connectrpc/connect";
	import { createConnectTransport } from "@connectrpc/connect-web";
  	import type { Meal } from '../../gen/meal_pb';

	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createClient(ShoppingListService, transport);

	let meals: Meal[] = $state([])

	async function refresh() {
		const rs = await client.getMeals({})
		meals = rs.meals
	}

	refresh()

</script>

<svelte:head>
	<title>Recipies</title>
	<meta name="description" content="About this app" />
</svelte:head>

<div class="text-column">
	<h1>Recipies</h1>

	<select>
		{#each meals as meal}
			<option value="{meal.id}">{meal.name}</option>
		{/each}
	</select>
</div>
