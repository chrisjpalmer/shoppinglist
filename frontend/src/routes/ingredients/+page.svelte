<script lang="ts">
  	import { CreateShoppingListService } from '$lib/shopping_list_service';
	import Button from '../../components/button.svelte';
	import Table from '../../components/table.svelte';
	import Td from '../../components/td.svelte';
	import TextInput from '../../components/text-input.svelte';
	import TrHeader from '../../components/tr-header.svelte';
	import TrTitle from '../../components/tr-title.svelte';
	import Tr from '../../components/tr.svelte';

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

<Table>
	<TrTitle><Td title={true}>Ingredients</Td><Td title={true}></Td></TrTitle>
	<TrHeader><Td header={true}>Name</Td><Td header={true}>Action</Td></TrHeader>
	{#each displayIngredients as dig (dig.id)}
		{#if dig.isEdit}
			<Tr><Td><TextInput bind:value={dig.name}/></Td><Td><Button onclick={() => saveIngredient(dig.id)}>Save</Button></Td></Tr>
		{:else}
			<Tr><Td>{dig.name}</Td><Td><Button onclick={() => editIngredient(dig.id)}>Edit</Button><Button classes="ml-1" onclick={() => deleteIngredient(dig.id)}>Delete</Button></Td></Tr>
		{/if}
	{/each}
	{#each displayNewIngredients as dig (psuedoIdCounter)}
		<Tr><Td><TextInput bind:value={dig.name}/></Td><Td><Button onclick={() => saveNewIngredient(dig.pseudoId)}>Save</Button></Td></Tr>
	{/each}
	<Tr><Td></Td><Td><Button onclick={addIngredient}>+</Button></Td></Tr>
</Table>
