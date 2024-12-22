import { hostname } from "./config";

export default async function getFromApi<T>(path: string): Promise<T> {
    const response = await fetch(`http://${hostname}/${path}`);

    if (!response.ok) {
        throw new Error(response.statusText);
    }

    return await response.json() as T;
}