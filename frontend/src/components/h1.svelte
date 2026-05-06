<script lang="ts">
    let { children, classes }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<h1 class="{classes} font-normal text-3xl text-left sm:text-4xl mb-5 font-sans">{@render children()}</h1>