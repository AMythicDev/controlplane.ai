<script lang="ts">
	import { mockPolicies } from '$lib/data/dummy';
	import PolicyCard from '$lib/components/ui/PolicyCard.svelte';
	
	let filter = $state('All');
	
	const filters = ['All', 'Performance', 'Cost', 'Responsibility'];
	
	let filteredPolicies = $derived(
		filter === 'All' 
			? mockPolicies 
			: mockPolicies.filter(p => p.engine.toLowerCase() === filter.toLowerCase())
	);

	let activeCount = $derived(mockPolicies.filter(p => p.enabled).length);
	let blockCount = $derived(mockPolicies.filter(p => p.action === 'block').length);
	let flagCount = $derived(mockPolicies.filter(p => p.action === 'flag').length);
	let editCount = $derived(mockPolicies.filter(p => p.action === 'edit').length);
	let escalateCount = $derived(mockPolicies.filter(p => p.action === 'escalate').length);
</script>

<svelte:head>
	<title>Policies | ControlPlane.ai</title>
</svelte:head>

<div class="max-w-7xl mx-auto p-6 md:p-10 space-y-8 min-h-screen">
	<!-- Header -->
	<header class="animate-fade-up space-y-2">
		<h1 class="font-serif text-5xl md:text-6xl text-primary tracking-tight">Policies</h1>
		<p class="text-secondary font-mono text-sm tracking-wide">
			{activeCount} active policies • {blockCount} blocking • {flagCount} flagging • {editCount} editing • {escalateCount} escalating
		</p>
	</header>

	<!-- Filters -->
	<div class="flex flex-wrap gap-2 animate-fade-up" style="animation-delay: 100ms;">
		{#each filters as f}
			<button 
				class="px-4 py-2 rounded-full text-sm font-medium transition-all {filter === f ? 'bg-primary text-base' : 'bg-surface text-secondary hover:bg-hover border border-elevated'}"
				onclick={() => filter = f}
			>
				{f}
			</button>
		{/each}
	</div>

	<!-- Grid -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 pb-20">
		{#each filteredPolicies as policy, i (policy.id)}
			<div class="animate-fade-up" style="animation-delay: {150 + i * 50}ms;">
				<PolicyCard {policy} />
			</div>
		{/each}
		
		{#if filteredPolicies.length === 0}
			<div class="col-span-1 lg:col-span-2 py-20 text-center text-tertiary font-mono">
				No policies found for this filter.
			</div>
		{/if}
	</div>
</div>
