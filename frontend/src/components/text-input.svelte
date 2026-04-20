<script lang="ts">
    let { value = $bindable(), classes }: {
        value: any,
        classes?: string,
    } = $props();
</script>
<input class="{classes} bg-white rounded-md h-7 px-2 border-solid border-gray-500 border-1 focus:border-gray-900 focus:outline-none" type="text" bind:value={value}>