<script lang="ts">
  import './layout.css';
  import FloatingNav from '$lib/components/layout/FloatingNav.svelte';
  import { page } from '$app/stores';

  let { children } = $props<{ children: any }>();

  let pageTitle = $derived.by(() => {
    const path = $page.url.pathname;
    if (path === '/') return 'Mission Control';
    if (path.startsWith('/requests')) return 'Request Telemetry';
    if (path.startsWith('/policies')) return 'Policy Engine';
    if (path.startsWith('/review')) return 'Review Queue';
    if (path.startsWith('/settings')) return 'System Settings';
    return 'ControlPlane';
  });
</script>

<svelte:head>
  <title>ControlPlane.ai - {pageTitle}</title>
</svelte:head>

<div class="noise-overlay"></div>

<div class="flex h-screen w-full overflow-hidden bg-gradient-to-b from-[#1A0A26] to-[#0D001A] relative text-primary">
  <FloatingNav />
  
  <div class="flex-1 flex flex-col relative z-10 h-full overflow-y-auto overflow-x-hidden pt-28 pb-12 px-6 lg:px-12">
    <main class="flex-1 max-w-[1600px] mx-auto w-full">
      {@render children()}
    </main>
  </div>
</div>
