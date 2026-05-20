<script lang="ts">
    let { children, classes }: {
        children?: any,
        classes?: string,
    } = $props();
</script>

<div class="{classes} overflow-y-auto rounded-lg shadow-md">
    <table class="w-full border-separate border-spacing-0 sm:text-base text-sm">
        <tbody>
            {@render children()}
        </tbody>
    </table>
</div>