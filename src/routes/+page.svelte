<script lang="ts">
  import { invoke } from "@tauri-apps/api/core";

  let name = $state("");
  let randomName = $state("");
  let greetMsg = $state("");

  async function greet(name: string): Promise<string> {
    return await invoke("greet", { name });
  }

  async function onsubmit(event: Event) {
    event.preventDefault();
    if (name === "" || name === randomName) {
      randomName = Math.random().toString(36).substring(2, 9);
      name = randomName;
    }
    greetMsg = await greet(name);
  }
</script>

<main class="container mx-auto p-4 text-center">
  <h1 class="text-4xl font-bold">Welcome to Tauri + Svelte</h1>

  <div class="flex justify-center">
    <a href="https://vitejs.dev" target="_blank" class="m-4">
      <img src="/vite.svg" class="h-16 w-auto" alt="Vite Logo" />
    </a>
    <a href="https://tauri.app" target="_blank" class="m-4">
      <img src="/tauri.svg" class="h-16 w-auto" alt="Tauri Logo" />
    </a>
    <a href="https://kit.svelte.dev" target="_blank" class="m-4">
      <img src="/svelte.svg" class="h-16 w-auto" alt="SvelteKit Logo" />
    </a>
  </div>
  <p class="mt-4">
    Click on the Tauri, Vite, and SvelteKit logos to learn more.
  </p>

  <form class="mt-8 flex flex-row gap-2 justify-center" {onsubmit}>
    <input
      id="greet-input"
      placeholder="Enter a name..."
      bind:value={name}
      class="input input-bordered"
    />
    <button type="submit" class="btn btn-primary">Greet</button>
  </form>
  <p class="mt-4">{greetMsg}</p>
</main>
