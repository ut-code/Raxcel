<script lang="ts">
  import type { Cell as CellType } from "$lib/types";
  import Cell from "./Cell.svelte";
  import { type Readable } from "svelte/store";
  import {
    createVirtualizer,
    type SvelteVirtualizer,
  } from "@tanstack/svelte-virtual";
  import type { SvelteSet } from "svelte/reactivity";

  interface Props {
    grid: Record<string, CellType>;
    selectedCells: SvelteSet<string>;
  }
  let { grid = $bindable(), selectedCells = $bindable() }: Props = $props();

  let virtualListEl: HTMLDivElement;
  let rowVirtualizer:
    | Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>>
    | undefined = $state();
  let columnVirtualizer:
    | Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>>
    | undefined = $state();

  let leftTopCell: CellType | null = $state(null);
  let isDragging = $state(false);

  const rowCount = 1000;
  const colCount = 1000;
  const cellHeight = 40;
  const cellWidth = 100;

  $effect(() => {
    if (virtualListEl) {
      rowVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
        count: rowCount,
        getScrollElement: () => virtualListEl,
        estimateSize: () => cellHeight,
        overscan: 5,
      });
      columnVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
        count: colCount,
        getScrollElement: () => virtualListEl,
        estimateSize: () => cellWidth,
        overscan: 5,
        horizontal: true,
      });
    }
  });

  // ✅ getCell内でgridを変更しない
  function getCell(x: number, y: number): CellType {
    const key = `${x}-${y}`;
    let cell = grid[key];
    if (!cell) {
      cell = {
        x,
        y,
        rawValue: "",
        displayValue: "",
        isSelected: false,
        isEditing: false,
      };
      // ✅ ここでは追加しない
    }
    return cell;
  }

  // ✅ セルを確実にgridに追加する関数
  function ensureCell(x: number, y: number): CellType {
    const key = `${x}-${y}`;
    if (!grid[key]) {
      grid[key] = {
        x,
        y,
        rawValue: "",
        displayValue: "",
        isSelected: false,
        isEditing: false,
      };
    }
    return grid[key];
  }

  function getCellKey(x: number, y: number): string {
    return `${x}-${y}`;
  }

  function convertEventLocToCell(event: Event): CellType | null {
    const target = event.target;
    if (!(target instanceof Element)) {
      return null;
    }
    const cellEl = target.closest("[data-cell-loc]");
    if (!cellEl) return null;
    const coords = cellEl.getAttribute("data-cell-loc");
    if (!coords) return null;
    const [x, y] = coords.split("-").map(Number);
    return ensureCell(x, y); // ✅ ensureCellを使用
  }

  function updateSelection(startCell: CellType, endCell: CellType) {
    for (const key of selectedCells) {
      const cell = grid[key];
      if (cell) {
        cell.isSelected = false;
      }
    }
    selectedCells.clear();

    const startX = Math.min(startCell.x, endCell.x);
    const endX = Math.max(startCell.x, endCell.x);
    const startY = Math.min(startCell.y, endCell.y);
    const endY = Math.max(startCell.y, endCell.y);

    for (let y = startY; y <= endY; y++) {
      for (let x = startX; x <= endX; x++) {
        const key = getCellKey(x, y);
        selectedCells.add(key);
        const cell = ensureCell(x, y); // ✅ ensureCellを使用
        cell.isSelected = true;
      }
    }
  }

  function handleMouseDown(event: MouseEvent) {
    for (const key of selectedCells) {
      const cell = grid[key];
      if (cell) {
        cell.isSelected = false;
      }
    }
    selectedCells.clear();

    const cell = convertEventLocToCell(event);

    if (cell) {
      leftTopCell = cell;
      isDragging = true;
      cell.isSelected = true;
      selectedCells.add(getCellKey(cell.x, cell.y));
    }
  }

  function handleMouseMove(event: MouseEvent) {
    if (isDragging && leftTopCell) {
      const cell = convertEventLocToCell(event);
      if (cell) {
        updateSelection(leftTopCell, cell);
      }
    }
  }

  function handleMouseUp(event: MouseEvent) {
    if (isDragging && leftTopCell) {
      const cell = convertEventLocToCell(event);
      if (cell) {
        updateSelection(leftTopCell, cell);
      }
    }
    isDragging = false;
    leftTopCell = null;
  }

  function handleEnterPress(event: KeyboardEvent) {
    const currentCell = convertEventLocToCell(event);
    if (currentCell) {
      currentCell.isEditing = false;
      currentCell.isSelected = false;
      selectedCells.delete(getCellKey(currentCell.x, currentCell.y));

      const nextY = currentCell.y + 1;
      if (nextY < rowCount) {
        for (const cell of Object.values(grid)) {
          cell.isSelected = false;
          cell.isEditing = false;
        }
        selectedCells.clear();
        const nextCell = ensureCell(currentCell.x, nextY); // ✅ ensureCellを使用
        nextCell.isSelected = true;
        nextCell.isEditing = true;
        selectedCells.add(getCellKey(currentCell.x, nextY));
      }
    }
  }

  function handleDelete() {
    for (const key of selectedCells) {
      const cell = grid[key];
      if (cell) {
        cell.displayValue = "";
        cell.rawValue = "";
        cell.isSelected = false;
      }
    }
  }
</script>

<div bind:this={virtualListEl} class="flex-1 w-full overflow-auto">
  {#if $rowVirtualizer && $columnVirtualizer}
    <div
      class="relative"
      style={`
        height: ${$rowVirtualizer.getTotalSize()}px;
        width: ${$columnVirtualizer.getTotalSize()}px;
      `}
      onmousemove={handleMouseMove}
      onkeydown={(event: KeyboardEvent) => {
        if (event.key === "Delete" || event.key === "Backspace") {
          handleDelete();
        }
      }}
      role="grid"
      tabindex="0"
    >
      {#each $rowVirtualizer.getVirtualItems() as row (row.index)}
        {#each $columnVirtualizer.getVirtualItems() as col (col.index)}
          <div
            class="absolute top-0 left-0"
            style={`
                width: ${col.size}px;
                height: ${row.size}px;
                transform: translateX(${col.start}px) translateY(${row.start}px);
             `}
            data-cell-loc={`${col.index}-${row.index}`}
          >
            <Cell
              cell={getCell(col.index, row.index)}
              bind:grid
              onMouseDown={handleMouseDown}
              onMouseUp={handleMouseUp}
              onEnterPress={handleEnterPress}
            />
          </div>
        {/each}
      {/each}
    </div>
  {/if}
</div>
