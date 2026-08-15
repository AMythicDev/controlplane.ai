<script lang="ts">
    let { score, label, dimension } = $props<{
        score: number;
        label: string;
        dimension: 'performance' | 'cost' | 'responsibility';
    }>();

    let mounted = $state(false);

    $effect(() => {
        mounted = true;
    });

    let color = $derived(
        dimension === 'performance' ? 'var(--accent-performance)' :
        dimension === 'cost' ? 'var(--accent-cost)' :
        'var(--accent-responsibility)'
    );

    // Score from 0 to 100
    // Semi-circle path: M 10 90 A 40 40 0 0 1 90 90
    // Length is roughly PI * r = 3.14 * 40 = 125.6
    let strokeDasharray = 125.6;
    let strokeDashoffset = $derived(
        mounted ? strokeDasharray - (score / 100) * strokeDasharray : strokeDasharray
    );
</script>

<div class="flex flex-col items-center">
    <div class="relative w-32 h-20 overflow-hidden">
        <svg viewBox="0 0 100 60" class="w-full h-full drop-shadow-lg">
            <!-- Background track -->
            <path
                d="M 10 50 A 40 40 0 0 1 90 50"
                fill="none"
                stroke="var(--bg-elevated)"
                stroke-width="8"
                stroke-linecap="round"
            />
            <!-- Foreground arc -->
            <path
                d="M 10 50 A 40 40 0 0 1 90 50"
                fill="none"
                stroke={color}
                stroke-width="8"
                stroke-linecap="round"
                stroke-dasharray={strokeDasharray}
                stroke-dashoffset={strokeDashoffset}
                class="transition-all duration-1000 ease-out"
            />
        </svg>
        <div class="absolute bottom-0 left-0 right-0 text-center font-mono text-2xl font-medium" style="color: {color}">
            {Math.round(score)}
        </div>
    </div>
    <div class="mt-2 text-xs uppercase tracking-widest text-secondary font-sans font-medium">
        {label}
    </div>
</div>
