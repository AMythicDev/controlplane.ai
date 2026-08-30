<script lang="ts">
	import type { ModelAnalyticsItem } from '$lib/types';
	import ShieldAlert from 'lucide-svelte/icons/shield-alert';
	import AlertCircle from 'lucide-svelte/icons/alert-circle';

	let { models = [] } = $props<{
		models: ModelAnalyticsItem[];
	}>();

	function getHallucinationRating(ratePct: number) {
		if (ratePct <= 10) return { label: 'Low Risk', color: '#98DA15', bg: 'rgba(152, 218, 21, 0.1)', border: 'rgba(152, 218, 21, 0.3)' };
		if (ratePct <= 30) return { label: 'Moderate', color: '#ff9100', bg: 'rgba(255, 145, 0, 0.1)', border: 'rgba(255, 145, 0, 0.3)' };
		return { label: 'High Risk', color: '#ff1744', bg: 'rgba(255, 23, 68, 0.1)', border: 'rgba(255, 23, 68, 0.3)' };
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
	<div class="pointer-events-none absolute top-0 right-0 -z-10 h-48 w-48 rounded-full bg-accent-resp/5 blur-3xl"></div>

	<div>
		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-2.5">
				<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-resp/10 text-accent-resp">
					<ShieldAlert class="h-4 w-4" />
				</div>
				<div>
					<h3 class="font-sans text-base font-semibold text-primary">Model Hallucination Score</h3>
					<p class="text-xs text-secondary">NLI contradiction probability (lower is safer)</p>
				</div>
			</div>
			<div class="flex items-center gap-1 text-[11px] font-mono text-tertiary">
				<span>DeBERTa NLI</span>
			</div>
		</div>

		<div class="mt-6 flex flex-col gap-3.5">
			{#if models.length === 0}
				<div class="flex flex-col items-center justify-center py-10 text-center">
					<div class="mb-2 rounded-full bg-elevated/30 p-3 text-secondary">
						<AlertCircle class="h-5 w-5 opacity-50" />
					</div>
					<p class="font-mono text-xs text-secondary">No NLI verification data available</p>
					<span class="mt-1 text-[11px] text-tertiary">Responses checked by NLI scanner will appear here</span>
				</div>
			{:else}
				{#each models as item}
					{@const ratePct = item.avg_hallucination !== null && item.avg_hallucination !== undefined ? Math.round((item.avg_hallucination > 1 ? item.avg_hallucination : item.avg_hallucination * 100)) : null}
					{@const rating = ratePct !== null ? getHallucinationRating(ratePct) : null}
					{@const name = cleanModelName(item.model)}
					<div class="group flex flex-col gap-2 rounded-xl border border-elevated/20 bg-base/30 p-3.5 transition-all duration-200 hover:border-elevated/50 hover:bg-base/50">
						<div class="flex items-center justify-between gap-2">
							<div class="flex items-center gap-2 truncate">
								<span class="truncate font-mono text-xs font-medium text-primary group-hover:text-accent-resp transition-colors">
									{name}
								</span>
								{#if item.nli_count > 0}
									<span class="text-[10px] text-tertiary font-mono">
										({item.nli_count} verified)
									</span>
								{/if}
							</div>

							<div class="flex items-center gap-2 shrink-0">
								{#if ratePct !== null && rating}
									<span
										class="rounded px-2 py-0.5 font-mono text-[11px] font-bold"
										style="color: {rating.color}; background-color: {rating.bg}; border: 1px solid {rating.border};"
									>
										{ratePct}%
									</span>
									<span class="text-[10px] uppercase font-mono tracking-wider font-semibold" style="color: {rating.color}">
										{rating.label}
									</span>
								{:else}
									<span class="rounded px-2 py-0.5 font-mono text-[10px] text-tertiary bg-elevated/40 border border-elevated/60">
										N/A
									</span>
								{/if}
							</div>
						</div>

						<!-- Risk Meter Bar -->
						<div class="relative h-2 w-full overflow-hidden rounded-full bg-elevated/40">
							{#if ratePct !== null && rating}
								<div
									class="h-full rounded-full transition-all duration-700 ease-out"
									style="width: {Math.max(ratePct, 2)}%; background: linear-gradient(90deg, {rating.color}99, {rating.color});"
								></div>
							{:else}
								<div class="h-full w-full rounded-full bg-elevated/20 opacity-40"></div>
							{/if}
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
