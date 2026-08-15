<script lang="ts">
    import { mockRequests } from '$lib/data/dummy';
    import { goto } from '$app/navigation';
    import type { Request } from '$lib/types';

    let selectedActions = $state<Set<string>>(new Set());
    let selectedScore = $state<string>('All');
    let selectedModel = $state<string>('All');
    let sortCol = $state<keyof Request>('timestamp');
    let sortAsc = $state<boolean>(false);
    let currentPage = $state(1);
    const ITEMS_PER_PAGE = 15;

    const availableModels = ['All', ...new Set(mockRequests.map(r => r.model))];

    function toggleAction(action: string) {
        if (selectedActions.has(action)) {
            selectedActions.delete(action);
        } else {
            selectedActions.add(action);
        }
        currentPage = 1;
    }

    let filteredRequests = $derived(
        mockRequests.filter(req => {
            if (selectedActions.size > 0 && !selectedActions.has(req.action)) return false;
            if (selectedModel !== 'All' && req.model !== selectedModel) return false;
            
            if (selectedScore !== 'All') {
                const minScore = Math.min(req.performanceScore, req.responsibilityScore);
                if (selectedScore === 'Critical' && minScore >= 40) return false;
                if (selectedScore === 'Warning' && (minScore < 40 || minScore >= 70)) return false;
                if (selectedScore === 'Good' && minScore < 70) return false;
            }
            return true;
        }).sort((a, b) => {
            const valA = a[sortCol];
            const valB = b[sortCol];
            if (valA < valB) return sortAsc ? -1 : 1;
            if (valA > valB) return sortAsc ? 1 : -1;
            return 0;
        })
    );

    let paginatedRequests = $derived(
        filteredRequests.slice((currentPage - 1) * ITEMS_PER_PAGE, currentPage * ITEMS_PER_PAGE)
    );
    let totalPages = $derived(Math.ceil(filteredRequests.length / ITEMS_PER_PAGE));

    function handleSort(col: keyof Request) {
        if (sortCol === col) {
            sortAsc = !sortAsc;
        } else {
            sortCol = col;
            sortAsc = false; // default desc for new sort
        }
        currentPage = 1;
    }

    function formatTime(isoStr: string) {
        const d = new Date(isoStr);
        return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    }

    function getScoreColor(score: number) {
        if (score >= 70) return 'text-[var(--accent-responsibility)] bg-[color-mix(in_srgb,var(--accent-responsibility)_15%,transparent)]';
        if (score >= 40) return 'text-[var(--accent-warning)] bg-[color-mix(in_srgb,var(--accent-warning)_15%,transparent)]';
        return 'text-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_15%,transparent)]';
    }

    function getActionStyle(action: string) {
        switch(action) {
            case 'pass': return 'text-[var(--accent-responsibility)] border-[var(--accent-responsibility)] bg-[color-mix(in_srgb,var(--accent-responsibility)_10%,transparent)]';
            case 'block': return 'text-[var(--accent-danger)] border-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_10%,transparent)]';
            case 'escalate': return 'text-[var(--accent-danger)] border-[var(--accent-danger)] bg-[color-mix(in_srgb,var(--accent-danger)_10%,transparent)]';
            case 'flag': return 'text-[var(--accent-warning)] border-[var(--accent-warning)] bg-[color-mix(in_srgb,var(--accent-warning)_10%,transparent)]';
            case 'edit': return 'text-[var(--accent-cost)] border-[var(--accent-cost)] bg-[color-mix(in_srgb,var(--accent-cost)_10%,transparent)]';
            default: return 'text-secondary border-secondary bg-surface';
        }
    }
</script>

<div class="p-8 max-w-7xl mx-auto animate-fade-up">
    <header class="mb-8 flex justify-between items-end">
        <div>
            <h1 class="font-serif text-5xl mb-2 text-primary tracking-wide">Request Explorer</h1>
            <p class="text-secondary font-sans text-lg">Investigate individual queries, intercepts, and scores.</p>
        </div>
    </header>

    <!-- Filters -->
    <div class="flex flex-wrap gap-6 mb-8 p-4 bg-surface rounded-lg border border-[var(--bg-elevated)]">
        <!-- Actions -->
        <div class="flex flex-col gap-2">
            <span class="text-xs uppercase tracking-widest text-tertiary font-bold">Actions</span>
            <div class="flex gap-2">
                {#each ['pass', 'flag', 'edit', 'block', 'escalate'] as act}
                    <button 
                        class="px-3 py-1 rounded-full text-xs font-mono font-medium border transition-colors {selectedActions.has(act) ? getActionStyle(act) : 'border-[var(--bg-elevated)] text-secondary hover:text-primary hover:border-secondary'}"
                        onclick={() => toggleAction(act)}
                    >
                        {act.toUpperCase()}
                    </button>
                {/each}
            </div>
        </div>

        <!-- Scores -->
        <div class="flex flex-col gap-2">
            <span class="text-xs uppercase tracking-widest text-tertiary font-bold">Score Range</span>
            <div class="flex gap-2">
                {#each ['All', 'Good', 'Warning', 'Critical'] as sc}
                    <button 
                        class="px-3 py-1 rounded-full text-xs font-mono font-medium border transition-colors {selectedScore === sc ? 'bg-elevated text-primary border-primary' : 'border-[var(--bg-elevated)] text-secondary hover:text-primary hover:border-secondary'}"
                        onclick={() => { selectedScore = sc; currentPage = 1; }}
                    >
                        {sc}
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
                onchange={() => currentPage = 1}
            >
                {#each availableModels as model}
                    <option value={model}>{model}</option>
                {/each}
            </select>
        </div>
    </div>

    <!-- Table -->
    <div class="w-full bg-surface border border-[var(--bg-elevated)] rounded-lg overflow-hidden">
        <table class="w-full text-left border-collapse">
            <thead>
                <tr class="bg-elevated border-b border-[var(--bg-hover)] text-xs uppercase tracking-widest text-secondary font-sans">
                    {#each [
                        { key: 'timestamp', label: 'Time' },
                        { key: 'userId', label: 'User' },
                        { key: 'endpoint', label: 'Endpoint' },
                        { key: 'model', label: 'Model' },
                        { key: 'performanceScore', label: 'Perf' },
                        { key: 'costAmount', label: 'Cost' },
                        { key: 'responsibilityScore', label: 'Resp' },
                        { key: 'action', label: 'Action' }
                    ] as col}
                        <th 
                            class="px-4 py-3 cursor-pointer hover:text-primary transition-colors {col.key === 'costAmount' || col.key.includes('Score') ? 'text-right' : ''}"
                            onclick={() => handleSort(col.key as keyof Request)}
                        >
                            <div class="flex items-center gap-1 {col.key === 'costAmount' || col.key.includes('Score') ? 'justify-end' : ''}">
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
                {#if paginatedRequests.length === 0}
                    <tr>
                        <td colspan="8" class="px-4 py-12 text-center text-secondary font-mono">No requests found matching criteria.</td>
                    </tr>
                {/if}
                {#each paginatedRequests as req}
                    <tr 
                        class="border-b border-[var(--bg-elevated)] hover:bg-hover cursor-pointer transition-colors group"
                        onclick={() => goto(`/requests/${req.id}`)}
                    >
                        <td class="px-4 py-3 font-mono text-xs text-secondary group-hover:text-primary transition-colors">{formatTime(req.timestamp)}</td>
                        <td class="px-4 py-3 font-sans text-sm text-primary truncate max-w-[100px]">{req.userName}</td>
                        <td class="px-4 py-3 font-mono text-xs text-tertiary truncate max-w-[120px]">{req.endpoint}</td>
                        <td class="px-4 py-3 font-mono text-xs text-secondary">{req.model}</td>
                        <td class="px-4 py-3 text-right">
                            <span class="inline-block px-2 py-0.5 rounded text-xs font-mono {getScoreColor(req.performanceScore)}">
                                {Math.round(req.performanceScore)}
                            </span>
                        </td>
                        <td class="px-4 py-3 font-mono text-xs text-right text-secondary">
                            ${req.costAmount.toFixed(4)}
                        </td>
                        <td class="px-4 py-3 text-right">
                            <span class="inline-block px-2 py-0.5 rounded text-xs font-mono {getScoreColor(req.responsibilityScore)}">
                                {Math.round(req.responsibilityScore)}
                            </span>
                        </td>
                        <td class="px-4 py-3">
                            <span class="inline-block px-2 py-1 rounded-sm text-[10px] uppercase font-bold font-mono border {getActionStyle(req.action)}">
                                {req.action}
                            </span>
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
        
        <!-- Pagination -->
        <div class="p-4 bg-elevated flex justify-between items-center text-sm font-mono text-secondary">
            <div>
                Showing {(currentPage - 1) * ITEMS_PER_PAGE + 1} to {Math.min(currentPage * ITEMS_PER_PAGE, filteredRequests.length)} of {filteredRequests.length}
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
</div>
