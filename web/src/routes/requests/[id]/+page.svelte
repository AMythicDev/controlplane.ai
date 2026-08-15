<script lang="ts">
    import { page } from '$app/stores';
    import { mockRequests, mockPolicies } from '$lib/data/dummy';
    import ScoreGauge from '$lib/components/ui/ScoreGauge.svelte';
    import TimelineBar from '$lib/components/ui/TimelineBar.svelte';
    
    let requestId = $derived($page.params.id);
    let request = $derived(mockRequests.find(r => r.id === requestId));
    let policy = $derived(request?.policyTriggered ? mockPolicies.find(p => p.name === request?.policyTriggered) : null);

    function formatTime(isoStr: string) {
        const d = new Date(isoStr);
        return d.toLocaleString();
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

    function formatText(text: string, isEdited: boolean) {
        if (!text) return '';
        // If action is edit, simulate redacted portions by wrapping dummy text if we want, or just wrapping [REDACTED] if it exists
        // Since dummy data might not have [REDACTED], we'll just return as is or do a simple replace if we inserted it
        return text.replace(/\[REDACTED\]/g, '<span class="bg-red-900/50 text-red-400 px-1 rounded border border-red-500/50 font-bold">[REDACTED]</span>');
    }
</script>

{#if !request}
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
                    <span class="inline-block px-3 py-1 rounded text-xs font-bold font-mono border uppercase tracking-wider {getActionStyle(request.action)}">
                        {request.action}
                    </span>
                </div>
                <div class="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm font-mono text-secondary mt-4">
                    <div class="flex items-center gap-2"><span class="text-tertiary">TIME:</span> {formatTime(request.timestamp)}</div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">USER:</span> {request.userName} ({request.userId})</div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">MODEL:</span> <span class="text-accent-perf">{request.model}</span></div>
                    <div class="flex items-center gap-2"><span class="text-tertiary">ENDPOINT:</span> {request.endpoint}</div>
                </div>
            </div>
            
            <div class="flex gap-4 items-center bg-surface p-3 rounded border border-elevated">
                <div class="text-center px-3 border-r border-elevated last:border-0">
                    <div class="text-xs text-tertiary mb-1">IN</div>
                    <div class="font-mono text-primary">{request.tokens.input}</div>
                </div>
                <div class="text-center px-3 border-r border-elevated last:border-0">
                    <div class="text-xs text-tertiary mb-1">OUT</div>
                    <div class="font-mono text-primary">{request.tokens.output}</div>
                </div>
                <div class="text-center px-3">
                    <div class="text-xs text-tertiary mb-1">CACHE</div>
                    <div class="font-mono {request.cacheHit ? 'text-[var(--accent-responsibility)]' : 'text-secondary'}">
                        {request.cacheHit ? 'HIT' : 'MISS'}
                    </div>
                </div>
            </div>
        </header>

        <!-- Evaluation Scores -->
        <section class="mb-12">
            <h2 class="font-serif text-3xl text-primary mb-6">Evaluation</h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="bg-surface p-6 rounded-xl border border-elevated flex items-center justify-center">
                    <ScoreGauge score={request.performanceScore} label="Performance" dimension="performance" />
                </div>
                <div class="bg-surface p-6 rounded-xl border border-elevated flex items-center justify-center">
                    <!-- Deriving a 0-100 score from cost for visual consistency, lower cost = higher score -->
                    <ScoreGauge score={Math.max(0, 100 - (request.costAmount * 1000))} label="Cost Efficiency" dimension="cost" />
                </div>
                <div class="bg-surface p-6 rounded-xl border border-elevated flex items-center justify-center">
                    <ScoreGauge score={request.responsibilityScore} label="Responsibility" dimension="responsibility" />
                </div>
            </div>
        </section>

        <!-- Action Log -->
        {#if request.action !== 'pass'}
            <section class="mb-12">
                <h2 class="font-serif text-3xl text-primary mb-6">Intervention Log</h2>
                <div class="bg-surface border border-elevated rounded-xl p-6 relative overflow-hidden">
                    <!-- Decorative background glow -->
                    <div class="absolute -right-20 -top-20 w-64 h-64 rounded-full blur-3xl opacity-10 pointer-events-none"
                         style="background-color: {getActionStyle(request.action).match(/var\((.*?)\)/)?.[0] || 'var(--accent-warning)'}">
                    </div>
                    
                    <div class="relative z-10 flex flex-col md:flex-row gap-6 md:items-center">
                        <div class="flex-shrink-0 flex items-center justify-center w-16 h-16 rounded-full bg-elevated border-2"
                             style="border-color: {getActionStyle(request.action).match(/var\((.*?)\)/)?.[0] || 'var(--accent-warning)'}">
                            <span class="text-2xl">{request.action === 'block' ? '🛑' : request.action === 'edit' ? '✂️' : request.action === 'flag' ? '⚠️' : '↗️'}</span>
                        </div>
                        <div>
                            <h3 class="text-xl font-mono text-primary mb-2 uppercase tracking-wide">Action: {request.action}</h3>
                            {#if policy}
                                <p class="text-secondary font-sans mb-1">
                                    Triggered by policy: <span class="text-primary font-medium">{policy.name}</span>
                                </p>
                                <p class="text-sm text-tertiary font-mono">
                                    {policy.description}
                                </p>
                            {:else}
                                <p class="text-secondary font-sans">Triggered by system guardrails.</p>
                            {/if}
                        </div>
                    </div>
                </div>
            </section>
        {/if}

        <!-- Check Pipeline -->
        <section class="mb-12">
            <h2 class="font-serif text-3xl text-primary mb-6">Execution Pipeline</h2>
            <TimelineBar checks={request.checks} />
        </section>

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
                        {@html formatText(request.prompt, false)}
                    </div>
                </div>

                <!-- Response -->
                <div class="bg-surface border border-elevated rounded-xl overflow-hidden flex flex-col">
                    <div class="bg-elevated px-4 py-3 border-b border-hover flex justify-between items-center">
                        <span class="font-mono text-xs font-bold text-tertiary uppercase tracking-widest">Response</span>
                        {#if request.action === 'edit'}
                            <span class="text-[10px] bg-accent-cost text-base px-2 py-0.5 rounded font-bold font-mono">MODIFIED</span>
                        {/if}
                    </div>
                    <div class="p-6 font-mono text-sm text-primary whitespace-pre-wrap flex-grow bg-[#111] relative">
                        <div class="noise-overlay" style="opacity: 0.05"></div>
                        {@html formatText(request.response, request.action === 'edit')}
                    </div>
                </div>
            </div>
        </section>
    </div>
{/if}
