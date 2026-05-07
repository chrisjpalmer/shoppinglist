<script lang="ts">
  	import type { Meal } from '../../gen/meal_pb';
  	import { CreateShoppingListService } from '$lib/shopping_list_service';

	const client = CreateShoppingListService()

	let meals: Meal[] = $state([])
	
	async function init() {
		const rs = await client.getMeals({})
		meals = rs.meals
	}

	init()
</script>
<svelte:head>
	<title>Meals</title>
</svelte:head>

<div class="flex flex-row flex-wrap w-full justify-center">
	{#each meals as m (m.id)}
	<a href="/recipies/{m.id}" class="flex flex-col justify-end relative overflow-hidden items-center w-25 h-25 m-1 sm:w-50 sm:h-50 sm:m-4 rounded-md bg-white shadow-md">
		<p class="bg-white/70 z-2 px-3 font-medium w-full text-center sm:text-base text-xs">{m.name}</p>
		<div class="absolute top-0 left-0 w-full h-full bg-cover" style="background-image: url({m.previewImageUrl || '/ramen.svg'})"></div>
	</a>
	{/each}
</div>