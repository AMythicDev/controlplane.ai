import type { RequestRecord } from '$lib/types';

export interface RequestsResponse {
    requests: RequestRecord[];
    total: number;
    limit: number;
    offset: number;
}

export async function fetchRequests(limit = 50, offset = 0): Promise<RequestsResponse> {
    const res = await fetch(`/v1/requests?limit=${limit}&offset=${offset}`);
    if (!res.ok) {
        throw new Error(`Failed to fetch requests: ${res.status}`);
    }
    return res.json();
}

export async function fetchRequestById(id: string): Promise<RequestRecord> {
    const res = await fetch(`/v1/requests/${id}`);
    if (!res.ok) {
        throw new Error(`Failed to fetch request: ${res.status}`);
    }
    return res.json();
}
