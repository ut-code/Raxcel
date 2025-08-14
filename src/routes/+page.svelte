<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import Cell from "$lib/components/Cell.svelte";
  import Chart from "$lib/components/Chart.svelte";

  //TODO: virtual scroll
  const rows = 10;
  const cols = 10;
  
  let grid: CellType[][] = $state(
    Array.from({ length: rows }, (_, y) =>
      Array.from({ length: cols }, (_, x) => ({
        x,
        y,
        rawValue: "",
        displayValue: "",
        isSelected: false,
        isEditing: false
      })),
    ),
  );

  let leftTopCell: CellType | null = $state(null);
  let isDragging = $state(false);
  
  function convertMouseLocToCell(event: MouseEvent) {
    const target = event.target as HTMLElement;
    const rowElem = target.parentElement;
    if (!rowElem || !rowElem.parentElement) return null;
    
    const y = Array.from(rowElem.parentElement.children).indexOf(rowElem);
    const x = Array.from(rowElem.children).indexOf(target);
    return grid[y][x];
  }

  function updateSelection(startCell: CellType, endCell: CellType) {
    // 全てのセルの選択状態をリセット
    grid.forEach((row) => row.forEach((cell) => (cell.isSelected = false)));
    
    // 選択範囲を計算
    const startX = Math.min(startCell.x, endCell.x);
    const endX = Math.max(startCell.x, endCell.x);
    const startY = Math.min(startCell.y, endCell.y);
    const endY = Math.max(startCell.y, endCell.y);
    
    // 選択範囲内のセルを選択状態にする
    for (let y = startY; y <= endY; y++) {
      for (let x = startX; x <= endX; x++) {
        grid[y][x].isSelected = true;
      }
    }
  }

  function handleMouseDown(event: MouseEvent) {
    grid.forEach((row) => row.forEach((cell) => (cell.isSelected = false)));
    const cell = convertMouseLocToCell(event);
    if (cell) {
      leftTopCell = cell;
      isDragging = true;
      cell.isSelected = true;
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
    grid[y][x].isEditing = false;
    grid[y][x].isSelected = false;

    const nextY = y + 1;
    if (nextY < rows) {
      grid.forEach((row) => row.forEach((cell) => {
        cell.isSelected = false;
        cell.isEditing = false;
      }));

      grid[nextY][x].isSelected = true;
      grid[nextY][x].isEditing = true;
    }
  }
</script>

<div class="flex flex-col" onmousemove={handleMouseMove} role="grid" tabindex=0>
  {#each grid as row}
    <div class="flex flex-row">
      {#each row as cell}
        <Cell 
          {cell} 
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
          onEnterPress={handleEnterPress}
        />
      {/each}
    </div>
  {/each}
</div>

<Chart {grid} />
