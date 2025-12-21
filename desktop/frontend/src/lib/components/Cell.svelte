<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import type { Action } from "svelte/action";
  import { evaluate } from "mathjs";
  import { resolveAll, updateCell } from "$lib/formula";

  interface Props {
    cell: CellType;
    grid: Record<string, CellType>;
    onMouseDown: (event: MouseEvent) => void;
    onMouseUp: (event: MouseEvent) => void;
    onEnterPress: (event: KeyboardEvent) => void;
  }

  let {
    cell,
    grid = $bindable(),
    onMouseDown,
    onMouseUp,
    onEnterPress,
  }: Props = $props();
  let startedWithTyping = $state(false)

  const focusInput: Action = (node) => {
    node.focus();
    if (node instanceof HTMLInputElement && !startedWithTyping) {
      node.select();
    }
    startedWithTyping = false
  };

  const processCell = () => {
    const cellKey = `${cell.x}-${cell.y}`;
    if (cell.rawValue[0] === "=") {
      try {
        const formula = cell.rawValue.slice(1);
        const resolvedFormula = resolveAll(formula, grid, cellKey);
        const result = evaluate(resolvedFormula);
        cell.displayValue = String(result);
      } catch (error) {
        cell.displayValue = "#ERROR";
      }
    } else {
      cell.displayValue = cell.rawValue;
    }

    grid[cellKey] = cell;

    updateCell(cellKey, grid);
  };

  function getCellHeader(x: number, y: number): string {
    if (x === 0 && y > 0) {
      return y.toString();
    } else if (y === 0 && x > 0) {
      let column = "";
      let columnIndex = x;
      while (columnIndex > 0) {
        columnIndex--;
        column = String.fromCharCode(65 + (columnIndex % 26)) + column;
        columnIndex = Math.floor(columnIndex / 26);
      }
      return column;
    } else {
      return "";
    }
  }
</script>

{#if cell.x === 0 || cell.y === 0}
  <div class="w-full h-full text-center border border-gray-300 box-border">
    {getCellHeader(cell.x, cell.y)}
  </div>
{:else if cell.isEditing}
  <input
    type="text"
    class="w-full h-full border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={cell.rawValue}
    use:focusInput
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        processCell();
        onEnterPress(event);
      }
      if (event.key === "Escape") {
        cell.isEditing = false;
        cell.isSelected = false;
      }
      if (event.key === "Delete" || event.key === "Backspace") {
        event.stopPropagation();
      }
    }}
    onblur={() => {
      processCell();
      cell.isEditing = false;
      cell.isSelected = false;
    }}
  />
{:else}
  <button
    class={[
      "w-full h-full border border-gray-300 box-border cursor-pointer flex-shrink-0",
      cell.isSelected ? "bg-gray-200" : "bg-white",
    ]}
    onmousedown={
      onMouseDown
    }
    onmouseup={
      onMouseUp
    }
    ondblclick={() => {
      cell.isEditing = true;
    }}
 onkeydown={(event: KeyboardEvent) => {
      if (event.key.length === 1 && !event.ctrlKey && !event.metaKey) {
        cell.isEditing = true
        cell.rawValue = event.key
        startedWithTyping = true
        event.preventDefault()
      }
    }}
  >
    {cell.displayValue}
  </button>
{/if}