<script lang="ts">
  import { Login } from "$lib/wailsjs/go/main/App";
  import { goto } from "$app/navigation";

  let email = $state("");
  let password = $state("");
  let error = $state("");
  let isLoading = $state(false);

  async function handleLogin() {
    if (!email || !password) {
      error = "Please fill in all fields";
      return;
    }

    isLoading = true;
    error = "";

    const result = await Login(email, password);

    isLoading = false;

    if (result.ok) {
      goto("/"); // Redirect to main page on success
    } else {
      error = result.message;
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === "Enter") {
      handleLogin();
    }
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-base-200">
  <div class="card w-96 bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">Sign In</h2>

      {#if error}
        <div class="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      {/if}

      <div class="form-control">
        <label class="label" for="email">
          <span class="label-text">Email</span>
        </label>
        <input
          id="email"
          type="email"
          placeholder="email@example.com"
          class="input input-bordered"
          bind:value={email}
          onkeypress={handleKeyPress}
          disabled={isLoading}
        />
      </div>

      <div class="form-control mt-4">
        <label class="label" for="password">
          <span class="label-text">Password</span>
        </label>
        <input
          id="password"
          type="password"
          placeholder="Enter your password"
          class="input input-bordered"
          bind:value={password}
          onkeypress={handleKeyPress}
          disabled={isLoading}
        />
      </div>

      <div class="card-actions mt-6">
        <button
          class="btn btn-primary w-full"
          onclick={handleLogin}
          disabled={isLoading}
        >
          {#if isLoading}
            <span class="loading loading-spinner"></span>
          {/if}
          Sign In
        </button>
      </div>

      <div class="text-center mt-4">
        <a href="/signup" class="link link-primary">
          Don't have an account? Sign up
        </a>
      </div>
    </div>
  </div>
</div>
