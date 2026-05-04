import { env } from '$env/dynamic/public';

export function ShoppingSiteUrl(): string {
    return env.PUBLIC_SHOPPING_SITE_URL
}