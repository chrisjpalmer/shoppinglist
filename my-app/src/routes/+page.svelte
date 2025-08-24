<script lang="ts">
  import { createClient } from "@connectrpc/connect";
  import { createConnectTransport } from "@connectrpc/connect-web";
  import { ShoppingListService } from "../gen/shopping_list_service_pb";
  import { Category, type Plan } from "../gen/plan_pb";
  import type { Meal } from "../gen/meal_pb";


	const transport = createConnectTransport({
		baseUrl: "http://localhost:8080",
	});

	const client = createClient(ShoppingListService, transport);

	let dirty = $state(false);

	let plan:Plan | null = $state(null);
	async function refresh() {
		const rs = await client.getPlan({})
		plan = <Plan> rs.plan
	}
	
	function getMealId(day: number, category: Category): bigint {
		if(!plan) {
			console.log("no plan but trying to call getMealId")
			return BigInt(0)
		}

		const catMeals = plan.days[day].categoryMeals
		const meal = catMeals.find(meal => meal.category = category)
		if(!meal) {
			console.log("could not find the meal for the given day and category")
			return BigInt(0)
		}

		return meal.mealId
	}
	function setMealId(day: string, category: Category, mealId: bigint) {
		if(!plan) {
			console.log("no plan but trying to call setMealId")
			return BigInt(0)
		}

		const catMeals = plan.days[day].categoryMeals
		if(!catMeals)
		chosenMeals.get(day)?.set(category, meal)
	}

	let categories = [
		Category.DINNER,
		Category.LUNCH,
		Category.SNACK,
	]

	let days = [
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	]

	let meals = new Map<string, string[]>()

	meals.set("Lunch", ["Spaghetti Bolognese", "Not Spaghetti Bolognese"])
	meals.set("Dinner", ["Meat Pie", "Not Meat Pie"])
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
				<td>{category}</td>
				{#each days as day}
				<td>
					<select bind:value={
							()=>getMeal(day, category), 
							(v) => setMeal(day, category, v)
						}>
						{#each allMeals as meal (meal.id)}
						<option value={meal.id}>{meal}</option>
						{/each}
					</select>
				</td>
				{/each}
			</tr>
			{/each}
		</tbody>
	</table>

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

	.welcome {
		display: block;
		position: relative;
		width: 100%;
		height: 0;
		padding: 0 0 calc(100% * 495 / 2048) 0;
	}

	.welcome img {
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		display: block;
	}
</style>
