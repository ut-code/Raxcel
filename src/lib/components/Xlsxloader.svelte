<script lang="ts">
  import * as XLSX from "xlsx";

  let currentFile = $state<File | null>(null);
  let jsonData = $state<Record<string, any[]> | null>(null);
  let error = $state<string | null>(null);

  async function processFile(file: File): Promise<void> {
    try {
      error = null;
      const arrayBuffer: ArrayBuffer = await file.arrayBuffer();

      // xlsxファイルを読み込み
      const workbook: XLSX.WorkBook = XLSX.read(arrayBuffer);

      // 全シートのデータを取得
      const result: Record<string, any[]> = {};
      workbook.SheetNames.forEach((sheetName: string) => {
        const worksheet: XLSX.WorkSheet = workbook.Sheets[sheetName];

        // シートをJSONに変換（ヘッダー行を自動認識）
        const jsonSheet: any[] = XLSX.utils.sheet_to_json(worksheet);
        result[sheetName] = jsonSheet;
      });

      jsonData = result;
      console.log("JSON Data Loaded");
    } catch (err) {
      error = `エラー: ${err instanceof Error ? err.message : "Unknown error"}`;
      jsonData = null;
    }
  }

  // ファイルアップロード処理
  export async function handleFileUpload(): Promise<void> {
    if (!currentFile) return;
    await processFile(currentFile);
  }

  // ファイル設定用の関数
  export function setFiles(newFiles: FileList): void {
    currentFile = newFiles[0] || null;
  }

  // データ取得用の関数
  export function getData() {
    return { jsonData, error };
  }
</script>
