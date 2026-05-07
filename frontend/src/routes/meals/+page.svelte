<script lang="ts">
  	import type { Meal } from '../../gen/meal_pb';
  	import { CreateShoppingListService } from '$lib/shopping_list_service';
	import TextInput from '../../components/text-input.svelte';
	import Button from '../../components/button.svelte';
	import TrHeader from '../../components/tr-header.svelte';
	import Td from '../../components/td.svelte';
	import Tr from '../../components/tr.svelte';
	import TrTitle from '../../components/tr-title.svelte';
	import Table from '../../components/table.svelte';

	const client = CreateShoppingListService()

	let psuedoIdCounter = $state(0)

	let displayMeals: DisplayMeal[] = $state([])
	let displayNewMeals: DisplayNewMeal[] = $state([])
	interface DisplayMeal {
		_meal: Meal
		id: bigint
		name: string
		isEdit: boolean
		recipeUrl: string
		previewImageUrl: string
		ingredientsImageUrl: string
		hasIngredients: boolean
	}

	interface DisplayNewMeal {
		pseudoId: number
		name: string
		recipeUrl: string
		previewImageUrl: string
		ingredientsImageUrl: string
	}

	async function refresh() {
		const rs = await client.getMeals({})
		displayMeals = rs.meals.map(m => ({
			id: m.id, 
			name: m.name, 
			_meal:m, 
			isEdit: false, 
			recipeUrl: m.recipeUrl,
			previewImageUrl: m.previewImageUrl,
			ingredientsImageUrl: m.ingredientsImageUrl,
			hasIngredients: m.ingredientRefs && m.ingredientRefs.length > 0
		}))
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
		_meal.previewImageUrl = m.previewImageUrl
		_meal.ingredientsImageUrl = m.ingredientsImageUrl

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
				previewImageUrl: m.previewImageUrl,
				ingredientsImageUrl: m.ingredientsImageUrl,
				ingredientRefs: [],
			},
		})

		refresh()
	}

	function addMeal() {
		displayNewMeals.push({pseudoId: psuedoIdCounter, name: "", recipeUrl: "", previewImageUrl: "", ingredientsImageUrl: ""})
		psuedoIdCounter++
	}

	refresh()
</script>
<svelte:head>
	<title>Meals</title>
</svelte:head>


<Table>
	<TrTitle>
		<Td title={true}>Meals</Td>
		<Td title={true}></Td>
		<Td title={true}></Td>
		<Td title={true}></Td>
		<Td title={true}></Td>
	</TrTitle>
	<TrHeader>
		<Td header={true}>Name</Td>
		<Td header={true}>Recipe Url</Td>
		<Td header={true}>Preview Image Url</Td>
		<Td header={true}>Ingredients Image Url</Td>
		<Td header={true}>Action</Td>
	</TrHeader>
	{#each displayMeals as dm (dm.id)}
		{#if dm.isEdit}
			<Tr>
				<Td>
					<TextInput bind:value={dm.name}/>
				</Td>
				<Td>
					<TextInput bind:value={dm.recipeUrl}/>
				</Td>
				<Td>
					<TextInput bind:value={dm.previewImageUrl}/>
				</Td>
				<Td>
					<TextInput bind:value={dm.ingredientsImageUrl}/>
				</Td>
				<Td>
					<Button onclick={() => saveMeal(dm.id)}>Save</Button>
				</Td>
			</Tr>
		{:else}
			<Tr>
				<Td classes={!dm.hasIngredients ? 'text-red-500' : ''}>{dm.name}</Td>
				<Td>{#if dm.recipeUrl != ''}<a href={dm.recipeUrl}>Link</a>{/if}</Td>
				<Td>{#if dm.previewImageUrl != ''}<a href={dm.previewImageUrl}>Link</a>{/if}</Td>
				<Td>{#if dm.ingredientsImageUrl != ''}<a href={dm.ingredientsImageUrl}>Link</a>{/if}</Td>
				<Td>
					<Button onclick={() => editMeal(dm.id)}>Edit</Button><Button classes="ml-1" onclick={() => deleteMeal(dm.id)}>Delete</Button>
				</Td>
			</Tr>
		{/if}
	{/each}
	{#each displayNewMeals as dm (psuedoIdCounter)}
		<Tr>
			<Td>
				<TextInput bind:value={dm.name}/>
			</Td>
			<Td>
				<TextInput bind:value={dm.recipeUrl}/>
			</Td>
			<Td>
				<TextInput bind:value={dm.previewImageUrl}/>
			</Td>
			<Td>
				<TextInput bind:value={dm.ingredientsImageUrl}/>
			</Td>
			<Td>
				<Button onclick={() => saveNewMeal(dm.pseudoId)}>Save</Button>
			</Td>
		</Tr>
	{/each}
	<Tr>
		<Td></Td>
		<Td></Td>
		<Td></Td>
		<Td></Td>
		<Td><Button onclick={addMeal}>+</Button></Td>
	</Tr>
</Table>
