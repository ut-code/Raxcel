<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import { clsx } from "clsx";

  interface Props {
    cell: CellType;
    onMouseDown: (event: MouseEvent) => void;
    onMouseUp: (event: MouseEvent) => void;
  }

  let { cell, onMouseDown, onMouseUp }: Props = $props();

  function focusInput(node: HTMLInputElement) {
    node.focus();
    node.select(); // テキストを全選択（オプション）
  }
</script>

{#if cell.isWritable}
  <input
    type="text"
    class="w-24 h-12 border border-gray-300 box-border cursor-pointer bg-white"
    bind:value={cell.value}
    use:focusInput
    onchange={(event: Event) => {
      cell.value = (event.target as HTMLInputElement).value;
    }}
    onkeydown={(event: KeyboardEvent) => {
      if (event.key === "Enter") {
        cell.isWritable = false;
        cell.isSelected = false;
      }
    }}
  />
{:else}
  <button
    class={clsx(
      "w-24 h-12 border border-gray-300 box-border cursor-pointer",
      {
        "bg-gray-200": cell.isSelected,
        "bg-white": !cell.isSelected,
      }
    )}
    onmousedown={onMouseDown}
    onmouseup={onMouseUp}
    ondblclick={() => {
      cell.isWritable = true;
    }}
  >
    {cell.value}
  </button>
{/if}
