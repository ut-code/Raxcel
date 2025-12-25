<script lang="ts">
  import { Signup } from "$lib/wailsjs/go/main/App";

  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let error = $state("");
  let isLoading = $state(false);
  let isSent = $state(false);

  async function handleRegister() {
    if (!email || !password || !confirmPassword) {
      error = "Please fill in all fields";
      return;
    }

    if (password !== confirmPassword) {
      error = "Passwords do not match";
      return;
    }

    if (password.length < 6) {
      error = "Password must be at least 6 characters";
      return;
    }

    isLoading = true;
    error = "";

    const result = await Signup(email, password);

    isLoading = false;

    if (result.ok) {
      isSent = true;
    } else {
      error = result.message;
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === "Enter") {
      handleRegister();
    }
  }
</script>

<div class="flex items-center justify-center min-h-screen bg-base-200">
  <div class="card w-96 bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title text-2xl mb-4">Sign Up</h2>

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

      <div class="form-control mt-4">
        <label class="label" for="confirmPassword">
          <span class="label-text">Confirm Password</span>
        </label>
        <input
          id="confirmPassword"
          type="password"
          placeholder="Confirm your password"
          class="input input-bordered"
          bind:value={confirmPassword}
          onkeypress={handleKeyPress}
          disabled={isLoading}
        />
      </div>

      {#if isSent}
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
            <h3 class="font-bold">Check your email!</h3>
            <div class="text-sm">
              We've sent a verification link to <strong>{email}</strong>. Please
              check your inbox and click the link to verify your account.
            </div>
          </div>
        </div>
        <div class="card-actions mt-4">
          <a href="/signin" class="btn btn-outline w-full"> Go to Sign In </a>
        </div>
      {:else}
        <div class="card-actions mt-6">
          <button
            class="btn btn-primary w-full"
            onclick={handleRegister}
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="loading loading-spinner"></span>
            {/if}
            Sign Up
          </button>
        </div>
      {/if}

      <div class="text-center mt-4">
        <a href="/signin" class="link link-primary">
          Already have an account? Sign in
        </a>
      </div>
    </div>
  </div>
</div>
