<script lang="ts">
  import { ShoppingSiteUrl } from '$lib/shopping_site_url';
  import Li from '../components/li.svelte';
  import Exit from '../lib/images/exit.svg';

  let { onclick }: { onclick?: () => void } = $props() 
  function _onclick() {
	if(onclick) {
		onclick()
		return
	}
  }
</script>

<nav>
	<ul>
		<Li onclick={_onclick} slug='/' title="Home"></Li>
		<Li onclick={_onclick} slug='/recipies' title="Recipies"></Li>
		<Li onclick={_onclick} slug='/meals' title="Meals"></Li>
		<Li onclick={_onclick} slug='/ingredients' title="Ingredients"></Li>
		<Li onclick={_onclick} slug={ShoppingSiteUrl()} title="Shopping" image={{source:Exit, alt: "shopping"}}></Li>
	</ul>
</nav>
