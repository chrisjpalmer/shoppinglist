<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import { ShoppingListService } from "../gen/shopping_list_service_pb";
  import { Category, type Plan } from "../gen/plan_pb";
  import type { Meal } from "../gen/meal_pb";
  import type { Ingredient } from "../gen/ingredient_pb";
  import { CreateShoppingListService } from "$lib/shopping_list_service";
  import Button from "../components/button.svelte";
  import Select from "../components/select.svelte";
  import Table from "../components/table.svelte";
  import Td from "../components/td.svelte";
  import TrHeader from "../components/tr-header.svelte";
  import TrTitle from "../components/tr-title.svelte";
  import Tr from "../components/tr.svelte";

	const categories = [
		Category.LUNCH,
		Category.DINNER,
		Category.SNACK,
	]

	const days = [
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	]

	const client = CreateShoppingListService()

	interface DisplaySummary {
		ingredients: DisplayIngredientCount[]
	}

	interface DisplayIngredientCount {
		name: string
		count: number
	}

	let dirty = $state(false);
	let plan:Plan | null = $state(null);
	let planSummary:DisplaySummary | null = $state(null);
	let allMeals:Meal[] | null = $state(null);

	async function refresh() {
		// plan
		const planRs = await client.getPlan({})
		plan = <Plan> planRs.plan

		// plan summary
		const igRs = await client.getIngredients({})

		const igmap = ingredientMap(igRs.ingredients)

		let ds:DisplaySummary = {
			ingredients: []
		}

		if(planRs.planSummary) {
			for(const igref of planRs.planSummary.ingredientRef) {
				if(igmap.has(igref.ingredientId)) {
					ds.ingredients.push({name: <string> igmap.get(igref.ingredientId), count: igref.number})
				} else {
					console.log(`ingredient id ${igref.ingredientId} was ignored when rendering plan summary`)
				}
			}
		}

		planSummary = ds

		// all meals
		const mealRs = await client.getMeals({})
		allMeals = <Meal[]> mealRs.meals
	}

	function ingredientMap(igs: Ingredient[]): Map<BigInt, string> {
		let igmap = new Map<BigInt, string>();
		for(const ig of igs) {
			igmap.set(ig.id, ig.name)
		}

		return igmap
	}
	
	function getMealId(day: number, category: Category): bigint {
		if(!plan) {
			console.log("no plan but trying to call getMealId")
			return BigInt(0)
		}

		const meals = plan.days[day].categoryMeals
		const meal = meals.find(meal => meal.category == category)
		if(!meal) {
			console.log("could not find the meal for the given day and category")
			return BigInt(0)
		}

		return meal.mealId
	}
	function setMealId(day: number, category: Category, mealId: bigint) {
		if(!plan) {
			console.log("no plan but trying to call setMealId")
			return BigInt(0)
		}

		const meals = plan.days[day].categoryMeals
		const meal = meals.find(meal => meal.category == category)
		if(!meal) {
			console.log("could not find the meal for the given day and category")
			return BigInt(0)
		}

		meal.mealId = mealId;
		dirty = true;
	}

	async function save() {
		if(!plan) {
			console.log("tried to update the plan, but plan is null")
			return
		}
		await client.updatePlan({plan: plan})
		dirty = false;
		refresh()
	}

	refresh()
	
</script>

<svelte:head>
	<title>Planner</title>
</svelte:head>

<div class="overflow-auto p-5 -m-5 w-[calc(100%_+_(var(--spacing)_*_10))]">
	<Table classes="min-w-300">
		<TrTitle><Td title={true}>Planner</Td>{#each days as day}<Td title={true}></Td>{/each}</TrTitle>
		<TrHeader>
			<Td header={true}></Td>
			{#each days as day}
			<Td header={true}>{day}</Td>
			{/each}
		</TrHeader>
		{#each categories as category}
		<Tr classes="h-20">
			<Td classes="font-bold">{Category[category]}</Td>
			{#each days as day, i}
			<Td>
				<Select classes="mx-1" bind:value={
						()=>getMealId(i, category),
						(v) => setMealId(i, category, v)
					}>
					{#if allMeals}
					{#each allMeals as meal (meal.id)}
					<option value={meal.id}>{meal.name}</option>
					{/each}
					{/if}
				</Select>
			</Td>
			{/each}
		</Tr>
		{/each}
	</Table>
</div>

{#if dirty}
<Button onclick={save}>Save</Button>
{/if}

{#if planSummary}
<Table classes="mt-10">
	<TrTitle><Td title={true}>Ingredients List</Td><Td title={true}></Td></TrTitle>
	<TrHeader><Td header={true}>Ingredient</Td><Td header={true}>Amount</Td></TrHeader>
	{#each planSummary.ingredients as ig}
	<Tr><Td>{ig.name}</Td><Td>{ig.count}</Td></Tr>
	{/each}
</Table>
{/if}
