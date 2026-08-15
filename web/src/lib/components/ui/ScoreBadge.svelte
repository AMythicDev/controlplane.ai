<script lang="ts">
  let { score, size = 'md', dimension = 'neutral' } = $props<{
    score: number;
    size?: 'sm' | 'md' | 'lg';
    dimension?: 'performance' | 'cost' | 'responsibility' | 'neutral';
  }>();

  let sizeClasses = $derived(
    size === 'sm' ? 'w-6 h-6 text-[10px]' :
    size === 'lg' ? 'w-12 h-12 text-lg' :
    'w-8 h-8 text-xs'
  );

  let colorClasses = $derived.by(() => {
    let baseColor = '';
    if (score >= 90) {
      if (dimension === 'performance') baseColor = 'text-[var(--accent-performance)] border-[var(--accent-performance)]/50';
      else if (dimension === 'cost') baseColor = 'text-[var(--accent-cost)] border-[var(--accent-cost)]/50';
      else if (dimension === 'responsibility') baseColor = 'text-[var(--accent-responsibility)] border-[var(--accent-responsibility)]/50';
      else baseColor = 'text-[var(--text-primary)] border-[var(--text-primary)]/50';
    } else if (score >= 70) {
      if (dimension === 'performance') baseColor = 'text-[var(--accent-performance)]/70 border-[var(--accent-performance)]/30';
      else if (dimension === 'cost') baseColor = 'text-[var(--accent-cost)]/70 border-[var(--accent-cost)]/30';
      else if (dimension === 'responsibility') baseColor = 'text-[var(--accent-responsibility)]/70 border-[var(--accent-responsibility)]/30';
      else baseColor = 'text-[var(--text-primary)]/70 border-[var(--text-primary)]/30';
    } else if (score >= 40) {
      baseColor = 'text-[var(--accent-warning)] border-[var(--accent-warning)]/50';
    } else {
      baseColor = 'text-[var(--accent-danger)] border-[var(--accent-danger)]/50';
    }
    return baseColor;
  });

  let shadowClass = $derived(
    score < 40 ? 'animate-[pulseGlow_2s_infinite]' :
    'hover:shadow-[0_0_8px_currentColor] transition-shadow duration-300'
  );
</script>

<div class="flex items-center justify-center rounded-full border border-solid font-mono font-medium {sizeClasses} {colorClasses} {shadowClass}" style="{score < 40 ? 'box-shadow: 0 0 10px var(--accent-danger);' : ''}">
  {Math.round(score)}
</div>
