<script lang="ts">
  import type { CellPosition, CellState } from "$lib/types.ts";
  import type { Action } from "svelte/action";
  import { evaluate } from "mathjs"
  import { resolveAll } from "$lib/formula";

  interface Props {
    state: CellState;
    pos: CellPosition;
    sheetData: Map<CellPosition, CellState>
    onMouseDown: (event: MouseEvent) => void;
    onMouseUp: (event: MouseEvent) => void;
    onEnterPress: (x: number, y: number) => void;
  }

  let { state = $bindable(), pos, sheetData,  onMouseDown, onMouseUp, onEnterPress }: Props = $props();

  const focusInput: Action = (node) => {
    node.focus();
    if (node instanceof HTMLInputElement) {
      node.select()
    }
  }

  const parseRawValue: Action = (_node) => {
    if (state.rawValue[0] === "=") {
      try {
      const formula = state.rawValue.slice(1);
      const resolvedFormula = resolveAll(formula, sheetData)
      state.displayValue = evaluate(resolvedFormula);
      } catch (error) {
        state.displayValue = "#ERROR";
      }
    } else {
      state.displayValue = state.rawValue;
    }
  }

</script>

{#if state.isEditing}
  <input
    type="text"
    class="w-full h-full border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={state.rawValue}
    use:focusInput
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        // cell.isEditing = false;
        // cell.isSelected = false;
        onEnterPress(pos.x, pos.y);
      }
    }}
    onblur = {() => {
      state.isEditing = false;
      state.isSelected = false;
    }}
  />
{:else}
  <button
    class={[
      "w-full h-full border border-gray-300 box-border cursor-pointer flex-shrink-0",
        state.isSelected ? 
        "bg-gray-200" : 
        "bg-white"
    ]}
    onmousedown={onMouseDown}
    onmouseup={onMouseUp}
    use:parseRawValue
    onclick={() => {
      state.isEditing = true;
    }}
  >
    {state.displayValue}
  </button>
{/if}
