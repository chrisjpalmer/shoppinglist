<script lang="ts">
    let { onclick, children, classes, disabled, image }: {
        onclick: () => void
        children?: any,
        classes?: string,
        disabled?: boolean,
        image?: string,
    } = $props();
</script>
<button disabled={disabled} class="{classes} min-h-9 rounded-md px-4 py-2 text-sm text-white font-semibold bg-blue-500 disabled:bg-blue-300 not-disabled:hover:bg-blue-400 not-disabled:active:bg-blue-600" onclick={onclick}>
    {#if children}
    {@render children()}
    {/if}
    {#if image}
        <img height=14 width=14 alt="icon" src={image}/>
    {/if}
</button>