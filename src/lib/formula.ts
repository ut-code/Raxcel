import type { Cell } from "./types";
import { sum } from "mathjs";

function columnToNumber(column: string): number {
  let result = 0;
  for (let i = 0; i < column.length; i++) {
    result = result * 26 + (column.charCodeAt(i) - "A".charCodeAt(0) + 1);
  }
  // Zero-based numbering
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

function parseA1Notation(cellRef: string): { x: number; y: number } | null {
  const match = cellRef.match(/^([A-Z]+)(\d+)$/);
  if (!match) return null;
  const [, column, row] = match;
  return {
    x: columnToNumber(column),
    // Zero-based numbering
    y: parseInt(row) - 1,
  };
}

function resolveSumFunction(
  formula: string,
  grid: Record<string, Cell>,
): string {
  const cellRangeRegex = /SUM\(([A-Z]+\d+:[A-Z]+\d+)\)/g;
  const rangeResolvedFormula = formula.replace(
    cellRangeRegex,
    (match, range) => {
      const [startCell, endCell] = range.split(":");
      const start = parseA1Notation(startCell);
      const end = parseA1Notation(endCell);

      if (!start || !end) {
        return "#ERROR";
      }

      const values = [];

      // Handle rectangular range (not just single row/column)
      for (
        let row = Math.min(start.y, end.y);
        row <= Math.max(start.y, end.y);
        row++
      ) {
        for (
          let col = Math.min(start.x, end.x);
          col <= Math.max(start.x, end.x);
          col++
        ) {
          const key = `${col}-${row}`;
          if (grid[key]) {
            const value = parseFloat(grid[key].displayValue);
            if (!isNaN(value)) {
              values.push(value);
            }
          }
        }
      }

      return String(sum(values));
    },
  );

  return rangeResolvedFormula;
}

function resolveCellReference(
  formula: string,
  grid: Record<string, Cell>,
): string {
  const cellRefRegex = /[A-Z]+\d+/g;
  const resolvedFormula = formula.replace(cellRefRegex, (match) => {
    const coords = parseA1Notation(match);
    if (!coords) return match;
    const key = `${coords.x}-${coords.y}`;
    const cell = grid[key];
    return cell ? cell.displayValue : "0";
  });
  return resolvedFormula;
}

export function resolveAll(
  formula: string,
  grid: Record<string, Cell>,
): string {
  const sumFuncResolvedFormula = resolveSumFunction(formula, grid);
  const cellRefResolvedFormula = resolveCellReference(
    sumFuncResolvedFormula,
    grid,
  );
  return cellRefResolvedFormula;
}
