<script lang="ts">
    import type { MockResult } from '$lib/data/models';
    import ScoreGauge from '$lib/components/ui/ScoreGauge.svelte';
    import Loader2 from 'lucide-svelte/icons/loader-2';
    import Clock from 'lucide-svelte/icons/clock';
    import DollarSign from 'lucide-svelte/icons/dollar-sign';
    import ShieldAlert from 'lucide-svelte/icons/shield-alert';
    import ShieldCheck from 'lucide-svelte/icons/shield-check';
    import Scale from 'lucide-svelte/icons/scale';

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

<div class="flex flex-col h-[520px] w-full rounded-2xl border border-elevated/40 bg-surface/30 backdrop-blur-sm shadow-xl overflow-hidden transition-all duration-300 hover:border-elevated">
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
    <div class="p-3.5 border-t border-elevated/30 bg-base/60 mt-auto">
        {#if isLoading}
            <div class="h-20 flex items-center justify-center">
                <span class="font-mono text-xs text-secondary opacity-50 uppercase tracking-widest animate-pulse">Generating...</span>
            </div>
        {:else if result}
            <div class="flex items-center justify-between gap-3">
                <!-- Confidence Gauge -->
                <div class="shrink-0 flex items-center">
                    {#if result.confidence !== null}
                        <ScoreGauge 
                            score={result.confidence} 
                            label="Confidence" 
                            dimension="performance"
                            size="sm"
                        />
                    {:else}
                        <div class="flex flex-col items-center opacity-50 grayscale w-20">
                            <div class="font-mono text-lg font-medium text-secondary">N/A</div>
                            <div class="mt-0.5 text-[9px] uppercase tracking-widest text-secondary font-sans font-medium text-center">No Logprobs</div>
                        </div>
                    {/if}
                </div>

                <!-- Right Metrics Column -->
                <div class="flex flex-col gap-1.5 min-w-0 flex-1 items-end">
                    <!-- Toxicity Score -->
                    <div class="flex flex-col items-end min-w-0">
                        <span class="text-[9px] tracking-wider text-secondary uppercase opacity-70 flex items-center gap-1 whitespace-nowrap">
                            {#if result.toxicity !== null && result.toxicity !== undefined && result.toxicity >= 0.5}
                                <ShieldAlert class="w-3 h-3 text-accent-danger shrink-0" />
                            {:else}
                                <ShieldCheck class="w-3 h-3 text-accent-resp shrink-0" />
                            {/if}
                            Toxicity
                        </span>
                        {#if result.toxicity !== null && result.toxicity !== undefined}
                            {@const toxVal = result.toxicity}
                            {@const toxPercent = toxVal * 100}
                            {@const isToxic = toxVal >= 0.5}
                            {@const isMedium = toxVal >= 0.1 && toxVal < 0.5}
                            <div class="flex items-center gap-1.5 mt-0.5 whitespace-nowrap">
                                <span class="font-mono text-xs sm:text-sm font-semibold {isToxic ? 'text-accent-danger' : isMedium ? 'text-accent-warning' : 'text-accent-resp'}">
                                    {toxPercent < 0.01 && toxPercent > 0 ? '<0.01%' : toxPercent < 1 ? toxPercent.toFixed(2) + '%' : toxPercent.toFixed(1) + '%'}
                                </span>
                                <span class="text-[9px] font-mono px-1.5 py-0.2 rounded uppercase font-bold tracking-wider {isToxic ? 'bg-accent-danger/20 text-accent-danger border border-accent-danger/30' : isMedium ? 'bg-accent-warning/20 text-accent-warning border border-accent-warning/30' : 'bg-accent-resp/20 text-accent-resp border border-accent-resp/30'}">
                                    {isToxic ? 'Toxic' : isMedium ? 'Flagged' : 'Clean'}
                                </span>
                            </div>
                        {:else}
                            <span class="font-mono text-xs text-secondary opacity-50">N/A</span>
                        {/if}
                    </div>

                    <!-- NLI / Hallucination Verification -->
                    {#if result.nli}
                        {@const nli = result.nli}
                        {@const isContradiction = nli.label === 'contradiction'}
                        {@const isEntailment = nli.label === 'entailment'}
                        <div class="flex flex-col items-end min-w-0">
                            <span class="text-[9px] tracking-wider text-secondary uppercase opacity-70 flex items-center gap-1 whitespace-nowrap">
                                <Scale class="w-3 h-3 {isContradiction ? 'text-accent-danger' : isEntailment ? 'text-accent-cost' : 'text-accent-warning'} shrink-0" />
                                NLI Check
                            </span>
                            <div class="flex items-center gap-1.5 mt-0.5 whitespace-nowrap">
                                <span class="font-mono text-xs sm:text-sm font-semibold {isContradiction ? 'text-accent-danger' : isEntailment ? 'text-accent-cost' : 'text-accent-warning'}">
                                    {(nli.score * 100).toFixed(0)}%
                                </span>
                                <span class="text-[9px] font-mono px-1.5 py-0.2 rounded uppercase font-bold tracking-wider {isContradiction ? 'bg-accent-danger/20 text-accent-danger border border-accent-danger/30' : isEntailment ? 'bg-accent-cost/20 text-accent-cost border border-accent-cost/30' : 'bg-accent-warning/20 text-accent-warning border border-accent-warning/30'}">
                                    {nli.label}
                                </span>
                            </div>
                        </div>
                    {/if}

                    <!-- Cost -->
                    <div class="flex flex-col items-end min-w-0">
                        <span class="text-[9px] tracking-wider text-secondary uppercase opacity-70 flex items-center gap-1 whitespace-nowrap">
                            <DollarSign class="w-3 h-3 text-accent-cost shrink-0" /> Cost
                        </span>
                        <span class="font-mono text-xs sm:text-sm text-accent-cost font-semibold whitespace-nowrap">{formatCost(result.cost)}</span>
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
