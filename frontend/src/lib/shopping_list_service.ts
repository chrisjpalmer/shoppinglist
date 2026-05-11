import { ShoppingListService } from '../gen/shopping_list_service_pb';
import { createClient} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { page } from '$app/state'
import { PUBLIC_BACKEND_PORT } from "$env/static/public";

export function CreateShoppingListService() {
	const url = BackendUrl()
    const transport = createConnectTransport({
		baseUrl: url,
	});

	const client = createClient(ShoppingListService, transport);

    return client
}

export function BackendUrl(): string {
	return `http://${page.url.hostname}:${PUBLIC_BACKEND_PORT}`
}