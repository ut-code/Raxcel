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

  // ✅ セルをローカルステートとして管理
  let localCell = $state({
    x: cell.x,
    y: cell.y,
    rawValue: cell.rawValue,
    displayValue: cell.displayValue,
    isSelected: cell.isSelected,
    isEditing: cell.isEditing,
  });

  // ✅ 親から変更があったら同期
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
    console.log(`processCell called for ${localCell.x}-${localCell.y}`);
    const cellKey = `${localCell.x}-${localCell.y}`;
    if (localCell.rawValue[0] === "=") {
      try {
        const formula = localCell.rawValue.slice(1);
        console.log(`Resolving formula: ${formula}`);
        const resolvedFormula = resolveAll(formula, grid, cellKey);
        console.log(`Resolved: ${resolvedFormula}`);
        const result = evaluate(resolvedFormula);
        console.log(`Result: ${result}`);
        localCell.displayValue = String(result);
      } catch (error) {
        console.error(`Error in cell ${cellKey}:`, error);
        localCell.displayValue = "#ERROR";
      }
    } else {
      localCell.displayValue = localCell.rawValue;
    }

    // ✅ gridに反映（参照は保持）
    grid[cellKey] = {
      x: localCell.x,
      y: localCell.y,
      rawValue: localCell.rawValue,
      displayValue: localCell.displayValue,
      isSelected: localCell.isSelected,
      isEditing: localCell.isEditing,
    };

    console.log(`Calling updateCell for ${cellKey}`);
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

  // デバッグ用
  $effect(() => {
    console.log(`Cell ${localCell.x}-${localCell.y} state:`, {
      isEditing: localCell.isEditing,
      isSelected: localCell.isSelected,
      rawValue: localCell.rawValue,
      displayValue: localCell.displayValue,
    });
  });
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
      console.log(`Key pressed: ${event.key}`);
      if (event.key === "Enter") {
        console.log(`Enter pressed in cell ${localCell.x}-${localCell.y}`);
        localCell.isEditing = false;
        processCell();
        onEnterPress(event);
      }
      if (event.key === "Escape") {
        console.log(`Escape pressed in cell ${localCell.x}-${localCell.y}`);
        localCell.isEditing = false;
        localCell.isSelected = false;
      }
      if (event.key === "Delete" || event.key === "Backspace") {
        event.stopPropagation();
      }
    }}
    onblur={() => {
      console.log(`Blur event in cell ${localCell.x}-${localCell.y}`);
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
      console.log(`MouseDown on cell ${localCell.x}-${localCell.y}`);
      onMouseDown(event);
    }}
    onmouseup={(event) => {
      console.log(`MouseUp on cell ${localCell.x}-${localCell.y}`);
      onMouseUp(event);
    }}
    ondblclick={() => {
      console.log(`DoubleClick on cell ${localCell.x}-${localCell.y}`);
      localCell.isEditing = true;
    }}
  >
    {localCell.displayValue}
  </button>
{/if}
