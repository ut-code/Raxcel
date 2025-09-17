<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import type { Action } from "svelte/action";
  import { evaluate } from "mathjs"
  import { resolveAll } from "$lib/formula";

  interface Props {
    cell: CellType;
    grid: Record<string, CellType>;
    onMouseDown: (event: MouseEvent) => void;
    onMouseUp: (event: MouseEvent) => void;
    onEnterPress: (event: KeyboardEvent) => void;
  }

  let { cell = $bindable(), grid, onMouseDown, onMouseUp, onEnterPress }: Props = $props();

  const focusInput: Action = (node) => {
    node.focus();
    if (node instanceof HTMLInputElement) {
      node.select()
    }
  }

  const parseRawValue: Action = (_node) => {
    if (cell.rawValue[0] === "=") {
      try {
      const formula = cell.rawValue.slice(1);
      const resolvedFormula = resolveAll(formula, grid)
      cell.displayValue = evaluate(resolvedFormula);
      } catch (error) {
        cell.displayValue = "#ERROR";
      }
    } else {
      cell.displayValue = cell.rawValue;
    }
  }

</script>

{#if cell.isEditing}
  <input
    type="text"
    class="w-full h-full border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={cell.rawValue}
    use:focusInput
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        // focus the next cell
        onEnterPress(event);
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
      "w-full h-full border border-gray-300 box-border cursor-pointer flex-shrink-0",
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
