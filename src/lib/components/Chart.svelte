<script lang="ts">
  import { setupPlot } from "$lib/chart";
  import type { Cell as CellType } from "$lib/types.ts";
  import {
    Chart,
    ScatterController,
    LinearScale,
    PointElement,
    Tooltip,
    Legend,
  } from "chart.js/auto";

  Chart.register(ScatterController, LinearScale, PointElement, Tooltip, Legend);

  interface Props {
    grid: CellType[][];
  }

  let { grid }: Props = $props();
  let chartInstance: Chart | null = null;
  let canvasRef: HTMLCanvasElement;

  function createPlot() {
    const selectedCells = grid.flat().filter((cell) => cell.isSelected);
    
    if (selectedCells.length % 2 !== 0) {
      alert("Please select an even number of cells.");
      return;
    }

    console.log("Selected cells:", selectedCells);
    const values = selectedCells.map((cell) => parseFloat(cell.value));
    
    if (values.some((value) => isNaN(value))) {
      alert("Invalid value in selected cells. Please enter valid numbers.");
      return;
    }

    console.log("Selected values:", values);
    const config = setupPlot(values);
    
    if (chartInstance) {
      chartInstance.destroy();
    }
    
    chartInstance = new Chart(canvasRef, config);
  }
</script>

<button onclick={createPlot}>Create Plot</button>
<div style="w-[500px]">
  <canvas bind:this={canvasRef}></canvas>
</div>
