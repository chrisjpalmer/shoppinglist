<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} group bg-white even:bg-gray-100">{#if children}{@render children()}{/if}</tr>