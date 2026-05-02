<script lang="ts">
    let { children, classes, header, title }: {
        children?: any,
        title?: boolean,
        header?: boolean,
        classes?: string,
    } = $props();
</script>
<td class="{classes} p-2 sm:pl-4 {title ? 'py-4':''} {header ? 'not-last:border-r-1 not-last:border-r-gray-200': ''} border-b-gray-200 border-b-1 group-first:border-t-1 group-first:border-t-gray-200 group-first:first:rounded-tl-lg group-first:last:rounded-tr-lg group-last:first:rounded-bl-lg group-last:last:rounded-br-lg">{#if children}{@render children()}{/if}</td>