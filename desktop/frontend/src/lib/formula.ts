import type { Cell } from "./types";
import * as math from "mathjs";

// Track which cells depend on which other cells
type DependencyMap = Record<string, Set<string>>;
const dependencyGraph: DependencyMap = {};

function columnToNumber(column: string): number {
  let result = 0;
  for (let i = 0; i < column.length; i++) {
    result = result * 26 + (column.charCodeAt(i) - "A".charCodeAt(0) + 1);
  }
  // Zero-based numbering
  return result;
}

function parseA1Notation(cellRef: string): { x: number; y: number } | null {
  const match = cellRef.match(/^([A-Z]+)(\d+)$/);
  if (!match) return null;
  const [, column, row] = match;
  return {
    x: columnToNumber(column),
    // Zero-based numbering
    y: parseInt(row),
  };
}

function getCellKey(x: number, y: number): string {
  return `${x}-${y}`;
}

function addDependency(dependentCell: string, dependsOnCell: string) {
  if (!dependencyGraph[dependsOnCell]) {
    dependencyGraph[dependsOnCell] = new Set();
  }
  dependencyGraph[dependsOnCell].add(dependentCell);
}

function clearDependencies(cellKey: string) {
  // Remove this cell as a dependent from all other cells
  Object.values(dependencyGraph).forEach((dependencies) => {
    dependencies.delete(cellKey);
  });
  // Clear this cell's dependents
  delete dependencyGraph[cellKey];
}

function getDependentCells(cellKey: string): Set<string> {
  return dependencyGraph[cellKey] || new Set();
}

/**
 * 範囲を配列に変換する関数
 * A1:A3 → [A1の値, A2の値, A3の値]
 * A1:B3 → [[A1の値, A2の値, A3の値], [B1の値, B2の値, B3の値]]
 */
function resolveRangeToArray(
  range: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
): string {
  const [startCell, endCell] = range.split(":");
  const start = parseA1Notation(startCell);
  const end = parseA1Notation(endCell);

  if (!start || !end) {
    return '"#ERROR"';
  }

  const minCol = Math.min(start.x, end.x);
  const maxCol = Math.max(start.x, end.x);
  const minRow = Math.min(start.y, end.y);
  const maxRow = Math.max(start.y, end.y);

  // 単一列の範囲（例: A1:A3）
  if (minCol === maxCol) {
    const values = [];
    for (let row = minRow; row <= maxRow; row++) {
      const key = getCellKey(minCol, row);
      addDependency(currentCellKey, key);

      if (grid[key]) {
        const value = parseFloat(grid[key].displayValue);
        values.push(isNaN(value) ? 0 : value);
      } else {
        values.push(0);
      }
    }
    return JSON.stringify(values);
  }

  // 単一行の範囲（例: A1:C1）
  if (minRow === maxRow) {
    const values = [];
    for (let col = minCol; col <= maxCol; col++) {
      const key = getCellKey(col, minRow);
      addDependency(currentCellKey, key);

      if (grid[key]) {
        const value = parseFloat(grid[key].displayValue);
        values.push(isNaN(value) ? 0 : value);
      } else {
        values.push(0);
      }
    }
    return JSON.stringify(values);
  }

  // 矩形範囲（例: A1:B3）→ 列ごとに配列を作成
  const columns = [];
  for (let col = minCol; col <= maxCol; col++) {
    const columnValues = [];
    for (let row = minRow; row <= maxRow; row++) {
      const key = getCellKey(col, row);
      addDependency(currentCellKey, key);

      if (grid[key]) {
        const value = parseFloat(grid[key].displayValue);
        columnValues.push(isNaN(value) ? 0 : value);
      } else {
        columnValues.push(0);
      }
    }
    columns.push(columnValues);
  }
  return JSON.stringify(columns);
}

/**
 * 範囲表記を配列に変換する
 * 関数内の範囲も含めてすべて処理する
 */
function resolveRangeNotation(
  formula: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
): string {
  // 範囲表記を配列に変換（A1:A3 または A1:B3 形式）
  const rangeRegex = /([A-Z]+\d+:[A-Z]+\d+)/g;
  const rangeResolved = formula.replace(rangeRegex, (match) => {
    return resolveRangeToArray(match, grid, currentCellKey);
  });

  return rangeResolved;
}

/**
 * 数式内の関数名をmathjsの関数呼び出しに変換
 */
function resolveFunctionCalls(formula: string): string {
  // 対応する関数のマッピング
  const functionMap: Record<string, string> = {
    SUM: "math.sum",
    MAX: "math.max",
    MIN: "math.min",
    MEAN: "math.mean",
    MEDIAN: "math.median",
    STD: "math.std",
    VARIANCE: "math.variance",
    // 相関係数は特殊処理が必要なため後述
  };

  let result = formula;

  // 各関数を置き換え
  Object.entries(functionMap).forEach(([excelFunc, mathFunc]) => {
    const regex = new RegExp(`\\b${excelFunc}\\(`, "g");
    result = result.replace(regex, `${mathFunc}(`);
  });

  return result;
}

function resolveCellReference(
  formula: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
): string {
  const cellRefRegex = /[A-Z]+\d+/g;
  const resolvedFormula = formula.replace(cellRefRegex, (match) => {
    const coords = parseA1Notation(match);
    if (!coords) return match;
    const key = getCellKey(coords.x, coords.y);
    // Add dependency
    addDependency(currentCellKey, key);

    const cell = grid[key];
    return cell ? cell.displayValue : "0";
  });
  return resolvedFormula;
}

export function resolveAll(
  formula: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
): string {
  // Clear existing dependencies for this cell
  clearDependencies(currentCellKey);

  // 1. 範囲表記を配列に変換
  const rangeResolvedFormula = resolveRangeNotation(
    formula,
    grid,
    currentCellKey,
  );

  // 2. 単一セル参照を解決
  const cellRefResolvedFormula = resolveCellReference(
    rangeResolvedFormula,
    grid,
    currentCellKey,
  );

  // 3. 関数名をmathjsの関数に変換
  const functionResolvedFormula = resolveFunctionCalls(cellRefResolvedFormula);

  return functionResolvedFormula;
}

// Function to get all cells that need to be updated when a cell changes
export function getAffectedCells(cellKey: string): Set<string> {
  const affected = new Set<string>();
  const queue = [cellKey];

  while (queue.length > 0) {
    const current = queue.shift()!;
    const dependents = getDependentCells(current);

    dependents.forEach((dependent) => {
      if (!affected.has(dependent)) {
        affected.add(dependent);
        queue.push(dependent);
      }
    });
  }

  return affected;
}

// Function to update a cell and all its dependents
export function updateCell(cellKey: string, grid: Record<string, Cell>): void {
  const affectedCells = getAffectedCells(cellKey);

  affectedCells.forEach((key) => {
    const cell = grid[key];
    if (cell && cell.rawValue.startsWith("=")) {
      const formula = cell.rawValue.slice(1);
      try {
        const resolvedFormula = resolveAll(formula, grid, key);
        grid[key] = {
          ...cell,
          displayValue: String(eval(resolvedFormula)),
        };
      } catch (error) {
        grid[key] = {
          ...cell,
          displayValue: "#ERROR",
        };
      }
    }
  });
}
