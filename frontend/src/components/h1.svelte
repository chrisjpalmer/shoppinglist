<script lang="ts">
    let { children, classes }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<h1 class="{classes} font-normal text-3xl text-center sm:text-4xl py-10 font-sans">{@render children()}</h1>