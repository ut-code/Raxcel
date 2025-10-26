<script lang="ts">
  import { browser } from "$app/environment";
  import "../app.css";
  type Props = {
    children: any;
  };
  const { children } = $props();
</script>

{@render children()}
