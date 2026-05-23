<script lang="ts">
  	import { CreateShoppingListService } from '$lib/shopping_list_service';
	import Button from '../../components/button.svelte';
	import ScrollableTable from '../../components/scrollableTable.svelte';
	import Select from '../../components/select.svelte';
	import StickyTrHeader from '../../components/sticky-tr-header.svelte';
	import StickyTrTitle from '../../components/sticky-tr-title.svelte';

	import Td from '../../components/td.svelte';
	import TextInput from '../../components/text-input.svelte';

	import Tr from '../../components/tr.svelte';

	const client = CreateShoppingListService()

	let psuedoIdCounter = $state(0)

	let displayIngredients: DisplayIngredient[] = $state([])
	let displayNewIngredients: DisplayNewIngredient[] = $state([])
	let displayCategories: { id: bigint; name: string }[] = $state([])

	interface DisplayIngredient {
		id: bigint
		name: string
		ingredientCategoryId: bigint
		isEdit: boolean
	}

	interface DisplayNewIngredient {
		pseudoId: number
		name: string
		ingredientCategoryId: bigint
	}

	async function refresh() {
		const igRes = await client.getIngredients({})

		const catRes = await client.getIngredientCategories({})

		displayIngredients = igRes.ingredients.map(ig => ({
			id: ig.id,
			name: ig.name,
			ingredientCategoryId: ig.ingredientCategoryId,
			isEdit: false
		}))

		displayCategories = catRes.ingredientCategories.map(c => ({ id: c.id, name: c.name }))

		displayNewIngredients = []
	}

	function findCategoryName(id: bigint): string {
		return displayCategories.find(c => c.id === id)?.name ?? '--'
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
				ingredientCategoryId: m.ingredientCategoryId,
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
				ingredientCategoryId: m.ingredientCategoryId,
			},
		})

		refresh()
	}

	function addIngredient() {
		displayNewIngredients.push({pseudoId: psuedoIdCounter, name: "", ingredientCategoryId: 0n})
		psuedoIdCounter++
	}

	refresh()
</script>
<svelte:head>
	<title>Ingredients</title>
</svelte:head>

<ScrollableTable classes="overflow-y-auto h-full">
	<StickyTrTitle><Td title={true} colspan={3}>Ingredients</Td></StickyTrTitle>
	<StickyTrHeader><Td header={true}>Name</Td><Td header={true}>Ingredient Category</Td><Td header={true}>Action</Td></StickyTrHeader>
	{#each displayIngredients as dig (dig.id)}
		{#if dig.isEdit}
			<Tr>
				<Td><TextInput bind:value={dig.name}/></Td>
				<Td>
					<Select bind:value={dig.ingredientCategoryId}>
						<option value={0n}>--</option>
						{#each displayCategories as cat}
							<option value={cat.id}>{cat.name}</option>
						{/each}
					</Select>
				</Td>
				<Td><Button onclick={() => saveIngredient(dig.id)}>Save</Button></Td>
			</Tr>
		{:else}
			<Tr>
				<Td>{dig.name}</Td>
				<Td>{findCategoryName(dig.ingredientCategoryId)}</Td>
				<Td><Button onclick={() => editIngredient(dig.id)}>Edit</Button><Button classes="ml-1" onclick={() => deleteIngredient(dig.id)}>Delete</Button></Td>
			</Tr>
		{/if}
	{/each}
	{#each displayNewIngredients as dig (psuedoIdCounter)}
		<Tr>
			<Td><TextInput bind:value={dig.name}/></Td>
			<Td>
				<Select bind:value={dig.ingredientCategoryId}>
					<option value={0n}>--</option>
					{#each displayCategories as cat}
						<option value={cat.id}>{cat.name}</option>
					{/each}
				</Select>
			</Td>
			<Td><Button onclick={() => saveNewIngredient(dig.pseudoId)}>Save</Button></Td>
		</Tr>
	{/each}
	<Tr><Td></Td><Td></Td><Td><Button onclick={addIngredient}>+</Button></Td></Tr>
</ScrollableTable>
