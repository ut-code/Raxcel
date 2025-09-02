import type { Cell } from "./types";

function columnToNumber(column: string): number {
  let result = 0;
  for (let i = 0; i < column.length; i++) {
    result = result * 26 + (column.charCodeAt(i) - "A".charCodeAt(0) + 1);
  }
  return result - 1;
}

// function numberToColumn(num: number): string {
//   let result = "";
//   while (num >= 0) {
//     result = String.fromCharCode((num % 26) + 'A'.charCodeAt(0)) + result;
//     num = Math.floor(num / 26) - 1;
//   }
//   return result;
// }

function parseA1Notation(cellRed: string): { x: number; y: number } | null {
  const match = cellRed.match(/^([A-Z]+)(\d+)$/);
  if (!match) return null;
  const [, column, row] = match;
  return {
    x: columnToNumber(column),
    y: parseInt(row) - 1,
  };
}

// 範囲参照には対応していない
export function resolveCellReference(
  formula: string,
  grid: Record<string, Cell>,
): string {
  const cellRefRegex = /\b([A-Z]+\d+)\b/g;

  return formula.replace(cellRefRegex, (match) => {
    const coords = parseA1Notation(match);
    if (!coords) return match;
    const key = `${coords.x}-${coords.y}`;
    const value = grid[key];

    return value.displayValue;
  });
}
