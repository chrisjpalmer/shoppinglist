<script lang="ts">
    let { onclick, children, classes, disabled }: {
        onclick: () => void
        children?: any,
        classes?: string,
        disabled?: boolean,
    } = $props();
</script>
<button disabled={disabled} class="{classes} rounded-md px-4 py-2 text-sm text-white font-semibold bg-blue-500 hover:bg-blue-400 active:bg-blue-600" onclick={onclick}>{@render children()}</button>