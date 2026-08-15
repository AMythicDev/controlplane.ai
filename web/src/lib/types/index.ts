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
