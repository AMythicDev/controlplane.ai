<script lang="ts">
  import { page } from '$app/stores';

  let isExpanded = $state(false);

  const navItems = [
    { name: 'Dashboard', path: '/', icon: '<path d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>' },
    { name: 'Requests', path: '/requests', icon: '<path d="M14 2H6c-1.1 0-1.99.9-1.99 2L4 20c0 1.1.89 2 1.99 2H18c1.1 0 2-.9 2-2V8l-6-6zm2 16H8v-2h8v2zm0-4H8v-2h8v2zm-3-5V3.5L18.5 9H13z"/>' },
    { name: 'Policies', path: '/policies', icon: '<path d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"/>' },
    { name: 'Review Queue', path: '/review', icon: '<path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>' },
    { name: 'Settings', path: '/settings', icon: '<path d="M19.14 12.94c.04-.3.06-.61.06-.94 0-.32-.02-.64-.06-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.56-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.05.3-.09.63-.09.94s.02.64.06.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .43-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.49-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/>' }
  ];
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<!-- svelte-ignore a11y_mouse_events_have_key_events -->
<div 
  class="fixed left-0 top-0 h-full bg-surface border-r border-elevated transition-all duration-300 ease-in-out flex flex-col z-50 {isExpanded ? 'w-60' : 'w-[72px]'}"
  onmouseover={() => isExpanded = true}
  onmouseleave={() => isExpanded = false}
>
  <div class="h-16 flex items-center px-4 border-b border-elevated overflow-hidden shrink-0">
    <div class="w-10 h-10 rounded bg-elevated flex items-center justify-center shrink-0 border border-hover shadow-[0_0_15px_rgba(0,229,255,0.15)]">
      <span class="font-serif text-xl text-primary leading-none mt-1">CP</span>
    </div>
    {#if isExpanded}
      <span class="ml-3 font-serif text-xl tracking-wide text-primary whitespace-nowrap animate-slide-in">ControlPlane</span>
    {/if}
  </div>

  <nav class="flex-1 py-4 flex flex-col gap-2 overflow-y-auto overflow-x-hidden px-3">
    {#each navItems as item}
      {@const isActive = $page.url.pathname === item.path || ($page.url.pathname.startsWith(item.path) && item.path !== '/')}
      <a 
        href={item.path}
        class="relative flex items-center h-10 rounded-md transition-colors group {isActive ? 'bg-elevated text-primary' : 'text-secondary hover:text-primary hover:bg-hover'}"
      >
        {#if isActive}
          <div class="absolute left-0 top-1 bottom-1 w-1 bg-accent-perf rounded-r-md"></div>
        {/if}
        <div class="w-12 h-10 flex items-center justify-center shrink-0 {isActive ? 'text-accent-perf' : ''}">
          <svg viewBox="0 0 24 24" fill="currentColor" class="w-5 h-5">
            {@html item.icon}
          </svg>
        </div>
        {#if isExpanded}
          <span class="whitespace-nowrap font-medium text-sm animate-slide-in">{item.name}</span>
        {/if}
      </a>
    {/each}
  </nav>

  <div class="p-3 border-t border-elevated overflow-hidden shrink-0">
    <div class="flex items-center h-10 rounded-md hover:bg-hover cursor-pointer px-1 transition-colors">
      <div class="w-8 h-8 rounded-full bg-accent-perf text-base flex items-center justify-center shrink-0 font-bold ml-1 text-black">
        A
      </div>
      {#if isExpanded}
        <div class="ml-3 flex flex-col animate-slide-in">
          <span class="text-sm font-medium text-primary whitespace-nowrap">Admin User</span>
          <span class="text-xs text-tertiary whitespace-nowrap">admin@controlplane.ai</span>
        </div>
      {/if}
    </div>
  </div>
</div>
