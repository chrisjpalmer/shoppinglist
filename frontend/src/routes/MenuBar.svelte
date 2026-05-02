<script lang="ts">
  import menuSvg from '$lib/images/menu.svg'
  import crossSvg from '$lib/images/cross.svg'
  import Nav from './Nav.svelte';

  let expanded = $state(false)
  function show() {
    expanded = true
  }
  function hide() {
    expanded = false
  }
</script>
<div class="sm:hidden absolute top-0 left-0 p-5 z-5 shadow-md w-full">
    <button class="cursor-pointer" onclick={() => show()}>
	    <img src={menuSvg} alt="menu" height="20" width="20" />
    </button>
</div>
<div aria-expanded={expanded} class="transition-[left] sm:hidden absolute top-0 -left-50 aria-expanded:left-0 p-5 bg-white h-screen z-6 w-50">
    <button class="cursor-pointer mb-2" onclick={() => hide()}>
        <img src={crossSvg} alt="exit" height="20" width="20" />
    </button>
    <Nav onclick={() => hide()}></Nav>
</div>