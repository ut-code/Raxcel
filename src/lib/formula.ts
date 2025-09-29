import type { Cell } from "./types";
import { sum } from "mathjs";

// Track which cells depend on which other cells
type DependencyMap = Record<string, Set<string>>;
const dependencyGraph: DependencyMap = {};

function columnToNumber(column: string): number {
  let result = 0;
  for (let i = 0; i < column.length; i++) {
    result = result * 26 + (column.charCodeAt(i) - "A".charCodeAt(0) + 1);
  }
  // Zero-based numbering
  return result - 1;
}

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

function resolveSumFunction(
  formula: string,
  grid: Record<string, Cell>,
  currentCellKey: string,
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

      // Handle rectangular range
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
          const key = getCellKey(col, row);
          // Add dependency
          addDependency(currentCellKey, key);

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

  const sumFuncResolvedFormula = resolveSumFunction(
    formula,
    grid,
    currentCellKey,
  );
  const cellRefResolvedFormula = resolveCellReference(
    sumFuncResolvedFormula,
    grid,
    currentCellKey,
  );
  return cellRefResolvedFormula;
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
