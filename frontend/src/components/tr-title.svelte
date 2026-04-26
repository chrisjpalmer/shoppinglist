<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} group bg-white font-bold">{#if children}{@render children()}{/if}</tr>