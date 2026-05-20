<script lang="ts">
  	import { type Meal, type ImageMeta, ImageMode } from '../../gen/meal_pb';
  	import { BackendUrl, CreateShoppingListService } from '$lib/shopping_list_service';
	import TextInput from '../../components/text-input.svelte';
	import Button from '../../components/button.svelte';
	import Td from '../../components/td.svelte';
	import Tr from '../../components/tr.svelte';
	import ScrollableTable from '../../components/scrollableTable.svelte';
	import StickyTrTitle from '../../components/sticky-tr-title.svelte';
	import StickyTrHeader from '../../components/sticky-tr-header.svelte';

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
		previewImage: DPImageMeta 
		ingredientsImage: DPImageMeta
		hasIngredients: boolean
	}

	interface DPImageMeta {
		id: string // dm.id + "preview/ingredients"
		mode: "none" | "external" | "internal"
		internalUrl?: string
		externalUrl?: string
		fileUpload?: HTMLInputElement
		fileBytes?: Uint8Array
	}

	function emptyImageMeta(id: string): DPImageMeta {
		return {
			id, 
			mode: "none",
		}
	}

	interface DisplayNewMeal {
		pseudoId: number
		name: string
		recipeUrl: string
		previewImage: DPImageMeta 
		ingredientsImage: DPImageMeta
	}

	function mapImageMeta(id: bigint, imgType: "preview" | "ingredients", img:ImageMeta | undefined): DPImageMeta {
		let mid = imgid(id, imgType)

		if (!img || img.mode == ImageMode.IM_NONE) {
			return emptyImageMeta(mid)
		}

		if (img.mode == ImageMode.IM_INTERNAL) {
			return {
				id: mid,
				mode: "internal",
				internalUrl: BackendUrl() + img.internalUrl,
			}
		}

		return {
			id: mid,
			mode: "external",
			externalUrl: img.externalUrl,
		}
	}

	function mapDPImageMeta(dpImg: DPImageMeta): ImageMeta {
		let mode: ImageMode = ImageMode.IM_NONE
		
		if (dpImg.mode == "internal") {
			mode = ImageMode.IM_INTERNAL
		} else if (dpImg.mode == "external") {
			mode = ImageMode.IM_EXTERNAL
		} else {
			// upgrade it to external if we are currently
			// set to none but an external link is provided
			if (!!dpImg.externalUrl) {
				mode = ImageMode.IM_EXTERNAL
			}
		}

		return {
			$typeName: "ImageMeta",
			mode,
			externalUrl: dpImg.externalUrl || "",
			internalUrl: "", // set to avoid typescript errors
		}
	} 

	async function refresh() {
		const rs = await client.getMeals({})
		displayMeals = rs.meals.map(m => ({
			id: m.id, 
			name: m.name, 
			_meal:m, 
			isEdit: false, 
			recipeUrl: m.recipeUrl,
			previewImage: mapImageMeta(m.id, "preview", m.previewImage),
			ingredientsImage: mapImageMeta(m.id, "ingredients", m.ingredientsImage),
			hasIngredients: m.ingredientRefs && m.ingredientRefs.length > 0,
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

		await client.updateMeal({
			meal: {
				id: m._meal.id,
				name: m.name,
				recipeUrl: m.recipeUrl,
				ingredientRefs: m._meal.ingredientRefs,
				previewImage: mapDPImageMeta(m.previewImage),
				ingredientsImage: mapDPImageMeta(m.ingredientsImage),
			},
		})

		if (m.previewImage.mode == "internal" && !!m.previewImage.fileBytes) {
			// a preview image was uploaded, we need to send this to the backend
			await client.updateMealPreviewImageRequest({id: m.id, imageBytes: m.previewImage.fileBytes})
		}

		if (m.ingredientsImage.mode == "internal" && !!m.ingredientsImage.fileBytes) {
			// an ingredients image was uploaded, we need to send this to the backend
			await client.updateMealIngredientsImageRequest({id: m.id, imageBytes: m.ingredientsImage.fileBytes})
		}

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
				previewImage: mapDPImageMeta(m.previewImage),
				ingredientsImage: mapDPImageMeta(m.ingredientsImage),
				ingredientRefs: [],
			},
		})

		refresh()
	}

	function addMeal() {
		displayNewMeals.push({
			pseudoId: psuedoIdCounter,
			name: "",
			recipeUrl: "",
			previewImage: emptyImageMeta(imgid(psuedoIdCounter, "preview")),
			ingredientsImage: emptyImageMeta(imgid(psuedoIdCounter, "ingredients")),
		})
		psuedoIdCounter++
	}

	function imgid(id: bigint | number, imgType: "preview" | "ingredients") {
		return `${id}_${imgType}`
	}

	function openFileDialog(img: DPImageMeta) {
		if(!img.fileUpload) {
			console.log("the fileUpload is undefined for dm: ", img.id)
			return
		}

		img.fileUpload.click()
	}

	async function fileChanged(img: DPImageMeta) {
		if(!img.fileUpload) {
			console.log("the fileUpload is undefined for dm: ", img.id)
			return
		}
		
		if (!img.fileUpload.files || img.fileUpload.files?.length == 0) {
			console.log("no files for dm: ", img.id)
			return
		}

		if (img.fileUpload.files.length > 1) {
			console.log("error: multiple files were selected for dm: ", img.id)
			return
		}

		const b = await img.fileUpload.files[0].arrayBuffer()
		img.fileBytes = new Uint8Array(b)
		img.mode = "internal"
	}

	function removeInternalImage(img: DPImageMeta) {
		// downgrade to either the external image or none
		if (!!img.externalUrl) {
			img.mode = "external"
		} else {
			img.mode = "none"
		}
	}

	refresh()
</script>
<svelte:head>
	<title>Meals</title>
</svelte:head>


<ScrollableTable classes="h-full">
	<StickyTrTitle>
		<Td title={true} colspan={5}>Meals</Td>
	</StickyTrTitle>
	<StickyTrHeader>
		<Td header={true}>Name</Td>
		<Td header={true}>Recipe Url</Td>
		<Td header={true}>Preview Image Url</Td>
		<Td header={true}>Ingredients Image Url</Td>
		<Td header={true}>Action</Td>
	</StickyTrHeader>
	{#each displayMeals as dm (dm.id)}
		{#if dm.isEdit}
			<Tr>
				<Td>
					<TextInput bind:value={dm.name}/>
				</Td>
				<Td>
					<TextInput bind:value={dm.recipeUrl}/>
				</Td>
				{#each [dm.previewImage, dm.ingredientsImage] as img}
				<Td>
					{#if img.mode == "internal" }
					<Button onclick={() => removeInternalImage(img)}>Remove Image</Button>
					{:else}
					<TextInput bind:value={img.externalUrl}/>
					<Button onclick={() => openFileDialog(img)}>Upload Image</Button>
					{/if}
					<input onchange={() => fileChanged(img)} bind:this={img.fileUpload} class="hidden" accept="image/png" type="file"/>
				</Td>
				{/each}
				<Td>
					<Button onclick={() => saveMeal(dm.id)}>Save</Button>
				</Td>
			</Tr>
		{:else}
			<Tr>
				<Td classes={!dm.hasIngredients ? 'text-red-500' : ''}>{dm.name}</Td>
				<Td>{#if dm.recipeUrl != ''}<a href={dm.recipeUrl}>Link</a>{/if}</Td>
				{#each [dm.previewImage, dm.ingredientsImage] as img}
				<Td>
					{#if img.mode == "external"}
					<img src={img.externalUrl} alt={img.id} height="50" width="50">
					{:else if img.mode == "internal"}
					<img src={img.internalUrl} alt={img.id} height="50" width="50">
					{/if}
				</Td>
				{/each}
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
			{#each [dm.previewImage, dm.ingredientsImage] as img}
				<Td>
					{#if img.mode == "internal" }
					<Button onclick={() => removeInternalImage(img)}>Remove Image</Button>
					{:else}
					<TextInput bind:value={img.externalUrl}/>
					<Button onclick={() => openFileDialog(img)}>Upload Image</Button>
					{/if}
					<input onchange={() => fileChanged(img)} bind:this={img.fileUpload} class="hidden" accept="image/png" type="file"/>
				</Td>
				{/each}
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
</ScrollableTable>
