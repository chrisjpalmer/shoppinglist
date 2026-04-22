<script lang="ts">
    let { value = $bindable(), children, classes }: {
        value: any
        children?: any,
        classes?: string,
    } = $props();
</script>
<select class="{classes} w-full bg-white rounded-md h-7 px-2 border-solid border-gray-500 border-1 focus:border-gray-900 focus:outline-none" bind:value={value}>
    {@render children()}				
</select>