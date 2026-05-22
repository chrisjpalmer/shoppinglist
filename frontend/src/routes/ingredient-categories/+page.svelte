<script lang="ts">
  	import { CreateShoppingListService } from '$lib/shopping_list_service';
	import Button from '../../components/button.svelte';
	import ScrollableTable from '../../components/scrollableTable.svelte';
	import StickyTrHeader from '../../components/sticky-tr-header.svelte';
	import StickyTrTitle from '../../components/sticky-tr-title.svelte';

	import Td from '../../components/td.svelte';
	import TextInput from '../../components/text-input.svelte';

	import Tr from '../../components/tr.svelte';

	import upArrow from '../../lib/images/up-arrow.svg';
	import downArrow from '../../lib/images/down-arrow.svg';

	const client = CreateShoppingListService()

	let psuedoIdCounter = $state(0)

	let displayCategories: DisplayIngredientCategory[] = $state([])
	let displayNewCategories: DisplayNewIngredientCategory[] = $state([])

	interface DisplayIngredientCategory {
		id: bigint
		name: string
		isEdit: boolean
		prevId?: bigint
		nextId?: bigint
	}

	interface DisplayNewIngredientCategory {
		pseudoId: number
		name: string
	}

	async function refresh() {
		const rs = await client.getIngredientCategories({})

		displayCategories = []
		for(let i = 0; i < rs.ingredientCategories.length; i++) {
			const c = rs.ingredientCategories[i]

			let cat:DisplayIngredientCategory = {
				id: c.id, 
				name: c.name,
				isEdit: false,
			}

			if (i > 0) {
				cat.prevId = rs.ingredientCategories[i-1].id
			}

			if (i != (rs.ingredientCategories.length - 1)) {
				cat.nextId = rs.ingredientCategories[i+1].id
			}

			displayCategories.push(cat)
		}
		displayNewCategories = []
	}

	function editCategory(id: bigint) {
		const c = displayCategories.find(dc => dc.id == id)
		if (!c) {
			console.log("edit clicked for a category whose id couldn't be found")
			return
		}
		c.isEdit = true
	}

	async function deleteCategory(id: bigint) {
		await client.deleteIngredientCategory({ingredientCategoryId: id})
		refresh()
	}

	async function saveCategory(id: bigint) {
		const c = displayCategories.find(dc => dc.id == id)
		if (!c) {
			console.log("save clicked for a category whose id couldn't be found")
			return;
		}

		await client.updateIngredientCategory({
			ingredientCategory: {
				id: id,
				name: c.name,
			},
		})

		refresh()
	}

	async function saveNewCategory(pid: number) {
		const c = displayNewCategories.find(dc => dc.pseudoId == pid)
		if (!c) {
			console.log("save new category clicked for a category whose pseudo id couldn't be found")
			return;
		}

		await client.createIngredientCategory({
			ingredientCategory: {
				name: c.name,
			},
		})

		refresh()
	}

	function addCategory() {
		displayNewCategories.push({pseudoId: psuedoIdCounter, name: ""})
		psuedoIdCounter++
	}

	async function swapCategories(aId: bigint, bId: bigint) {
		await client.swapIngredientCategories({
			ingredientCategoryA: aId,
			ingredientCategoryB: bId,
		})

		refresh()
	}

	refresh()
</script>
<svelte:head>
	<title>Ingredient Categories</title>
</svelte:head>

<ScrollableTable classes="overflow-y-auto h-full">
	<StickyTrTitle><Td title={true} colspan={2}>Ingredient Categories</Td></StickyTrTitle>
	<StickyTrHeader><Td header={true}>Name</Td><Td header={true}>Action</Td></StickyTrHeader>
	{#each displayCategories as dc (dc.id)}
		{#if dc.isEdit}
			<Tr>
				<Td>
					<TextInput bind:value={dc.name}/>
				</Td>
				<Td>
					<Button onclick={() => saveCategory(dc.id)}>Save</Button>
				</Td>
			</Tr>
		{:else}
			<Tr>
				<Td>{dc.name}</Td>
				<Td>
					<div class="flex flex-row justify-start items-center">
						<Button onclick={() => editCategory(dc.id)}>Edit</Button>
						<Button classes="ml-1" onclick={() => deleteCategory(dc.id)}>Delete</Button>
						{#if dc.prevId}
						<Button classes="ml-1" onclick={() => swapCategories(dc.id, <bigint> dc.prevId)} image={upArrow}></Button>
						{/if}
						{#if dc.nextId}
						<Button classes="ml-1" onclick={() => swapCategories(dc.id, <bigint> dc.nextId)} image={downArrow}></Button>
						{/if}
					</div>
				</Td>
			</Tr>
		{/if}
	{/each}
	{#each displayNewCategories as dc (psuedoIdCounter)}
		<Tr>
			<Td>
				<TextInput bind:value={dc.name}/></Td>
			<Td>
				<Button onclick={() => saveNewCategory(dc.pseudoId)}>Save</Button>
			</Td>
		</Tr>
	{/each}
	<Tr>
		<Td></Td>
		<Td>
			<Button onclick={addCategory}>+</Button>
		</Td>
	</Tr>
</ScrollableTable>
