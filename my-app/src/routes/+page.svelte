<script lang="ts">

	function getKeys(m:Map<any, any>) {
		let keys = []
		for (let k of m.keys()) {
			keys.push(k)
		}
		return keys
	}

	let chosenMeals = $state(new Map<string, Map<string, string>>()) // day, category, meal
	
	function getMeal(day: string, category: string): string {
		return <string> chosenMeals.get(day)?.get(category)
	}
	function setMeal(day: string, category: string, meal: string) {
		chosenMeals.get(day)?.set(category, meal)
	}
	
	let monday = new Map();
	monday.set("Lunch", "Not Spaghetti Bolognese")
	chosenMeals.set("Monday", monday)

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

	let categories = $derived(getKeys(meals))
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
				{#if meals.has(category)}
				{@const categorymeals = <string[]> meals.get(category)}
				<td>{category}</td>
				{#each days as day}
				<td>
					<select bind:value={
							()=>getMeal(day, category), 
							(v) => setMeal(day, category, v)
						}>
						{#each categorymeals as meal}
						<option>{meal}</option>
						{/each}
					</select>
				</td>
				{/each}
				{/if}
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
