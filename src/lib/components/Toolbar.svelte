<script lang="ts">
  import Chart from "$lib/components/Chart.svelte";
  import Xlsxloader from "$lib/components/Xlsxloader.svelte";

  let isDragOver = $state(false);
  let fileInput = $state<HTMLInputElement>();
  let files = $state<FileList | undefined>();
  let jsonData = $state<Record<string, any[]> | null>(null);
  let error = $state<string | null>(null);

  interface Props {
    chartComponent: Chart | null;
  }
  let { chartComponent }: Props = $props();

  function handlerCreateChart() {
    if (chartComponent) {
      const success = chartComponent.drawChart();
      if (!success) {
        alert("select valid numbers");
      }
    }
  }

  let xlsxLoaderComponent: any;

  async function handleFileUpload(): Promise<void> {
    await xlsxLoaderComponent.handleFileUpload();
    const result = xlsxLoaderComponent.getData();
    jsonData = result.jsonData;
    error = result.error;
  }

  // ドラッグオーバー時
  function handleDragOver(e: DragEvent): void {
    e.preventDefault();
    isDragOver = true;
  }

  // ドラッグ離脱時
  function handleDragLeave(e: DragEvent): void {
    e.preventDefault();
    isDragOver = false;
  }

  // ドロップ時
  async function handleDrop(e: DragEvent): Promise<void> {
    e.preventDefault();
    isDragOver = false;

    const droppedFiles = e.dataTransfer?.files;

    if (droppedFiles && droppedFiles.length > 0) {
      const dt = new DataTransfer();
      dt.items.add(droppedFiles[0]);
      if (fileInput) {
        fileInput.files = dt.files;
        files = dt.files;
      }
      xlsxLoaderComponent.setFiles(dt.files);
    }
  }
</script>

<Xlsxloader bind:this={xlsxLoaderComponent} />

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

  <div class="ml-auto flex gap-2">
    <button onclick={handleFileUpload} class="btn"
      >Excelファイルをアップロード</button
    >

    <!-- ドラッグ&ドロップエリア -->
    <div
      role="button"
      tabindex="0"
      ondragover={handleDragOver}
      ondragleave={handleDragLeave}
      ondrop={handleDrop}
      style="border: 2px dashed {isDragOver
        ? '#0066cc'
        : '#ccc'}; padding: 10px; border-radius: 4px;"
    >
      <input
        bind:this={fileInput}
        type="file"
        accept=".xlsx,.xls"
        bind:files
        onchange={() => {
          if (files) {
            xlsxLoaderComponent.setFiles(files);
          }
        }}
      />
      {#if isDragOver}
        <p style="margin: 5px 0; font-size: 0.9em; color: #0066cc;">
          ファイルをドロップしてください
        </p>
      {:else}
        <p style="margin: 5px 0; font-size: 0.9em; color: #666;">
          または、ここにドラッグ＆ドロップ
        </p>
      {/if}
    </div>
  </div>
</div>

{#if error}
  <p style="color: red;">{error}</p>
{/if}

{#if jsonData}
  <p style="color: green;">ファイルが正常に読み込まれました</p>
{/if}
