<script lang="ts">
	import type { Policy } from '$lib/types';
	import { slide } from 'svelte/transition';

	let { policy } = $props<{ policy: Policy }>();
	
	let expanded = $state(false);
	
	const engineColors: Record<string, string> = {
		performance: 'text-accent-perf bg-accent-perf/10 border-accent-perf/20',
		cost: 'text-accent-cost bg-accent-cost/10 border-accent-cost/20',
		responsibility: 'text-accent-resp bg-accent-resp/10 border-accent-resp/20'
	};
	
	const actionColors: Record<string, string> = {
		pass: 'text-accent-resp bg-accent-resp/10 border-accent-resp/20',
		flag: 'text-accent-cost bg-accent-cost/10 border-accent-cost/20',
		edit: 'text-accent-perf bg-accent-perf/10 border-accent-perf/20',
		block: 'text-accent-danger bg-accent-danger/10 border-accent-danger/20',
		escalate: 'text-accent-warning bg-accent-warning/10 border-accent-warning/20'
	};
	
	const severityColors: Record<string, string> = {
		low: 'text-accent-resp',
		medium: 'text-accent-cost',
		high: 'text-accent-warning',
		critical: 'text-accent-danger'
	};

	let enabled = $state(policy.enabled);

	function toggleExpand() {
		expanded = !expanded;
	}
</script>

<div class="bg-surface rounded-xl border border-elevated overflow-hidden relative shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)] transition-colors hover:border-hover group">
	<div class="absolute inset-0 opacity-[0.015] pointer-events-none" style="background-image: url(&quot;data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E&quot;);"></div>
	
	<div 
		class="p-5 flex flex-col md:flex-row md:items-center justify-between gap-4 cursor-pointer relative z-10"
		onclick={toggleExpand}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Enter' && toggleExpand()}
	>
		<div class="flex items-center gap-4">
			<button 
				class="w-6 h-6 flex items-center justify-center text-secondary group-hover:text-primary transition-colors"
				aria-label={expanded ? "Collapse" : "Expand"}
			>
				<svg class="w-4 h-4 transition-transform duration-300 {expanded ? 'rotate-180' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
				</svg>
			</button>
			
			<div class="flex flex-col gap-1">
				<h3 class="font-bold text-primary text-lg leading-none">{policy.name}</h3>
				<div class="flex items-center gap-2 text-xs font-mono text-secondary">
					<span>{policy.check}</span>
				</div>
			</div>
		</div>

		<div class="flex items-center gap-3">
			<span class="px-2 py-1 rounded-md text-xs font-medium uppercase tracking-wider border {engineColors[policy.engine]}">
				{policy.engine}
			</span>
			<span class="px-2 py-1 rounded-full text-xs font-bold uppercase tracking-wider border {actionColors[policy.action]}">
				{policy.action}
			</span>
			
			<div class="w-px h-6 bg-elevated mx-2 hidden md:block"></div>
			
			<button 
				class="relative w-11 h-6 rounded-full transition-colors duration-200 ease-in-out focus-ring z-20 {enabled ? 'bg-accent-perf' : 'bg-elevated'}"
				onclick={(e) => { e.stopPropagation(); enabled = !enabled; }}
				role="switch"
				aria-checked={enabled}
			>
				<span class="absolute top-[2px] left-[2px] bg-white w-5 h-5 rounded-full transition-transform duration-200 {enabled ? 'translate-x-5' : 'translate-x-0'} shadow-sm"></span>
			</button>
		</div>
	</div>

	{#if expanded}
		<div class="px-5 pb-5 pt-2 border-t border-elevated bg-elevated/20 relative z-10" transition:slide={{ duration: 300 }}>
			<div class="mb-4">
				<p class="text-secondary text-sm">{policy.description}</p>
			</div>
			
			<div class="flex items-center gap-2 mb-4 text-sm">
				<span class="text-tertiary">Severity:</span>
				<span class="font-bold uppercase text-xs tracking-wider {severityColors[policy.severity]}">{policy.severity}</span>
			</div>
			
			<div class="bg-[#09080a] p-4 rounded-lg border border-elevated font-mono text-xs overflow-x-auto shadow-inner">
				<pre class="text-secondary m-0 leading-relaxed">{@html policy.yaml.replace(/\n/g, '<br/>').replace(/ /g, '&nbsp;')}</pre>
			</div>
		</div>
	{/if}
</div>
