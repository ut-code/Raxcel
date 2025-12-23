import type { Cell } from "./types";
import * as math from "mathjs";

export type DependencyMap = Record<string, Set<string>>;
const dependencyGraph: DependencyMap = {};
const evaluationStack = new Set<string>();

export function resetFormulaState(): void {
  Object.keys(dependencyGraph).forEach((key) => delete dependencyGraph[key]);
  evaluationStack.clear();
}

function columnToNumber(column: string): number {
  let result = 0;
  for (let i = 0; i < column.length; i++) {
    result = result * 26 + (column.charCodeAt(i) - "A".charCodeAt(0) + 1);
  }
  return result;
}

function parseA1Notation(cellRef: string): { x: number; y: number } | null {
  const match = cellRef.match(/^([A-Z]+)(\d+)$/);
  if (!match) return null;
  const [, column, row] = match;
  return {
    x: columnToNumber(column),
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
  // Remove cellKey from all dependency lists (cellKey no longer references those cells)
  Object.values(dependencyGraph).forEach((dependencies) => {
    dependencies.delete(cellKey);
  });
  // Note: We don't delete dependencyGraph[cellKey] because other cells may still reference this cell
}

function getDependentCells(cellKey: string): Set<string> {
  return dependencyGraph[cellKey] || new Set();
}

function detectCircularReference(
  startCell: string,
  visited = new Set<string>(),
  path: string[] = [],
): string[] | null {
  if (path.includes(startCell)) {
    const cycleStart = path.indexOf(startCell);
    return [...path.slice(cycleStart), startCell];
  }

  if (visited.has(startCell)) {
    return null;
  }

  visited.add(startCell);
  path.push(startCell);

  const dependents = getDependentCells(startCell);

  for (const dependent of dependents) {
    const cycle = detectCircularReference(dependent, visited, [...path]);
    if (cycle) {
      return cycle;
    }
  }

  return null;
}

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

  const columns = [];
  for (let col = minCol; col <= maxCol; col++) {
    const columnValues = [];
    for (let row = minRow; row <= maxRow; row++) {
      const key = getCellKey(col, row);

      if (evaluationStack.has(key)) {
        throw new Error("#CIRCULAR");
      }

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

function resolveRangeNotation(
  formula: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
): string {
  const rangeRegex = /([A-Z]+\d+:[A-Z]+\d+)/g;
  let processed = formula.replace(rangeRegex, (match) => {
    return resolveRangeToArray(match, grid, currentCellKey);
  });

  return processed;
}

const COLUMN_SPLIT_FUNCTIONS = ["CORR"];

function resolveFunctionCalls(formula: string): string {
  const functionMap: Record<string, string> = {
    SUM: "sum",
    MAX: "max",
    MIN: "min",
    MEAN: "mean",
    MEDIAN: "median",
    STD: "std",
    VARIANCE: "variance",
    CORR: "corr",
  };

  let result = formula;

  COLUMN_SPLIT_FUNCTIONS.forEach((funcName) => {
    const pattern = new RegExp(
      `\\b${funcName}\\s*\\(\\s*(\\[\\[.+?\\]\\])\\s*\\)`,
      "gi",
    );

    result = result.replace(pattern, (match, array2d) => {
      const flattened = array2d.slice(1, -1);
      return `${funcName}(${flattened})`;
    });
  });

  Object.entries(functionMap).forEach(([excelFunc, mathFunc]) => {
    const regex = new RegExp(`\\b${excelFunc}\\(`, "gi");
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

    if (evaluationStack.has(key)) {
      throw new Error("#CIRCULAR");
    }

    addDependency(currentCellKey, key);

    // Check for circular reference: if the referenced cell depends on currentCellKey
    const cycle = detectCircularReference(key);
    if (cycle && cycle.includes(currentCellKey)) {
      throw new Error("#CIRCULAR");
    }

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
  clearDependencies(currentCellKey);

  const rangeResolvedFormula = resolveRangeNotation(
    formula,
    grid,
    currentCellKey,
  );

  const cellRefResolvedFormula = resolveCellReference(
    rangeResolvedFormula,
    grid,
    currentCellKey,
  );

  const functionResolvedFormula = resolveFunctionCalls(cellRefResolvedFormula);

  return functionResolvedFormula;
}

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

export function updateCell(cellKey: string, grid: Record<string, Cell>): void {
  const affectedCells = getAffectedCells(cellKey);

  // Include the cell itself if it has a formula
  const cellsToUpdate = new Set(affectedCells);
  if (grid[cellKey]?.rawValue.startsWith("=")) {
    cellsToUpdate.add(cellKey);
  }

  cellsToUpdate.forEach((key) => {
    const cell = grid[key];

    if (cell && cell.rawValue.startsWith("=")) {
      evaluationStack.add(key);

      const formula = cell.rawValue.slice(1);

      try {
        const resolvedFormula = resolveAll(formula, grid, key);
        const result = math.evaluate(resolvedFormula);

        grid[key] = {
          ...cell,
          displayValue: String(result),
        };
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : "#ERROR";
        grid[key] = {
          ...cell,
          displayValue: errorMessage.includes("#CIRCULAR")
            ? "#CIRCULAR"
            : "#ERROR",
        };
      } finally {
        evaluationStack.delete(key);
      }
    }
  });
}
