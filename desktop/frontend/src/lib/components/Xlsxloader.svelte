<script lang="ts">
  import * as XLSX from "xlsx";
  import type { Cell } from "$lib/types";

  interface Props {
    grid: Record<string, Cell>;
  }

  let { grid = $bindable() }: Props = $props();

  let fileInput: HTMLInputElement;
  let error = $state<string | null>(null);

  async function processFile(file: File): Promise<void> {
    try {
      error = null;
      const arrayBuffer: ArrayBuffer = await file.arrayBuffer();
      const workbook: XLSX.WorkBook = XLSX.read(arrayBuffer);

      // We'll use the first sheet by default
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];

      // Get the range of the worksheet
      const range = XLSX.utils.decode_range(worksheet["!ref"] || "A1");

      // Clear existing grid
      grid = {};

      // Iterate through all cells in the worksheet
      for (let row = range.s.r; row <= range.e.r; row++) {
        for (let col = range.s.c; col <= range.e.c; col++) {
          const cellAddress = XLSX.utils.encode_cell({ r: row, c: col });
          const cell = worksheet[cellAddress];

          // Create a key for our grid (format: "col-row")
          const gridKey = `${col}-${row}`;

          // Get the display value of the cell
          let displayValue = "";
          if (cell) {
            displayValue = cell.w || XLSX.utils.format_cell(cell) || "";
          }

          // Update grid with new cell data
          grid[gridKey] = {
            x: col,
            y: row,
            rawValue: cell?.v?.toString() || "",
            displayValue: displayValue,
            isSelected: false,
            isEditing: false,
          };
        }
      }

      console.log("Grid updated with Excel data:", $state.snapshot(grid));
    } catch (err) {
      error = `エラー: ${err instanceof Error ? err.message : "Unknown error"}`;
    }
  }

  async function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      await processFile(input.files[0]);
    }
  }

  function openFileDialog() {
    fileInput.click();
  }
</script>

<!-- Hidden file input -->
<input
  bind:this={fileInput}
  type="file"
  class="hidden"
  accept=".xlsx,.xls"
  onchange={handleFileChange}
/>

<!-- Styled button -->
<button onclick={openFileDialog} class="btn btn-sm btn-ghost gap-2">
  <svg
    xmlns="http://www.w3.org/2000/svg"
    class="h-4 w-4"
    fill="none"
    viewBox="0 0 24 24"
    stroke="currentColor"
  >
    <path
      stroke-linecap="round"
      stroke-linejoin="round"
      stroke-width="2"
      d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
    />
  </svg>
  Open File
</button>

{#if error}
  <div class="toast toast-top toast-end">
    <div class="alert alert-error">
      <span>{error}</span>
    </div>
  </div>
{/if}
