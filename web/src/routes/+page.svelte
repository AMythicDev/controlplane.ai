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
	const formatCost = (val: number) =>
		new Intl.NumberFormat('en-US', {
			style: 'currency',
			currency: 'USD',
			minimumFractionDigits: 3
		}).format(val);
	const formatNum = (val: number) =>
		new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(val);

	// Time format for cards
	const formatTime = (ts: string) => {
		const d = new Date(ts);
		return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
	};

	let liveCost = $state<number>(0);
	let userDailyLimit = $state<number>(50);
	let userMonthlyLimit = $state<number>(1500);

	let burnPercent = $derived(userMonthlyLimit > 0 ? Math.min(100, Math.round((liveCost / userMonthlyLimit) * 100)) : (liveCost > 0 ? 100 : 0));

	onMount(async () => {
		try {
			const [costRes, configRes] = await Promise.all([
				fetch('/v1/cost'),
				fetch('/v1/config')
			]);

			if (costRes.ok) {
				const costData = await costRes.json();
				liveCost = costData.cost_dollars;
			}

			if (configRes.ok) {
				const configData = await configRes.json();
				if (configData.per_user_daily_limit !== undefined) {
					userDailyLimit = configData.per_user_daily_limit / 1000000;
				}
				if (configData.per_user_monthly_limit !== undefined) {
					userMonthlyLimit = configData.per_user_monthly_limit / 1000000;
				}
			}
		} catch (e) {
			console.error('Failed to fetch dashboard data', e);
		}
	});
</script>

<div class="relative z-10 flex min-h-full w-full flex-col gap-24 pb-32">
	<!-- Radial gradient background behind hero -->
	<div
		class="pointer-events-none absolute top-0 left-1/2 -z-10 h-[80vh] w-[80vw] -translate-x-1/2 rounded-full bg-[radial-gradient(ellipse_at_top,rgba(var(--color-elevated),0.15)_0%,transparent_70%)] blur-3xl"
	></div>

	<!-- HERO ZONE -->
	<section class="flex flex-col items-center pt-8">
		<div class="relative grid w-full max-w-6xl grid-cols-1 gap-8 md:grid-cols-3 md:gap-0">
			<!-- Vertical Separators -->
			<div
				class="via-border absolute top-12 bottom-12 left-1/3 z-0 hidden w-px bg-gradient-to-b from-transparent to-transparent md:block"
			></div>
			<div
				class="via-border absolute top-12 right-1/3 bottom-12 z-0 hidden w-px bg-gradient-to-b from-transparent to-transparent md:block"
			></div>

			<!-- Performance -->
			<div class="group relative z-10 flex flex-col items-center px-4 md:px-12">
				<div
					class="cursor-default font-mono text-[80px] leading-none font-bold tracking-tighter text-accent-perf tabular-nums shadow-accent-perf/20 drop-shadow-[0_0_25px_rgba(var(--color-accent-perf),0.4)] transition-transform duration-500 hover:scale-105 lg:text-[110px]"
				>
					{formatNum(mockDashboardStats.avgPerformance)}
				</div>
				<div class="mt-4 mb-8 text-xs font-medium tracking-[0.2em] text-secondary uppercase">
					Performance Score
				</div>
				<div class="h-8 w-full opacity-70 transition-opacity group-hover:opacity-100">
					<Sparkline data={mockTimeSeriesPerf} color="var(--color-accent-perf)" />
				</div>
			</div>

			<!-- Cost -->
			<div class="group relative z-10 flex flex-col items-center px-4 md:px-12">
				<div
					class="flex cursor-default items-start font-mono text-[80px] leading-none font-bold tracking-tighter text-accent-cost tabular-nums drop-shadow-[0_0_25px_rgba(var(--color-accent-cost),0.4)] transition-transform duration-500 hover:scale-105 lg:text-[110px]"
				>
					<span class="mt-4 text-4xl opacity-50 lg:text-6xl">$</span>
					<span>{formatNum(liveCost)}</span>
				</div>
				<div class="mt-4 mb-8 text-xs font-medium tracking-[0.2em] text-secondary uppercase">
					Total Spend
				</div>
				<div class="h-8 w-full opacity-70 transition-opacity group-hover:opacity-100">
					<Sparkline data={mockTimeSeriesCost} color="var(--color-accent-cost)" />
				</div>
			</div>

			<!-- Responsibility -->
			<div class="group relative z-10 flex flex-col items-center px-4 md:px-12">
				<div
					class="cursor-default font-mono text-[80px] leading-none font-bold tracking-tighter text-accent-resp tabular-nums drop-shadow-[0_0_25px_rgba(var(--color-accent-resp),0.4)] transition-transform duration-500 hover:scale-105 lg:text-[110px]"
				>
					{formatNum(mockDashboardStats.avgResponsibility)}
				</div>
				<div class="mt-4 mb-8 text-xs font-medium tracking-[0.2em] text-secondary uppercase">
					Responsibility Score
				</div>
				<div class="h-8 w-full opacity-70 transition-opacity group-hover:opacity-100">
					<Sparkline data={mockTimeSeriesResp} color="var(--color-accent-resp)" />
				</div>
			</div>
		</div>

		<!-- Status Ribbon -->
		<div class="mt-16 flex items-center justify-center">
			<div
				class="flex items-center gap-4 rounded-full border border-elevated/50 bg-surface/50 px-6 py-2 font-mono text-xs tracking-widest text-secondary shadow-xl backdrop-blur-md"
			>
				<span class="flex items-center gap-2 text-primary"
					><div class="h-2 w-2 animate-pulse rounded-full bg-accent-perf"></div>
					{mockDashboardStats.requestsPerMinute} REQ/MIN</span
				>
				<span class="opacity-50">·</span>
				<span>$0.18 AVG COST</span>
				<span class="opacity-50">·</span>
				<span class="text-accent-danger">3 FLAGS/HR</span>
			</div>
		</div>
	</section>

	<!-- ACTIVITY ZONE -->
	<section class="-mx-6 flex flex-col gap-6 lg:-mx-12">
		<div class="mx-auto flex w-full max-w-[1600px] items-baseline justify-between px-6 lg:px-12">
			<h2 class="font-serif text-3xl tracking-wide text-primary italic md:text-4xl">
				Recent Activity
			</h2>
			<a
				href="/requests"
				class="group flex items-center gap-2 font-mono text-sm text-accent-resp transition-colors hover:text-accent-resp"
			>
				View All <span class="transition-transform group-hover:translate-x-1">→</span>
			</a>
		</div>

		<!-- Horizontal Scroll Container -->
		<div
			class="hide-scrollbar w-full cursor-grab overflow-x-auto px-6 pb-8 active:cursor-grabbing lg:px-12"
		>
			<div class="flex w-max gap-4 pr-12">
				{#each mockRequests.slice(0, 15) as req}
					<div
						class="group flex w-72 shrink-0 cursor-pointer flex-col gap-4 rounded-2xl border border-elevated/40 bg-surface/30 p-5 backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-accent-resp/30 hover:bg-surface/60 hover:shadow-2xl hover:shadow-elevated/20"
						onclick={() => goto(`/requests/${req.id}`)}
					>
						<div class="flex items-start justify-between">
							<span class="font-mono text-xs text-secondary">{formatTime(req.timestamp)}</span>
							<ActionBadge action={req.action} />
						</div>

						<div class="flex flex-col gap-1">
							<span
								class="truncate text-sm font-medium text-primary transition-colors group-hover:text-accent-resp"
								>{req.endpoint}</span
							>
							<span class="truncate text-xs text-secondary">{req.model}</span>
						</div>

						<div class="mt-auto flex items-end justify-between border-t border-elevated/30 pt-2">
							<div class="flex flex-col">
								<span class="text-[10px] tracking-widest text-secondary uppercase opacity-60"
									>Score</span
								>
								<span class="font-mono text-sm text-primary"
									>{formatNum(req.responsibilityScore)}</span
								>
							</div>
							<div class="flex flex-col text-right">
								<span class="text-[10px] tracking-widest text-secondary uppercase opacity-60"
									>Cost</span
								>
								<span class="font-mono text-sm text-primary">{formatCost(req.costAmount)}</span>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<!-- INSIGHTS ZONE -->
	<section class="mx-auto grid w-full max-w-[1600px] grid-cols-1 gap-12 lg:grid-cols-12 lg:gap-8">
		<!-- Left: Performance (Wide) -->
		<div class="relative flex flex-col gap-8 p-6 lg:col-span-5 lg:p-0">
			<h3 class="font-mono text-sm tracking-widest text-secondary uppercase">
				Performance Insights
			</h3>

			<div class="flex items-end gap-4">
				<div class="font-mono text-6xl font-bold tracking-tighter text-accent-danger tabular-nums">
					0.4%
				</div>
				<div class="max-w-[150px] pb-2 text-sm text-secondary">
					Hallucination rate across all models (7d)
				</div>
			</div>

			<div class="flex flex-col gap-4">
				<div class="flex justify-between text-xs">
					<span class="text-primary">Confidence Distribution</span>
					<span class="text-accent-perf">92% High</span>
				</div>
				<!-- Thin horizontal bar -->
				<div class="flex h-1.5 w-full overflow-hidden rounded-full bg-surface">
					<div class="h-full w-[92%] bg-accent-perf"></div>
					<div class="h-full w-[6%] bg-accent-cost"></div>
					<div class="h-full w-[2%] bg-accent-danger"></div>
				</div>
			</div>

			<div class="mt-4">
				<h4 class="mb-4 text-xs tracking-wider text-secondary uppercase">Top Flagged Categories</h4>
				<div class="flex flex-col gap-3">
					<div class="group flex items-center justify-between">
						<span class="text-sm text-primary transition-colors group-hover:text-accent-resp"
							>Factual Inconsistency</span
						>
						<span class="rounded bg-surface px-2 py-1 font-mono text-xs">241</span>
					</div>
					<div class="bg-border/40 h-px w-full"></div>
					<div class="group flex items-center justify-between">
						<span class="text-sm text-primary transition-colors group-hover:text-accent-resp"
							>Formatting Errors</span
						>
						<span class="rounded bg-surface px-2 py-1 font-mono text-xs">189</span>
					</div>
					<div class="bg-border/40 h-px w-full"></div>
					<div class="group flex items-center justify-between">
						<span class="text-sm text-primary transition-colors group-hover:text-accent-resp"
							>Context Misses</span
						>
						<span class="rounded bg-surface px-2 py-1 font-mono text-xs">56</span>
					</div>
				</div>
			</div>
		</div>

		<!-- Center: Cost (Narrow) -->
		<div
			class="relative flex flex-col items-center gap-8 border-elevated/30 p-6 lg:col-span-3 lg:items-start lg:border-r lg:border-l lg:p-0 lg:px-8"
		>
			<h3
				class="w-full text-center font-mono text-sm tracking-widest text-secondary uppercase lg:text-left"
			>
				Cost Pulse
			</h3>

			<div class="group relative flex h-full w-full flex-col items-center justify-center py-4">
				<!-- Vertical thermometer/gauge -->
				<div
					class="relative flex h-48 w-4 flex-col justify-end overflow-hidden rounded-full border border-elevated bg-surface/50 shadow-inner"
				>
					<div class="relative w-full bg-gradient-to-t from-accent-cost to-accent-cost transition-all duration-1000" style="height: {burnPercent}%;">
						<div class="absolute top-0 right-0 left-0 h-1 bg-white/50"></div>
					</div>
				</div>

				<div
					class="absolute top-1/2 left-1/2 z-10 flex -translate-x-1/2 -translate-y-1/2 flex-col items-center gap-1 rounded-xl border border-elevated bg-base/80 p-3 whitespace-nowrap opacity-0 shadow-2xl backdrop-blur-md transition-opacity group-hover:opacity-100"
				>
					<span class="text-xs tracking-widest text-secondary uppercase">Burn Rate</span>
					<span class="font-mono text-xl font-bold text-accent-cost">${formatNum(liveCost)} / ${formatNum(userMonthlyLimit)}</span>
					<span class="text-[10px] text-secondary">{burnPercent}% of Monthly Budget</span>
				</div>
			</div>

			<div
				class="mt-auto flex w-full flex-col items-center text-center lg:items-start lg:text-left"
			>
				<span class="mb-1 font-mono text-2xl text-accent-perf">$342.50</span>
				<span class="max-w-[200px] text-xs text-secondary">Saved this month via Semantic Cache</span
				>
			</div>
		</div>

		<!-- Right: Responsibility (Wide) -->
		<div class="relative flex flex-col gap-8 p-6 lg:col-span-4 lg:p-0">
			<h3 class="font-mono text-sm tracking-widest text-secondary uppercase">
				Responsibility Status
			</h3>

			<div class="flex items-end gap-8">
				<div class="flex flex-col gap-2">
					<span class="text-xs tracking-wider text-secondary uppercase">PII Incidents</span>
					<div class="flex items-center gap-3">
						<span class="font-mono text-5xl font-bold text-primary">12</span>
						<div
							class="flex items-center rounded bg-accent-perf/10 px-2 py-1 font-mono text-xs text-accent-perf"
						>
							<TrendingDown class="mr-1 h-3 w-3" /> 24%
						</div>
					</div>
				</div>
			</div>

			<div class="mt-4 flex flex-col gap-4">
				<h4 class="text-xs tracking-wider text-secondary uppercase">Toxicity Breakdown</h4>
				<div class="flex flex-wrap gap-1.5">
					<!-- Represents severities -->
					{#each Array(45) as _, i}
						<div
							class="h-3 w-3 rounded-sm {i < 3
								? 'animate-pulse bg-accent-danger shadow-[0_0_8px_rgba(var(--color-accent-danger),0.8)]'
								: i < 12
									? 'bg-accent-cost opacity-80'
									: 'border border-elevated/50 bg-surface'}"
						></div>
					{/each}
				</div>
				<div class="mt-1 flex gap-4 font-mono text-[10px] text-secondary uppercase">
					<span class="flex items-center gap-1"
						><div class="h-2 w-2 rounded-sm bg-accent-danger"></div>
						 Critical (3)</span
					>
					<span class="flex items-center gap-1"
						><div class="h-2 w-2 rounded-sm bg-accent-cost"></div>
						 Flagged (9)</span
					>
				</div>
			</div>

			<div class="mt-auto flex flex-wrap gap-3 pt-6">
				<div
					class="flex items-center gap-2 rounded-full border border-accent-perf/30 bg-accent-perf/5 px-3 py-1.5 font-mono text-xs text-accent-perf"
				>
					<CheckCircle2 class="h-3.5 w-3.5" /> HIPAA Active
				</div>
				<div
					class="flex items-center gap-2 rounded-full border border-accent-perf/30 bg-accent-perf/5 px-3 py-1.5 font-mono text-xs text-accent-perf"
				>
					<CheckCircle2 class="h-3.5 w-3.5" /> SOC2 Monitored
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
