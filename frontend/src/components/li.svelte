<script lang="ts">
    import { page } from '$app/state';
    export interface Image {
        source: string
        alt: string
    }
    let { title, slug, onclick, image }: {
        title: string,
        slug: string,
        onclick: () => void,
        image?: Image
    } = $props();
</script>
<li class="text-md py-1 text-gray-600 aria-current:font-medium aria-current:text-blue-500" aria-current={page.url.pathname === slug ? 'page' : undefined}>
    <a onclick={onclick} href={slug}>{title}</a>
    {#if image}
    <img class="ml-5" src={image.source} alt={image.alt} height="20" width="20">
    {/if}
</li>