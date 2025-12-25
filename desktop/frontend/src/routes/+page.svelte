<script lang="ts">
  import type { Cell as CellType } from "$lib/types.ts";
  import Chart from "$lib/components/Chart.svelte";
  import Toolbar from "$lib/components/Toolbar.svelte";
  import { SvelteSet } from "svelte/reactivity";
  import Chat from "$lib/components/Chat.svelte";
  import { GetCurrentUser } from "$lib/wailsjs/go/main/App";
  import { onMount } from "svelte";
  import { authState } from "$lib/stores/auth.svelte";
  import Sheet from "$lib/components/Sheet.svelte";

  let grid = $state<Record<string, CellType>>({});
  onMount(() => {
    const checkUser = async () => {
      const result = await GetCurrentUser();
      if (result.error === "") {
        authState.login();
      }
    };
    checkUser();
  });

  let selectedCells = $state(new SvelteSet<string>());

  let chartComponent: Chart | null = $state(null);

  let isChatOpen = $state(false);
</script>

<div class="h-screen flex flex-col">
  <Toolbar {chartComponent} bind:grid bind:isChatOpen />

  <Sheet bind:grid bind:selectedCells />

  {#if isChatOpen}
    <Chat bind:isChatOpen {grid} />
  {/if}

  <Chart {grid} {selectedCells} bind:this={chartComponent} />
</div>
