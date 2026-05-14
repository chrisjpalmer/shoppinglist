import { ShoppingListService } from '../gen/shopping_list_service_pb';
import { createClient} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { env } from '$env/dynamic/public';

export function CreateShoppingListService() {
	const url = BackendUrl()
    const transport = createConnectTransport({
		baseUrl: url,
	});

	const client = createClient(ShoppingListService, transport);

    return client
}

export function BackendUrl(): string {
	return env.PUBLIC_BACKEND_URL
}