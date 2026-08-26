<script lang="ts">
    import type { MockResult } from '$lib/data/models';
    import ScoreGauge from '$lib/components/ui/ScoreGauge.svelte';
    import Loader2 from 'lucide-svelte/icons/loader-2';
    import Clock from 'lucide-svelte/icons/clock';
    import DollarSign from 'lucide-svelte/icons/dollar-sign';
    import Activity from 'lucide-svelte/icons/activity';

    let { result, isLoading, providerColor, providerName, modelName } = $props<{
        result: MockResult | null;
        isLoading: boolean;
        providerColor: string;
        providerName: string;
        modelName: string;
    }>();

    // Formatters
    const formatCost = (val: number) =>
        new Intl.NumberFormat('en-US', {
            style: 'currency',
            currency: 'USD',
            minimumFractionDigits: 3,
            maximumFractionDigits: 4
        }).format(val);
</script>

<div class="flex flex-col h-[500px] w-full rounded-2xl border border-elevated/40 bg-surface/30 backdrop-blur-sm shadow-xl overflow-hidden transition-all duration-300 hover:border-elevated">
    <!-- Card Header -->
    <div class="flex items-center justify-between p-4 border-b border-elevated/30 bg-surface/50">
        <div class="flex items-center gap-3">
            <!-- Provider Badge -->
            <div 
                class="w-3 h-3 rounded-full shadow-[0_0_10px_rgba(255,255,255,0.2)]" 
                style="background-color: {providerColor}; box-shadow: 0 0 10px {providerColor}40;"
            ></div>
            <div class="flex flex-col">
                <span class="text-[10px] tracking-widest text-secondary uppercase opacity-80">{providerName}</span>
                <span class="font-mono text-sm font-medium text-primary">{modelName}</span>
            </div>
        </div>
        
        <!-- Status Indicator -->
        {#if isLoading}
            <Loader2 class="w-4 h-4 text-accent-perf animate-spin" />
        {:else if result}
            <div class="flex items-center gap-1.5 px-2 py-1 rounded bg-surface border border-elevated text-xs font-mono">
                <Clock class="w-3 h-3 text-secondary" />
                <span class="text-primary">{result.latency_ms}ms</span>
            </div>
        {/if}
    </div>

    <!-- Card Body (Content) -->
    <div class="flex-1 p-4 overflow-y-auto custom-scrollbar relative">
        {#if isLoading}
            <div class="flex flex-col gap-3 opacity-50 animate-pulse">
                <div class="h-3 w-3/4 rounded bg-elevated"></div>
                <div class="h-3 w-full rounded bg-elevated"></div>
                <div class="h-3 w-5/6 rounded bg-elevated"></div>
                <div class="h-3 w-full rounded bg-elevated"></div>
                <div class="h-3 w-2/3 rounded bg-elevated"></div>
            </div>
        {:else if result}
            <div class="font-sans text-sm leading-relaxed text-secondary whitespace-pre-wrap">
                {result.content}
            </div>
        {:else}
            <div class="h-full flex items-center justify-center text-secondary opacity-50 font-mono text-sm">
                Waiting to run...
            </div>
        {/if}
    </div>

    <!-- Card Footer (Metrics) -->
    <div class="p-4 border-t border-elevated/30 bg-base/50 mt-auto">
        {#if isLoading}
            <div class="h-20 flex items-center justify-center">
                <span class="font-mono text-xs text-secondary opacity-50 uppercase tracking-widest animate-pulse">Generating...</span>
            </div>
        {:else if result}
            <div class="flex items-center justify-between gap-4">
                <!-- Confidence Gauge -->
                <div class="scale-75 origin-left">
                    {#if result.confidence !== null}
                        <ScoreGauge 
                            score={result.confidence} 
                            label="Confidence" 
                            dimension="performance" 
                        />
                    {:else}
                        <div class="flex flex-col items-center opacity-50 grayscale pt-2">
                            <div class="font-mono text-2xl font-medium text-secondary">N/A</div>
                            <div class="mt-2 text-xs uppercase tracking-widest text-secondary font-sans font-medium">No Logprobs</div>
                        </div>
                    {/if}
                </div>

                <!-- Other Metrics -->
                <div class="flex flex-col gap-3 flex-1 items-end pt-2">
                    <!-- Cost -->
                    <div class="flex flex-col items-end">
                        <span class="text-[9px] tracking-widest text-secondary uppercase opacity-60 flex items-center gap-1">
                            <DollarSign class="w-3 h-3" /> Cost
                        </span>
                        <span class="font-mono text-sm text-accent-cost">{formatCost(result.cost)}</span>
                    </div>
                    
                    <!-- Perplexity -->
                    <div class="flex flex-col items-end">
                        <span class="text-[9px] tracking-widest text-secondary uppercase opacity-60 flex items-center gap-1">
                            <Activity class="w-3 h-3" /> Perplexity
                        </span>
                        <span class="font-mono text-sm text-primary">
                            {#if result.perplexity !== null}
                                {result.perplexity.toFixed(3)}
                            {:else}
                                <span class="opacity-50">N/A</span>
                            {/if}
                        </span>
                    </div>
                </div>
            </div>
        {:else}
            <div class="h-20 flex items-center justify-center opacity-0">
                <!-- Empty placeholder to maintain height -->
            </div>
        {/if}
    </div>
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: var(--bg-elevated);
        border-radius: 2px;
    }
    .custom-scrollbar:hover::-webkit-scrollbar-thumb {
        background: var(--bg-hover);
    }
</style>
