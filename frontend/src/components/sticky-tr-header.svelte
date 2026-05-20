<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} sticky top-18 z-2 relative shadow-sm group bg-white font-medium">{#if children}{@render children()}{/if}</tr>