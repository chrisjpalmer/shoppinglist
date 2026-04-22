<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} group h-10 bg-white font-bold">{#if children}{@render children()}{/if}</tr>