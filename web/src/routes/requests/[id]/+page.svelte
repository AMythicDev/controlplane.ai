<script lang="ts">
    import { page } from '$app/stores';
    import { fetchRequestById } from '$lib/data/api';
    import type { RequestRecord } from '$lib/types';
    import ScoreGauge from '$lib/components/ui/ScoreGauge.svelte';
    import Zap from 'lucide-svelte/icons/zap';
    
    let requestId = $derived($page.params.id);
    let request = $state<RequestRecord | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    $effect(() => {
        if (requestId) {
            loadData(requestId);
        }
    });

    async function loadData(id: string) {
        if (!id) return;
        try {
            loading = true;
            error = null;
            request = await fetchRequestById(id);
        } catch (e: any) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    function formatTime(isoStr: string) {
        const d = new Date(isoStr);
        return d.toLocaleString();
    }

    function formatText(text: string) {
        if (!text) return '';
        return text;
    }

    function formatConfidenceScore(conf: number | null): number {
        if (conf === null || conf === undefined) return 0;
        return Math.round(conf > 1 ? conf : conf * 100);
    }

    function formatToxicityScore(tox: number | null): number {
        if (tox === null || tox === undefined) return 0;
        return Math.round(tox > 1 ? tox : tox * 100);
    }

    function formatCost(microcents: number): string {
        const dollars = (microcents || 0) / 1000000;
        if (dollars === 0) return '$0.00';
        if (dollars < 0.01) return `$${dollars.toFixed(4)}`;
        return `$${dollars.toFixed(2)}`;
    }

    function getNliStyle(label: string) {
        switch(label) {
            case 'entailment': return 'text-[var(--accent-responsibility)] border-[var(--accent-responsibility)] bg-[color-mix(in_srgb,var(--accent-responsibility)_10%,transparent)]';
            case 'neutral': return 'text-[var(--accent-warning)] border-[var(--accent-warning)] bg-[color-mix(in_srgb,var(--accent-warning)_10%,transparent)]';
            case 'contradiction': return 'text-[var(--accent-danger)] border-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_10%,transparent)]';
            default: return 'text-secondary border-secondary bg-surface';
        }
    }
</script>

{#if loading}
    <div class="p-12 text-center font-mono text-secondary">
        <p class="text-2xl mb-4">Loading request...</p>
    </div>
{:else if error}
    <div class="p-12 text-center font-mono text-[var(--accent-danger)]">
        <p class="text-2xl mb-4">Error: {error}</p>
        <a href="/requests" class="text-accent-perf hover:underline">&larr; Back to Requests</a>
    </div>
{:else if !request}
    <div class="p-12 text-center font-mono text-secondary">
        <p class="text-2xl mb-4">Request Not Found</p>
        <a href="/requests" class="text-accent-perf hover:underline">&larr; Back to Requests</a>
    </div>
{:else}
    <div class="p-8 max-w-7xl mx-auto animate-fade-up pb-24">
        <!-- Header -->
        <a href="/requests" class="inline-flex items-center text-sm font-mono text-secondary hover:text-primary mb-6 transition-colors">
            &larr; Back to Requests
        </a>
        
        <header class="flex flex-col md:flex-row md:justify-between md:items-start gap-4 mb-10 pb-6 border-b border-elevated">
            <div>
                <div class="flex items-center gap-3 mb-2">
                    <h1 class="font-mono text-3xl text-primary tracking-tight">{request.id}</h1>
                    {#if request.cached}
                        <span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-accent-perf/20 border border-accent-perf/40 text-xs font-mono text-accent-perf font-bold tracking-wider uppercase shadow-[0_0_12px_rgba(161,0,255,0.25)]">
                            <Zap class="w-3.5 h-3.5 text-accent-perf fill-accent-perf" />
                            Cache Hit
                        </span>
                    {/if}
                </div>
                <div class="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm font-mono text-secondary mt-4">
                    <div class="flex items-center gap-2"><span class="text-tertiary">TIME:</span> {formatTime(request.timestamp)}</div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">PROVIDER:</span> {request.provider}</div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">MODEL:</span> <span class="text-accent-perf">{request.model}</span></div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">ENDPOINT:</span> {request.endpoint}</div>
                    <div class="flex items-center gap-2">
                        <span class="text-tertiary">CACHE:</span>
                        {#if request.cached}
                            <span class="inline-flex items-center gap-1 text-accent-perf font-semibold"><Zap class="w-3 h-3 fill-accent-perf" /> Hit</span>
                        {:else}
                            <span class="text-secondary">Miss</span>
                        {/if}
                    </div>
                </div>
            </div>
        </header>

        <!-- Evaluation Scores -->
        <section class="mb-12">
            <h2 class="font-serif text-3xl text-primary mb-6">Evaluation</h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="bg-surface p-6 rounded-xl border border-elevated flex items-center justify-center">
                    {#if request.confidence !== null}
                        <ScoreGauge score={formatConfidenceScore(request.confidence)} label="Confidence" dimension="performance" />
                    {:else}
                        <div class="text-secondary font-mono">No Confidence Data</div>
                    {/if}
                </div>
                <div class="bg-surface p-6 rounded-xl border border-elevated flex items-center justify-center">
                    <ScoreGauge score={formatToxicityScore(request.toxicity)} label="Toxicity" dimension="responsibility" />
                </div>
                <div class="bg-surface p-6 rounded-xl border border-elevated flex flex-col items-center justify-center">
                    <div class="text-3xl font-mono text-[var(--accent-cost)] mb-2">
                        {formatCost(request.cost_microcents)}
                    </div>
                    <div class="text-xs uppercase tracking-widest text-secondary font-sans font-medium text-center">
                        Cost
                    </div>
                </div>
            </div>
        </section>

        <!-- NLI Section -->
        {#if request.nli}
            <section class="mb-12">
                <h2 class="font-serif text-3xl text-primary mb-6">Hallucination (NLI) Analysis</h2>
                <div class="bg-surface border border-elevated rounded-xl p-6 relative overflow-hidden">
                    <div class="flex items-center gap-4 mb-6">
                        <span class="text-sm font-mono text-secondary uppercase tracking-widest">Result:</span>
                        <span class="inline-block px-4 py-1.5 rounded-sm text-sm uppercase font-bold font-mono border {getNliStyle(request.nli.label)}">
                            {request.nli.label}
                        </span>
                    </div>
                    
                    <div class="space-y-4">
                        <div>
                            <div class="flex justify-between text-xs font-mono text-secondary mb-1">
                                <span>Entailment</span>
                                <span>{(request.nli.entailment_prob * 100).toFixed(1)}%</span>
                            </div>
                            <div class="w-full bg-elevated rounded-full h-2">
                                <div class="bg-[var(--accent-responsibility)] h-2 rounded-full" style="width: {request.nli.entailment_prob * 100}%"></div>
                            </div>
                        </div>
                        <div>
                            <div class="flex justify-between text-xs font-mono text-secondary mb-1">
                                <span>Neutral</span>
                                <span>{(request.nli.neutral_prob * 100).toFixed(1)}%</span>
                            </div>
                            <div class="w-full bg-elevated rounded-full h-2">
                                <div class="bg-[var(--accent-warning)] h-2 rounded-full" style="width: {request.nli.neutral_prob * 100}%"></div>
                            </div>
                        </div>
                        <div>
                            <div class="flex justify-between text-xs font-mono text-secondary mb-1">
                                <span>Contradiction</span>
                                <span>{(request.nli.contradiction_prob * 100).toFixed(1)}%</span>
                            </div>
                            <div class="w-full bg-elevated rounded-full h-2">
                                <div class="bg-[var(--accent-danger)] h-2 rounded-full" style="width: {request.nli.contradiction_prob * 100}%"></div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        {/if}

        <!-- Prompt & Response -->
        <section>
            <h2 class="font-serif text-3xl text-primary mb-6">Payload</h2>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <!-- Prompt -->
                <div class="bg-surface border border-elevated rounded-xl overflow-hidden flex flex-col">
                    <div class="bg-elevated px-4 py-3 border-b border-hover flex justify-between items-center">
                        <span class="font-mono text-xs font-bold text-tertiary uppercase tracking-widest">Prompt</span>
                    </div>
                    <div class="p-6 font-mono text-sm text-primary whitespace-pre-wrap flex-grow bg-[#111] relative">
                        <div class="noise-overlay" style="opacity: 0.05"></div>
                        {@html formatText(request.prompt)}
                    </div>
                </div>

                <!-- Response -->
                <div class="bg-surface border border-elevated rounded-xl overflow-hidden flex flex-col {request.cached ? 'ring-1 ring-accent-perf/30' : ''}">
                    <div class="bg-elevated px-4 py-3 border-b border-hover flex justify-between items-center">
                        <span class="font-mono text-xs font-bold text-tertiary uppercase tracking-widest">Response</span>
                        {#if request.cached}
                            <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-accent-perf/20 border border-accent-perf/30 text-[10px] font-mono text-accent-perf uppercase font-bold tracking-wider">
                                <Zap class="w-3 h-3 fill-accent-perf" /> Cache Hit
                            </span>
                        {/if}
                    </div>
                    <div class="p-6 font-mono text-sm text-primary whitespace-pre-wrap flex-grow bg-[#111] relative">
                        <div class="noise-overlay" style="opacity: 0.05"></div>
                        {@html formatText(request.response)}
                    </div>
                </div>
            </div>
        </section>
    </div>
{/if}
