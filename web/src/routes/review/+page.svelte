<script lang="ts">
	import { mockReviewQueue, mockRequests } from '$lib/data/dummy';
	import { slide } from 'svelte/transition';
	
	let items = $state(mockReviewQueue.map(item => ({
		...item,
		request: mockRequests.find(r => r.id === item.requestId)
	})));
	
	let sortMode = $state<'severity' | 'time'>('severity');
	
	const severityOrder = { critical: 3, high: 2, medium: 1, low: 0 };
	
	let sortedItems = $derived(
		[...items].sort((a, b) => {
			if (sortMode === 'severity') {
				return severityOrder[b.severity] - severityOrder[a.severity];
			} else {
				return new Date(a.slaDeadline).getTime() - new Date(b.slaDeadline).getTime();
			}
		})
	);
	
	let pendingCount = $derived(items.filter(i => i.status === 'pending').length);
	
	let expandedItemId = $state<string | null>(null);
	
	const severityColors: Record<string, string> = {
		medium: 'text-accent-warning border-accent-warning',
		high: 'text-[#ff6d00] border-[#ff6d00]', // orange
		critical: 'text-accent-danger border-accent-danger'
	};
	
	function getTimeRemaining(deadline: string) {
		const ms = new Date(deadline).getTime() - Date.now();
		const mins = Math.max(0, Math.floor(ms / 60000));
		return mins;
	}
	
	function getTimeSince(timestamp: string) {
		const ms = Date.now() - new Date(timestamp).getTime();
		const hrs = Math.floor(ms / 3600000);
		const mins = Math.floor((ms % 3600000) / 60000);
		if (hrs > 0) return `${hrs}h ${mins}m ago`;
		return `${mins}m ago`;
	}

	function handleAction(id: string, action: 'approved' | 'rejected' | 'pending', e: Event) {
		e.stopPropagation();
		items = items.map(i => i.id === id ? { ...i, status: action } : i);
		if (action !== 'pending') {
			expandedItemId = null;
		}
	}
	
	function toggleExpand(id: string) {
		if (expandedItemId === id) {
			expandedItemId = null;
		} else {
			expandedItemId = id;
		}
	}
</script>

<svelte:head>
	<title>Review Queue | ControlPlane.ai</title>
</svelte:head>

<div class="max-w-7xl mx-auto p-6 md:p-10 space-y-8 min-h-screen">
	<!-- Header -->
	<header class="animate-fade-up space-y-2 flex flex-col md:flex-row md:items-end justify-between gap-4">
		<div>
			<h1 class="font-serif text-5xl md:text-6xl text-primary tracking-tight">Review Queue</h1>
			<p class="text-secondary font-mono text-sm tracking-wide">
				{pendingCount} items pending human review
			</p>
		</div>
		
		<div class="flex items-center gap-2 bg-surface p-1 rounded-lg border border-elevated">
			<button 
				class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors {sortMode === 'severity' ? 'bg-elevated text-primary' : 'text-secondary hover:text-primary'}"
				onclick={() => sortMode = 'severity'}
			>
				Sort by Severity
			</button>
			<button 
				class="px-3 py-1.5 rounded-md text-sm font-medium transition-colors {sortMode === 'time' ? 'bg-elevated text-primary' : 'text-secondary hover:text-primary'}"
				onclick={() => sortMode = 'time'}
			>
				Sort by SLA
			</button>
		</div>
	</header>

	<!-- Queue -->
	<div class="space-y-4 pb-20">
		{#each sortedItems as item, i (item.id)}
			{#if item.status === 'pending'}
				{@const minsLeft = getTimeRemaining(item.slaDeadline)}
				{@const pctLeft = Math.min(100, (minsLeft / 120) * 100)} <!-- assuming 2h max SLA -->
				{@const progressColor = pctLeft > 50 ? 'bg-accent-resp' : pctLeft > 25 ? 'bg-accent-warning' : 'bg-accent-danger'}
				
				<div class="bg-surface rounded-xl border border-elevated overflow-hidden relative shadow-[inset_0_1px_1px_rgba(255,255,255,0.05)] transition-all animate-fade-up" style="animation-delay: {100 + i * 50}ms;">
					<!-- Header Row -->
					<div 
						class="p-5 flex flex-col md:flex-row md:items-center justify-between gap-4 cursor-pointer hover:bg-hover/50 transition-colors"
						onclick={() => toggleExpand(item.id)}
						role="button"
						tabindex="0"
						onkeydown={(e) => e.key === 'Enter' && toggleExpand(item.id)}
					>
						<div class="flex items-center gap-4 flex-1">
							<span class="px-2.5 py-1 rounded border uppercase text-[10px] font-bold tracking-widest {severityColors[item.severity]}">
								{item.severity}
							</span>
							
							<div class="flex flex-col">
								<h3 class="font-bold text-primary">{item.reason}</h3>
								<div class="flex items-center gap-3 text-xs font-mono text-secondary">
									<span>{item.request?.model || 'Unknown Model'}</span>
									<span class="w-1 h-1 rounded-full bg-elevated"></span>
									<span>Escalated {getTimeSince(item.escalatedAt)}</span>
								</div>
							</div>
						</div>
						
						<div class="flex items-center gap-6">
							<div class="flex flex-col items-end gap-1.5 min-w-[120px]">
								<span class="text-xs font-mono {minsLeft < 30 ? 'text-accent-danger' : 'text-secondary'}">
									{minsLeft} min remaining
								</span>
								<div class="w-full h-1.5 bg-elevated rounded-full overflow-hidden">
									<div class="h-full rounded-full transition-all duration-1000 {progressColor}" style="width: {pctLeft}%"></div>
								</div>
							</div>
							
							<button 
								class="px-4 py-2 bg-elevated hover:bg-hover text-primary text-sm font-medium rounded-lg border border-transparent hover:border-elevated transition-colors"
							>
								{expandedItemId === item.id ? 'Close' : 'Review'}
							</button>
						</div>
					</div>
					
					<!-- Expanded Area -->
					{#if expandedItemId === item.id}
						<div class="border-t border-elevated p-5 bg-elevated/10" transition:slide={{duration: 300}}>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
								<!-- Original -->
								<div class="space-y-2">
									<h4 class="text-xs uppercase tracking-widest text-secondary font-bold">Original Response</h4>
									<div class="bg-[#09080a] p-4 rounded-lg border border-elevated border-dashed h-48 overflow-y-auto">
										<p class="text-secondary text-sm whitespace-pre-wrap">{item.request?.response || item.originalResponse}</p>
									</div>
								</div>
								
								<!-- Suggested -->
								<div class="space-y-2">
									<h4 class="text-xs uppercase tracking-widest text-secondary font-bold flex justify-between">
										<span>Suggested Edit</span>
										<span class="text-accent-perf">AI Generated</span>
									</h4>
									<div class="bg-accent-perf/5 p-4 rounded-lg border border-accent-perf/20 h-48 overflow-y-auto">
										<!-- Just mockup diff visualization -->
										<p class="text-primary text-sm whitespace-pre-wrap"><span class="bg-accent-perf/20 text-accent-perf line-through mr-1 opacity-70">original bad text</span> <span class="bg-accent-resp/20 text-accent-resp font-medium">redacted or safer text</span> continuing the response here...</p>
									</div>
								</div>
							</div>
							
							<div class="flex justify-end gap-3 pt-4 border-t border-elevated">
								<button 
									class="px-5 py-2 rounded-lg font-medium text-sm border border-accent-danger text-accent-danger hover:bg-accent-danger hover:text-base transition-colors"
									onclick={(e) => handleAction(item.id, 'rejected', e)}
								>
									Reject (Block)
								</button>
								<button 
									class="px-5 py-2 rounded-lg font-medium text-sm border border-accent-perf text-accent-perf hover:bg-accent-perf hover:text-base transition-colors"
									onclick={(e) => handleAction(item.id, 'pending', e)}
								>
									Manual Edit
								</button>
								<button 
									class="px-5 py-2 rounded-lg font-medium text-sm bg-accent-resp text-base hover:bg-[#00c968] transition-colors shadow-[0_0_15px_rgba(0,230,118,0.2)]"
									onclick={(e) => handleAction(item.id, 'approved', e)}
								>
									Approve Edit
								</button>
							</div>
						</div>
					{/if}
				</div>
			{/if}
		{/each}
		
		{#if pendingCount === 0}
			<div class="py-20 text-center space-y-4 animate-fade-up">
				<div class="w-16 h-16 mx-auto rounded-full bg-accent-resp/10 flex items-center justify-center text-accent-resp">
					<svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
					</svg>
				</div>
				<h3 class="text-xl font-serif text-primary">Queue is clear</h3>
				<p class="text-secondary font-mono text-sm">All pending items have been reviewed.</p>
			</div>
		{/if}
	</div>
</div>
