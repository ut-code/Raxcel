<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import Cell from "$lib/components/Cell.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import Toolbar from "$lib/components/Toolbar.svelte";
  import {
    createVirtualizer,
    type SvelteVirtualizer,
  } from "@tanstack/svelte-virtual";
  import { SvelteSet } from "svelte/reactivity";
  import { type Readable } from "svelte/store";
  import Chat from "$lib/components/Chat.svelte";
  import { CheckUser } from "$lib/wailsjs/go/main/App";
  import { onMount } from "svelte";

  const rowCount = 1000;
  const colCount = 1000;
  const cellHeight = 40;
  const cellWidth = 100;

  // global state
  let grid = $state<Record<string, CellType>>({});
  let isLoggedIn = $state(false);
  onMount(() => {
    const checkUser = async () => {
      const result = await CheckUser();
      if (result.ok) {
        isLoggedIn = true;
      }
    };
    checkUser();
  });

  // Function to update the grid
  function updateGrid(x: number, y: number, cellData: Partial<CellType>) {
    const key = `${x}-${y}`;
    grid[key] = {
      ...grid[key],
      ...cellData,
    };
  }

  // to select cells
  let leftTopCell: CellType | null = $state(null);
  let isDragging = $state(false);
  let selectedCells = new SvelteSet<string>();
  let selectedValues: string[] = $derived(
    Array.from(selectedCells)
      .map((key) => {
        const cell = grid[key];
        return cell?.displayValue || "";
      })
      .filter((val) => val !== ""),
  );

  // to show the graph
  let chartComponent: Chart | null = $state(null);

  let isChatOpen = $state(false);

  // virtual scroll
  let virtualListEl: HTMLDivElement;
  let rowVirtualizer:
    | Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>>
    | undefined = $state();
  let columnVirtualizer:
    | Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>>
    | undefined = $state();

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

  // cell utility
  function getCell(x: number, y: number): CellType {
    const key = `${x}-${y}`;
    let cell = grid[key];
    // If cell is undefined, initialize it.
    if (!cell) {
      cell = {
        x,
        y,
        rawValue: "",
        displayValue: "",
        isSelected: false,
        isEditing: false,
      };
    }
    return cell;
  }

  function setCell(newCell: CellType) {
    const key = `${newCell.x}-${newCell.y}`;
    grid[key] = newCell;
  }

  function getCellKey(x: number, y: number): string {
    return `${x}-${y}`;
  }

  function convertEventLocToCell(event: Event): CellType | null {
    const target = event.target;
    if (!(target instanceof Element)) {
      return null;
    }
    // get the element by data-cell-loc attribute
    const cellEl = target.closest("[data-cell-loc]");
    if (!cellEl) return null;
    const coords = cellEl.getAttribute("data-cell-loc");
    if (!coords) return null;
    const [x, y] = coords.split("-").map(Number);
    return getCell(x, y);
  }

  function updateSelection(startCell: CellType, endCell: CellType) {
    selectedCells.clear();
    for (const cell of Object.values(grid)) {
      cell.isSelected = false;
    }

    const startX = Math.min(startCell.x, endCell.x);
    const endX = Math.max(startCell.x, endCell.x);
    const startY = Math.min(startCell.y, endCell.y);
    const endY = Math.max(startCell.y, endCell.y);

    for (let y = startY; y <= endY; y++) {
      for (let x = startX; x <= endX; x++) {
        const key = getCellKey(x, y);
        selectedCells.add(key);
        const cell = getCell(x, y);
        cell.isSelected = true;
      }
    }
  }

  function handleMouseDown(event: MouseEvent) {
    selectedCells.clear();
    for (const cell of Object.values(grid)) {
      cell.isSelected = false;
    }

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
        const nextCell = getCell(currentCell.x, nextY);
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
        grid[key] = { ...cell };
      }
    }
  }
</script>

<div class="h-screen flex flex-col">
  <Toolbar {chartComponent} bind:grid bind:isChatOpen bind:isLoggedIn />

  <div bind:this={virtualListEl} class="flex-1 w-full overflow-auto">
    <div
      class="relative"
      style={`
      height: ${$rowVirtualizer?.getTotalSize()}px;
      width: ${$columnVirtualizer?.getTotalSize()}px;
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
      {#each $rowVirtualizer?.getVirtualItems() ?? [] as row (row.index)}
        {#each $columnVirtualizer?.getVirtualItems() ?? [] as col (col.index)}
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
              bind:cell={
                () => getCell(col.index, row.index), (cell) => setCell(cell)
              }
              bind:grid
              onMouseDown={handleMouseDown}
              onMouseUp={handleMouseUp}
              onEnterPress={handleEnterPress}
            />
          </div>
        {/each}
      {/each}
    </div>
  </div>

  {#if isChatOpen}
    <Chat />
  {/if}
  <Chart {selectedValues} bind:this={chartComponent} />
</div>
