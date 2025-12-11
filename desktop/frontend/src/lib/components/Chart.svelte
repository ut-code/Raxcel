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
  import type { SvelteSet } from "svelte/reactivity";

  Chart.register(ScatterController, LinearScale, PointElement, Tooltip, Legend);

  interface Props {
    grid: Record<string, CellType>;
    selectedCells: SvelteSet<string>;
  }

  let { grid, selectedCells }: Props = $props();

  let selectedValues: string[] = $derived(
    Array.from(selectedCells)
      .map((key) => {
        const cell = grid[key];
        return cell?.displayValue || "";
      })
      .filter((val) => val !== ""),
  );
  let chartInstance: Chart | null = null;
  let canvasRef: HTMLCanvasElement;

  let dialogRef: HTMLDialogElement;

  function validateValues(selectedValues: string[]): {
    validatedValues: number[];
    isValid: boolean;
  } {
    const validatedValues = selectedValues.map((value) => Number(value));
    for (const value of validatedValues) {
      if (isNaN(value)) {
        return {
          validatedValues: [],
          isValid: false,
        };
      }
    }
    return {
      validatedValues,
      isValid: true,
    };
  }

  export function drawChart() {
    const { validatedValues, isValid } = validateValues(selectedValues);
    if (isValid && validatedValues.length > 0) {
      if (chartInstance) {
        chartInstance.destroy();
        chartInstance = null;
      }
      const config = setupPlot(validatedValues);
      dialogRef.showModal();
      chartInstance = new Chart(canvasRef, config);
      return true;
    }
    return false;
  }
</script>

<dialog class="modal" bind:this={dialogRef}>
  <div class="relative h-96 w-128">
    <form method="dialog">
      <button class="btn">Close</button>
    </form>
    <canvas bind:this={canvasRef} class="bg-white"></canvas>
  </div>
</dialog>
