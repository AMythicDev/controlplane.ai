<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchCostData, fetchConfig, fetchAnalytics } from '$lib/data/api';
	import type { ModelAnalyticsItem, AnalyticsResponse } from '$lib/types';
	import DollarSign from 'lucide-svelte/icons/dollar-sign';
	import Calculator from 'lucide-svelte/icons/calculator';
	import RefreshCw from 'lucide-svelte/icons/refresh-cw';
	import ArrowUpRight from 'lucide-svelte/icons/arrow-up-right';
	import Activity from 'lucide-svelte/icons/activity';

	import ModelUsageChart from '$lib/components/dashboard/ModelUsageChart.svelte';
	import ModelConfidenceChart from '$lib/components/dashboard/ModelConfidenceChart.svelte';
	import ModelHallucinationChart from '$lib/components/dashboard/ModelHallucinationChart.svelte';
	import ModelToxicityChart from '$lib/components/dashboard/ModelToxicityChart.svelte';
	import CostPulseCard from '$lib/components/dashboard/CostPulseCard.svelte';

	// Live state from API
	let liveCost = $state<number>(0);
	let liveSavings = $state<number>(0);
	let liveAvgCost = $state<number>(0);
	let userDailyLimit = $state<number>(0);
	let userMonthlyLimit = $state<number>(0);

	let analyticsData = $state<AnalyticsResponse | null>(null);
	let modelsList = $state<ModelAnalyticsItem[]>([]);
	let weeklyTotalReqs = $state<number>(0);

	let loading = $state<boolean>(true);
	let refreshing = $state<boolean>(false);
	let lastUpdated = $state<string>('');

	// Number formatters
	const formatNum = (val: number) =>
		new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(val);

	const formatCostDisplay = (val: number) => {
		if (val === 0) return '$0.00';
		if (val < 0.01) return `$${val.toFixed(4)}`;
		return `$${val.toFixed(2)}`;
	};

	async function loadDashboardData() {
		try {
			const [costData, configData, analyticsRes] = await Promise.allSettled([
				fetchCostData(),
				fetchConfig(),
				fetchAnalytics()
			]);

			if (costData.status === 'fulfilled' && costData.value) {
				const c = costData.value;
				liveCost = c.cost_dollars ?? 0;
				liveSavings = c.semantic_cache_savings ?? 0;
				liveAvgCost = c.avg_cost ?? c.average_cost_dollars ?? 0;
			}

			if (configData.status === 'fulfilled' && configData.value) {
				const cfg = configData.value;
				if (cfg.per_user_daily_limit !== undefined) {
					userDailyLimit = cfg.per_user_daily_limit / 1000000;
				}
				if (cfg.per_user_monthly_limit !== undefined) {
					userMonthlyLimit = cfg.per_user_monthly_limit / 1000000;
				}
			}

			if (analyticsRes.status === 'fulfilled' && analyticsRes.value) {
				analyticsData = analyticsRes.value;
				modelsList = analyticsRes.value.models || [];
				weeklyTotalReqs = analyticsRes.value.weekly_requests || analyticsRes.value.total_requests || 0;
			}

			const now = new Date();
			lastUpdated = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
		} catch (e) {
			console.error('Failed to load observatory dashboard data', e);
		} finally {
			loading = false;
			refreshing = false;
		}
	}

	async function handleRefresh() {
		refreshing = true;
		await loadDashboardData();
	}

	onMount(() => {
		loadDashboardData();
		const interval = setInterval(loadDashboardData, 15000);
		return () => clearInterval(interval);
	});
</script>

<div class="relative z-10 flex min-h-full w-full flex-col gap-10 pb-24 animate-fade-up">
	<!-- Radial ambient glow -->
	<div
		class="pointer-events-none absolute top-0 left-1/2 -z-10 h-[600px] w-[90vw] -translate-x-1/2 rounded-full bg-[radial-gradient(ellipse_at_top,rgba(161,0,255,0.08)_0%,rgba(152,218,21,0.04)_40%,transparent_70%)] blur-3xl"
	></div>

	<!-- HEADER & TITLE SECTION -->
	<header class="flex flex-col sm:flex-row sm:items-end justify-between gap-4 border-b border-elevated/30 pb-6">
		<div>
			<div class="flex items-center gap-2 text-xs font-mono tracking-widest text-accent-perf uppercase">
				<span class="flex h-2 w-2 rounded-full bg-accent-perf animate-pulse"></span>
				Live Telemetry
			</div>
			<h1 class="font-serif text-4xl sm:text-5xl text-primary italic tracking-wide mt-1">
				Observatory
			</h1>
			<p class="text-sm text-secondary font-sans mt-1">
				Real-time LLM cost governance, quality scoring & model intelligence
			</p>
		</div>

		<!-- Action bar / Refresh info -->
		<div class="flex items-center gap-3 shrink-0">
			{#if lastUpdated}
				<span class="font-mono text-xs text-tertiary">
					Updated {lastUpdated}
				</span>
			{/if}
			<button
				onclick={handleRefresh}
				class="flex items-center gap-2 rounded-full border border-elevated/60 bg-surface/60 px-4 py-1.5 font-mono text-xs text-secondary shadow-sm backdrop-blur-md transition-all hover:border-accent-perf/40 hover:text-primary hover:bg-surface active:scale-95 disabled:opacity-50"
				disabled={refreshing}
				title="Refresh telemetry metrics"
			>
				<RefreshCw class="h-3.5 w-3.5 {refreshing ? 'animate-spin text-accent-perf' : ''}" />
				<span>{refreshing ? 'Refreshing...' : 'Refresh'}</span>
			</button>
			<a
				href="/requests"
				class="flex items-center gap-1 rounded-full border border-elevated/40 bg-elevated/30 px-4 py-1.5 font-mono text-xs text-secondary hover:text-primary hover:border-elevated transition-colors"
			>
				<span>Explorer</span>
				<ArrowUpRight class="h-3.5 w-3.5" />
			</a>
		</div>
	</header>

	<!-- SECTION 1: WORKING COST METRICS & COST PULSE -->
	<section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 gap-6">
		<!-- Total Cost Card -->
		<div class="lg:col-span-3 relative flex flex-col justify-between overflow-hidden rounded-2xl border border-elevated/40 bg-surface/40 p-6 backdrop-blur-md transition-all duration-300 hover:border-elevated/70 group">
			<div class="pointer-events-none absolute -top-12 -right-12 h-32 w-32 rounded-full bg-accent-cost/10 blur-2xl"></div>

			<div>
				<div class="flex items-center justify-between">
					<span class="font-mono text-xs font-medium uppercase tracking-widest text-secondary">
						Total Spend
					</span>
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-cost/10 text-accent-cost">
						<DollarSign class="h-4 w-4" />
					</div>
				</div>

				<div class="mt-4 flex items-baseline font-mono text-4xl sm:text-5xl font-bold tracking-tight text-accent-cost">
					<span>${formatNum(liveCost)}</span>
				</div>
			</div>

			<div class="mt-6 flex items-center justify-between border-t border-elevated/30 pt-3 font-mono text-xs text-secondary">
				<span class="text-tertiary text-[11px]">Cumulative usage</span>
				<span class="text-accent-cost flex items-center gap-1 font-semibold">
					<span class="h-1.5 w-1.5 rounded-full bg-accent-cost"></span> Active
				</span>
			</div>
		</div>

		<!-- Avg Cost Card -->
		<div class="lg:col-span-3 relative flex flex-col justify-between overflow-hidden rounded-2xl border border-elevated/40 bg-surface/40 p-6 backdrop-blur-md transition-all duration-300 hover:border-elevated/70 group">
			<div class="pointer-events-none absolute -top-12 -right-12 h-32 w-32 rounded-full bg-accent-perf/10 blur-2xl"></div>

			<div>
				<div class="flex items-center justify-between">
					<span class="font-mono text-xs font-medium uppercase tracking-widest text-secondary">
						Avg Cost / Req
					</span>
					<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-perf/10 text-accent-perf">
						<Calculator class="h-4 w-4" />
					</div>
				</div>

				<div class="mt-4 flex items-baseline font-mono text-4xl sm:text-5xl font-bold tracking-tight text-primary">
					<span>{formatCostDisplay(liveAvgCost)}</span>
				</div>
			</div>

			<div class="mt-6 flex items-center justify-between border-t border-elevated/30 pt-3 font-mono text-xs text-secondary">
				<span class="text-tertiary text-[11px]">Per completion</span>
				<span class="text-accent-perf font-semibold">Micro-metered</span>
			</div>
		</div>

		<!-- Cost Pulse & Budget Thermometer (Wide) -->
		<div class="lg:col-span-6">
			<CostPulseCard
				{liveCost}
				{liveSavings}
				{userMonthlyLimit}
				{userDailyLimit}
			/>
		</div>
	</section>

	<!-- SECTION 2: MODEL QUALITY & USAGE CHARTS -->
	<section class="flex flex-col gap-6">
		<div class="flex items-center justify-between">
			<div>
				<h2 class="font-serif text-2xl sm:text-3xl text-primary italic tracking-wide">
					Model Intelligence & Quality Analytics
				</h2>
				<p class="text-xs sm:text-sm text-secondary font-sans mt-0.5">
					Aggregated telemetry across all active LLM inference endpoints
				</p>
			</div>
		</div>

		<!-- 2x2 Grid of Charts -->
		<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
			<!-- Chart 1: Most Used Models Over Last Week -->
			<ModelUsageChart
				models={modelsList}
				weeklyTotal={weeklyTotalReqs}
			/>

			<!-- Chart 2: Model Wise Confidence Score -->
			<ModelConfidenceChart
				models={modelsList}
			/>

			<!-- Chart 3: Model Wise Hallucination Score -->
			<ModelHallucinationChart
				models={modelsList}
			/>

			<!-- Chart 4: Model Wise Toxicity Score -->
			<ModelToxicityChart
				models={modelsList}
			/>
		</div>
	</section>
</div>
