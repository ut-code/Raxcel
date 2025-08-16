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
    selectedValues: string[];
  }

  let { selectedValues = [] }: Props = $props();
  let chartInstance: Chart | null = null;
  let canvasRef: HTMLCanvasElement;

  function validateValues(selectedValues: string[]): {
    validatedValues: number[];
    isValid: boolean
  } {
    const validatedValues = selectedValues.map(value => Number(value));
    for (const value of validatedValues) {
        if (isNaN(value)) {
          return {
            validatedValues: [],
            isValid: false,
          }
        }
    }
    return {
      validatedValues,
      isValid: true,
    }
  }

  export function drawChart() {
    const { validatedValues, isValid } = validateValues(selectedValues)
    if (isValid && validatedValues.length > 0) {
      if (chartInstance) {
        chartInstance.destroy();
        chartInstance = null;
      }
      const config = setupPlot(validatedValues)
      chartInstance = new Chart(canvasRef, config)
      return true
    }
    return false
  }
</script>

<div style="w-[500px]">
  <canvas bind:this={canvasRef}></canvas>
</div>
