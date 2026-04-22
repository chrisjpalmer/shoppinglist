<script lang="ts">
    let { children, classes, }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<tr class="{classes} group bg-white h-10 even:bg-gray-100">{#if children}{@render children()}{/if}</tr>