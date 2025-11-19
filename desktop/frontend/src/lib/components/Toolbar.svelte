<script lang="ts">
  import { goto } from "$app/navigation";
  import Chart from "$lib/components/Chart.svelte";
  import Xlsxloader from "$lib/components/Xlsxloader.svelte";
  import type { Cell } from "$lib/types";

  interface Props {
    chartComponent: Chart | null;
    grid: Record<string, Cell>;
    isChatOpen: boolean;
    isLoggedIn: boolean;
  }
  let {
    chartComponent,
    grid = $bindable(),
    isChatOpen = $bindable(),
    isLoggedIn,
  }: Props = $props();

  function handlerCreateChart() {
    if (chartComponent) {
      const success = chartComponent.drawChart();
      if (!success) {
        alert("select valid numbers");
      }
    }
  }
</script>

<Xlsxloader bind:grid />

<div class="flex flex-row">
  <div>
    <button onclick={handlerCreateChart} class="btn">Create Chart</button>
  </div>
  {#if isChatOpen}
    <button
      class="btn"
      onclick={() => {
        isChatOpen = false;
      }}>Close AI Chat</button
    >
  {:else}
    <button
      class="btn"
      onclick={() => {
        if (isLoggedIn) {
          isChatOpen = true;
        } else {
          alert("sign in to chat with AI");
        }
      }}>Open AI Chat</button
    >
  {/if}

  <div>
    <button class="btn" onclick={() => goto("/signin")}
      >Go to SignIn page</button
    >
  </div>

  <div>
    <button class="btn" onclick={() => goto("/signup")}
      >Go to SignUp page</button
    >
  </div>
</div>
