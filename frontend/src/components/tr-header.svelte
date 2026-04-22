<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} z-2 relative shadow-sm group h-10 bg-white font-medium">{#if children}{@render children()}{/if}</tr>