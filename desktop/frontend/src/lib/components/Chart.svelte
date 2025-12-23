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
  function exportImage() {
    const imageData = canvasRef.toDataURL("image/png");
    const link = document.createElement("a");
    link.href = imageData;
    link.download = `chart-${new Date().getTime()}.png`;
    link.click();
  }
</script>

<dialog class="modal" bind:this={dialogRef}>
  <div class="modal-box">
    <button aria-label="export as image" class="btn btn-sm absolute top-2 right-6" onclick={exportImage}>
      <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M5.25589 16C3.8899 15.0291 3 13.4422 3 11.6493C3 9.20008 4.8 6.9375 7.5 6.5C8.34694 4.48637 10.3514 3 12.6893 3C15.684 3 18.1317 5.32251 18.3 8.25C19.8893 8.94488 21 10.6503 21 12.4969C21 14.0582 20.206 15.4339 19 16.2417M12 21V11M12 21L9 18M12 21L15 18" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg>
    </button>
    <canvas bind:this={canvasRef} class="bg-white"></canvas>
  </div>
    <form method="dialog" class="modal-backdrop">
      <button>close</button>
    </form>
</dialog>
