<script lang="ts">
	import type { ModelAnalyticsItem } from '$lib/types';
	import TrendingUp from 'lucide-svelte/icons/trending-up';
	import Layers from 'lucide-svelte/icons/layers';

	let { models = [], weeklyTotal = 0 } = $props<{
		models: ModelAnalyticsItem[];
		weeklyTotal: number;
	}>();

	function getProviderColor(provider: string) {
		const p = (provider || '').toLowerCase();
		if (p.includes('openai')) return '#10a37f';
		if (p.includes('anthropic')) return '#d97706';
		if (p.includes('google')) return '#4285f4';
		if (p.includes('nvidia')) return '#76b900';
		if (p.includes('openrouter')) return '#6366f1';
		return '#A100FF';
	}

	function cleanModelName(modelSpec: string) {
		if (modelSpec && modelSpec.includes('/')) {
			const parts = modelSpec.split('/');
			return parts.slice(1).join('/');
		}
		return modelSpec || 'Unknown';
	}
</script>

<div class="relative flex flex-col justify-between overflow-hidden rounded-2xl border border-elevated/40 bg-surface/40 p-6 backdrop-blur-md transition-all duration-300 hover:border-elevated/70">
	<!-- Background glow -->
	<div class="pointer-events-none absolute top-0 right-0 -z-10 h-48 w-48 rounded-full bg-accent-perf/5 blur-3xl"></div>

	<div>
		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-2.5">
				<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-perf/10 text-accent-perf">
					<Layers class="h-4 w-4" />
				</div>
				<div>
					<h3 class="font-sans text-base font-semibold text-primary">Most Used Models</h3>
					<p class="text-xs text-secondary">Request volume distribution over the last 7 days</p>
				</div>
			</div>
			<div class="flex items-center gap-1.5 rounded-full border border-elevated/60 bg-base/60 px-3 py-1 font-mono text-xs text-secondary">
				<span class="text-primary font-bold">{weeklyTotal}</span>
				<span>total reqs</span>
			</div>
		</div>

		<div class="mt-6 flex flex-col gap-3.5">
			{#if models.length === 0}
				<div class="flex flex-col items-center justify-center py-10 text-center">
					<div class="mb-2 rounded-full bg-elevated/30 p-3 text-secondary">
						<TrendingUp class="h-5 w-5 opacity-50" />
					</div>
					<p class="font-mono text-xs text-secondary">No requests recorded in the last 7 days</p>
					<span class="mt-1 text-[11px] text-tertiary">Run prompts in the Playground or API to populate data</span>
				</div>
			{:else}
				{#each models as item, index}
					{@const color = getProviderColor(item.provider || item.model)}
					{@const name = cleanModelName(item.model)}
					<div class="group flex flex-col gap-2 rounded-xl border border-elevated/20 bg-base/30 p-3.5 transition-all duration-200 hover:border-elevated/50 hover:bg-base/50">
						<div class="flex items-center justify-between gap-2">
							<div class="flex items-center gap-2.5 min-w-0">
								<span class="flex h-5 w-5 shrink-0 items-center justify-center rounded bg-elevated/50 font-mono text-[11px] font-bold text-secondary">
									#{index + 1}
								</span>
								<div class="flex items-center gap-2 truncate">
									<span class="truncate font-mono text-xs font-medium text-primary group-hover:text-accent-perf transition-colors">
										{name}
									</span>
									<span
										class="shrink-0 rounded px-1.5 py-0.5 font-mono text-[10px] font-semibold uppercase tracking-wider"
										style="background-color: {color}15; color: {color}; border: 1px solid {color}30;"
									>
										{item.provider || 'AI'}
									</span>
								</div>
							</div>

							<div class="flex items-center gap-3 shrink-0">
								<span class="font-mono text-xs text-secondary">
									{item.request_count} <span class="text-[10px] text-tertiary">reqs</span>
								</span>
								<span class="font-mono text-xs font-bold text-primary w-12 text-right">
									{item.percentage.toFixed(1)}%
								</span>
							</div>
						</div>

						<!-- Progress Bar -->
						<div class="relative h-2 w-full overflow-hidden rounded-full bg-elevated/40">
							<div
								class="h-full rounded-full transition-all duration-700 ease-out"
								style="width: {Math.max(item.percentage, 2)}%; background: linear-gradient(90deg, {color}cc, {color});"
							></div>
						</div>

						<!-- Mini 7-day spark bars if daily counts available -->
						{#if item.daily_counts && item.daily_counts.length > 0}
							{@const maxDay = Math.max(...item.daily_counts.map((d: { count: number }) => d.count), 1)}
							<div class="mt-1 flex items-end justify-between gap-1 pt-1 border-t border-elevated/20">
								<span class="text-[9px] uppercase tracking-wider text-tertiary">7-day trend</span>
								<div class="flex items-end gap-1 h-3.5">
									{#each item.daily_counts as day}
										{@const heightPct = Math.max((day.count / maxDay) * 100, 15)}
										<div
											class="w-1.5 rounded-xs transition-all duration-300"
											title="{day.date}: {day.count} requests"
											style="height: {day.count > 0 ? heightPct : 15}%; background-color: {day.count > 0 ? color : 'var(--bg-elevated)'}; opacity: {day.count > 0 ? 0.85 : 0.3};"
										></div>
									{/each}
								</div>
							</div>
						{/if}
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
