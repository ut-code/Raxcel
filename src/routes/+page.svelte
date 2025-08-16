<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import Cell from "$lib/components/Cell.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import Toolbar from "$lib/components/Toolbar.svelte";
  import { createVirtualizer, type SvelteVirtualizer } from "@tanstack/svelte-virtual";
  import { SvelteMap, SvelteSet } from "svelte/reactivity";
  import { type Readable } from "svelte/store";

  const rowCount = 1000;
  const colCount = 1000;
  const sheetHeight = 20;
  const sheetWidth = 20;

  let cellData = new SvelteMap<string, CellType>();

  function getCell(x: number, y: number): CellType {
    const key = `${x}-${y}`;
    let cell = cellData.get(key);
    if (!cell){
     cell = {
       x,
       y,
       rawValue: "",
       displayValue: "",
       isSelected: false,
       isEditing: false,
     }
      cellData.set(key, cell)
    }
    return cell
  }

  function setCell(x: number, y: number, newCell: CellType) {
    const key=`${x}-${y}`;
    cellData.set(key,newCell);
  }

  let leftTopCell: CellType | null = $state(null);
  let isDragging = $state(false);

  let selectedCells = new SvelteSet<string>();
  let selectedValues: string[] = $derived(
    Array.from(selectedCells).map(key => {
      const cell = cellData.get(key);
      return cell?.displayValue || "";
    }).filter(val => val !== "")
  )

  let chartComponent: Chart | null = $state(null);
  let virtualListEl: HTMLDivElement;

  let rowVirtualizer: Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>> | undefined = $state()
  let columnVirtualizer: Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>> | undefined =  $state()

  $effect(() => {
    if (virtualListEl) {
      rowVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: rowCount,
    getScrollElement: () => virtualListEl,
    estimateSize: () => sheetHeight,
    overscan: 5,
      })
      columnVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: rowCount,
    getScrollElement: () => virtualListEl,
    estimateSize: () => sheetWidth,
    overscan: 5,
    horizontal: true
      })
    }
  })

  function getCellKey(x: number, y: number): string {
    return `${x}-${y}`
  }

  function convertMouseLocToCell(event: MouseEvent) : CellType | null{
    const target = event.target as HTMLElement;
    const cellEl = target.closest("[data-cell-coords]")
    if (!cellEl) return null;
    const coords = cellEl.getAttribute("data-cell-coords")
    if (!coords) return null;
    const [x, y] = coords.split("-").map(Number)
    return getCell(x, y)
  }

  function updateSelection(startCell: CellType, endCell: CellType) {
    selectedCells.clear()    
    for (const [_key, cell] of cellData) {
      cell.isSelected = false
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
        cell.isSelected = true
      }
    }
  }

  function handleMouseDown(event: MouseEvent) {
    selectedCells.clear()
    for (const [_key, cell] of cellData) {
      cell.isSelected = false;
    }

    const cell = convertMouseLocToCell(event);
    if (cell) {
      leftTopCell = cell;
      isDragging = true;
      cell.isSelected = true;
      selectedCells.add(getCellKey(cell.x, cell.y))
    }
  }

  function handleMouseMove(event: MouseEvent) {
    if (isDragging && leftTopCell) {
      const cell = convertMouseLocToCell(event);
      if (cell) {
        updateSelection(leftTopCell, cell);
      }
    }
  }

  function handleMouseUp(event: MouseEvent) {
    if (isDragging && leftTopCell) {
      const cell = convertMouseLocToCell(event);
      if (cell) {
        updateSelection(leftTopCell, cell);
      }
    }
    isDragging = false;
    leftTopCell = null;
  }

  function handleEnterPress(x: number, y: number) {
    const currentCell = getCell(x, y);
    currentCell.isEditing = false;
    currentCell.isSelected = false;
    selectedCells.delete(getCellKey(x, y));

    const nextY = y + 1;
    if (nextY < rowCount) {
      for (const [_key, cell] of cellData) {
        cell.isSelected = false;
        cell.isEditing = false;
      }
      selectedCells.clear()
      const nextCell = getCell(x, nextY)
      nextCell.isSelected = true;
      nextCell.isEditing = true;
      selectedCells.add(getCellKey(x, nextY))
    }
  }

</script>

<Toolbar {chartComponent} />

<div bind:this={virtualListEl} class="h-[600px] w-full overflow-auto">
  <div
    class="relative"
    style={`
      height: ${$rowVirtualizer?.getTotalSize()}px;
      width: ${$columnVirtualizer?.getTotalSize()}px;
    `}
    onmousemove={handleMouseMove}
    role="grid"
    tabindex=0
  >

    {#each $rowVirtualizer?.getVirtualItems() ?? []as row (row.index)}
      {#each $columnVirtualizer?.getVirtualItems() ?? [] as col (col.index)}
        <div
          class="absolute top-0 left-0"
          style={`
              width: ${col.size}px;
              height: ${row.size}px;
              transform: translateX(${col.start}px) translateY(${row.start}px);
           `}
          data-cell-coords={`${col.index}-${row.index}`}>
         <Cell 
          bind:cell={
            () => getCell(col.index, row.index),
            (newCell) => setCell(col.index, row.index, newCell)
          } 
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
          onEnterPress={handleEnterPress}
        />
        </div>
      {/each}
    {/each}
  </div>
</div>

<Chart {selectedValues} bind:this={chartComponent}/>
