<script lang="ts">
  import * as XLSX from "xlsx";
  import type { Cell } from "$lib/types";

  interface Props {
    grid: Record<string, Cell>;
  }

  let { grid = $bindable() }: Props = $props();

  let currentFile = $state<File | null>(null);
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

  function handleFileInput(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      currentFile = input.files[0];
    }
  }
</script>

<form
  onsubmit={async (e) => {
    if (currentFile) await processFile(currentFile);
  }}
>
  <input type="file" accept=".xlsx,.xls" onchange={handleFileInput} />
  <button type="submit">upload Excel file</button>
</form>

{#if error}
  <p class="text-red-500">{error}</p>
{/if}
