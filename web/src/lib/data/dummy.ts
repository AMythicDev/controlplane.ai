import type { Request, Policy, ReviewItem, DashboardStats, TimeSeriesPoint, CostBreakdown } from '$lib/types';

export const mockPolicies: Policy[] = [
    {
        id: 'pol_1', name: 'Block Toxic Content', engine: 'responsibility', check: 'ToxicityCheck', severity: 'critical', action: 'block', enabled: true,
        yaml: `name: Block Toxic Content\nengine: responsibility\naction: block\nconditions:\n  - check: toxicity\n    threshold: > 0.8`, description: 'Blocks any output with toxicity score > 0.8'
    },
    {
        id: 'pol_2', name: 'Flag Low Confidence', engine: 'performance', check: 'ConfidenceScore', severity: 'medium', action: 'flag', enabled: true,
        yaml: `name: Flag Low Confidence\nengine: performance\naction: flag\nconditions:\n  - check: confidence\n    threshold: < 0.6`, description: 'Flags responses where the model exhibits low confidence'
    },
    {
        id: 'pol_3', name: 'Auto-Redact PII', engine: 'responsibility', check: 'PIIDetector', severity: 'high', action: 'edit', enabled: true,
        yaml: `name: Auto-Redact PII\nengine: responsibility\naction: edit\nconditions:\n  - check: pii\n    types: [SSN, EMAIL, PHONE]`, description: 'Automatically redacts PII before returning response'
    },
    {
        id: 'pol_4', name: 'Escalate Bias Concerns', engine: 'responsibility', check: 'BiasAnalyzer', severity: 'high', action: 'escalate', enabled: true,
        yaml: `name: Escalate Bias Concerns\nengine: responsibility\naction: escalate\nconditions:\n  - check: bias\n    threshold: > 0.7`, description: 'Escalates highly biased responses to human review'
    },
    {
        id: 'pol_5', name: 'Enforce Budget limits', engine: 'cost', check: 'BudgetTracker', severity: 'high', action: 'block', enabled: true,
        yaml: `name: Enforce Budget\nengine: cost\naction: block\nconditions:\n  - check: monthly_spend\n    threshold: > 5000`, description: 'Blocks requests if monthly budget exceeded'
    },
    {
        id: 'pol_6', name: 'Block Prompt Injection', engine: 'responsibility', check: 'JailbreakDetector', severity: 'critical', action: 'block', enabled: true,
        yaml: `name: Block Prompt Injection\nengine: responsibility\naction: block\nconditions:\n  - check: jailbreak\n    threshold: > 0.9`, description: 'Blocks requests containing potential jailbreak attempts'
    },
    {
        id: 'pol_7', name: 'Flag Hallucination', engine: 'performance', check: 'FactualConsistency', severity: 'high', action: 'flag', enabled: true,
        yaml: `name: Flag Hallucination\nengine: performance\naction: flag\nconditions:\n  - check: factual_consistency\n    threshold: < 0.5`, description: 'Flags output that contradicts established facts'
    },
    {
        id: 'pol_8', name: 'Cache Duplicate Queries', engine: 'cost', check: 'SemanticCache', severity: 'low', action: 'pass', enabled: true,
        yaml: `name: Cache Duplicate Queries\nengine: cost\naction: pass\nconditions:\n  - check: semantic_similarity\n    threshold: > 0.95`, description: 'Serves identical/highly similar queries from cache'
    }
];

export const mockRequests: Request[] = Array.from({ length: 50 }).map((_, i) => {
    const isBad = i % 15 === 0;
    const isModerate = i % 7 === 0 && !isBad;
    
    let action: Request['action'] = 'pass';
    if (isBad) action = i % 2 === 0 ? 'block' : 'escalate';
    else if (isModerate) action = i % 2 === 0 ? 'flag' : 'edit';
    
    const models = ['gpt-4o', 'gpt-4o-mini', 'claude-3.5-sonnet', 'gemini-1.5-pro'];
    const model = models[i % models.length];

    const endpoints = ['/v1/chat/completions', '/v1/completions', '/v1/assistants'];

    return {
        id: `req_${Math.random().toString(36).substring(2, 9)}`,
        timestamp: new Date(Date.now() - Math.floor(Math.random() * 86400000)).toISOString(),
        endpoint: endpoints[i % endpoints.length],
        userId: `usr_${Math.random().toString(36).substring(2, 8)}`,
        userName: `User ${i + 1}`,
        model,
        performanceScore: isBad ? 30 + Math.random() * 20 : isModerate ? 60 + Math.random() * 20 : 85 + Math.random() * 14,
        costAmount: Math.random() * 0.05 + 0.001,
        responsibilityScore: isBad ? 20 + Math.random() * 30 : isModerate ? 70 + Math.random() * 10 : 90 + Math.random() * 9,
        action,
        prompt: `This is a sample prompt ${i} regarding a business inquiry.`,
        response: `This is the generated AI response ${i} addressing the business inquiry in detail.`,
        tokens: {
            input: Math.floor(Math.random() * 500) + 50,
            output: Math.floor(Math.random() * 800) + 100,
            reasoning: Math.floor(Math.random() * 100)
        },
        cacheHit: Math.random() > 0.8,
        checks: [
            { name: 'Toxicity', layer: 1, duration: 12, status: isBad ? 'fail' : 'pass', score: isBad ? 0.9 : 0.1, details: 'Checked for toxic language.' }
        ],
        policyTriggered: action !== 'pass' ? mockPolicies[Math.floor(Math.random() * mockPolicies.length)].name : undefined
    };
});

export const mockReviewQueue: ReviewItem[] = [
    { id: 'rev_1', requestId: mockRequests[0].id, reason: 'High bias score detected', severity: 'high', escalatedAt: new Date().toISOString(), slaDeadline: new Date(Date.now() + 3600000).toISOString(), originalResponse: '...', suggestedEdit: '...', status: 'pending' },
    { id: 'rev_2', requestId: mockRequests[15].id, reason: 'Potential prompt injection', severity: 'critical', escalatedAt: new Date().toISOString(), slaDeadline: new Date(Date.now() + 1800000).toISOString(), originalResponse: '...', suggestedEdit: '...', status: 'pending' },
    { id: 'rev_3', requestId: mockRequests[30].id, reason: 'PII exposure', severity: 'high', escalatedAt: new Date().toISOString(), slaDeadline: new Date(Date.now() + 3600000).toISOString(), originalResponse: '...', suggestedEdit: '...', status: 'approved' },
    { id: 'rev_4', requestId: mockRequests[45].id, reason: 'Toxicity threshold exceeded', severity: 'critical', escalatedAt: new Date().toISOString(), slaDeadline: new Date(Date.now() + 1800000).toISOString(), originalResponse: '...', suggestedEdit: '...', status: 'rejected' },
    { id: 'rev_5', requestId: mockRequests[7].id, reason: 'Low factual consistency', severity: 'medium', escalatedAt: new Date().toISOString(), slaDeadline: new Date(Date.now() + 7200000).toISOString(), originalResponse: '...', suggestedEdit: '...', status: 'pending' },
];

export const mockDashboardStats: DashboardStats = {
    totalRequests: 14592,
    avgPerformance: 94.2,
    totalSpend: 1245.60,
    avgResponsibility: 98.1,
    requestsPerMinute: 42
};

export const mockTimeSeriesPerf: TimeSeriesPoint[] = Array.from({ length: 24 }).map((_, i) => ({
    timestamp: new Date(Date.now() - (23 - i) * 3600000).toISOString(),
    value: 80 + Math.random() * 20
}));

export const mockTimeSeriesCost: TimeSeriesPoint[] = Array.from({ length: 24 }).map((_, i) => ({
    timestamp: new Date(Date.now() - (23 - i) * 3600000).toISOString(),
    value: 10 + Math.random() * 50
}));

export const mockTimeSeriesResp: TimeSeriesPoint[] = Array.from({ length: 24 }).map((_, i) => ({
    timestamp: new Date(Date.now() - (23 - i) * 3600000).toISOString(),
    value: 90 + Math.random() * 10
}));

export const mockCostBreakdown: CostBreakdown[] = [
    { model: 'gpt-4o', amount: 800, percentage: 64.2 },
    { model: 'claude-3.5-sonnet', amount: 300, percentage: 24.1 },
    { model: 'gpt-4o-mini', amount: 100, percentage: 8.0 },
    { model: 'gemini-1.5-pro', amount: 45.6, percentage: 3.7 }
];
