export interface CheckResult {
    name: string;
    layer: 1 | 2 | 3 | 4;
    duration: number; // ms
    status: 'pass' | 'fail' | 'warn';
    score: number;
    details: string;
}

export interface Request {
    id: string;
    timestamp: string;
    endpoint: string;
    userId: string;
    userName: string;
    model: string;
    performanceScore: number;
    costAmount: number;
    responsibilityScore: number;
    action: 'pass' | 'flag' | 'edit' | 'block' | 'escalate';
    prompt: string;
    response: string;
    tokens: {
        input: number;
        output: number;
        reasoning: number;
    };
    cacheHit: boolean;
    checks: CheckResult[];
    policyTriggered?: string;
}

export interface Policy {
    id: string;
    name: string;
    engine: 'performance' | 'cost' | 'responsibility';
    check: string;
    severity: string;
    action: string;
    enabled: boolean;
    yaml: string;
    description: string;
}

export interface ReviewItem {
    id: string;
    requestId: string;
    reason: string;
    severity: 'critical' | 'high' | 'medium';
    escalatedAt: string;
    slaDeadline: string;
    originalResponse: string;
    suggestedEdit: string;
    status: 'pending' | 'approved' | 'rejected';
}

export interface DashboardStats {
    totalRequests: number;
    avgPerformance: number;
    totalSpend: number;
    avgResponsibility: number;
    requestsPerMinute: number;
}

export interface TimeSeriesPoint {
    timestamp: string;
    value: number;
}

export interface CostBreakdown {
    model: string;
    amount: number;
    percentage: number;
}

export interface NLIResult {
    label: string;
    score: number;
    contradiction_prob: number;
    neutral_prob: number;
    entailment_prob: number;
}

export interface RequestRecord {
    id: string;
    endpoint: string;
    model: string;
    provider: string;
    prompt: string;
    messages?: { Role: string; Content: string }[];
    response: string;
    confidence: number | null;
    toxicity: number;
    nli: NLIResult | null;
    latency_ms: number;
    cost_microcents: number;
    cached?: boolean;
    timestamp: string;
}

export interface DailyModelCount {
    date: string;
    count: number;
}

export interface ModelAnalyticsItem {
    model: string;
    provider: string;
    request_count: number;
    percentage: number;
    avg_confidence: number | null;
    confidence_count: number;
    avg_hallucination: number | null;
    nli_count: number;
    avg_toxicity: number;
    toxicity_count: number;
    total_cost_microcents: number;
    total_cost_dollars: number;
    daily_counts: DailyModelCount[];
}

export interface AnalyticsResponse {
    total_requests: number;
    weekly_requests: number;
    weekly_start_date: string;
    weekly_end_date: string;
    models: ModelAnalyticsItem[];
    daily_totals: DailyModelCount[];
}

