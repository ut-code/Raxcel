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
    isDrawing: boolean;
  }

  let { selectedValues = [], isDrawing = false }: Props = $props();
  let chartInstance: Chart | null = null;
  let canvasRef: HTMLCanvasElement;

  $effect(() => {
    const { validatedValues, isValid } = validateValues(selectedValues)
    if (isDrawing && isValid) {
      createChart(validatedValues);
    }
    return () => {
      if (chartInstance) {
        chartInstance.destroy();
      }
    }
  });

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

  function createChart(values: number[]) {
    const config = setupPlot(values)
    
    if (chartInstance) {
      chartInstance.destroy();
    }

    chartInstance = new Chart(canvasRef, config)
    isDrawing = false;
  }

</script>

<div style="w-[500px]">
  <canvas bind:this={canvasRef}></canvas>
</div>
