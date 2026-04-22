<script lang="ts">
    let { children, classes }: {
        children?: any,
        classes?: string,
    } = $props();
</script>
<table class="{classes} w-full table-fixed border-separate border-spacing-0 rounded-lg shadow-md">
    <tbody>
        {@render children()}
    </tbody>
</table>