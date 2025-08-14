<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import type { Action } from "svelte/action";

  interface Props {
    cell: CellType;
    onMouseDown: (event: MouseEvent) => void;
    onMouseUp: (event: MouseEvent) => void;
    onEnterPress: (x: number, y: number) => void;
  }

  let { cell, onMouseDown, onMouseUp, onEnterPress }: Props = $props();

  const focusInput: Action = (node) => {
    node.focus();
    if (node instanceof HTMLInputElement) {
      node.select()
    }
  }

  const parseRawValue: Action = (_node) => {
   // Parse here
    cell.displayValue = cell.rawValue
  }

// TODO: focus the lower cell when the user clicks Enter
</script>

{#if cell.isEditing}
  <input
    type="text"
    class="w-24 h-12 border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={cell.rawValue}
    use:focusInput
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        // cell.isEditing = false;
        // cell.isSelected = false;
        onEnterPress(cell.x, cell.y);
      }
    }}
    onblur = {() => {
      cell.isEditing = false;
      cell.isSelected = false;
    }}
  />
{:else}
  <button
    class={[
      "w-24 h-12 border border-gray-300 box-border cursor-pointer flex-shrink-0",
        cell.isSelected ? 
        "bg-gray-200" : 
        "bg-white"
    ]}
    onmousedown={onMouseDown}
    onmouseup={onMouseUp}
    use:parseRawValue
    onclick={() => {
      cell.isEditing = true;
    }}
  >
    {cell.displayValue}
  </button>
{/if}
