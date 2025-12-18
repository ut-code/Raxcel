<script lang="ts">
  import { onMount } from "svelte";
  import { ChatWithAI, GetMessages } from "../wailsjs/go/main/App";
  interface Props {
    isChatOpen: boolean;
  }
  let { isChatOpen = $bindable() }: Props = $props()
  type Message = {
    author: "ai" | "user";
    message: string;
  };
  let messages = $state<Message[]>([]);
  let isLoading = $state(false);
  let userMessage = $state("");

  onMount(async () => {
    // Load chat history when component mounts
    const res = await GetMessages();
    if (res.ok && res.messages) {
      messages = res.messages.map((msg) => ({
        author: msg.role === "user" ? "user" : "ai",
        message: msg.content,
      }));
    } else if (res.error) {
      console.error("Failed to load messages:", res.error);
    }
  });

  async function sendMessage() {
    const newUserMessage: Message = {
      author: "user",
      message: userMessage,
    };
    messages.push(newUserMessage);
    isLoading = true;
    const res = await ChatWithAI(userMessage);
    userMessage = "";
    if (!res.ok) {
      alert(`Error: ${res.message}`);
    }
    const newAiMessage: Message = {
      author: "ai",
      message: res.message,
    };
    messages.push(newAiMessage);
    isLoading = false;
  }
</script>

<div
  class="absolute top-0 right-0 w-80 h-full bg-base-100 border-l border-base-300 shadow-xl flex flex-col"
>
  <!-- Header -->
  <div class="flex items-center justify-between p-4 border-b border-base-300">
    <h2 class="text-lg font-bold flex items-center gap-2">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
      </svg>
      AI Chat
    </h2>
    <button 
      onclick={() => isChatOpen = false} 
      class="btn btn-sm btn-ghost btn-circle"
      aria-label="Close chat"
    >
      <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  </div>

  <!-- Messages Container -->
  <div class="flex-1 overflow-y-auto p-4 space-y-4">
  {#each messages as message}
    {#if message.author === "ai"}
      <div class="chat chat-start">
        <div class="chat-header">AI</div>
        <div class="chat-bubble">{message.message}</div>
      </div>
    {:else}
      <div class="chat chat-end">
        <div class="chat-header">You</div>
        <div class="chat-bubble">{message.message}</div>
      </div>
    {/if}
  {/each}
  {#if isLoading}
    <div class="flex justify-center py-2">
      <span class="loading loading-dots loading-md"></span>
    </div>
  {/if}
  </div>

  <!-- Input Area -->
  <div class="p-4 border-t border-base-300">
    <div class="flex gap-2">
      <input 
        class="input input-bordered flex-1" 
        bind:value={userMessage} 
        placeholder="Type a message..." 
        onkeydown={(e) => {
          if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
          }
        }}
      />
      <button 
        class="btn btn-primary" 
        onclick={sendMessage}
        aria-label="Send message button"
        disabled={!userMessage.trim() || isLoading}
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
        </svg>
      </button>
    </div>
  </div>
</div>
