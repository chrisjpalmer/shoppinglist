<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import { ShoppingListService } from "../gen/shopping_list_service_pb";
  import { Category, type Plan } from "../gen/plan_pb";
  import type { Meal } from "../gen/meal_pb";

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

	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createClient(ShoppingListService, transport);

	let dirty = $state(false);

	let plan:Plan | null = $state(null);
	let allMeals:Meal[] | null = $state(null);

	async function refresh() {
		const planRs = await client.getPlan({})
		plan = <Plan> planRs.plan

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

<section>
	<table>
		<thead>
			<tr><td>Ingredient</td><td>Amount</td></tr>
		</thead>
		<tbody>
			<tr><td>Onion</td><td>2</td></tr>
		</tbody>
	</table>
</section>

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
