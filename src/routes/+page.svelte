<script lang="ts">
  import type { CellState, CellPosition} from "$lib/types.ts";
  import Cell from "$lib/components/Cell.svelte";
  import Chart from "$lib/components/Chart.svelte";
  import Toolbar from "$lib/components/Toolbar.svelte";
  import { createVirtualizer, type SvelteVirtualizer } from "@tanstack/svelte-virtual";
  import { SvelteMap, SvelteSet } from "svelte/reactivity";
  import { type Readable } from "svelte/store";

  const rowCount = 1000;
  const colCount = 1000;
  const cellHeight = 32;
  const cellWidth = 100;

  // Glocal State. This is everything.
  let sheetData = new SvelteMap<CellPosition, CellState>();

  // To select cells
  let leftTopPos: CellPosition | null = $state(null);
  let isDragging = $state(false);
  let selectedPoss = new SvelteSet<CellPosition>();

  // chart
  let chartComponent: Chart | null = $state(null);

  // For virtual scroll
  let virtualListEl: HTMLDivElement;
  let rowVirtualizer: Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>> | undefined = $state()
  let columnVirtualizer: Readable<SvelteVirtualizer<HTMLDivElement, HTMLDivElement>> | undefined =  $state()

  $effect(() => {
    if (virtualListEl) {
      rowVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: rowCount,
    getScrollElement: () => virtualListEl,
    estimateSize: () => cellHeight,
    overscan: 5,
      })
      columnVirtualizer = createVirtualizer<HTMLDivElement, HTMLDivElement>({
    count: colCount,
    getScrollElement: () => virtualListEl,
    estimateSize: () => cellWidth,
    overscan: 5,
    horizontal: true
      })
    }
  })

  function convertMousePosToCellPos(event: MouseEvent) : CellPosition | null{
    //TODO: write here
    // This is a mock data.
    const pos: CellPosition = {
      x: 1,
      y: 1,
    }
    return pos
  }

  function updateSelection(startPos: CellPosition, endPos: CellPosition) {
    selectedPoss.clear()    
    for (const state of sheetData.values()) {
    //TODO: is this the right way to update SvelteMap?
      state.isSelected = false
    }

    const startX = Math.min(startPos.x, endPos.x);
    const endX = Math.max(startPos.x, endPos.x);
    const startY = Math.min(startPos.y, endPos.y);
    const endY = Math.max(startPos.y, endPos.y);
    
    for (let y = startY; y <= endY; y++) {
      for (let x = startX; x <= endX; x++) {
        const key: CellPosition = {x, y};
        selectedPoss.add(key);
        const cell = sheetData.get(key);
        if (cell) {
          const newCellState: CellState = {...cell, isSelected: true};
          sheetData.set(key, newCellState);
        }
      }
    }
  }

  function handleMouseDown(event: MouseEvent) {
    selectedPoss.clear()
    for (const state of sheetData.values()) {
      state.isSelected = false;
    }

    const pos = convertMousePosToCellPos(event);
    if (pos) {
      leftTopPos = pos;
      isDragging = true;
      const state = sheetData.get(pos);
      if (state) {
        sheetData.set(pos, {
          ...state,
          isSelected: true
        })
      }
      selectedPoss.add(pos)
    }
  }

  function handleMouseMove(event: MouseEvent) {
    if (isDragging && leftTopPos) {
      const pos = convertMousePosToCellPos(event);
      if (pos) {
        updateSelection(leftTopPos, pos);
      }
    }
  }

  function handleMouseUp(event: MouseEvent) {
    if (isDragging && leftTopPos) {
      const pos = convertMousePosToCellPos(event);
      if (pos) {
        updateSelection(leftTopPos, pos);
      }
    }
    isDragging = false;
    leftTopPos = null;
  }

  //TODO: why handleEnterPress takes numbers as arguments??
  function handleEnterPress(x: number, y: number) {
    const pos: CellPosition = {x, y};
    const state = sheetData.get(pos);
    if (state) {
      sheetData.set(pos, {
        ...state,
        isEditing: false,
        isSelected: false
      })
    }
    selectedPoss.delete(pos);

    const nextY = y + 1;
    if (nextY < rowCount) {
      for (const state of sheetData.values()) {
        state.isSelected = false;
        state.isEditing = false;
      }
      selectedPoss.clear()
      const pos: CellPosition = {x, y: nextY}
      const state = sheetData.get(pos);
      if (state) {
        state.isSelected = true;
        state.isEditing = true;
      }
      selectedPoss.add(pos)
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
            () => {
              const pos: CellPosition = {
                x: col.index,
                y: row.index,
              }
              const state = sheetData.get(pos);
              if (state) {
                return state;
              } else {
                //TODO: what to return when state is undefined
                const mockState: CellState = {
                  isEditing: false,
                  isSelected: false,
                  rawValue: "",
                  displayValue: ""
                }
                return mockState
              }
            },
            (newState: CellState) => {
              const pos: CellPosition = {
                x: col.index,
                y: row.index,
              }
              sheetData.set(pos, newState)
            }
          } 
          grid={sheetData}
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
          onEnterPress={handleEnterPress}
        />
        </div>
      {/each}
    {/each}
  </div>
</div>

<Chart {sheetData} bind:this={chartComponent}/>
