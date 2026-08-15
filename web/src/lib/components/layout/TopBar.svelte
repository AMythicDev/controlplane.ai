<script lang="ts">
  let { title, subtitle = '' } = $props<{ title: string; subtitle?: string }>();
  
  let timeStr = $state('');

  $effect(() => {
    const updateTime = () => {
      const now = new Date();
      timeStr = now.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false });
    };
    updateTime();
    const interval = setInterval(updateTime, 1000);
    return () => clearInterval(interval);
  });
</script>

<header class="h-16 sticky top-0 z-40 bg-base/80 backdrop-blur-md flex items-center justify-between px-8 w-full border-b border-transparent" style="border-image: linear-gradient(to right, var(--bg-elevated), transparent) 1; border-bottom: 1px solid;">
  <div class="flex flex-col justify-center">
    <h1 class="font-serif text-3xl text-primary tracking-wide leading-tight">{title}</h1>
    {#if subtitle}
      <span class="text-sm text-secondary font-sans">{subtitle}</span>
    {/if}
  </div>

  <div class="flex items-center gap-6">
    <div class="relative group">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 text-tertiary group-focus-within:text-accent-perf transition-colors">
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
      </div>
      <input 
        type="text" 
        placeholder="Search requests, policies..." 
        class="w-64 bg-surface border border-elevated rounded-full py-1.5 pl-10 pr-4 text-sm text-primary placeholder-tertiary focus:outline-none focus:border-accent-perf/50 focus:ring-1 focus:ring-accent-perf/50 transition-all"
      >
      <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
        <span class="text-[10px] text-tertiary font-mono bg-elevated px-1.5 py-0.5 rounded border border-hover">⌘K</span>
      </div>
    </div>

    <button class="relative text-secondary hover:text-primary transition-colors focus-ring rounded-full p-1">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-5 h-5">
        <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
        <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
      </svg>
      <span class="absolute top-0 right-0 w-2 h-2 bg-accent-danger rounded-full border border-base"></span>
    </button>

    <div class="font-mono text-sm text-accent-perf opacity-80 select-none">
      {timeStr}
    </div>
  </div>
</header>
