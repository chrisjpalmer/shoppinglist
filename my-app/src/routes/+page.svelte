<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import { ShoppingListService } from "../gen/shopping_list_service_pb";
  import { Category, type Plan } from "../gen/plan_pb";
  import type { Meal } from "../gen/meal_pb";
  import type { Ingredient } from "../gen/ingredient_pb";
  import { PUBLIC_BACKEND_URL } from "$env/static/public";
  import { CreateShoppingListService } from "$lib/shopping_list_service";

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

	const client = CreateShoppingListService(PUBLIC_BACKEND_URL)

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
		refresh()
	}

	refresh()
	
</script>

<svelte:head>
	<title>Planner</title>
</svelte:head>

<section>
	<h1>Planner</h1>
	<table>
		<thead>
			<tr>
				<td></td>
				{#each days as day}
				<td>{day}</td>
				{/each}
			</tr>
		</thead>
		<tbody>
			{#each categories as category}
			<tr>
				<td>{Category[category]}</td>
				{#each days as day, i}
				<td>
					<select bind:value={
							()=>getMealId(i, category), 
							(v) => setMealId(i, category, v)
						}>
						{#if allMeals}
						{#each allMeals as meal (meal.id)}
						<option value={meal.id}>{meal.name}</option>
						{/each}
						{/if}
					</select>
				</td>
				{/each}
			</tr>
			{/each}
		</tbody>
	</table>

	{#if dirty}
	<button onclick={save}>Save</button>
	{/if}
</section>

{#if planSummary}
<section>
	<table>
		<thead>
			<tr><td>Ingredient</td><td>Amount</td></tr>
		</thead>
		<tbody>
			{#each planSummary.ingredients as ig}
			<tr><td>{ig.name}</td><td>{ig.count}</td></tr>
			{/each}
		</tbody>
	</table>
</section>
{/if}

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}

	h1 {
		width: 100%;
	}

	
</style>
