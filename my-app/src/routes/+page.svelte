<script lang="ts">

	function getKeys(m:Map<any, any>) {
		let keys = []
		for (let k of m.keys()) {
			keys.push(k)
		}
		return keys
	}

	let chosenMeals = new Map<string, Map<string, string>>() // day, category, meal
	let mondayLunch = new Map
	chosenMeals.set("Monday", new Map())

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
				<td>Monday</td>
				<td>Tuesday</td>
				<td>Wednesday</td>
				<td>Thursday</td>
				<td>Friday</td>
				<td>Saturday</td>
				<td>Sunday</td>
			</tr>
		</thead>
		<tbody>
			{#each categories as category}
			<tr>
				<td>{category}</td>
				{#each days as day}
				<td>
					<select>
						{#each meals.get(category) as meal}
						<option>{meal}</option>
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
