<script lang="ts">
    import type { CheckResult } from '$lib/types';

    let { checks } = $props<{ checks: CheckResult[] }>();

    function getStatusColor(status: string) {
        switch(status) {
            case 'pass': return 'var(--accent-responsibility)';
            case 'fail': return 'var(--accent-danger)';
            case 'warn': return 'var(--accent-warning)';
            default: return 'var(--text-tertiary)';
        }
    }
</script>

<div class="flex items-center w-full bg-surface p-4 rounded border border-[var(--bg-elevated)]">
    {#each checks as check, i (check.name)}
        <div class="flex-1 flex flex-col group relative">
            <!-- Connecting line -->
            <div class="absolute top-1/2 left-0 w-full h-1 -translate-y-1/2 bg-elevated z-0"
                 class:rounded-l-full={i === 0}
                 class:rounded-r-full={i === checks.length - 1}>
            </div>
            
            <!-- Progress portion of line -->
            <div class="absolute top-1/2 left-0 h-1 -translate-y-1/2 z-0 transition-all duration-500"
                 style="width: 100%; background-color: {getStatusColor(check.status)}; opacity: 0.8;"
                 class:rounded-l-full={i === 0}
                 class:rounded-r-full={i === checks.length - 1}>
            </div>

            <!-- Node -->
            <div class="relative z-10 w-4 h-4 rounded-full mx-auto shadow-[0_0_10px_rgba(0,0,0,0.5)] transition-transform group-hover:scale-125"
                 style="background-color: {getStatusColor(check.status)}; border: 2px solid var(--bg-surface)">
            </div>

            <!-- Labels -->
            <div class="relative z-10 text-center mt-3">
                <div class="text-xs font-mono font-medium truncate px-1 text-primary">
                    {check.name}
                </div>
                <div class="text-[10px] text-tertiary font-mono mt-1">
                    {check.duration}ms
                </div>
            </div>

            <!-- Tooltip (CSS only) -->
            <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 w-48 p-2 bg-elevated text-xs border border-hover rounded shadow-xl opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-20">
                <div class="font-bold text-primary mb-1">{check.name} (L{check.layer})</div>
                <div class="text-secondary">{check.details}</div>
                <div class="mt-1 font-mono text-[10px]" style="color: {getStatusColor(check.status)}">
                    Status: {check.status.toUpperCase()} | Score: {check.score}
                </div>
            </div>
        </div>
    {/each}
</div>
