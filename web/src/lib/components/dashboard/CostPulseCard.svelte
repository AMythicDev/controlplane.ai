<script lang="ts">
	import Zap from 'lucide-svelte/icons/zap';
	import Activity from 'lucide-svelte/icons/activity';
	import Flame from 'lucide-svelte/icons/flame';
	import ShieldCheck from 'lucide-svelte/icons/shield-check';

	let { liveCost = 0, liveSavings = 0, userMonthlyLimit = 0, userDailyLimit = 0 } = $props<{
		liveCost: number;
		liveSavings: number;
		userMonthlyLimit: number;
		userDailyLimit: number;
	}>();

	let burnPercent = $derived(
		userMonthlyLimit > 0 ? Math.min(100, Math.round((liveCost / userMonthlyLimit) * 100)) : 0
	);

	function formatNum(val: number) {
		return new Intl.NumberFormat('en-US', { maximumFractionDigits: 2 }).format(val);
	}
</script>

<div class="relative flex flex-col justify-between overflow-hidden rounded-2xl border border-elevated/40 bg-surface/40 p-6 backdrop-blur-md transition-all duration-300 hover:border-elevated/70">
	<!-- Background glow -->
	<div class="pointer-events-none absolute top-0 right-0 -z-10 h-48 w-48 rounded-full bg-accent-cost/5 blur-3xl"></div>

	<div>
		<div class="flex items-center justify-between gap-3">
			<div class="flex items-center gap-2.5">
				<div class="flex h-8 w-8 items-center justify-center rounded-lg bg-accent-cost/10 text-accent-cost">
					<Activity class="h-4 w-4" />
				</div>
				<div>
					<h3 class="font-sans text-base font-semibold text-primary">Cost Pulse & Budget Burn</h3>
					<p class="text-xs text-secondary">Real-time spend pacing against configured caps</p>
				</div>
			</div>
			<div class="flex items-center gap-1.5 rounded-full border border-elevated/60 bg-base/60 px-3 py-1 font-mono text-xs text-secondary">
				{#if userMonthlyLimit > 0}
					<span class="text-accent-cost font-bold">{burnPercent}%</span>
					<span>budget used</span>
				{:else}
					<span class="text-accent-cost font-bold">Uncapped</span>
					<span>no limit</span>
				{/if}
			</div>
		</div>

		<!-- Pulse Meter Zone -->
		<div class="mt-6 flex flex-col gap-5">
			<!-- Burn gauge progress bar with markers -->
			<div class="flex flex-col gap-2">
				<div class="flex justify-between items-baseline text-xs font-mono">
					<span class="text-secondary">Monthly Spend Pacing</span>
					<span class="font-bold text-primary">
						{#if userMonthlyLimit > 0}
							${formatNum(liveCost)} <span class="text-tertiary">/ ${formatNum(userMonthlyLimit)}</span>
						{:else}
							${formatNum(liveCost)} <span class="text-tertiary">/ Unlimited</span>
						{/if}
					</span>
				</div>

				<div class="relative h-3 w-full overflow-hidden rounded-full bg-elevated/50 p-0.5 border border-elevated/40">
					{#if userMonthlyLimit > 0}
						<div
							class="h-full rounded-full transition-all duration-1000 ease-out"
							style="width: {Math.max(burnPercent, 1)}%; background: linear-gradient(90deg, var(--accent-cost), #c0eb38);"
						></div>
					{:else}
						<div class="h-full w-full rounded-full bg-accent-cost/20 opacity-50"></div>
					{/if}
				</div>

				<div class="flex justify-between text-[10px] font-mono text-tertiary">
					<span>$0.00</span>
					{#if userMonthlyLimit > 0}
						<span>${formatNum(userMonthlyLimit / 2)} (50%)</span>
						<span>${formatNum(userMonthlyLimit)} Cap</span>
					{:else}
						<span>No Monthly Cap Configured</span>
						<span>Unlimited</span>
					{/if}
				</div>
			</div>

			<!-- Metrics Grid -->
			<div class="grid grid-cols-2 gap-3 pt-2">
				<!-- Cache Savings -->
				<div class="flex flex-col gap-1 rounded-xl border border-elevated/30 bg-base/40 p-3.5">
					<div class="flex items-center gap-1.5 text-xs text-secondary">
						<Zap class="h-3.5 w-3.5 text-accent-perf" />
						<span>Semantic Cache Savings</span>
					</div>
					<div class="font-mono text-xl font-bold text-accent-perf">
						${formatNum(liveSavings)}
					</div>
					<span class="text-[10px] text-tertiary font-mono">Saved via zero-cost cache hits</span>
				</div>

				<!-- Daily Limit -->
				<div class="flex flex-col gap-1 rounded-xl border border-elevated/30 bg-base/40 p-3.5">
					<div class="flex items-center gap-1.5 text-xs text-secondary">
						<ShieldCheck class="h-3.5 w-3.5 text-accent-cost" />
						<span>Daily Cap</span>
					</div>
					<div class="font-mono text-xl font-bold text-primary">
						{userDailyLimit > 0 ? `$${formatNum(userDailyLimit)}` : 'Unlimited'}
					</div>
					<span class="text-[10px] text-tertiary font-mono">{userDailyLimit > 0 ? 'Per-user daily limit' : 'No daily limit set'}</span>
				</div>
			</div>
		</div>
	</div>
</div>
