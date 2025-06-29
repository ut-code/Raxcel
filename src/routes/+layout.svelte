<script lang="ts">
    import "../app.css";
    type Props = {
        children: any;
    };
    const { children } = $props();
</script>

{@render children()}
