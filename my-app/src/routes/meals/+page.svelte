<script lang="ts">
  	import type { Meal } from '../../gen/meal_pb';
  	import { PUBLIC_BACKEND_URL } from '$env/static/public';
  	import { CreateShoppingListService } from '$lib/shopping_list_service';

	const client = CreateShoppingListService(PUBLIC_BACKEND_URL)

	let psuedoIdCounter = $state(0)

	let displayMeals: DisplayMeal[] = $state([])
	let displayNewMeals: DisplayNewMeal[] = $state([])
	interface DisplayMeal {
		_meal: Meal
		id: bigint
		name: string
		isEdit: boolean
		recipeUrl: string
	}

	interface DisplayNewMeal {
		pseudoId: number
		name: string
		recipeUrl: string
	}

	async function refresh() {
		const rs = await client.getMeals({})
		displayMeals = rs.meals.map(m => ({id: m.id, name: m.name, _meal:m, isEdit: false, recipeUrl: m.recipeUrl}))
		displayNewMeals = []
	}

	function editMeal(id: bigint) {
		const m = displayMeals.find(dm => dm.id == id)
		if (!m) {
			console.log("edit clicked for a meal whose id couldn't be found")
			return
		}
		m.isEdit = true
	}

	async function deleteMeal(id: bigint) {
		await client.deleteMeal({mealId: id})
		refresh()
	}

	async function saveMeal(id: bigint) {
		const m = displayMeals.find(dm => dm.id == id)
		if (!m) {
			console.log("save clicked for a meal whose id couldn't be found")
			return;
		}

		let _meal = m._meal
		_meal.name = m.name
		_meal.recipeUrl = m.recipeUrl

		await client.updateMeal({
			meal: _meal,
		})

		refresh()
	}

	async function saveNewMeal(pid: number) {
		const m = displayNewMeals.find(dm => dm.pseudoId == pid)
		if (!m) {
			console.log("save new meal clicked for a meal whose psuedo id couldn't be found")
			return;
		}

		await client.createMeal({
			meal: {
				name: m.name,
				recipeUrl: m.recipeUrl,
				ingredientRefs: [],
			},
		})

		refresh()
	}

	function addMeal() {
		displayNewMeals.push({pseudoId: psuedoIdCounter, name: "", recipeUrl: ""})
		psuedoIdCounter++
	}

	refresh()
</script>
<svelte:head>
	<title>Meals</title>
</svelte:head>

<div class="text-column">
	<h1>Meals</h1>

	<table>
		<thead>
			<tr><td>Name</td><td>Recipe Url</td><td>Action</td></tr>
		</thead>
		<tbody>
			{#each displayMeals as dm (dm.id)}
				{#if dm.isEdit}
					<tr>
						<td>
							<input type="text" bind:value={dm.name}>
						</td>
						<td>
							<input type="text" bind:value={dm.recipeUrl}>
						</td>
						<td>
							<button onclick={() => saveMeal(dm.id)}>Save</button>
						</td>
					</tr>
				{:else}
					<tr>
						<td>{dm.name}</td>
						<td>{#if dm.recipeUrl != ''}<a href={dm.recipeUrl}>Link</a>{/if}</td>
						<td>
							<button onclick={() => editMeal(dm.id)}>Edit</button><button onclick={() => deleteMeal(dm.id)}>Delete</button>
						</td>
					</tr>
				{/if}
			{/each}
			{#each displayNewMeals as dm (psuedoIdCounter)}
				<tr>
					<td>
						<input type="text" bind:value={dm.name}>
					</td>
					<td>
						<input type="text" bind:value={dm.recipeUrl}>
					</td>
					<td>
						<button onclick={() => saveNewMeal(dm.pseudoId)}>Save</button>
					</td>
				</tr>
			{/each}
			<tr><td></td><td></td><td><button onclick={addMeal}>+</button></td></tr>
		</tbody>
	</table>
</div>
