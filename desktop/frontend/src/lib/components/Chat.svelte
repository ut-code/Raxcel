<script lang="ts">
  import { ChatWithAI } from "../wailsjs/go/main/App";
  type Message = {
    author: "ai" | "user";
    message: string;
  };
  let messages = $state<Message[]>([]);
  let isLoading = $state(false);
  let userMessage = $state("");
  async function sendMessage() {
    const newUserMessage: Message = {
      author: "user",
      message: userMessage,
    };
    messages.push(newUserMessage);
    isLoading = true;
    const aiRes = await ChatWithAI(userMessage);
    userMessage = "";
    if (!aiRes.ok) {
      alert(`Error: ${aiRes.message}`);
    }
    const newAiMessage: Message = {
      author: "ai",
      message: aiRes.message,
    };
    messages.push(newAiMessage);
    isLoading = false;
  }
</script>

<div
  class="absolute top-0 right-0 w-60 h-full p-2 bg-white border-3 border-solid"
>
  <h1>AI Chat</h1>
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
    <div class="loading loading-spinner loading-xs"></div>
  {/if}
  <input class="input" bind:value={userMessage} placeholder="say something" />
  <button class="btn" onclick={sendMessage}>Send!</button>
</div>
