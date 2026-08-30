import type { RequestRecord, AnalyticsResponse } from '$lib/types';

export interface RequestsResponse {
    requests: RequestRecord[];
    total: number;
    limit: number;
    offset: number;
}

export interface CostResponse {
    user_id: string;
    cost_dollars: number;
    cost_microcents: number;
    semantic_cache_savings: number;
    avg_cost: number;
    average_cost_dollars?: number;
}

export interface ConfigResponse {
    per_user_daily_limit: number;
    per_user_monthly_limit: number;
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

export async function fetchAnalytics(): Promise<AnalyticsResponse> {
    const res = await fetch('/v1/analytics');
    if (!res.ok) {
        throw new Error(`Failed to fetch analytics: ${res.status}`);
    }
    return res.json();
}

export async function fetchCostData(): Promise<CostResponse> {
    const res = await fetch('/v1/cost');
    if (!res.ok) {
        throw new Error(`Failed to fetch cost: ${res.status}`);
    }
    return res.json();
}

export async function fetchConfig(): Promise<ConfigResponse> {
    const res = await fetch('/v1/config');
    if (!res.ok) {
        throw new Error(`Failed to fetch config: ${res.status}`);
    }
    return res.json();
}

