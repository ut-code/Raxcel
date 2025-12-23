<script lang="ts">
  import * as XLSX from "xlsx";
  import type { Cell as CellType } from "$lib/types";

  interface Props {
    grid: Record<string, CellType>;
  }

  let { grid = $bindable() }: Props = $props();
  let error = $state<string | null>(null);
  let fileName = $state("spreadsheet");
  let dialogElement: HTMLDialogElement;
  
  function focusOnMount(node: HTMLElement) {
    setTimeout(() => node.focus(), 0);
  }
  
  function exportToXlsx(grid: Record<string, CellType>, fileName: string): void {
    try {
      const workbook = XLSX.utils.book_new();
      
      let maxRow = 0;
      let maxCol = 0;
      for (const key in grid) {
        const cell = grid[key];
        if (cell.rawValue || cell.displayValue) {
          maxRow = Math.max(maxRow, cell.y);
          maxCol = Math.max(maxCol, cell.x);
        }
      }

      const worksheetData: (string | number)[][] = [];
      for (let row = 0; row <= maxRow; row++) {
        worksheetData[row] = [];
        for (let col = 0; col <= maxCol; col++) {
          const key = `${col}-${row}`;
          const cell = grid[key];
          worksheetData[row][col] = cell?.rawValue || "";
        }
      }

      const worksheet = XLSX.utils.aoa_to_sheet(worksheetData);
      XLSX.utils.book_append_sheet(workbook, worksheet, "Sheet1");
      
      const fullFileName = fileName.endsWith('.xlsx') ? fileName : `${fileName}.xlsx`;
      XLSX.writeFile(workbook, fullFileName);
      
      dialogElement?.close();
      fileName = "spreadsheet";
      error = null;
    } catch (err) {
      error = `Export error: ${err instanceof Error ? err.message : "Unknown error"}`;
    }
  }

  function handleSubmit(e: Event) {
    e.preventDefault();
    if (fileName.trim()) {
      exportToXlsx(grid, fileName.trim());
    }
  }
  
  function closeModal() {
    dialogElement?.close();
    fileName = "spreadsheet";
  }
</script>

<button 
  onclick={() => dialogElement?.showModal()} 
  class="btn btn-sm btn-ghost gap-2"
>
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
      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" 
    />
  </svg>
  Export to Excel
</button>

<!-- dialog  -->
<dialog bind:this={dialogElement} class="modal">
  <div class="modal-box">
    <h3 class="font-bold text-lg">Enter filename</h3>
    
    <form onsubmit={handleSubmit} class="py-4">
      <div class="form-control">
        <label class="label" for="filename-input">
          <span class="label-text">filename</span>
        </label>
        <div class="input-group">
          <input
            id="filename-input"
            type="text"
            bind:value={fileName}
            placeholder="spreadsheet"
            class="input input-bordered w-full"
            use:focusOnMount
          />
          <span class="bg-base-200 px-3 flex items-center">.xlsx</span>
        </div>
      </div>

      <div class="modal-action">
        <button 
          type="button" 
          class="btn btn-ghost"
          onclick={closeModal}
        >
          cancel
        </button>
        <button 
          type="submit" 
          class="btn btn-primary"
          disabled={!fileName.trim()}
        >
          export
        </button>
      </div>
    </form>
  </div>
  
  <form method="dialog" class="modal-backdrop">
    <button type="submit" aria-label="close modal">close</button>
  </form>
</dialog>

{#if error}
  <div class="toast toast-top toast-end">
    <div class="alert alert-error">
      <span>{error}</span>
    </div>
  </div>
{/if}