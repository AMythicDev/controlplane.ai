<script lang="ts">
  import { mockRequests } from '$lib/data/dummy';
  import ScoreBadge from '$lib/components/ui/ScoreBadge.svelte';
  import ActionBadge from '$lib/components/ui/ActionBadge.svelte';
  import { onMount } from 'svelte';
  
  let requests = $derived(mockRequests.slice(0, 15));
  
  function formatTime(iso: string) {
    const d = new Date(iso);
    return d.toLocaleTimeString([], { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' });
  }

  let highlightIndex = $state(-1);
  
  onMount(() => {
    const interval = setInterval(() => {
      highlightIndex = 0;
      setTimeout(() => highlightIndex = -1, 1000);
    }, 5000);
    return () => clearInterval(interval);
  });
</script>

<div class="flex flex-col h-full bg-[var(--bg-surface)] rounded-2xl border border-white/5 overflow-hidden shadow-lg animate-slide-in">
  <div class="p-5 border-b border-white/5 flex items-center justify-between bg-[var(--bg-elevated)]/30">
    <div class="flex items-center gap-3">
      <div class="w-2 h-2 rounded-full bg-[var(--accent-danger)] animate-pulse"></div>
      <h3 class="font-sans font-medium text-[var(--text-primary)] text-sm uppercase tracking-widest">Live Stream</h3>
    </div>
    <div class="text-xs font-mono text-[var(--text-secondary)]">42 req/s</div>
  </div>
  
  <div class="flex-1 overflow-y-auto">
    <table class="w-full text-left border-collapse">
      <thead class="sticky top-0 bg-[var(--bg-surface)]/95 backdrop-blur z-10 border-b border-white/5 text-[10px] uppercase tracking-widest text-[var(--text-tertiary)] font-sans">
        <tr>
          <th class="px-4 py-3 font-medium">Time</th>
          <th class="px-4 py-3 font-medium">Endpoint</th>
          <th class="px-4 py-3 font-medium">Perf</th>
          <th class="px-4 py-3 font-medium">Resp</th>
          <th class="px-4 py-3 font-medium">Cost</th>
          <th class="px-4 py-3 font-medium text-right">Action</th>
        </tr>
      </thead>
      <tbody class="text-sm font-mono">
        {#each requests as req, i}
          <tr 
            class="border-b border-white/5 hover:bg-[var(--bg-elevated)]/50 transition-colors cursor-pointer group"
            class:bg-[var(--bg-elevated)]={i === highlightIndex}
          >
            <td class="px-4 py-3 text-[var(--text-secondary)] whitespace-nowrap">{formatTime(req.timestamp)}</td>
            <td class="px-4 py-3 text-[var(--text-primary)] truncate max-w-[100px]" title={req.endpoint}>{req.endpoint}</td>
            <td class="px-4 py-3">
              <ScoreBadge score={req.performanceScore} size="sm" dimension="performance" />
            </td>
            <td class="px-4 py-3">
              <ScoreBadge score={req.responsibilityScore} size="sm" dimension="responsibility" />
            </td>
            <td class="px-4 py-3 text-[var(--text-secondary)]">${req.costAmount.toFixed(4)}</td>
            <td class="px-4 py-3 text-right">
              <ActionBadge action={req.action} />
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
