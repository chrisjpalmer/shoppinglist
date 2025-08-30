<script lang="ts">
  	import type { IngredientRef, Meal } from '../../gen/meal_pb';
  	import type { Ingredient } from '../../gen/ingredient_pb';
  	import { CreateShoppingListService } from '$lib/shopping_list_service';

	const client = CreateShoppingListService()

	let psuedoId = 10000

	let meals: Meal[] = $state([])
	let ingredients: Ingredient[] = $state([])

	async function refresh() {
		const gmRs = await client.getMeals({})
		meals = gmRs.meals
		
		const giRs = await client.getIngredients({});
		ingredients = giRs.ingredients

		if(selectedMeal) {
			await setSelectedMealId(selectedMeal.id)
		} else {
			if(meals.length > 0) {
				await setSelectedMealId(meals[0].id)
			}
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
	let selectedMealPristine: SelectedMeal | null = null;

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
			id: meal.id,
			name: meal.name,
			ingredients: smig,
		}
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
	}

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

	function valid(): boolean {
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
	}

	refresh()

</script>

<svelte:head>
	<title>Recipies</title>
	<meta name="description" content="Build the recipes" />
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
			<tr><td>Ingredient</td><td>Number</td><td>Action</td></tr>
		</thead>
		<tbody>
			{#each selectedMeal.ingredients as ing (ing.id)}
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
				<td>
					<button onclick={() => deleteIngredient(ing.id)}>Delete</button>
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
