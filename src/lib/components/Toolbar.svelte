<script lang="ts">
  import Chart from "$lib/components/Chart.svelte";
  import Xlsxloader from "$lib/components/Xlsxloader.svelte";
  import type { Cell } from "$lib/types";

  interface Props {
    chartComponent: Chart | null;
    grid: Record<string, Cell>;
  }
  let { chartComponent, grid = $bindable() }: Props = $props();

  function handlerCreateChart() {
    if (chartComponent) {
      const success = chartComponent.drawChart();
      if (!success) {
        alert("select valid numbers");
      }
    }
  }
</script>

<Xlsxloader bind:grid />

<div class="flex flex-row">
  <div>
    <button onclick={handlerCreateChart} class="btn">Create Chart</button>
  </div>

  <div>
    <button class="btn">Do nothing</button>
  </div>

  <div>
    <button class="btn">Do nothing</button>
  </div>
</div>
