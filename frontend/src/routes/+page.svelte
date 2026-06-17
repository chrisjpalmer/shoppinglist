<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import { ShoppingListService } from "../gen/shopping_list_service_pb";
  import { Category, type Plan } from "../gen/plan_pb";
  import type { Meal } from "../gen/meal_pb";
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
		Category.BABY,
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

	interface DisplayIngredientCount {
		name: string
		count: number
	}

	let dirty = $state(false);
	let plan:Plan | null = $state(null);
	let allMeals:Meal[] | null = $state(null);
	let isEditing = $state(false)

	async function refresh() {
		// plan
		const planRs = await client.getPlan({})
		plan = <Plan> planRs.plan

		normalizePlan(plan)

		// all meals
		const mealRs = await client.getMeals({})
		allMeals = <Meal[]> mealRs.meals
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
	function getMealName(day: number, category: Category): string {
		const mealId = getMealId(day, category)

		if(mealId == BigInt(0)) {
			return ""
		}

		const meal = allMeals?.find(m => m.id == mealId)

		if(!meal) {
			console.log("could not find the meal id in the list of meals")
			return ""
		}

		return meal.name
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

	function edit() {
		isEditing = true;
	}

	async function save() {
		if(!plan) {
			console.log("tried to update the plan, but plan is null")
			return
		}

		if(!dirty) {
			console.log("not saving the plan since nothing changed")
			isEditing = false;
			return
		}

		await client.updatePlan({plan: plan})
		dirty = false;
		isEditing = false;
		refresh()
	}

	async function reset() {
		if(!plan) {
			console.log("tried to reset the plan, but plan is null")
			return
		}

		for(let d of plan.days) {
			for(let c of d.categoryMeals) {
				c.mealId = 0n
			}
		}

		dirty = true
	}

	function normalizePlan(p:Plan) {
		for (let day of p.days) {
			for(let c = 0; c < categories.length; c++) {

				// fill category
				if (day.categoryMeals.length < (c+1)) {
					console.log(`filled category ${c}`)

					day.categoryMeals.push({
						category: categories[c],
						mealId: 0n,
						$typeName: 'CategoryMeal'
					})
				}
			}
		}
	}

	refresh()
	
</script>

<svelte:head>
	<title>Planner</title>
</svelte:head>

<div class="overflow-auto p-5 -m-5 w-[calc(100%_+_(var(--spacing)_*_10))]">
	<Table classes="min-w-300">
		<TrTitle>
			<Td title={true} colspan={days.length+1}>
				<p>Planner</p>
				{#if isEditing}
					<div class="flex flex-row justify-end">
						<Button classes="mr-2" onclick={() => save()}>Save</Button>
						<Button classes="mr-2" onclick={reset}>Reset</Button>
					</div>
				{:else}
					<div class="flex flex-row justify-end">
						<Button classes="mr-2" onclick={() => edit()}>Edit</Button>
					</div>
				{/if}
			</Td>
		</TrTitle>
		<TrHeader>
			<Td header={true}></Td>
			{#each days as day}
			<Td header={true}>{day}</Td>
			{/each}
		</TrHeader>
		{#each categories as category}
		<Tr classes="h-25">
			<Td classes="font-bold">{Category[category]}</Td>
			{#each days as day, i}
			<Td>
				{#if isEditing}
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
				{:else}
				<a class="text-blue-500" href="/recipies/{getMealId(i, category)}">{getMealName(i, category)}</a>
				{/if}
			</Td>
			{/each}
		</Tr>
		{/each}
	</Table>
</div>