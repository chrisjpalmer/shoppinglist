<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} sticky top-0 h-18 group bg-white font-bold">{#if children}{@render children()}{/if}</tr>