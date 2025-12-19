import { describe, test, expect, beforeEach } from "bun:test";
import { resolveAll, getAffectedCells, updateCell, resetFormulaState } from "../src/lib/formula";
import type { Cell } from "../src/lib/types";

describe("formula.ts", () => {
    let grid: Record<string, Cell>;

    beforeEach(() => {
        grid = {};
        resetFormulaState();
    });

    describe("resolveAll", () => {
        test("resolves simple cell reference", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "10", displayValue: "10", isSelected: false, isEditing: false };
            const result = resolveAll("A1", grid, "1-2");
            expect(result).toBe("10");
        });

        test("resolves multiple cell references", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "5", displayValue: "5", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "3", displayValue: "3", isSelected: false, isEditing: false };
            const result = resolveAll("A1 + A2", grid, "1-3");
            expect(result).toBe("5 + 3");
        });

        test("resolves SUM function", () => {
            const result = resolveAll("SUM(1, 2, 3)", grid, "1-1");
            expect(result).toBe("sum(1, 2, 3)");
        });

        test("resolves range notation A1:A3", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "1", displayValue: "1", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "2", displayValue: "2", isSelected: false, isEditing: false };
            grid["1-3"] = { x: 1, y: 3, rawValue: "3", displayValue: "3", isSelected: false, isEditing: false };
            const result = resolveAll("A1:A3", grid, "1-4");
            expect(result).toBe("[[1,2,3]]");
        });

        test("resolves SUM with range", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "1", displayValue: "1", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "2", displayValue: "2", isSelected: false, isEditing: false };
            grid["1-3"] = { x: 1, y: 3, rawValue: "3", displayValue: "3", isSelected: false, isEditing: false };
            const result = resolveAll("SUM(A1:A3)", grid, "1-4");
            expect(result).toBe("sum([[1,2,3]])");
        });

        test("resolves multi-column range", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "1", displayValue: "1", isSelected: false, isEditing: false };
            grid["2-1"] = { x: 2, y: 1, rawValue: "2", displayValue: "2", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "3", displayValue: "3", isSelected: false, isEditing: false };
            grid["2-2"] = { x: 2, y: 2, rawValue: "4", displayValue: "4", isSelected: false, isEditing: false };
            const result = resolveAll("A1:B2", grid, "3-1");
            expect(result).toBe("[[1,3],[2,4]]");
        });

        test("resolves CORR function with 2D array", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "1", displayValue: "1", isSelected: false, isEditing: false };
            grid["2-1"] = { x: 2, y: 1, rawValue: "2", displayValue: "2", isSelected: false, isEditing: false };
            const result = resolveAll("CORR(A1:B1)", grid, "3-1");
            expect(result).toBe("corr([1],[2])");
        });

        test("resolves empty cells as 0", () => {
            const result = resolveAll("A1", grid, "1-2");
            expect(result).toBe("0");
        });

        test("detects circular reference via updateCell", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "=A1", displayValue: "0", isSelected: false, isEditing: false };
            updateCell("1-1", grid);
            expect(grid["1-1"].displayValue).toBe("#CIRCULAR");
        });
    });

    describe("getAffectedCells", () => {
        test("returns empty set when no dependencies", () => {
            const affected = getAffectedCells("1-1");
            expect(affected.size).toBe(0);
        });

        test("returns direct dependents", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "10", displayValue: "10", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1", displayValue: "10", isSelected: false, isEditing: false };
            resolveAll("A1", grid, "1-2");
            
            const affected = getAffectedCells("1-1");
            expect(affected.has("1-2")).toBe(true);
        });

        test("returns transitive dependents", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "10", displayValue: "10", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1", displayValue: "10", isSelected: false, isEditing: false };
            grid["1-3"] = { x: 1, y: 3, rawValue: "=A2", displayValue: "10", isSelected: false, isEditing: false };
            
            resolveAll("A1", grid, "1-2");
            resolveAll("A2", grid, "1-3");
            
            const affected = getAffectedCells("1-1");
            expect(affected.has("1-2")).toBe(true);
            expect(affected.has("1-3")).toBe(true);
        });
    });

    describe("updateCell", () => {
        test("updates dependent cells", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "10", displayValue: "10", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1 + 5", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("A1 + 5", grid, "1-2");
            updateCell("1-1", grid);
            
            expect(grid["1-2"].displayValue).toBe("15");
        });

        test("handles formula errors gracefully", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "10", displayValue: "10", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1 / 0", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("A1 / 0", grid, "1-2");
            updateCell("1-1", grid);
            
            expect(grid["1-2"].displayValue).toBe("Infinity");
        });

        test("detects circular reference with two cells", () => {
            // Start with A1 = 5, A2 = A1
            grid["1-1"] = { x: 1, y: 1, rawValue: "5", displayValue: "5", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1", displayValue: "0", isSelected: false, isEditing: false };
            
            // Establish dependency A2 -> A1
            updateCell("1-2", grid);
            expect(grid["1-2"].displayValue).toBe("5");
            
            // Now change A1 to =A2+1, creating a circular reference
            grid["1-1"] = { x: 1, y: 1, rawValue: "=A2+1", displayValue: "5", isSelected: false, isEditing: false };
            updateCell("1-1", grid);
            
            // Should detect circular reference
            expect(grid["1-1"].displayValue).toBe("#CIRCULAR");
        });

        test("updates chain of dependent cells", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "5", displayValue: "5", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=A1 * 2", displayValue: "0", isSelected: false, isEditing: false };
            grid["1-3"] = { x: 1, y: 3, rawValue: "=A2 + 1", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("A1 * 2", grid, "1-2");
            resolveAll("A2 + 1", grid, "1-3");
            updateCell("1-1", grid);
            
            expect(grid["1-2"].displayValue).toBe("10");
            expect(grid["1-3"].displayValue).toBe("11");
        });

        test("handles MAX function", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "5", displayValue: "5", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=MAX(A1, 10)", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("MAX(A1, 10)", grid, "1-2");
            updateCell("1-1", grid);
            
            expect(grid["1-2"].displayValue).toBe("10");
        });

        test("handles MIN function", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "5", displayValue: "5", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "=MIN(A1, 10)", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("MIN(A1, 10)", grid, "1-2");
            updateCell("1-1", grid);
            
            expect(grid["1-2"].displayValue).toBe("5");
        });

        test("handles MEAN function with range", () => {
            grid["1-1"] = { x: 1, y: 1, rawValue: "2", displayValue: "2", isSelected: false, isEditing: false };
            grid["1-2"] = { x: 1, y: 2, rawValue: "4", displayValue: "4", isSelected: false, isEditing: false };
            grid["1-3"] = { x: 1, y: 3, rawValue: "6", displayValue: "6", isSelected: false, isEditing: false };
            grid["1-4"] = { x: 1, y: 4, rawValue: "=MEAN(A1:A3)", displayValue: "0", isSelected: false, isEditing: false };
            
            resolveAll("MEAN(A1:A3)", grid, "1-4");
            updateCell("1-1", grid);
            
            expect(grid["1-4"].displayValue).toBe("4");
        });
    });
});