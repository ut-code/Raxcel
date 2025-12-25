<script lang="ts">
  import { Signin } from "$lib/wailsjs/go/main/App";
  import { authState } from "$lib/stores/auth.svelte";

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

    const result = await Signin(email, password);

    isLoading = false;

    if (result.error === "") {
      authState.login();
    } else {
      error = result.error;
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

      {#if authState.isLoggedIn}
        <div class="alert alert-success mt-6">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div>
            <h3 class="font-bold">Welcome back!</h3>
            <div class="text-sm">You're now signed in.</div>
          </div>
        </div>
        <div class="card-actions mt-4">
          <a href="/" class="btn btn-primary w-full"> Go to Home </a>
        </div>
      {:else}
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
      {/if}

      <div class="text-center mt-4">
        <a href="/signup" class="link link-primary">
          Don't have an account? Sign up
        </a>
      </div>
    </div>
  </div>
</div>
