<script lang="ts">
  let { action } = $props<{ action: 'pass' | 'flag' | 'edit' | 'block' | 'escalate' }>();

  let config = $derived.by(() => {
    switch (action) {
      case 'pass': return { icon: '✓', label: 'Pass', classes: 'text-[var(--accent-responsibility)] bg-[var(--accent-responsibility)]/10 border-[var(--accent-responsibility)]/20' };
      case 'flag': return { icon: '⚑', label: 'Flag', classes: 'text-[var(--accent-warning)] bg-[var(--accent-warning)]/10 border-[var(--accent-warning)]/20' };
      case 'edit': return { icon: '✎', label: 'Edit', classes: 'text-[var(--accent-performance)] bg-[var(--accent-performance)]/10 border-[var(--accent-performance)]/20' };
      case 'block': return { icon: '✕', label: 'Block', classes: 'text-[var(--accent-danger)] bg-[var(--accent-danger)]/10 border-[var(--accent-danger)]/20' };
      case 'escalate': return { icon: '⇧', label: 'Escalate', classes: 'text-[#d500f9] bg-[#d500f9]/10 border-[#d500f9]/20' };
    }
  });
</script>

<span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full border text-xs font-mono font-medium whitespace-nowrap {config.classes}">
  <span>{config.icon}</span>
  <span>{config.label}</span>
</span>
