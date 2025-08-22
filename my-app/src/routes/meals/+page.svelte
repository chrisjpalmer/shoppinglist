<script lang="ts">
	import { ShoppingListService } from '../../gen/shopping_list_service_pb';
	import { createClient } from "@connectrpc/connect";
	import { createConnectTransport } from "@connectrpc/connect-web";
  	import type { Meal } from '../../gen/meal_pb';

	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createClient(ShoppingListService, transport);

	let psuedoIdCounter = $state(0)

	let displayMeals: DisplayMeal[] = $state([])
	let displayNewMeals: DisplayNewMeal[] = $state([])
	interface DisplayMeal {
		_meal: Meal
		id: bigint
		name: string
		isEdit: boolean
	}

	interface DisplayNewMeal {
		pseudoId: number
		name: string
	}

	async function refresh() {
		const rs = await client.getMeals({})
		displayMeals = rs.meals.map(m => ({id: m.id, name: m.name, _meal:m, isEdit: false}))
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
				ingredientRefs: [],
			},
		})

		refresh()
	}

	function addMeal() {
		displayNewMeals.push({pseudoId: psuedoIdCounter, name: ""})
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
			<tr><td>Name</td><td>Action</td></tr>
		</thead>
		<tbody>
			{#each displayMeals as dm (dm.id)}
				{#if dm.isEdit}
					<tr><td><input type="text" bind:value={dm.name}></td><td><button onclick={() => saveMeal(dm.id)}>Save</button></td></tr>
				{:else}
					<tr><td>{dm.name}</td><td><button onclick={() => editMeal(dm.id)}>Edit</button><button onclick={() => deleteMeal(dm.id)}>Delete</button></td></tr>
				{/if}
			{/each}
			{#each displayNewMeals as dm (psuedoIdCounter)}
				<tr><td><input type="text" bind:value={dm.name}></td><td><button onclick={() => saveNewMeal(dm.pseudoId)}>Save</button></td></tr>
			{/each}
			<tr><td></td><td><button onclick={addMeal}>+</button></td></tr>
		</tbody>
	</table>
</div>
