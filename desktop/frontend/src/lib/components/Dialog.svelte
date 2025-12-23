<script lang="ts">
  interface Props {
    isOpen: boolean;
    title?: string;
    message: string;
    type?: "info" | "error" | "warning" | "success";
  }

  let {
    isOpen = $bindable(),
    title = "Notification",
    message,
    type = "info",
  }: Props = $props();

  function close() {
    isOpen = false;
  }

  // アイコンとスタイルをタイプに応じて変更
  const config = $derived(
    {
      info: {
        icon: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
        colorClass: "text-info",
      },
      error: {
        icon: "M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z",
        colorClass: "text-error",
      },
      warning: {
        icon: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z",
        colorClass: "text-warning",
      },
      success: {
        icon: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
        colorClass: "text-success",
      },
    }[type],
  );
</script>

{#if isOpen}
  <dialog class="modal modal-open">
    <div class="modal-box">
      <div class="flex items-start gap-3">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-6 w-6 flex-shrink-0 {config.colorClass}"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d={config.icon}
          />
        </svg>
        <div class="flex-1">
          <h3 class="font-bold text-lg">{title}</h3>
          <p class="py-4">{message}</p>
        </div>
      </div>
      <div class="modal-action">
        <button class="btn btn-sm" onclick={close}>Close</button>
      </div>
    </div>
    <!-- 背景をクリックしてもモーダルが閉じる -->
    <button
      type="button"
      class="modal-backdrop"
      onclick={close}
      aria-label="Close dialog"
    >
    </button>
  </dialog>
{/if}
