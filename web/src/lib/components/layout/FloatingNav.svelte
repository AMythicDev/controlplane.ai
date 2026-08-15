<script lang="ts">
  import { page } from '$app/stores';
  import LayoutDashboard from 'lucide-svelte/icons/layout-dashboard';
  import Activity from 'lucide-svelte/icons/activity';
  import ShieldAlert from 'lucide-svelte/icons/shield-alert';
  import CheckSquare from 'lucide-svelte/icons/check-square';
  import Settings from 'lucide-svelte/icons/settings';
  import Bell from 'lucide-svelte/icons/bell';

  let currentPath = $derived($page.url.pathname);

  const navItems = [
    { href: '/', label: 'Observatory', icon: LayoutDashboard },
    { href: '/requests', label: 'Activity', icon: Activity },
    { href: '/policies', label: 'Directives', icon: ShieldAlert },
    { href: '/review', label: 'Review', icon: CheckSquare },
    { href: '/settings', label: 'Config', icon: Settings }
  ];
</script>

<div class="fixed top-6 left-1/2 -translate-x-1/2 z-50 transition-all duration-300">
  <nav class="flex items-center gap-1 p-2 bg-surface/40 backdrop-blur-xl border border-surface/50 rounded-full shadow-2xl shadow-black/50 overflow-hidden group">
    
    <!-- Logo area -->
    <div class="flex items-center gap-2 pr-6 pl-4 border-r border-elevated/50">
      <div class="w-6 h-6 rounded-md bg-gradient-to-br from-accent-perf to-accent-perf flex items-center justify-center shrink-0 shadow-[0_0_15px_rgba(var(--color-accent-perf),0.5)]">
        <div class="w-3 h-3 bg-base rounded-sm rotate-45 transform transition-transform duration-700 group-hover:rotate-[225deg]"></div>
      </div>
      <span class="font-display font-medium text-primary text-sm tracking-tight hidden sm:block">CP.ai</span>
    </div>

    <!-- Nav Items -->
    <div class="flex items-center px-2 gap-1">
      {#each navItems as item}
        <a 
          href={item.href}
          class="relative flex items-center gap-2 px-4 py-2 rounded-full transition-all duration-300 group/item overflow-hidden"
          class:bg-surface={currentPath === item.href}
          class:text-accent-perf={currentPath === item.href}
          class:text-primary-muted={currentPath !== item.href}
        >
          <!-- Active Background Indicator (Glow) -->
          {#if currentPath === item.href}
            <div class="absolute inset-0 bg-accent-perf/10 opacity-0 group-hover/item:opacity-100 transition-opacity"></div>
          {/if}

          <item.icon class="w-4 h-4 shrink-0 transition-transform duration-300 group-hover/item:scale-110" />
          
          <span class="text-xs font-mono font-medium tracking-wide uppercase transition-all duration-300 w-0 md:w-auto opacity-0 md:opacity-100 overflow-hidden md:overflow-visible">
            {item.label}
          </span>
        </a>
      {/each}
    </div>

    <!-- Actions -->
    <div class="flex items-center gap-3 pl-4 pr-2 border-l border-elevated/50">
      <button class="relative p-2 text-primary-muted hover:text-primary transition-colors rounded-full hover:bg-surface/50">
        <Bell class="w-4 h-4" />
        <span class="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-accent-danger rounded-full shadow-[0_0_5px_rgba(var(--color-accent-danger),1)] animate-pulse"></span>
      </button>
      <div class="w-8 h-8 rounded-full bg-surface border border-elevated/50 overflow-hidden flex items-center justify-center relative cursor-pointer hover:border-accent-perf/50 transition-colors">
        <!-- Generic avatar pattern -->
        <div class="absolute inset-0 bg-gradient-to-tr from-accent-cost/20 to-accent-perf/20"></div>
        <span class="font-mono text-xs text-primary relative z-10">AK</span>
      </div>
    </div>
  </nav>
</div>

<style>
  /* Optional specific styles */
</style>
