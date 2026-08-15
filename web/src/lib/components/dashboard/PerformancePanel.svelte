<script lang="ts">
  import { mockDashboardStats, mockTimeSeriesPerf } from '$lib/data/dummy';
  import Sparkline from '$lib/components/ui/Sparkline.svelte';
  
  let data = $derived(mockTimeSeriesPerf.map(d => d.value));
</script>

<div class="relative bg-[var(--bg-surface)] rounded-2xl border border-white/5 p-6 md:p-8 flex flex-col gap-6 overflow-hidden shadow-lg">
  <div class="absolute inset-0 opacity-[0.015] pointer-events-none mix-blend-overlay" style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noise%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noise)%22/%3E%3C/svg%3E');"></div>
  
  <div class="relative z-10 flex flex-col sm:flex-row sm:items-start justify-between gap-4">
    <div>
      <h2 class="font-serif text-3xl text-[var(--text-primary)] italic tracking-wide">Performance</h2>
      <p class="text-sm text-[var(--text-secondary)] font-sans mt-1">System latency & response quality</p>
    </div>
    <div class="flex flex-col items-end">
      <div class="text-5xl md:text-6xl lg:text-7xl font-mono font-bold text-[var(--accent-performance)] tracking-tighter leading-none">
        {mockDashboardStats.avgPerformance.toFixed(1)}
      </div>
      <div class="mt-2 opacity-80 hover:opacity-100 transition-opacity">
        <Sparkline data={data} color="var(--accent-performance)" width={120} height={32} />
      </div>
    </div>
  </div>

  <div class="relative z-10 grid grid-cols-3 gap-4 border-t border-white/5 pt-6 mt-2">
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Hallucination</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">1.2%</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Avg Confidence</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">0.94</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Citation Acc.</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">98.5%</div>
    </div>
  </div>

  <div class="relative z-10 bg-[var(--bg-elevated)]/50 rounded-xl p-4 border border-white/5 mt-2">
    <div class="text-xs uppercase tracking-widest text-[var(--text-secondary)] mb-3">Top Ungrounded Claims</div>
    <ul class="flex flex-col gap-3">
      <li class="flex items-start gap-3 text-sm">
        <span class="text-[var(--accent-warning)] font-mono mt-0.5">!</span>
        <span class="text-[var(--text-secondary)] font-mono">"Our SLA guarantees 100% uptime..."</span>
      </li>
      <li class="flex items-start gap-3 text-sm">
        <span class="text-[var(--accent-warning)] font-mono mt-0.5">!</span>
        <span class="text-[var(--text-secondary)] font-mono">"The API supports rate limits of 10k/sec..."</span>
      </li>
    </ul>
  </div>
</div>
