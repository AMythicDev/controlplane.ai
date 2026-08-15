<script lang="ts">
	// State for inputs
	let budgetGlobalDaily = $state(500);
	let budgetGlobalMonthly = $state(10000);
	let budgetUserDaily = $state(50);
	let budgetUserMonthly = $state(1000);
	
	let providers = $state([
		{ id: 'openai', name: 'OpenAI', connected: true, models: 12, enabled: true, color: 'bg-emerald-500' },
		{ id: 'anthropic', name: 'Anthropic', connected: true, models: 4, enabled: true, color: 'bg-orange-500' },
		{ id: 'google', name: 'Google', connected: false, models: 3, enabled: false, color: 'bg-blue-500' },
		{ id: 'mistral', name: 'Mistral', connected: true, models: 5, enabled: true, color: 'bg-yellow-500' }
	]);
	
	let notifications = $state([
		{ id: 'slack', name: 'Slack', type: 'webhook', value: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX', enabled: true },
		{ id: 'pagerduty', name: 'PagerDuty', type: 'integration key', value: 'R03X2B1A9...', enabled: true },
		{ id: 'email', name: 'Email Alerts', type: 'email addresses', value: 'security@controlplane.ai', enabled: true },
		{ id: 'webhook', name: 'Custom Webhook', type: 'url', value: '', enabled: false }
	]);
</script>

<svelte:head>
	<title>Settings | ControlPlane.ai</title>
</svelte:head>

<div class="max-w-5xl mx-auto p-6 md:p-10 space-y-12 min-h-screen pb-20">
	<!-- Header -->
	<header class="animate-fade-up">
		<h1 class="font-serif text-5xl md:text-6xl text-primary tracking-tight">Settings</h1>
		<p class="text-secondary font-mono text-sm tracking-wide mt-2">
			Configure engine parameters and integrations
		</p>
	</header>

	<!-- API Config -->
	<section class="animate-fade-up" style="animation-delay: 100ms;">
		<h2 class="font-serif text-3xl text-primary mb-6">API Configuration</h2>
		<div class="bg-surface border border-elevated rounded-xl p-6 shadow-inner relative overflow-hidden">
			<!-- Noise -->
			<div class="absolute inset-0 opacity-[0.01] pointer-events-none" style="background-image: url(&quot;data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.65' numOctaves='3' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E&quot;);"></div>
			
			<div class="relative z-10 space-y-6">
				<div>
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary mb-2">API Key</label>
					<div class="flex gap-3">
						<div class="flex-1 bg-[#09080a] border border-elevated rounded-lg px-4 py-2.5 font-mono text-primary flex items-center justify-between">
							<span>sk-cp-xxxxxxxxxxxxxxxxxxxxxxxx</span>
							<button class="text-accent-perf hover:text-primary transition-colors text-sm">Copy</button>
						</div>
						<button class="px-4 py-2 bg-elevated hover:bg-hover border border-elevated rounded-lg text-primary text-sm font-medium transition-colors">
							Regenerate
						</button>
					</div>
				</div>
				
				<div>
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary mb-2">Endpoint URL</label>
					<div class="flex-1 bg-[#09080a] border border-elevated rounded-lg px-4 py-2.5 font-mono text-primary">
						https://checker.controlplane.ai/v1
					</div>
				</div>
				
				<div class="flex items-center gap-2 pt-2">
					<span class="w-2 h-2 rounded-full bg-accent-resp shadow-[0_0_8px_rgba(0,230,118,0.6)]"></span>
					<span class="text-sm font-mono text-secondary">Connected & Active</span>
				</div>
			</div>
		</div>
	</section>

	<!-- LLM Providers -->
	<section class="animate-fade-up" style="animation-delay: 200ms;">
		<h2 class="font-serif text-3xl text-primary mb-6">LLM Providers</h2>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			{#each providers as provider}
				<div class="bg-surface border border-elevated rounded-xl p-5 flex items-center justify-between hover:border-hover transition-colors">
					<div class="flex items-center gap-4">
						<div class="w-10 h-10 rounded-full flex items-center justify-center font-bold text-white shadow-inner {provider.color}">
							{provider.name.charAt(0)}
						</div>
						<div>
							<h3 class="font-bold text-primary">{provider.name}</h3>
							<div class="flex items-center gap-2 text-xs font-mono text-secondary mt-1">
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
						class="relative w-11 h-6 rounded-full transition-colors duration-200 ease-in-out focus-ring {provider.enabled ? 'bg-accent-perf' : 'bg-elevated'}"
						onclick={() => provider.enabled = !provider.enabled}
						role="switch"
						aria-checked={provider.enabled}
					>
						<span class="absolute top-[2px] left-[2px] bg-white w-5 h-5 rounded-full transition-transform duration-200 {provider.enabled ? 'translate-x-5' : 'translate-x-0'} shadow-sm"></span>
					</button>
				</div>
			{/each}
		</div>
	</section>

	<!-- Budget Limits -->
	<section class="animate-fade-up" style="animation-delay: 300ms;">
		<h2 class="font-serif text-3xl text-primary mb-6">Budget Limits</h2>
		<div class="bg-surface border border-elevated rounded-xl p-6">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div class="space-y-2">
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary">Global Daily Limit</label>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 text-secondary font-mono">$</span>
						<input type="number" bind:value={budgetGlobalDaily} class="w-full bg-[#09080a] border border-elevated rounded-lg pl-8 pr-4 py-2.5 font-mono text-primary focus-ring transition-shadow" />
					</div>
				</div>
				<div class="space-y-2">
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary">Global Monthly Limit</label>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 text-secondary font-mono">$</span>
						<input type="number" bind:value={budgetGlobalMonthly} class="w-full bg-[#09080a] border border-elevated rounded-lg pl-8 pr-4 py-2.5 font-mono text-primary focus-ring transition-shadow" />
					</div>
				</div>
				<div class="space-y-2">
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary">Per-User Daily Limit</label>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 text-secondary font-mono">$</span>
						<input type="number" bind:value={budgetUserDaily} class="w-full bg-[#09080a] border border-elevated rounded-lg pl-8 pr-4 py-2.5 font-mono text-primary focus-ring transition-shadow" />
					</div>
				</div>
				<div class="space-y-2">
					<label class="block text-xs font-mono uppercase tracking-widest text-secondary">Per-User Monthly Limit</label>
					<div class="relative">
						<span class="absolute left-4 top-1/2 -translate-y-1/2 text-secondary font-mono">$</span>
						<input type="number" bind:value={budgetUserMonthly} class="w-full bg-[#09080a] border border-elevated rounded-lg pl-8 pr-4 py-2.5 font-mono text-primary focus-ring transition-shadow" />
					</div>
				</div>
			</div>
			
			<div class="mt-6 flex justify-end">
				<button class="px-6 py-2.5 bg-primary text-base font-bold rounded-lg hover:bg-white transition-colors">
					Save Limits
				</button>
			</div>
		</div>
	</section>

	<!-- Notifications -->
	<section class="animate-fade-up" style="animation-delay: 400ms;">
		<h2 class="font-serif text-3xl text-primary mb-6">Notifications</h2>
		<div class="bg-surface border border-elevated rounded-xl overflow-hidden">
			{#each notifications as notif, i}
				<div class="p-5 flex flex-col md:flex-row md:items-center gap-4 {i !== notifications.length - 1 ? 'border-b border-elevated' : ''}">
					<div class="w-32">
						<span class="font-bold text-primary">{notif.name}</span>
					</div>
					
					<div class="flex-1">
						<input 
							type="text" 
							bind:value={notif.value} 
							placeholder={`Enter ${notif.type}...`}
							class="w-full bg-[#09080a] border border-elevated rounded-lg px-4 py-2 font-mono text-sm text-primary focus-ring transition-shadow placeholder:text-tertiary" 
						/>
					</div>
					
					<div class="flex items-center justify-end w-12">
						<button 
							class="relative w-11 h-6 rounded-full transition-colors duration-200 ease-in-out focus-ring {notif.enabled ? 'bg-accent-perf' : 'bg-elevated'}"
							onclick={() => notif.enabled = !notif.enabled}
							role="switch"
							aria-checked={notif.enabled}
						>
							<span class="absolute top-[2px] left-[2px] bg-white w-5 h-5 rounded-full transition-transform duration-200 {notif.enabled ? 'translate-x-5' : 'translate-x-0'} shadow-sm"></span>
						</button>
					</div>
				</div>
			{/each}
		</div>
	</section>
</div>
