<script lang="ts">
    let { children, classes }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<table class="{classes} w-full border-separate border-spacing-0 rounded-lg shadow-md sm:text-base text-sm">
    <tbody>
        {@render children()}
    </tbody>
</table>