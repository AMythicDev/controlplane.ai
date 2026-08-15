<script lang="ts">
  let { data, color, width = 120, height = 32 } = $props<{
    data: number[];
    color: string;
    width?: number;
    height?: number;
  }>();

  let min = $derived(Math.min(...data));
  let max = $derived(Math.max(...data));
  let range = $derived(max - min || 1);
  
  let points = $derived(data.map((val, i) => {
    const x = (i / (data.length - 1)) * width;
    const y = height - ((val - min) / range) * (height - 4) - 2; // pad top/bottom
    return `${x},${y}`;
  }).join(' '));

  let pathData = $derived(`M 0,${height} L ${points} L ${width},${height} Z`);
  let lineData = $derived(`M ${points}`);
  
  let safeColorName = $derived(color.replace(/[^a-zA-Z0-9]/g, ''));
</script>

<svg {width} {height} class="overflow-visible" viewBox="0 0 {width} {height}">
  <defs>
    <linearGradient id="grad-{safeColorName}" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0%" stop-color={color} stop-opacity="0.2" />
      <stop offset="100%" stop-color={color} stop-opacity="0" />
    </linearGradient>
  </defs>
  
  <path d={pathData} fill="url(#grad-{safeColorName})" stroke="none" />
  <path d={lineData} fill="none" stroke={color} stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
</svg>
