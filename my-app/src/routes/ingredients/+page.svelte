<script lang="ts">
  	import { CreateShoppingListService } from '$lib/shopping_list_service';

	const client = CreateShoppingListService()

	let psuedoIdCounter = $state(0)

	let displayIngredients: DisplayIngredient[] = $state([])
	let displayNewIngredients: DisplayNewIngredient[] = $state([])
	interface DisplayIngredient {
		id: bigint
		name: string
		isEdit: boolean
	}

	interface DisplayNewIngredient {
		pseudoId: number
		name: string
	}

	async function refresh() {
		const rs = await client.getIngredients({})
		displayIngredients = rs.ingredients.map(ig => ({id: ig.id, name: ig.name, isEdit: false}))
		displayNewIngredients = []
	}

	function editIngredient(id: bigint) {
		const m = displayIngredients.find(dm => dm.id == id)
		if (!m) {
			console.log("edit clicked for a meal whose id couldn't be found")
			return
		}
		m.isEdit = true
	}

	async function deleteIngredient(id: bigint) {
		await client.deleteIngredient({ingredientId: id})
		refresh()
	}

	async function saveIngredient(id: bigint) {
		const m = displayIngredients.find(dm => dm.id == id)
		if (!m) {
			console.log("save clicked for a meal whose id couldn't be found")
			return;
		}

		await client.updateIngredient({
			ingredient: {
				id: id,
				name: m.name,
			},
		})

		refresh()
	}

	async function saveNewIngredient(pid: number) {
		const m = displayNewIngredients.find(dm => dm.pseudoId == pid)
		if (!m) {
			console.log("save new meal clicked for a meal whose psuedo id couldn't be found")
			return;
		}

		await client.createIngredient({
			ingredient: {
				name: m.name,
			},
		})

		refresh()
	}

	function addIngredient() {
		displayNewIngredients.push({pseudoId: psuedoIdCounter, name: ""})
		psuedoIdCounter++
	}

	refresh()
</script>
<svelte:head>
	<title>Ingredients</title>
</svelte:head>

<div class="text-column">
	<h1>Ingredients</h1>

	<table>
		<thead>
			<tr><td>Name</td><td>Action</td></tr>
		</thead>
		<tbody>
			{#each displayIngredients as dig (dig.id)}
				{#if dig.isEdit}
					<tr><td><input type="text" bind:value={dig.name}></td><td><button onclick={() => saveIngredient(dig.id)}>Save</button></td></tr>
				{:else}
					<tr><td>{dig.name}</td><td><button onclick={() => editIngredient(dig.id)}>Edit</button><button onclick={() => deleteIngredient(dig.id)}>Delete</button></td></tr>
				{/if}
			{/each}
			{#each displayNewIngredients as dig (psuedoIdCounter)}
				<tr><td><input type="text" bind:value={dig.name}></td><td><button onclick={() => saveNewIngredient(dig.pseudoId)}>Save</button></td></tr>
			{/each}
			<tr><td></td><td><button onclick={addIngredient}>+</button></td></tr>
		</tbody>
	</table>
</div>
