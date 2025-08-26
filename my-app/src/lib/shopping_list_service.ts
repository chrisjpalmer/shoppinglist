import { ShoppingListService } from '../gen/shopping_list_service_pb';
import { createClient} from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

export function CreateShoppingListService(url: string) {
    const transport = createConnectTransport({
		baseUrl: url,
	});

	const client = createClient(ShoppingListService, transport);

    return client
}