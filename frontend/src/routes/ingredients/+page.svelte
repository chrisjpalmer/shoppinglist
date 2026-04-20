<script lang="ts">
  	import { CreateShoppingListService } from '$lib/shopping_list_service';
	import Button from '../../components/button.svelte';
  import H1 from '../../components/h1.svelte';
  import TextInput from '../../components/text-input.svelte';

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

<div class="px-8">
	<H1>Ingredients</H1>

	<table class="w-full table-fixed">
		<thead>
			<tr class="font-bold"><td>Name</td><td>Action</td></tr>
		</thead>
		<tbody>
			{#each displayIngredients as dig (dig.id)}
				{#if dig.isEdit}
					<tr class="h-10"><td><TextInput bind:value={dig.name}/></td><td><Button onclick={() => saveIngredient(dig.id)}>Save</Button></td></tr>
				{:else}
					<tr class="h-10"><td>{dig.name}</td><td><Button onclick={() => editIngredient(dig.id)}>Edit</Button><Button classes="ml-1" onclick={() => deleteIngredient(dig.id)}>Delete</Button></td></tr>
				{/if}
			{/each}
			{#each displayNewIngredients as dig (psuedoIdCounter)}
				<tr class="h-10"><td><TextInput bind:value={dig.name}/></td><td><Button onclick={() => saveNewIngredient(dig.pseudoId)}>Save</Button></td></tr>
			{/each}
			<tr class="h-10"><td></td><td><Button onclick={addIngredient}>+</Button></td></tr>
		</tbody>
	</table>
</div>
