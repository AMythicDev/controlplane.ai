<script lang="ts">
	// State for inputs
	let budgetGlobalDaily = $state(0);
	let budgetGlobalMonthly = $state(0);
	let budgetUserDaily = $state(0);
	let budgetUserMonthly = $state(0);
	let isSaving = $state(false);
	let saveMessage = $state('');

	async function saveLimits() {
		isSaving = true;
		saveMessage = '';
		try {
			const res = await fetch('/v1/config', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					per_user_daily_limit: Math.round(budgetUserDaily * 1000000),
					per_user_monthly_limit: Math.round(budgetUserMonthly * 1000000)
				})
			});
			if (res.ok) {
				saveMessage = 'Limits saved successfully';
			} else {
				saveMessage = 'Failed to save limits';
			}
		} catch (e) {
			console.error('Failed to save limits', e);
			saveMessage = 'Failed to save limits';
		} finally {
			isSaving = false;
			setTimeout(() => {
				saveMessage = '';
			}, 3000);
		}
	}

	let providers = $state([
		{
			id: 'openai',
			name: 'OpenAI',
			connected: true,
			models: 12,
			enabled: true,
			color: 'bg-emerald-500'
		},
		{
			id: 'anthropic',
			name: 'Anthropic',
			connected: true,
			models: 4,
			enabled: true,
			color: 'bg-orange-500'
		},
		{
			id: 'google',
			name: 'Google',
			connected: false,
			models: 3,
			enabled: false,
			color: 'bg-blue-500'
		},
		{
			id: 'mistral',
			name: 'Mistral',
			connected: true,
			models: 5,
			enabled: true,
			color: 'bg-yellow-500'
		}
	]);

	let notifications = $state([
		{
			id: 'slack',
			name: 'Slack',
			type: 'webhook',
			value: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX',
			enabled: true
		},
		{
			id: 'pagerduty',
			name: 'PagerDuty',
			type: 'integration key',
			value: 'R03X2B1A9...',
			enabled: true
		},
		{
			id: 'email',
			name: 'Email Alerts',
			type: 'email addresses',
			value: 'security@controlplane.ai',
			enabled: true
		},
		{ id: 'webhook', name: 'Custom Webhook', type: 'url', value: '', enabled: false }
	]);

	import { onMount } from 'svelte';
	onMount(async () => {
		try {
			const res = await fetch('/v1/config');
			if (res.ok) {
				const data = await res.json();
				if (data.per_user_daily_limit !== undefined) {
					budgetUserDaily = data.per_user_daily_limit / 1000000;
				}
				if (data.per_user_monthly_limit !== undefined) {
					budgetUserMonthly = data.per_user_monthly_limit / 1000000;
				}
			}
		} catch (e) {
			console.error('Failed to fetch config', e);
		}
	});
</script>

<svelte:head>
	<title>Settings | ControlPlane.ai</title>
</svelte:head>

<div class="mx-auto min-h-screen max-w-5xl space-y-12 p-6 pb-20 md:p-10">
	<!-- Header -->
	<header class="animate-fade-up">
		<h1 class="font-serif text-5xl tracking-tight text-primary md:text-6xl">Settings</h1>
		<p class="mt-2 font-mono text-sm tracking-wide text-secondary">
			Configure engine parameters and integrations
		</p>
	</header>

	<!-- API Config -->
	<section class="animate-fade-up" style="animation-delay: 100ms;">
		<h2 class="mb-6 font-serif text-3xl text-primary">API Configuration</h2>
		<div
			class="relative overflow-hidden rounded-xl border border-elevated bg-surface p-6 shadow-inner"
		>
			<!-- Noise -->
			<div
				class="pointer-events-none absolute inset-0 opacity-[0.01]"
				style="background-image: url(&quot;data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E&quot;);"
			></div>

			<div class="relative z-10 space-y-6">
				<div>
					<label class="mb-2 block font-mono text-xs tracking-widest text-secondary uppercase"
						>API Key</label
					>
					<div class="flex gap-3">
						<div
							class="flex flex-1 items-center justify-between rounded-lg border border-elevated bg-[#09080a] px-4 py-2.5 font-mono text-primary"
						>
							<span>sk-cp-xxxxxxxxxxxxxxxxxxxxxxxx</span>
							<button class="text-sm text-accent-perf transition-colors hover:text-primary"
								>Copy</button
							>
						</div>
						<button
							class="rounded-lg border border-elevated bg-elevated px-4 py-2 text-sm font-medium text-primary transition-colors hover:bg-hover"
						>
							Regenerate
						</button>
					</div>
				</div>

				<div>
					<label class="mb-2 block font-mono text-xs tracking-widest text-secondary uppercase"
						>Endpoint URL</label
					>
					<div
						class="flex-1 rounded-lg border border-elevated bg-[#09080a] px-4 py-2.5 font-mono text-primary"
					>
						https://checker.controlplane.ai/v1
					</div>
				</div>

				<div class="flex items-center gap-2 pt-2">
					<span class="h-2 w-2 rounded-full bg-accent-resp shadow-[0_0_8px_rgba(0,230,118,0.6)]"
					></span>
					<span class="font-mono text-sm text-secondary">Connected & Active</span>
				</div>
			</div>
		</div>
	</section>

	<!-- LLM Providers -->
	<section class="animate-fade-up" style="animation-delay: 200ms;">
		<h2 class="mb-6 font-serif text-3xl text-primary">LLM Providers</h2>
		<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
			{#each providers as provider}
				<div
					class="flex items-center justify-between rounded-xl border border-elevated bg-surface p-5 transition-colors hover:border-hover"
				>
					<div class="flex items-center gap-4">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-full font-bold text-white shadow-inner {provider.color}"
						>
							{provider.name.charAt(0)}
						</div>
						<div>
							<h3 class="font-bold text-primary">{provider.name}</h3>
							<div class="mt-1 flex items-center gap-2 font-mono text-xs text-secondary">
								{#if provider.connected}
									<span class="text-accent-resp">● Connected</span>
									<span>• {provider.models} models</span>
								{:else}
									<span class="text-tertiary">○ Disconnected</span>
								{/if}
							</div>
						</div>
					</div>

					<button
						class="focus-ring relative h-6 w-11 rounded-full transition-colors duration-200 ease-in-out {provider.enabled
							? 'bg-accent-perf'
							: 'bg-elevated'}"
						onclick={() => (provider.enabled = !provider.enabled)}
						role="switch"
						aria-checked={provider.enabled}
					>
						<span
							class="absolute top-[2px] left-[2px] h-5 w-5 rounded-full bg-white transition-transform duration-200 {provider.enabled
								? 'translate-x-5'
								: 'translate-x-0'} shadow-sm"
						></span>
					</button>
				</div>
			{/each}
		</div>
	</section>

	<!-- Budget Limits -->
	<section class="animate-fade-up" style="animation-delay: 300ms;">
		<h2 class="mb-6 font-serif text-3xl text-primary">Budget Limits</h2>
		<div class="rounded-xl border border-elevated bg-surface p-6">
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				<div class="space-y-2">
					<label class="block font-mono text-xs tracking-widest text-secondary uppercase"
						>Global Daily Limit</label
					>
					<div class="relative">
						<span class="absolute top-1/2 left-4 -translate-y-1/2 font-mono text-secondary">$</span>
						<input
							type="number"
							bind:value={budgetGlobalDaily}
							class="focus-ring w-full rounded-lg border border-elevated bg-[#09080a] py-2.5 pr-4 pl-8 font-mono text-primary transition-shadow"
						/>
					</div>
				</div>
				<div class="space-y-2">
					<label class="block font-mono text-xs tracking-widest text-secondary uppercase"
						>Global Monthly Limit</label
					>
					<div class="relative">
						<span class="absolute top-1/2 left-4 -translate-y-1/2 font-mono text-secondary">$</span>
						<input
							type="number"
							bind:value={budgetGlobalMonthly}
							class="focus-ring w-full rounded-lg border border-elevated bg-[#09080a] py-2.5 pr-4 pl-8 font-mono text-primary transition-shadow"
						/>
					</div>
				</div>
				<div class="space-y-2">
					<label class="block font-mono text-xs tracking-widest text-secondary uppercase"
						>Per-User Daily Limit</label
					>
					<div class="relative">
						<span class="absolute top-1/2 left-4 -translate-y-1/2 font-mono text-secondary">$</span>
						<input
							type="number"
							bind:value={budgetUserDaily}
							class="focus-ring w-full rounded-lg border border-elevated bg-[#09080a] py-2.5 pr-4 pl-8 font-mono text-primary transition-shadow"
						/>
					</div>
				</div>
				<div class="space-y-2">
					<label class="block font-mono text-xs tracking-widest text-secondary uppercase"
						>Per-User Monthly Limit</label
					>
					<div class="relative">
						<span class="absolute top-1/2 left-4 -translate-y-1/2 font-mono text-secondary">$</span>
						<input
							type="number"
							bind:value={budgetUserMonthly}
							class="focus-ring w-full rounded-lg border border-elevated bg-[#09080a] py-2.5 pr-4 pl-8 font-mono text-primary transition-shadow"
						/>
					</div>
				</div>
			</div>

			<div class="mt-6 flex items-center justify-end gap-4">
				{#if saveMessage}
					<span
						class="font-mono text-xs {saveMessage.includes('success')
							? 'text-accent-perf'
							: 'text-accent-danger'}">{saveMessage}</span
					>
				{/if}
				<button
					onclick={saveLimits}
					disabled={isSaving}
					class="rounded-lg bg-primary px-6 py-2.5 font-bold text-base transition-colors hover:bg-white disabled:opacity-50"
				>
					{isSaving ? 'Saving...' : 'Save Limits'}
				</button>
			</div>
		</div>
	</section>

	<!-- Notifications -->
	<section class="animate-fade-up" style="animation-delay: 400ms;">
		<h2 class="mb-6 font-serif text-3xl text-primary">Notifications</h2>
		<div class="overflow-hidden rounded-xl border border-elevated bg-surface">
			{#each notifications as notif, i}
				<div
					class="flex flex-col gap-4 p-5 md:flex-row md:items-center {i !== notifications.length - 1
						? 'border-b border-elevated'
						: ''}"
				>
					<div class="w-32">
						<span class="font-bold text-primary">{notif.name}</span>
					</div>

					<div class="flex-1">
						<input
							type="text"
							bind:value={notif.value}
							placeholder={`Enter ${notif.type}...`}
							class="focus-ring w-full rounded-lg border border-elevated bg-[#09080a] px-4 py-2 font-mono text-sm text-primary transition-shadow placeholder:text-tertiary"
						/>
					</div>

					<div class="flex w-12 items-center justify-end">
						<button
							class="focus-ring relative h-6 w-11 rounded-full transition-colors duration-200 ease-in-out {notif.enabled
								? 'bg-accent-perf'
								: 'bg-elevated'}"
							onclick={() => (notif.enabled = !notif.enabled)}
							role="switch"
							aria-checked={notif.enabled}
						>
							<span
								class="absolute top-[2px] left-[2px] h-5 w-5 rounded-full bg-white transition-transform duration-200 {notif.enabled
									? 'translate-x-5'
									: 'translate-x-0'} shadow-sm"
							></span>
						</button>
					</div>
				</div>
			{/each}
		</div>
	</section>
</div>
