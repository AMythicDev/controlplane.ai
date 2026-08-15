<script lang="ts">
  import { mockDashboardStats, mockTimeSeriesCost, mockCostBreakdown } from '$lib/data/dummy';
  import Sparkline from '$lib/components/ui/Sparkline.svelte';
  
  let data = $derived(mockTimeSeriesCost.map(d => d.value));
  let dailyBudget = 500;
  let spentToday = 142.50;
  let progress = $derived((spentToday / dailyBudget) * 100);
</script>

<div class="relative bg-[var(--bg-surface)] rounded-2xl border border-white/5 p-6 md:p-8 flex flex-col gap-6 overflow-hidden shadow-lg">
  <div class="absolute inset-0 opacity-[0.015] pointer-events-none mix-blend-overlay" style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noise%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noise)%22/%3E%3C/svg%3E');"></div>
  
  <div class="relative z-10 flex flex-col sm:flex-row sm:items-start justify-between gap-4">
    <div>
      <h2 class="font-serif text-3xl text-[var(--text-primary)] italic tracking-wide">Cost</h2>
      <p class="text-sm text-[var(--text-secondary)] font-sans mt-1">Resource utilization & billing</p>
    </div>
    <div class="flex flex-col items-end">
      <div class="text-5xl md:text-6xl lg:text-7xl font-mono font-bold text-[var(--accent-cost)] tracking-tighter leading-none">
        ${mockDashboardStats.totalSpend.toLocaleString(undefined, {minimumFractionDigits: 0, maximumFractionDigits: 0})}
      </div>
      <div class="mt-2 opacity-80 hover:opacity-100 transition-opacity">
        <Sparkline data={data} color="var(--accent-cost)" width={120} height={32} />
      </div>
    </div>
  </div>

  <div class="relative z-10 border-t border-white/5 pt-6 mt-2">
    <div class="flex justify-between items-end mb-2">
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)]">Daily Budget</div>
      <div class="font-mono text-sm text-[var(--text-secondary)]"><span class="text-[var(--text-primary)]">${spentToday.toFixed(2)}</span> / ${dailyBudget}</div>
    </div>
    <div class="h-1.5 w-full bg-[var(--bg-elevated)] rounded-full overflow-hidden">
      <div class="h-full bg-[var(--accent-cost)] rounded-full transition-all duration-1000" style="width: {progress}%"></div>
    </div>
  </div>

  <div class="relative z-10 grid grid-cols-3 gap-4 border-t border-white/5 pt-6">
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Cache Hits</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">42.8%</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Token Eff.</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">89%</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Proj. Monthly</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">$4,250</div>
    </div>
  </div>

  <div class="relative z-10 bg-[var(--bg-elevated)]/50 rounded-xl p-4 border border-white/5 mt-2 flex flex-col gap-3">
    {#each mockCostBreakdown as item}
      <div class="flex flex-col gap-1.5">
        <div class="flex justify-between text-xs font-mono">
          <span class="text-[var(--text-secondary)]">{item.model}</span>
          <span class="text-[var(--text-primary)]">${item.amount.toFixed(0)}</span>
        </div>
        <div class="h-1 w-full bg-black/40 rounded-full overflow-hidden">
          <div class="h-full bg-[var(--accent-cost)]/70 rounded-full" style="width: {item.percentage}%"></div>
        </div>
      </div>
    {/each}
  </div>
</div>
