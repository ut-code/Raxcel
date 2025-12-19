import type { Cell } from "./types";

export function gridToMarkdownTable(grid: Record<string, Cell>): string {
  // グリッドの範囲を取得
  const cells = Object.keys(grid);
  if (cells.length === 0) return "";
  
  const coords = cells.map(key => {
    const [x, y] = key.split('-').map(Number);
    return { x, y };
  });
  
  const maxX = Math.max(...coords.map(c => c.x));
  const maxY = Math.max(...coords.map(c => c.y));
  
  // Markdownテーブルを構築
  let markdown = "\nCurrent spreadsheet content:\n\n";
  
  // ヘッダー行（列名: A, B, C...）
  const headers = [];
  for (let x = 0; x <= maxX; x++) {
    headers.push(String.fromCharCode(65 + x));
  }
  markdown += "| " + headers.join(" | ") + " |\n";
  markdown += "|" + headers.map(() => "---").join("|") + "|\n";
  
  // データ行
  for (let y = 0; y <= maxY; y++) {
    const row = [];
    for (let x = 0; x <= maxX; x++) {
      const cell = grid[`${x}-${y}`];
      row.push(cell?.displayValue || "");
    }
    markdown += "| " + row.join(" | ") + " |\n";
  }
  
  return markdown;
}
