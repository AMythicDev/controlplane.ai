export interface ModelOption {
    id: string;
    name: string;
    provider: string;
    spec: string; // e.g. "openai/gpt-4o"
}

export interface ProviderGroup {
    id: string;
    name: string;
    color: string;
    models: ModelOption[];
}

export const modelCatalog: ProviderGroup[] = [
    {
        id: 'openai',
        name: 'OpenAI',
        color: '#10a37f',
        models: [
            { id: 'gpt-4o', name: 'GPT-4o', provider: 'openai', spec: 'openai/gpt-4o' },
            { id: 'gpt-4o-mini', name: 'GPT-4o Mini', provider: 'openai', spec: 'openai/gpt-4o-mini' },
            { id: 'o3', name: 'o3', provider: 'openai', spec: 'openai/o3' },
            { id: 'o3-mini', name: 'o3 Mini', provider: 'openai', spec: 'openai/o3-mini' },
        ]
    },
    {
        id: 'anthropic',
        name: 'Anthropic',
        color: '#d97706',
        models: [
            { id: 'claude-sonnet-4', name: 'Claude Sonnet 4', provider: 'anthropic', spec: 'anthropic/claude-sonnet-4-20250514' },
            { id: 'claude-haiku-3-5', name: 'Claude 3.5 Haiku', provider: 'anthropic', spec: 'anthropic/claude-haiku-3-5-20241022' },
        ]
    },
    {
        id: 'google',
        name: 'Google',
        color: '#4285f4',
        models: [
            { id: 'gemini-2.5-pro', name: 'Gemini 2.5 Pro', provider: 'google', spec: 'google/gemini-2.5-pro' },
            { id: 'gemini-2.5-flash', name: 'Gemini 2.5 Flash', provider: 'google', spec: 'google/gemini-2.5-flash' },
        ]
    },
    {
        id: 'openrouter',
        name: 'OpenRouter',
        color: '#6366f1',
        models: [
            { id: 'meta-llama/llama-4-maverick', name: 'Llama 4 Maverick', provider: 'openrouter', spec: 'openrouter/meta-llama/llama-4-maverick' },
            { id: 'deepseek/deepseek-r1', name: 'DeepSeek R1', provider: 'openrouter', spec: 'openrouter/deepseek/deepseek-r1' },
        ]
    },
    {
        id: 'nvidia',
        name: 'NVIDIA',
        color: '#76b900',
        models: [
            { id: 'google/diffusiongemma-26b-a4b-it', name: 'Diffusion Gemma 20B', provider: 'nvidia', spec: 'nvidia/google/diffusiongemma-26b-a4b-it' },
            { id: 'nvidia/nemotron-3.5-lightning-30b-a3b', name: 'Nemotron 5.3 Lightning 30B', provider: 'nvidia', spec: 'nvidia/nvidia/nemotron-3.5-lightning-30b-a3b' },
        ]
    }
];

export interface NLIResult {
    label: 'contradiction' | 'entailment' | 'neutral' | string;
    score: number;
    contradiction_prob: number;
    neutral_prob: number;
    entailment_prob: number;
}

export interface MockResult {
    model: string;
    provider: string;
    content: string;
    confidence: number | null;
    toxicity: number | null;
    nli: NLIResult | null;
    latency_ms: number;
    cost: number;
}

export const runPlaygroundRequest = async (prompt: string, modelSpec: string): Promise<MockResult> => {
    const res = await fetch('/v1/playground', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            prompt: prompt,
            model_spec: modelSpec
        })
    });

    if (!res.ok) {
        const errorData = await res.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error ${res.status}`);
    }

    const data = await res.json();
    return {
        model: data.model,
        provider: data.provider,
        content: data.content,
        confidence: data.confidence,
        toxicity: data.toxicity ?? null,
        nli: data.nli ?? null,
        latency_ms: data.latency_ms,
        cost: data.cost
    };
};


