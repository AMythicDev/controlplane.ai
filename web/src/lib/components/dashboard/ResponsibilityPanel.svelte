<script lang="ts">
  import { mockDashboardStats, mockTimeSeriesResp } from '$lib/data/dummy';
  import Sparkline from '$lib/components/ui/Sparkline.svelte';
  
  let data = $derived(mockTimeSeriesResp.map(d => d.value));
</script>

<div class="relative bg-[var(--bg-surface)] rounded-2xl border border-white/5 p-6 md:p-8 flex flex-col gap-6 overflow-hidden shadow-lg">
  <div class="absolute inset-0 opacity-[0.015] pointer-events-none mix-blend-overlay" style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noise%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noise)%22/%3E%3C/svg%3E');"></div>
  
  <div class="relative z-10 flex flex-col sm:flex-row sm:items-start justify-between gap-4">
    <div>
      <h2 class="font-serif text-3xl text-[var(--text-primary)] italic tracking-wide">Responsibility</h2>
      <p class="text-sm text-[var(--text-secondary)] font-sans mt-1">Safety, bias & compliance</p>
    </div>
    <div class="flex flex-col items-end">
      <div class="text-5xl md:text-6xl lg:text-7xl font-mono font-bold text-[var(--accent-responsibility)] tracking-tighter leading-none">
        {mockDashboardStats.avgResponsibility.toFixed(1)}
      </div>
      <div class="mt-2 opacity-80 hover:opacity-100 transition-opacity">
        <Sparkline data={data} color="var(--accent-responsibility)" width={120} height={32} />
      </div>
    </div>
  </div>

  <div class="relative z-10 flex gap-2 flex-wrap border-t border-white/5 pt-6 mt-2">
    <div class="px-2.5 py-1 rounded-full bg-[var(--accent-responsibility)]/10 border border-[var(--accent-responsibility)]/20 text-[var(--accent-responsibility)] text-xs font-mono flex items-center gap-1.5">
      <span>✓</span> GDPR
    </div>
    <div class="px-2.5 py-1 rounded-full bg-[var(--accent-responsibility)]/10 border border-[var(--accent-responsibility)]/20 text-[var(--accent-responsibility)] text-xs font-mono flex items-center gap-1.5">
      <span>✓</span> HIPAA
    </div>
    <div class="px-2.5 py-1 rounded-full bg-[var(--accent-responsibility)]/10 border border-[var(--accent-responsibility)]/20 text-[var(--accent-responsibility)] text-xs font-mono flex items-center gap-1.5">
      <span>✓</span> SOC2
    </div>
  </div>

  <div class="relative z-10 grid grid-cols-3 gap-4 border-t border-white/5 pt-6 mt-2">
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">PII Redacted</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">142</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Toxicity Flags</div>
      <div class="font-mono text-lg text-[var(--text-primary)]">12</div>
    </div>
    <div>
      <div class="text-xs uppercase tracking-widest text-[var(--text-tertiary)] mb-1">Bias Status</div>
      <div class="font-mono text-lg text-[var(--accent-responsibility)]">PASS</div>
    </div>
  </div>
</div>
