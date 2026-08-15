<script lang="ts">
  let { label, value, trend, trendLabel, accentColor = 'var(--accent-performance)' } = $props<{
    label: string;
    value: string;
    trend: number;
    trendLabel: string;
    accentColor?: string;
  }>();

  let isPositive = $derived(trend >= 0);
</script>

<div class="relative overflow-hidden bg-[var(--bg-surface)] rounded-xl border border-white/5 p-5 animate-fade-up shadow-[inset_0_1px_0_rgba(255,255,255,0.05)] backdrop-blur-sm">
  <div class="absolute inset-0 opacity-[0.02] mix-blend-overlay pointer-events-none" style="background-image: url('data:image/svg+xml,%3Csvg viewBox=%220 0 200 200%22 xmlns=%22http://www.w3.org/2000/svg%22%3E%3Cfilter id=%22noise%22%3E%3CfeTurbulence type=%22fractalNoise%22 baseFrequency=%220.65%22 numOctaves=%223%22 stitchTiles=%22stitch%22/%3E%3C/filter%3E%3Crect width=%22100%25%22 height=%22100%25%22 filter=%22url(%23noise)%22/%3E%3C/svg%3E');"></div>
  
  <div class="relative z-10 flex flex-col gap-1">
    <div class="text-3xl sm:text-4xl font-mono font-medium tracking-tight" style="color: {accentColor}">
      {value}
    </div>
    
    <div class="text-xs uppercase tracking-widest text-[var(--text-secondary)] font-sans font-medium mt-1">
      {label}
    </div>
    
    <div class="flex items-center gap-2 mt-3">
      <span class="text-xs font-mono font-medium px-1.5 py-0.5 rounded bg-black/20 {isPositive ? 'text-[var(--accent-responsibility)]' : 'text-[var(--accent-danger)]'}">
        {isPositive ? '↑' : '↓'} {Math.abs(trend)}%
      </span>
      <span class="text-xs text-[var(--text-tertiary)] font-sans">{trendLabel}</span>
    </div>
  </div>
</div>
