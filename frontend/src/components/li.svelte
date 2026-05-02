<script lang="ts">
    import { page } from '$app/state';
    let { title, slug, onclick }: {
        title: string,
        slug: string,
        onclick: () => void,
    } = $props();
</script>
<li class="text-md py-1 text-gray-600 aria-current:font-medium aria-current:text-blue-500" aria-current={page.url.pathname === slug ? 'page' : undefined}>
    <a onclick={onclick} href={slug}>{title}</a>
</li>