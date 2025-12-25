<script lang="ts">
  import { goto } from "$app/navigation";
  import Chart from "$lib/components/Chart.svelte";
  import Xlsxloader from "$lib/components/Xlsxloader.svelte";
  import XlsxExporter from "./XlsxExporter.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import { SignOut } from "$lib/wailsjs/go/main/App";
  import type { Cell } from "$lib/types";
  import { authState } from "$lib/stores/auth.svelte";

  interface Props {
    chartComponent: Chart | null;
    grid: Record<string, Cell>;
    isChatOpen: boolean;
  }
  let {
    chartComponent,
    grid = $bindable(),
    isChatOpen = $bindable(),
  }: Props = $props();

  // Dialog state
  let dialogOpen = $state(false);
  let dialogMessage = $state("");
  let dialogTitle = $state("");
  let dialogType = $state<"info" | "error" | "warning" | "success">("info");

  function showDialog(
    message: string,
    title: string = "Notification",
    type: "info" | "error" | "warning" | "success" = "info",
  ) {
    dialogMessage = message;
    dialogTitle = title;
    dialogType = type;
    dialogOpen = true;
  }

  function handlerCreateChart() {
    if (chartComponent) {
      const success = chartComponent.drawChart();
      if (!success) {
        showDialog(
          "Please select valid numbers",
          "Chart Creation Failed",
          "warning",
        );
      }
    }
  }

  async function handleSignOut() {
    const result = await SignOut();
    if (result.error === "") {
      authState.logout();
    } else {
      showDialog(`Sign out failed: ${result.error}`, "Sign Out Error", "error");
    }
  }
</script>

<div class="navbar bg-base-200 shadow-lg">
  <div class="flex-1 gap-2">
    <!-- Divider -->
    <div class="divider divider-horizontal"></div>

    <!-- File Operations -->
    <div class="flex gap-1">
      <Xlsxloader bind:grid />
      <XlsxExporter bind:grid />
    </div>

    <!-- Divider -->
    <div class="divider divider-horizontal"></div>

    <!-- Chart Operations -->
    <button onclick={handlerCreateChart} class="btn btn-sm btn-primary gap-2">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="h-4 w-4"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
        />
      </svg>
      Create Chart
    </button>
  </div>

  <div class="flex-none gap-2">
    {#if !isChatOpen && authState.isLoggedIn}
      <button
        class="btn btn-sm btn-accent gap-2"
        onclick={() => {
          isChatOpen = true;
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
          />
        </svg>
        AI Chat
      </button>
    {/if}

    <!-- User Menu -->
    {#if authState.isLoggedIn}
      <button class="btn btn-sm btn-outline" onclick={handleSignOut}>
        Sign Out
      </button>
    {:else}
      <button class="btn btn-sm btn-outline" onclick={() => goto("/signin")}>
        Sign In
      </button>
      <button class="btn btn-sm btn-primary" onclick={() => goto("/signup")}>
        Sign Up
      </button>
    {/if}
  </div>
</div>

<Dialog
  bind:isOpen={dialogOpen}
  message={dialogMessage}
  title={dialogTitle}
  type={dialogType}
/>
