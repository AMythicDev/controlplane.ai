<script lang="ts">
    import { modelCatalog, runPlaygroundRequest, type MockResult, type ModelOption } from '$lib/data/models';
    import ResultCard from '$lib/components/playground/ResultCard.svelte';
    import Play from 'lucide-svelte/icons/play';
    import ChevronDown from 'lucide-svelte/icons/chevron-down';
    import Check from 'lucide-svelte/icons/check';
    import FlaskConical from 'lucide-svelte/icons/flask-conical';

    // State
    let prompt = $state('');
    let selectedModels = $state<Set<string>>(new Set(['openai/gpt-4o', 'anthropic/claude-sonnet-4-20250514']));
    let isRunning = $state(false);
    let globalError = $state<string | null>(null);
    
    // Results store per model spec
    let results = $state<Record<string, MockResult | null>>({});
    let loadingStates = $state<Record<string, boolean>>({});

    // Provider expansion state for accordion
    let expandedProviders = $state<Set<string>>(new Set(['openai', 'anthropic']));

    // Derived
    let selectedCount = $derived(selectedModels.size);
    
    // Get all selected model objects for rendering the grid
    let activeModels = $derived.by(() => {
        const active: { model: ModelOption, providerColor: string, providerName: string }[] = [];
        for (const provider of modelCatalog) {
            for (const model of provider.models) {
                if (selectedModels.has(model.spec)) {
                    active.push({
                        model,
                        providerColor: provider.color,
                        providerName: provider.name
                    });
                }
            }
        }
        return active;
    });

    // Actions
    function toggleModel(spec: string) {
        const newSet = new Set(selectedModels);
        if (newSet.has(spec)) {
            newSet.delete(spec);
        } else {
            newSet.add(spec);
        }
        selectedModels = newSet;
    }

    function toggleProvider(providerId: string) {
        const newSet = new Set(expandedProviders);
        if (newSet.has(providerId)) {
            newSet.delete(providerId);
        } else {
            newSet.add(providerId);
        }
        expandedProviders = newSet;
    }

    async function runPlayground() {
        if (!prompt.trim() || selectedModels.size === 0) return;
        
        isRunning = true;
        globalError = null;
        
        // Reset state for selected models
        const specs = Array.from(selectedModels);
        for (const spec of specs) {
            results[spec] = null;
            loadingStates[spec] = true;
        }

        // Run all selected models in parallel
        const promises = specs.map(async (spec) => {
            try {
                const res = await runPlaygroundRequest(prompt, spec);
                results[spec] = res;
            } catch (err: any) {
                console.error(`Error with ${spec}:`, err);
                if (err.message && err.message.toLowerCase().includes('budget')) {
                    globalError = err.message;
                } else if (!globalError) {
                    globalError = `Error with ${spec}: ${err.message}`;
                }
            } finally {
                loadingStates[spec] = false;
            }
        });

        await Promise.all(promises);
        isRunning = false;
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter' && (e.ctrlKey || e.metaKey)) {
            e.preventDefault();
            runPlayground();
        }
    }
</script>

<div class="relative z-10 flex min-h-full w-full flex-col gap-8 pb-32">
    <!-- Radial gradient background -->
    <div
        class="pointer-events-none absolute top-0 left-1/2 -z-10 h-[60vh] w-[80vw] -translate-x-1/2 rounded-full bg-[radial-gradient(ellipse_at_top,rgba(var(--color-accent-perf),0.1)_0%,transparent_70%)] blur-3xl"
    ></div>

    <!-- Header -->
    <header class="flex flex-col gap-6 px-6 lg:px-12 pt-8">
        <div class="flex items-end justify-between">
            <div>
                <h1 class="font-serif text-4xl tracking-wide text-primary italic mb-2">LLM Playground</h1>
                <p class="font-mono text-xs text-secondary uppercase tracking-widest">Multi-Model Orchestration & Evaluation</p>
            </div>
        </div>

        {#if globalError}
            <div class="bg-accent-danger/10 border border-accent-danger/30 rounded-xl p-4 flex items-center gap-4 animate-fade-up">
                <div class="w-8 h-8 rounded-full bg-accent-danger/20 flex items-center justify-center shrink-0">
                    <span class="text-accent-danger font-bold text-lg">!</span>
                </div>
                <div>
                    <h3 class="text-accent-danger font-bold text-sm">Request Failed</h3>
                    <p class="text-secondary text-xs mt-1">{globalError}</p>
                </div>
            </div>
        {/if}
    </header>

    <div class="px-6 lg:px-12 grid grid-cols-1 lg:grid-cols-12 gap-8">
        
        <!-- Left Column: Controls (Span 4) -->
        <div class="lg:col-span-4 flex flex-col gap-6">
            
            <!-- Prompt Area -->
            <div class="flex flex-col gap-3 p-1 rounded-2xl border border-elevated bg-surface/50 backdrop-blur shadow-lg overflow-hidden focus-within:border-accent-perf/50 focus-within:ring-1 focus-within:ring-accent-perf/50 transition-all">
                <textarea
                    bind:value={prompt}
                    onkeydown={handleKeydown}
                    placeholder="Enter your prompt here... (Ctrl+Enter to run)"
                    class="w-full h-40 bg-transparent resize-none p-4 font-sans text-primary placeholder:text-secondary/50 focus:outline-none custom-scrollbar"
                ></textarea>
                
                <div class="flex items-center justify-between p-3 border-t border-elevated/50 bg-base/30">
                    <span class="font-mono text-[10px] text-secondary opacity-60">System prompts not yet supported</span>
                    <button
                        onclick={runPlayground}
                        disabled={isRunning || selectedCount === 0 || !prompt.trim()}
                        class="flex items-center gap-2 px-5 py-2 rounded-full bg-accent-perf text-white font-mono text-xs font-bold uppercase tracking-wider shadow-[0_0_15px_rgba(161,0,255,0.4)] transition-all hover:scale-105 hover:shadow-[0_0_20px_rgba(161,0,255,0.6)] disabled:opacity-50 disabled:pointer-events-none disabled:hover:scale-100 disabled:shadow-none"
                    >
                        {#if isRunning}
                            <div class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                            Running
                        {:else}
                            <Play class="w-4 h-4 fill-white" />
                            Run Models
                        {/if}
                    </button>
                </div>
            </div>

            <!-- Model Selector -->
            <div class="flex flex-col rounded-2xl border border-elevated bg-surface/50 backdrop-blur shadow-lg overflow-hidden">
                <div class="p-4 border-b border-elevated/50 bg-base/30 flex items-center justify-between">
                    <span class="font-mono text-xs text-primary uppercase tracking-widest font-medium">Model Matrix</span>
                    <div class="flex items-center gap-1.5 px-2 py-1 rounded bg-elevated/50 font-mono text-[10px] text-secondary">
                        <span class="text-accent-perf font-bold">{selectedCount}</span> selected
                    </div>
                </div>
                
                <div class="flex flex-col max-h-[500px] overflow-y-auto custom-scrollbar">
                    {#each modelCatalog as provider}
                        {@const isExpanded = expandedProviders.has(provider.id)}
                        {@const selectedInProvider = provider.models.filter(m => selectedModels.has(m.spec)).length}
                        
                        <div class="flex flex-col border-b border-elevated/30 last:border-0">
                            <!-- Provider Header -->
                            <button 
                                onclick={() => toggleProvider(provider.id)}
                                class="flex items-center justify-between p-3 hover:bg-elevated/20 transition-colors w-full text-left"
                            >
                                <div class="flex items-center gap-3">
                                    <ChevronDown class="w-4 h-4 text-secondary transition-transform duration-200 {isExpanded ? '' : '-rotate-90'}" />
                                    <span class="font-mono text-sm text-primary">{provider.name}</span>
                                </div>
                                {#if selectedInProvider > 0}
                                    <div class="w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-mono text-base font-bold" style="background-color: {provider.color}">
                                        {selectedInProvider}
                                    </div>
                                {/if}
                            </button>
                            
                            <!-- Models List -->
                            {#if isExpanded}
                                <div class="flex flex-col bg-base/20 py-1">
                                    {#each provider.models as model}
                                        {@const isSelected = selectedModels.has(model.spec)}
                                        <button 
                                            onclick={() => toggleModel(model.spec)}
                                            class="flex items-center gap-3 py-2 px-10 hover:bg-elevated/30 transition-colors w-full text-left group"
                                        >
                                            <div class="w-4 h-4 rounded border flex items-center justify-center transition-colors {isSelected ? 'bg-accent-perf border-accent-perf' : 'border-elevated/80 group-hover:border-primary/50'}">
                                                {#if isSelected}
                                                    <Check class="w-3 h-3 text-white" />
                                                {/if}
                                            </div>
                                            <span class="font-sans text-sm {isSelected ? 'text-primary font-medium' : 'text-secondary'}">
                                                {model.name}
                                            </span>
                                        </button>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
            
        </div>

        <!-- Right Column: Results Grid (Span 8) -->
        <div class="lg:col-span-8">
            {#if activeModels.length === 0}
                <div class="h-full min-h-[400px] flex flex-col items-center justify-center border border-dashed border-elevated rounded-2xl bg-surface/10">
                    <div class="w-16 h-16 rounded-full bg-elevated/50 flex items-center justify-center mb-4">
                        <FlaskConical class="w-8 h-8 text-secondary opacity-50" />
                    </div>
                    <span class="font-mono text-sm text-secondary uppercase tracking-widest mb-2">Awaiting Selection</span>
                    <span class="text-xs text-secondary/60 font-sans max-w-xs text-center">Select models from the matrix on the left to begin comparison.</span>
                </div>
            {:else}
                <!-- Grid columns fixed to static Tailwind classes to avoid parsing issues -->
                <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6 w-full items-start">
                    {#each activeModels as item (item.model.spec)}
                        <div class="w-full min-w-0">
                            <ResultCard 
                                result={results[item.model.spec] || null}
                                isLoading={loadingStates[item.model.spec] || false}
                                providerColor={item.providerColor}
                                providerName={item.providerName}
                                modelName={item.model.name}
                            />
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
        
    </div>
</div>

<style>
    .custom-scrollbar::-webkit-scrollbar {
        width: 6px;
    }
    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }
    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: var(--bg-elevated);
        border-radius: 3px;
    }
    .custom-scrollbar:hover::-webkit-scrollbar-thumb {
        background: var(--bg-hover);
    }
</style>
