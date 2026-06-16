<script lang="ts">
  	import { type Meal, type IngredientRef, type ImageMeta, ImageMode } from '../../../gen/meal_pb';
  	import type { Ingredient } from '../../../gen/ingredient_pb';
  	import { BackendUrl, CreateShoppingListService } from '$lib/shopping_list_service';
	import Select from '../../../components/select.svelte';
	import Table from '../../../components/table.svelte';
	import Td from '../../../components/td.svelte';
	import TrHeader from '../../../components/tr-header.svelte';
	import TrTitle from '../../../components/tr-title.svelte';
	import Tr from '../../../components/tr.svelte';

	import type { PageProps } from './$types';
  	import H1 from '../../../components/h1.svelte';
  	import Button from '../../../components/button.svelte';

	let { params }: PageProps = $props();

	const client = CreateShoppingListService()

	let ingredients: Ingredient[] = $state([])
	let isEditing: boolean = $state(false)
	let meals: Meal[];
	let psuedoId = 10000

	async function refresh() {
		const gmRs = await client.getMeals({})
		meals = gmRs.meals

		const giRs = await client.getIngredients({});
		ingredients = giRs.ingredients

		await setSelectedMealId(BigInt(params.meal))
		
	}
	interface SelectedMeal {
		_meal: Meal
		id: bigint
		name: string
		recipeUrl: string
		ingredientsImageUrl: string
		ingredients: SelectedMealIngredient[]
	}

	interface SelectedMealIngredient {
		id: bigint
		name: string
		number: number
		isNew: boolean
	}

	let selectedMeal: SelectedMeal | null = $state(null)
	let selectedMealPristine: SelectedMeal | null = $state(null)

	async function setSelectedMealId(id:bigint | undefined) {
		if (!id) {
			selectedMeal = null
			return
		}
			
		const meal = meals.find(m => m.id == id) || null
		if(!meal) {
			selectedMeal = null
			return
		}

		selectedMeal = toSelectedMeal(meal)
		selectedMealPristine = toSelectedMeal(meal)
	}

	function toSelectedMeal(meal:Meal) : SelectedMeal {
		let smig: SelectedMealIngredient[] = []

		for(const igRef of meal.ingredientRefs) {
			const ig = ingredients.find(ig => ig.id === igRef.ingredientId)
			if(!ig) {
				console.log("failed to resolve ingredient id: ", igRef.ingredientId)
				continue
			}

			smig.push({
				id: ig.id,
				name: ig.name,
				number: igRef.number,
				isNew: false,
			})
		}

		return {
			_meal: meal,
			id: meal.id,
			name: meal.name,
			recipeUrl: meal.recipeUrl,
			ingredientsImageUrl: mapImageSource(meal.ingredientsImage),
			ingredients: smig,
		}
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

	function addNewIngredient() {
		selectedMeal?.ingredients.push({id: BigInt(psuedoId), name:"", number: 1, isNew: true})
		psuedoId++
	}

	function deleteIngredient(id:bigint) {
		if(!selectedMeal) {
			console.log("trying to delete an ingredient, but no selected meal")
			return
		}
		const ing = selectedMeal.ingredients;
		ing.splice(ing.findIndex(ml => ml.id == id), 1)
	}

	async function save() {
		if(!selectedMeal) {
			console.log("tried to save a selected meal when no selected meal is set")
			return
		}

		isEditing = false

		let igRefs: IngredientRef[] = []
		for(const ig of selectedMeal.ingredients) {
			igRefs.push(<IngredientRef>{
				ingredientId: ig.id,
				number: ig.number,
			})
		}

		await client.updateMeal({
			meal: {
				id: selectedMeal.id,
				name: selectedMeal.name,
				recipeUrl: selectedMeal.recipeUrl,
				previewImage: selectedMeal._meal.previewImage,
				ingredientsImage: selectedMeal._meal.ingredientsImage,
				ingredientRefs: igRefs,
			},
		})

		refresh()
	}

	function edit() {
		isEditing = true
	}

	let dirty = $derived.by(() => {
		if(!selectedMeal) {
			return false;
		}
		if(!selectedMealPristine) {
			console.log("selectedMealPristine was null but shouldnt be")
			return false;
		}

		if(selectedMeal.ingredients.length != selectedMealPristine.ingredients.length) {
			return true
		}

		for(let i = 0; i < selectedMeal.ingredients.length; i++) {
			const ig = selectedMeal.ingredients[i]
			const igprist = selectedMealPristine.ingredients[i]
			if(!ingredientsEqual(ig, igprist)) {
				return true
			}
		}
		
		return false
	})

	let valid = $derived.by(() => {
		if(!selectedMeal) {
			return true;
		}

		const newIngredients = selectedMeal.ingredients.filter(ig => ig.isNew)

		for(const ning of newIngredients) {
			if(ning.id == BigInt(0)) {
				return false
			}

			if(ning.number == 0) {
				return false
			}
		}
		return true
	})

	function ingredientsEqual(a: SelectedMealIngredient, b:SelectedMealIngredient): boolean {
		if(a.id != b.id) {
			return false
		}

		if (a.isNew != b.isNew) {
			return false
		}

		if (a.name != b.name) {
			return false
		}

		if (a.number != b.number) {
			return false
		}

		return true
	}

	refresh()

</script>

<svelte:head>
	<title>{selectedMeal?.name}</title>
	<meta name="description" content="Build the recipes" />
</svelte:head>

<div class="overflow-y-auto h-full p-3 -m-3">
	{#if selectedMeal}
	<H1>{selectedMeal.name}</H1>
	<div class="flex flex-col justify-end w-full sm:h-100 h-70 rounded-xl shadow-md mb-5 bg-white bg-contain bg-no-repeat bg-center" style="background-image: url({selectedMeal.ingredientsImageUrl || "/carrot.svg"});">
		<div class="flex flex-row justify-center w-full p-2 bg-white/70">
			{#if selectedMeal.recipeUrl}
			<a class="text-blue-500" href={selectedMeal.recipeUrl} target="_blank">{selectedMeal.name}</a>
			{:else}
			<p>{selectedMeal.name}</p>
			{/if}
		</div>
	</div>
	<Table>
		<TrTitle>
			<Td colspan={3} title={true}>
				<div class="flex flex-row items-center justify-between">
					<p>Ingredients</p>
					{#if isEditing}
						{#if dirty}
						<div class="flex flex-row justify-end">
							<Button disabled={!valid} classes="mr-2" onclick={() => save()}>Save</Button>
						</div>
						{/if}
					{:else}
						<div class="flex flex-row justify-end">
							<Button classes="mr-2" onclick={() => edit()}>Edit</Button>
						</div>
					{/if}
				</div>
			</Td>
		</TrTitle>
		<TrHeader>
			<Td header={true}>Ingredient</Td>
			<Td header={true}>Number</Td>
			{#if isEditing}
			<Td header={true}>Action</Td>
			{/if}
		</TrHeader>
		{#each selectedMeal.ingredients as ing (ing.id)}
		<Tr>
			<Td>
				{#if isEditing}
					<Select bind:value={ing.id}>
					{#each ingredients as ing}
						<option value={ing.id}>{ing.name}</option>
					{/each}
					</Select>
				{:else}
					{ing.name}
				{/if}
			</Td>
			<Td>
				{#if isEditing}
					<input class="w-12 bg-white rounded-md h-7 px-2 border-solid border-gray-500 border-1 focus:border-gray-900 focus:outline-none" bind:value={ing.number} type="number" min="1" max="20">
				{:else}
					{ing.number}
				{/if}
			</Td>
			{#if isEditing}
			<Td>
				<Button onclick={() => deleteIngredient(ing.id)}>Delete</Button>
			</Td>
			{/if}
		</Tr>
		{/each}
		{#if isEditing}
		<Tr><Td></Td><Td></Td><Td><Button onclick={addNewIngredient}>+</Button></Td></Tr>
		{/if}
	</Table>
	{/if}
</div>
