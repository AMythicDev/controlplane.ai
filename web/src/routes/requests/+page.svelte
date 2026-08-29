<script lang="ts">
    import { goto } from '$app/navigation';
    import type { RequestRecord } from '$lib/types';
    import { fetchRequests } from '$lib/data/api';
    import Zap from 'lucide-svelte/icons/zap';

    let requests = $state<RequestRecord[]>([]);
    let totalRequests = $state(0);
    let loading = $state(true);
    let error = $state<string | null>(null);

    let selectedEndpoints = $state<Set<string>>(new Set());
    let selectedModel = $state<string>('All');
    let sortCol = $state<keyof RequestRecord>('timestamp');
    let sortAsc = $state<boolean>(false);
    
    let currentPage = $state(1);
    const ITEMS_PER_PAGE = 50;

    $effect(() => {
        loadData(currentPage);
    });

    async function loadData(page: number) {
        try {
            loading = true;
            error = null;
            const offset = (page - 1) * ITEMS_PER_PAGE;
            const result = await fetchRequests(ITEMS_PER_PAGE, offset);
            requests = result.requests;
            totalRequests = result.total;
        } catch (e: any) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    let availableModels = $derived(['All', ...new Set(requests.map(r => r.model))]);

    function toggleEndpoint(endpoint: string) {
        if (selectedEndpoints.has(endpoint)) {
            selectedEndpoints.delete(endpoint);
        } else {
            selectedEndpoints.add(endpoint);
        }
    }

    let filteredRequests = $derived(
        requests.filter(req => {
            if (selectedEndpoints.size > 0 && !selectedEndpoints.has(req.endpoint)) return false;
            if (selectedModel !== 'All' && req.model !== selectedModel) return false;
            return true;
        }).sort((a, b) => {
            let valA = a[sortCol];
            let valB = b[sortCol];
            
            if (sortCol === 'confidence') {
                valA = a.confidence ?? 0;
                valB = b.confidence ?? 0;
            } else if (sortCol === 'nli' as any) {
                valA = a.nli?.label ?? '';
                valB = b.nli?.label ?? '';
            }

            if (valA! < valB!) return sortAsc ? -1 : 1;
            if (valA! > valB!) return sortAsc ? 1 : -1;
            return 0;
        })
    );

    let totalPages = $derived(Math.ceil(totalRequests / ITEMS_PER_PAGE));

    function handleSort(col: keyof RequestRecord | 'nli') {
        if (sortCol === col as any) {
            sortAsc = !sortAsc;
        } else {
            sortCol = col as keyof RequestRecord;
            sortAsc = false; // default desc for new sort
        }
    }

    function formatTime(isoStr: string) {
        const d = new Date(isoStr);
        return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    }

    function formatConfidence(conf: number | null): number | null {
        if (conf === null || conf === undefined) return null;
        const pct = conf > 1 ? conf : conf * 100;
        return Math.round(pct);
    }

    function formatToxicity(tox: number | null): number {
        if (tox === null || tox === undefined) return 0;
        const pct = tox > 1 ? tox : tox * 100;
        return Math.round(pct);
    }

    function formatCost(microcents: number): string {
        const dollars = (microcents || 0) / 1000000;
        if (dollars === 0) return '$0.00';
        if (dollars < 0.01) return `$${dollars.toFixed(4)}`;
        return `$${dollars.toFixed(2)}`;
    }

    function getConfidenceColor(score: number) {
        if (score >= 70) return 'text-[var(--accent-responsibility)] bg-[color-mix(in_srgb,var(--accent-responsibility)_15%,transparent)]';
        if (score >= 40) return 'text-[var(--accent-warning)] bg-[color-mix(in_srgb,var(--accent-warning)_15%,transparent)]';
        return 'text-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_15%,transparent)]';
    }

    function getToxicityColor(score: number) {
        // INVERT the scale: low is green, high is red
        if (score >= 60) return 'text-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_15%,transparent)]';
        if (score >= 30) return 'text-[var(--accent-warning)] bg-[color-mix(in_srgb,var(--accent-warning)_15%,transparent)]';
        return 'text-[var(--accent-responsibility)] bg-[color-mix(in_srgb,var(--accent-responsibility)_15%,transparent)]';
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

<div class="p-8 max-w-7xl mx-auto animate-fade-up">
    <header class="mb-8 flex justify-between items-end">
        <div>
            <h1 class="font-serif text-5xl mb-2 text-primary tracking-wide">Request Explorer</h1>
            <p class="text-secondary font-sans text-lg">Investigate individual queries, endpoints, and scores.</p>
        </div>
    </header>

    <!-- Filters -->
    <div class="flex flex-wrap gap-6 mb-8 p-4 bg-surface rounded-lg border border-[var(--bg-elevated)]">
        <!-- Endpoints -->
        <div class="flex flex-col gap-2">
            <span class="text-xs uppercase tracking-widest text-tertiary font-bold">Endpoints</span>
            <div class="flex gap-2">
                {#each ['/v1/chat/completions', '/v1/playground'] as ep}
                    <button 
                        class="px-3 py-1 rounded-full text-xs font-mono font-medium border transition-colors {selectedEndpoints.has(ep) ? 'bg-elevated text-primary border-primary' : 'border-[var(--bg-elevated)] text-secondary hover:text-primary hover:border-secondary'}"
                        onclick={() => toggleEndpoint(ep)}
                    >
                        {ep}
                    </button>
                {/each}
            </div>
        </div>

        <!-- Models -->
        <div class="flex flex-col gap-2">
            <span class="text-xs uppercase tracking-widest text-tertiary font-bold">Model</span>
            <select 
                class="bg-base border border-elevated text-primary text-sm rounded-md px-3 py-1.5 font-mono focus:outline-none focus:ring-1 focus:ring-accent-perf"
                bind:value={selectedModel}
            >
                {#each availableModels as model}
                    <option value={model}>{model}</option>
                {/each}
            </select>
        </div>
    </div>

    {#if loading}
        <div class="text-center py-12 text-secondary font-mono">Loading requests...</div>
    {:else if error}
        <div class="text-center py-12 text-[var(--accent-danger)] font-mono">Error: {error}</div>
    {:else}
        <!-- Table -->
        <div class="w-full bg-surface border border-[var(--bg-elevated)] rounded-lg overflow-hidden">
            <table class="w-full text-left border-collapse">
                <thead>
                    <tr class="bg-elevated border-b border-[var(--bg-hover)] text-xs uppercase tracking-widest text-secondary font-sans">
                        {#each [
                            { key: 'timestamp', label: 'Time' },
                            { key: 'endpoint', label: 'Endpoint' },
                            { key: 'model', label: 'Model' },
                            { key: 'confidence', label: 'Confidence' },
                            { key: 'toxicity', label: 'Toxicity' },
                            { key: 'nli', label: 'Hallucination' },
                            { key: 'cost_microcents', label: 'Cost' }
                        ] as col}
                            <th 
                                class="px-4 py-3 cursor-pointer hover:text-primary transition-colors {['confidence', 'toxicity', 'cost_microcents'].includes(col.key) ? 'text-right' : ''}"
                                onclick={() => handleSort(col.key as any)}
                            >
                                <div class="flex items-center gap-1 {['confidence', 'toxicity', 'cost_microcents'].includes(col.key) ? 'justify-end' : ''}">
                                    {col.label}
                                    {#if sortCol === col.key}
                                        <span class="text-[10px] text-accent-perf">{sortAsc ? '▲' : '▼'}</span>
                                    {/if}
                                </div>
                            </th>
                        {/each}
                    </tr>
                </thead>
                <tbody>
                    {#if filteredRequests.length === 0}
                        <tr>
                            <td colspan="7" class="px-4 py-12 text-center text-secondary font-mono">No requests found matching criteria.</td>
                        </tr>
                    {/if}
                    {#each filteredRequests as req}
                        {@const conf = formatConfidence(req.confidence)}
                        {@const tox = formatToxicity(req.toxicity)}
                        <tr 
                            class="border-b border-[var(--bg-elevated)] hover:bg-hover cursor-pointer transition-colors group"
                            onclick={() => goto(`/requests/${req.id}`)}
                        >
                            <td class="px-4 py-3 font-mono text-xs text-secondary group-hover:text-primary transition-colors">{formatTime(req.timestamp)}</td>
                            <td class="px-4 py-3 font-mono text-xs text-tertiary truncate max-w-[200px]">
                                <div class="flex items-center gap-1.5">
                                    <span>{req.endpoint}</span>
                                    {#if req.cached}
                                        <span class="inline-flex items-center gap-0.5 px-1.5 py-0.2 rounded bg-accent-perf/20 text-accent-perf text-[10px] font-bold uppercase tracking-wider">
                                            <Zap class="w-2.5 h-2.5 fill-accent-perf" /> Cached
                                        </span>
                                    {/if}
                                </div>
                            </td>
                            <td class="px-4 py-3 font-mono text-xs text-secondary">{req.model}</td>
                            <td class="px-4 py-3 text-right">
                                {#if conf !== null}
                                    <span class="inline-block px-2 py-0.5 rounded text-xs font-mono {getConfidenceColor(conf)}">
                                        {conf}%
                                    </span>
                                {:else}
                                    <span class="text-secondary">-</span>
                                {/if}
                            </td>
                            <td class="px-4 py-3 text-right">
                                <span class="inline-block px-2 py-0.5 rounded text-xs font-mono {getToxicityColor(tox)}">
                                    {tox}%
                                </span>
                            </td>
                            <td class="px-4 py-3">
                                {#if req.nli}
                                    <span class="inline-block px-2 py-1 rounded-sm text-[10px] uppercase font-bold font-mono border {getNliStyle(req.nli.label)}">
                                        {req.nli.label}
                                    </span>
                                {:else}
                                    <span class="text-secondary font-mono text-xs">-</span>
                                {/if}
                            </td>
                            <td class="px-4 py-3 font-mono text-xs text-right text-secondary">
                                {formatCost(req.cost_microcents)}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
            
            <!-- Pagination -->
            <div class="p-4 bg-elevated flex justify-between items-center text-sm font-mono text-secondary">
                <div>
                    Total {totalRequests} | Page {currentPage} of {totalPages || 1}
                </div>
                <div class="flex gap-2">
                    <button 
                        class="px-3 py-1 bg-surface border border-hover rounded hover:text-primary disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={currentPage === 1}
                        onclick={() => currentPage--}
                    >
                        Prev
                    </button>
                    <button 
                        class="px-3 py-1 bg-surface border border-hover rounded hover:text-primary disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={currentPage === totalPages || totalPages === 0}
                        onclick={() => currentPage++}
                    >
                        Next
                    </button>
                </div>
            </div>
        </div>
    {/if}
</div>
