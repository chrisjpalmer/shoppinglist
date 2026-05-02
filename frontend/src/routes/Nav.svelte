<script lang="ts">
  import Li from '../components/li.svelte';
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
	</ul>
</nav>
