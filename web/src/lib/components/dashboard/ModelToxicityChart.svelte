<script lang="ts">
	import type { ModelAnalyticsItem } from '$lib/types';
	import Shield from 'lucide-svelte/icons/shield';
	import CheckCircle2 from 'lucide-svelte/icons/check-circle-2';

	let { models = [] } = $props<{
		models: ModelAnalyticsItem[];
	}>();

	function getToxicityRating(toxPct: number) {
		if (toxPct <= 5) return { label: 'Clean', color: '#98DA15', bg: 'rgba(152, 218, 21, 0.1)', border: 'rgba(152, 218, 21, 0.3)' };
		if (toxPct <= 15) return { label: 'Guarded', color: '#ff9100', bg: 'rgba(255, 145, 0, 0.1)', border: 'rgba(255, 145, 0, 0.3)' };
		return { label: 'Flagged', color: '#ff1744', bg: 'rgba(255, 23, 68, 0.1)', border: 'rgba(255, 23, 68, 0.3)' };
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
	<div class="pointer-events-none absolute top-0 right-0 -z-10 h-48 w-48 rounded-full bg-accent-danger/5 blur-3xl"></div>

	<div>
		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-2.5">
				<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-danger/10 text-accent-danger">
					<Shield class="h-4 w-4" />
				</div>
				<div>
					<h3 class="font-sans text-base font-semibold text-primary">Model Toxicity Score</h3>
					<p class="text-xs text-secondary">Detoxify safety assessment (lower is cleaner)</p>
				</div>
			</div>
			<div class="flex items-center gap-1 text-[11px] font-mono text-tertiary">
				<span>Detoxify ONNX</span>
			</div>
		</div>

		<div class="mt-6 flex flex-col gap-3.5">
			{#if models.length === 0}
				<div class="flex flex-col items-center justify-center py-10 text-center">
					<div class="mb-2 rounded-full bg-elevated/30 p-3 text-secondary">
						<CheckCircle2 class="h-5 w-5 opacity-50" />
					</div>
					<p class="font-mono text-xs text-secondary">No toxicity telemetry available</p>
					<span class="mt-1 text-[11px] text-tertiary">Responses scanned for toxicity will appear here</span>
				</div>
			{:else}
				{#each models as item}
					{@const toxPct = item.avg_toxicity !== null && item.avg_toxicity !== undefined ? Math.round((item.avg_toxicity > 1 ? item.avg_toxicity : item.avg_toxicity * 100)) : 0}
					{@const rating = getToxicityRating(toxPct)}
					{@const name = cleanModelName(item.model)}
					<div class="group flex flex-col gap-2 rounded-xl border border-elevated/20 bg-base/30 p-3.5 transition-all duration-200 hover:border-elevated/50 hover:bg-base/50">
						<div class="flex items-center justify-between gap-2">
							<div class="flex items-center gap-2 truncate">
								<span class="truncate font-mono text-xs font-medium text-primary group-hover:text-accent-perf transition-colors">
									{name}
								</span>
								{#if item.toxicity_count > 0}
									<span class="text-[10px] text-tertiary font-mono">
										({item.toxicity_count} scanned)
									</span>
								{/if}
							</div>

							<div class="flex items-center gap-2 shrink-0">
								<span
									class="rounded px-2 py-0.5 font-mono text-[11px] font-bold"
									style="color: {rating.color}; background-color: {rating.bg}; border: 1px solid {rating.border};"
								>
									{toxPct}%
								</span>
								<span class="text-[10px] uppercase font-mono tracking-wider font-semibold" style="color: {rating.color}">
									{rating.label}
								</span>
							</div>
						</div>

						<!-- Safety Meter Bar -->
						<div class="relative h-2 w-full overflow-hidden rounded-full bg-elevated/40">
							<div
								class="h-full rounded-full transition-all duration-700 ease-out"
								style="width: {Math.max(toxPct, 2)}%; background: linear-gradient(90deg, {rating.color}99, {rating.color});"
							></div>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
