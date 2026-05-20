<script lang="ts">
  	import { ImageMode, type ImageMeta } from '../../gen/meal_pb';
  	import { BackendUrl, CreateShoppingListService } from '$lib/shopping_list_service';

	const client = CreateShoppingListService()

	interface DisplayMeal {
		 id: bigint
		 name: string
		 previewImageUrl: string
	}

	let meals: DisplayMeal[] = $state([])

	async function init() {
		const rs = await client.getMeals({})
		meals = rs.meals.map(m => {
			return {
				id: m.id,
				name: m.name,
				previewImageUrl: mapImageSource(m.previewImage)
			}
		})
	}

	function mapImageSource(img:ImageMeta | undefined): string {
		if (!img || img.mode == ImageMode.IM_NONE) {
			return ''
		}

		if (img.mode == ImageMode.IM_INTERNAL) {
			return BackendUrl() + img.internalUrl
		} 
		
		return img.externalUrl
	}

	init()
</script>
<svelte:head>
	<title>Meals</title>
</svelte:head>

<div class="flex flex-row flex-wrap w-full justify-center overflow-y-auto h-full">
	{#each meals as m (m.id)}
	<a href="/recipies/{m.id}" class="flex flex-col justify-end relative overflow-hidden items-center w-25 h-25 m-1 sm:w-50 sm:h-50 sm:m-4 rounded-md bg-white shadow-md">
		<p class="bg-white/70 z-2 px-3 font-medium w-full text-center sm:text-base text-xs">{m.name}</p>
		<div class="absolute top-0 left-0 w-full h-full bg-cover" style="background-image: url({m.previewImageUrl || '/ramen.svg'})"></div>
	</a>
	{/each}
</div>