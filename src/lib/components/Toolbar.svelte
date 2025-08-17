<script lang="ts">
  import Chart from "$lib/components/Chart.svelte"
  interface Props {
    chartComponent: Chart | null
  }
  let { chartComponent }: Props = $props();

  function handlerCreateChart() {
    if (chartComponent) {
      const success = chartComponent.drawChart();
      if (!success) {
        alert("select valid numbers")
      }
    }
  }
</script>


<div class="flex flex-row">
  <div>

<button onclick={handlerCreateChart} class="btn">Create Chart</button>
  </div>

  <div>

<button class="btn">Do nothing</button>
  </div>
  <div>

<button class="btn">Do nothing</button>
  </div>
</div>
