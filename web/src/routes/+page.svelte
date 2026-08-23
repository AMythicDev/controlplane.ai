<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    mockDashboardStats, 
    mockTimeSeriesPerf,
    mockTimeSeriesCost,
    mockTimeSeriesResp,
    mockRequests
  } from '$lib/data/dummy';
  import Sparkline from '$lib/components/ui/Sparkline.svelte';
  import ActionBadge from '$lib/components/ui/ActionBadge.svelte';
  import ScoreBadge from '$lib/components/ui/ScoreBadge.svelte';
  import { goto } from '$app/navigation';
  import ShieldAlert from 'lucide-svelte/icons/shield-alert';
  import AlertTriangle from 'lucide-svelte/icons/alert-triangle';
  import CheckCircle2 from 'lucide-svelte/icons/check-circle-2';
  import TrendingUp from 'lucide-svelte/icons/trending-up';
  import TrendingDown from 'lucide-svelte/icons/trending-down';
  import DollarSign from 'lucide-svelte/icons/dollar-sign';
  import Activity from 'lucide-svelte/icons/activity';

  // Formatters
  const formatCost = (val: number) => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 3 }).format(val);
  const formatNum = (val: number) => new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(val);

  // Time format for cards
  const formatTime = (ts: string) => {
    const d = new Date(ts);
    return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  let liveCost = $state<number | null>(null);

  onMount(async () => {
    try {
      const res = await fetch('/v1/cost');
      if (res.ok) {
        const data = await res.json();
        liveCost = data.cost_dollars;
      }
    } catch (e) {
      console.error('Failed to fetch live cost', e);
    }
  });
</script>

<div class="flex flex-col gap-24 relative z-10 w-full min-h-full pb-32">
  <!-- Radial gradient background behind hero -->
  <div class="absolute top-0 left-1/2 -translate-x-1/2 w-[80vw] h-[80vh] bg-[radial-gradient(ellipse_at_top,rgba(var(--color-elevated),0.15)_0%,transparent_70%)] pointer-events-none -z-10 rounded-full blur-3xl"></div>

  <!-- HERO ZONE -->
  <section class="flex flex-col items-center pt-8">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-0 w-full max-w-6xl relative">
      <!-- Vertical Separators -->
      <div class="hidden md:block absolute top-12 bottom-12 left-1/3 w-px bg-gradient-to-b from-transparent via-border to-transparent z-0"></div>
      <div class="hidden md:block absolute top-12 bottom-12 right-1/3 w-px bg-gradient-to-b from-transparent via-border to-transparent z-0"></div>

      <!-- Performance -->
      <div class="flex flex-col items-center px-4 md:px-12 relative z-10 group">
        <div class="text-[80px] lg:text-[110px] font-mono font-bold leading-none tracking-tighter text-accent-perf tabular-nums hover:scale-105 transition-transform duration-500 cursor-default shadow-accent-perf/20 drop-shadow-[0_0_25px_rgba(var(--color-accent-perf),0.4)]">
          {formatNum(mockDashboardStats.avgPerformance)}
        </div>
        <div class="uppercase tracking-[0.2em] text-xs text-secondary mt-4 mb-8 font-medium">Performance Score</div>
        <div class="w-full h-8 opacity-70 group-hover:opacity-100 transition-opacity">
          <Sparkline data={mockTimeSeriesPerf} color="var(--color-accent-perf)" />
        </div>
      </div>

      <!-- Cost -->
      <div class="flex flex-col items-center px-4 md:px-12 relative z-10 group">
        <div class="text-[80px] lg:text-[110px] font-mono font-bold leading-none tracking-tighter text-accent-cost tabular-nums hover:scale-105 transition-transform duration-500 cursor-default drop-shadow-[0_0_25px_rgba(var(--color-accent-cost),0.4)] flex items-start">
          <span class="text-4xl lg:text-6xl mt-4 opacity-50">$</span>
          {#if liveCost !== null}
            <span>{formatNum(liveCost)}</span>
          {:else}
            <span>{formatNum(mockDashboardStats.totalSpend / 1000)}<span class="text-4xl lg:text-6xl text-accent-cost/50">k</span></span>
          {/if}
        </div>
        <div class="uppercase tracking-[0.2em] text-xs text-secondary mt-4 mb-8 font-medium">Total Spend</div>
        <div class="w-full h-8 opacity-70 group-hover:opacity-100 transition-opacity">
          <Sparkline data={mockTimeSeriesCost} color="var(--color-accent-cost)" />
        </div>
      </div>

      <!-- Responsibility -->
      <div class="flex flex-col items-center px-4 md:px-12 relative z-10 group">
        <div class="text-[80px] lg:text-[110px] font-mono font-bold leading-none tracking-tighter text-accent-resp tabular-nums hover:scale-105 transition-transform duration-500 cursor-default drop-shadow-[0_0_25px_rgba(var(--color-accent-resp),0.4)]">
          {formatNum(mockDashboardStats.avgResponsibility)}
        </div>
        <div class="uppercase tracking-[0.2em] text-xs text-secondary mt-4 mb-8 font-medium">Responsibility Score</div>
        <div class="w-full h-8 opacity-70 group-hover:opacity-100 transition-opacity">
          <Sparkline data={mockTimeSeriesResp} color="var(--color-accent-resp)" />
        </div>
      </div>
    </div>

    <!-- Status Ribbon -->
    <div class="mt-16 flex items-center justify-center">
      <div class="px-6 py-2 rounded-full bg-surface/50 border border-elevated/50 backdrop-blur-md flex items-center gap-4 text-xs font-mono text-secondary tracking-widest shadow-xl">
        <span class="flex items-center gap-2 text-primary"><div class="w-2 h-2 rounded-full bg-accent-perf animate-pulse"></div> {mockDashboardStats.requestsPerMinute} REQ/MIN</span>
        <span class="opacity-50">·</span>
        <span>$0.18 AVG COST</span>
        <span class="opacity-50">·</span>
        <span class="text-accent-danger">3 FLAGS/HR</span>
      </div>
    </div>
  </section>

  <!-- ACTIVITY ZONE -->
  <section class="flex flex-col gap-6 -mx-6 lg:-mx-12">
    <div class="px-6 lg:px-12 flex items-baseline justify-between w-full max-w-[1600px] mx-auto">
      <h2 class="font-serif text-3xl md:text-4xl text-primary italic tracking-wide">Recent Activity</h2>
      <a href="/requests" class="text-sm font-mono text-accent-resp hover:text-accent-resp transition-colors flex items-center gap-2 group">
        View All <span class="group-hover:translate-x-1 transition-transform">→</span>
      </a>
    </div>

    <!-- Horizontal Scroll Container -->
    <div class="w-full overflow-x-auto pb-8 hide-scrollbar px-6 lg:px-12 cursor-grab active:cursor-grabbing">
      <div class="flex gap-4 w-max pr-12">
        {#each mockRequests.slice(0, 15) as req}
          <div 
            class="w-72 shrink-0 bg-surface/30 hover:bg-surface/60 border border-elevated/40 hover:border-accent-resp/30 rounded-2xl p-5 flex flex-col gap-4 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl hover:shadow-elevated/20 backdrop-blur-sm cursor-pointer group"
            onclick={() => goto(`/requests/${req.id}`)}
          >
            <div class="flex justify-between items-start">
              <span class="text-xs font-mono text-secondary">{formatTime(req.timestamp)}</span>
              <ActionBadge action={req.action} />
            </div>
            
            <div class="flex flex-col gap-1">
              <span class="text-sm font-medium text-primary truncate group-hover:text-accent-resp transition-colors">{req.endpoint}</span>
              <span class="text-xs text-secondary truncate">{req.model}</span>
            </div>
            
            <div class="flex items-end justify-between mt-auto pt-2 border-t border-elevated/30">
              <div class="flex flex-col">
                <span class="text-[10px] uppercase tracking-widest text-secondary opacity-60">Score</span>
                <span class="text-sm font-mono text-primary">{formatNum(req.responsibilityScore)}</span>
              </div>
              <div class="flex flex-col text-right">
                <span class="text-[10px] uppercase tracking-widest text-secondary opacity-60">Cost</span>
                <span class="text-sm font-mono text-primary">{formatCost(req.costAmount)}</span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- INSIGHTS ZONE -->
  <section class="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-8 max-w-[1600px] mx-auto w-full">
    
    <!-- Left: Performance (Wide) -->
    <div class="lg:col-span-5 flex flex-col gap-8 relative p-6 lg:p-0">
      <h3 class="text-sm font-mono tracking-widest uppercase text-secondary">Performance Insights</h3>
      
      <div class="flex items-end gap-4">
        <div class="text-6xl font-mono text-accent-danger font-bold tabular-nums tracking-tighter">0.4%</div>
        <div class="text-sm text-secondary pb-2 max-w-[150px]">Hallucination rate across all models (7d)</div>
      </div>

      <div class="flex flex-col gap-4">
        <div class="flex justify-between text-xs">
          <span class="text-primary">Confidence Distribution</span>
          <span class="text-accent-perf">92% High</span>
        </div>
        <!-- Thin horizontal bar -->
        <div class="w-full h-1.5 rounded-full bg-surface flex overflow-hidden">
          <div class="h-full bg-accent-perf w-[92%]"></div>
          <div class="h-full bg-accent-cost w-[6%]"></div>
          <div class="h-full bg-accent-danger w-[2%]"></div>
        </div>
      </div>

      <div class="mt-4">
        <h4 class="text-xs uppercase tracking-wider text-secondary mb-4">Top Flagged Categories</h4>
        <div class="flex flex-col gap-3">
          <div class="flex items-center justify-between group">
            <span class="text-sm text-primary group-hover:text-accent-resp transition-colors">Factual Inconsistency</span>
            <span class="font-mono text-xs bg-surface px-2 py-1 rounded">241</span>
          </div>
          <div class="w-full h-px bg-border/40"></div>
          <div class="flex items-center justify-between group">
            <span class="text-sm text-primary group-hover:text-accent-resp transition-colors">Formatting Errors</span>
            <span class="font-mono text-xs bg-surface px-2 py-1 rounded">189</span>
          </div>
          <div class="w-full h-px bg-border/40"></div>
          <div class="flex items-center justify-between group">
            <span class="text-sm text-primary group-hover:text-accent-resp transition-colors">Context Misses</span>
            <span class="font-mono text-xs bg-surface px-2 py-1 rounded">56</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Center: Cost (Narrow) -->
    <div class="lg:col-span-3 flex flex-col gap-8 relative p-6 lg:p-0 items-center lg:items-start lg:border-l lg:border-r border-elevated/30 lg:px-8">
      <h3 class="text-sm font-mono tracking-widest uppercase text-secondary w-full text-center lg:text-left">Cost Pulse</h3>
      
      <div class="flex flex-col items-center justify-center h-full w-full py-4 relative group">
        <!-- Vertical thermometer/gauge -->
        <div class="h-48 w-4 rounded-full bg-surface/50 border border-elevated relative overflow-hidden flex flex-col justify-end shadow-inner">
          <div class="w-full bg-gradient-to-t from-accent-cost to-accent-cost h-[68%] relative">
            <div class="absolute top-0 left-0 right-0 h-1 bg-white/50"></div>
          </div>
        </div>
        
        <div class="absolute left-1/2 -translate-x-1/2 top-1/2 -translate-y-1/2 flex flex-col items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity bg-base/80 backdrop-blur-md p-3 rounded-xl border border-elevated shadow-2xl z-10 whitespace-nowrap">
          <span class="text-xs uppercase tracking-widest text-secondary">Burn Rate</span>
          <span class="font-mono text-xl text-accent-cost font-bold">$1.2k / $1.8k</span>
          <span class="text-[10px] text-secondary">68% of Monthly Budget</span>
        </div>
      </div>

      <div class="flex flex-col items-center lg:items-start text-center lg:text-left w-full mt-auto">
        <span class="text-2xl font-mono text-accent-perf mb-1">$342.50</span>
        <span class="text-xs text-secondary max-w-[200px]">Saved this month via Semantic Cache</span>
      </div>
    </div>

    <!-- Right: Responsibility (Wide) -->
    <div class="lg:col-span-4 flex flex-col gap-8 relative p-6 lg:p-0">
      <h3 class="text-sm font-mono tracking-widest uppercase text-secondary">Responsibility Status</h3>
      
      <div class="flex gap-8 items-end">
        <div class="flex flex-col gap-2">
          <span class="text-xs uppercase tracking-wider text-secondary">PII Incidents</span>
          <div class="flex items-center gap-3">
            <span class="text-5xl font-mono text-primary font-bold">12</span>
            <div class="flex items-center text-accent-perf text-xs font-mono bg-accent-perf/10 px-2 py-1 rounded">
              <TrendingDown class="w-3 h-3 mr-1" /> 24%
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 flex flex-col gap-4">
        <h4 class="text-xs uppercase tracking-wider text-secondary">Toxicity Breakdown</h4>
        <div class="flex gap-1.5 flex-wrap">
          <!-- Represents severities -->
          {#each Array(45) as _, i}
            <div class="w-3 h-3 rounded-sm {i < 3 ? 'bg-accent-danger animate-pulse shadow-[0_0_8px_rgba(var(--color-accent-danger),0.8)]' : i < 12 ? 'bg-accent-cost opacity-80' : 'bg-surface border border-elevated/50'}"></div>
          {/each}
        </div>
        <div class="flex gap-4 text-[10px] uppercase font-mono text-secondary mt-1">
          <span class="flex items-center gap-1"><div class="w-2 h-2 bg-accent-danger rounded-sm"></div> Critical (3)</span>
          <span class="flex items-center gap-1"><div class="w-2 h-2 bg-accent-cost rounded-sm"></div> Flagged (9)</span>
        </div>
      </div>

      <div class="mt-auto pt-6 flex gap-3 flex-wrap">
        <div class="px-3 py-1.5 rounded-full border border-accent-perf/30 bg-accent-perf/5 text-accent-perf text-xs font-mono flex items-center gap-2">
          <CheckCircle2 class="w-3.5 h-3.5" /> HIPAA Active
        </div>
        <div class="px-3 py-1.5 rounded-full border border-accent-perf/30 bg-accent-perf/5 text-accent-perf text-xs font-mono flex items-center gap-2">
          <CheckCircle2 class="w-3.5 h-3.5" /> SOC2 Monitored
        </div>
      </div>
    </div>
  </section>
</div>

<style>
  .hide-scrollbar::-webkit-scrollbar {
    display: none;
  }
  .hide-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>
