<script lang="ts">
	import { ShoppingListService } from '../../gen/shopping_list_service_pb';
	import { createClient } from "@connectrpc/connect";
	import { createConnectTransport } from "@connectrpc/connect-web";
  	import type { IngredientRef, Meal } from '../../gen/meal_pb';
  import type { Ingredient } from '../../gen/ingredient_pb';

	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createClient(ShoppingListService, transport);

	let meals: Meal[] = $state([])
	let ingredients: Ingredient[] = $state([])

	async function refresh() {
		const gmRs = await client.getMeals({})
		meals = gmRs.meals
		
		const giRs = await client.getIngredients({});
		ingredients = giRs.ingredients

		if(selectedMeal) {
			await setSelectedMealId(selectedMeal.id)
		}
	}
	interface SelectedMeal {
		id: bigint
		name: string
		ingredients: SelectedMealIngredient[]
	}

	interface SelectedMealIngredient {
		id: bigint
		name: string
		number: number
		isNew: boolean
	}

	let selectedMeal: SelectedMeal | null = $state(null)

	function getSelectedMealId() {
		if(selectedMeal) {
			return selectedMeal.id
		}
	}

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

		selectedMeal = {
			id: id,
			name: meal.name,
			ingredients: smig,
		}
	}

	function addNewIngredient() {
		selectedMeal?.ingredients.push({id: BigInt(0), name:"", number: 1, isNew: true})
	}

	async function saveSelectedMeal() {
		if(!selectedMeal) {
			console.log("tried to save a selected meal when no selected meal is set")
			return
		}

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
				ingredientRefs: igRefs,
			},
		})

		refresh()
	}

	function dirty(): boolean {
		if(!selectedMeal) {
			return false;
		}
		return selectedMeal.newIngredients.length > 0
	}

	function valid(): boolean {
		if(!selectedMeal) {
			return true;
		}

		for(const ning of selectedMeal.newIngredients) {
			if(ning.id == BigInt(0)) {
				return false
			}

			if(ning.number == 0) {
				return false
			}
		}
		return true
	}

	refresh()

</script>

<svelte:head>
	<title>Recipies</title>
	<meta name="description" content="About this app" />
</svelte:head>

<div class="text-column">
	<h1>Recipies</h1>

	<select bind:value={()=>getSelectedMealId(), (v) => setSelectedMealId(v)}>
		{#each meals as meal (meal.id)}
			<option value={meal.id}>{meal.name}</option>
		{/each}
	</select>

	{#if selectedMeal}
	<table>
		<thead>
			<tr><td>Ingredient</td><td>Number</td></tr>
		</thead>
		<tbody>
			{#each selectedMeal.ingredients as ing}
			<tr>
				<td>
					<select bind:value={ing.id}>
					{#each ingredients as ing}
						<option value={ing.id}>{ing.name}</option>
					{/each}
					</select>
				</td>
				<td>
					<input bind:value={ing.number} type="number" min="1" max="20">
				</td>
			</tr>
			{/each}

			<tr><td></td><td><button onclick={addNewIngredient}>+</button></td></tr>
		</tbody>
	</table>

	{#if dirty()}
	<button onclick={saveSelectedMeal} disabled={!valid()}>Save</button>
	{/if}
	{/if}

</div>
