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

  let localCell = $state({
    x: cell.x,
    y: cell.y,
    rawValue: cell.rawValue,
    displayValue: cell.displayValue,
    isSelected: cell.isSelected,
    isEditing: cell.isEditing,
  });

  $effect(() => {
    localCell = {
      x: cell.x,
      y: cell.y,
      rawValue: cell.rawValue,
      displayValue: cell.displayValue,
      isSelected: cell.isSelected,
      isEditing: cell.isEditing,
    };
  });

  const focusInput: Action = (node) => {
    node.focus();
    if (node instanceof HTMLInputElement) {
      node.select();
    }
  };

  const processCell = () => {
    const cellKey = `${localCell.x}-${localCell.y}`;
    if (localCell.rawValue[0] === "=") {
      try {
        const formula = localCell.rawValue.slice(1);
        const resolvedFormula = resolveAll(formula, grid, cellKey);
        const result = evaluate(resolvedFormula);
        localCell.displayValue = String(result);
      } catch (error) {
        localCell.displayValue = "#ERROR";
      }
    } else {
      localCell.displayValue = localCell.rawValue;
    }

    grid[cellKey] = {
      x: localCell.x,
      y: localCell.y,
      rawValue: localCell.rawValue,
      displayValue: localCell.displayValue,
      isSelected: localCell.isSelected,
      isEditing: localCell.isEditing,
    };

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

{#if localCell.x === 0 || localCell.y === 0}
  <div class="w-full h-full text-center border border-gray-300 box-border">
    {getCellHeader(localCell.x, localCell.y)}
  </div>
{:else if localCell.isEditing}
  <input
    type="text"
    class="w-full h-full border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={localCell.rawValue}
    use:focusInput
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        localCell.isEditing = false;
        processCell();
        onEnterPress(event);
      }
      if (event.key === "Escape") {
        localCell.isEditing = false;
        localCell.isSelected = false;
      }
      if (event.key === "Delete" || event.key === "Backspace") {
        event.stopPropagation();
      }
    }}
    onblur={() => {
      localCell.isEditing = false;
      processCell();
      localCell.isSelected = false;
    }}
  />
{:else}
  <button
    class={[
      "w-full h-full border border-gray-300 box-border cursor-pointer flex-shrink-0",
      localCell.isSelected ? "bg-gray-200" : "bg-white",
    ]}
    onmousedown={(event) => {
      onMouseDown(event);
    }}
    onmouseup={(event) => {
      onMouseUp(event);
    }}
    ondblclick={() => {
      localCell.isEditing = true;
    }}
  >
    {localCell.displayValue}
  </button>
{/if}