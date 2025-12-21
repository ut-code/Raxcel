import Cell from "$lib/components/Cell.svelte";
import { beforeEach, describe, test, type Mock, vi, expect } from "vitest";
import type { Cell as CellType } from "$lib/types";
import { flushSync, mount, unmount } from "svelte";
import userEvent from "@testing-library/user-event";
import { fireEvent } from "@testing-library/svelte"

type MouseEventCallback = (event: MouseEvent) => void;
type KeyboardEventCallback = (event: KeyboardEvent) => void;

describe("Cell", () => {
    let grid: Record<string, CellType>;
    let mockCell: CellType;
    let handleEnterPress: Mock<KeyboardEventCallback>;
    let handleMouseDown: Mock<MouseEventCallback>;
    let handleMouseUp: Mock<MouseEventCallback>;

    beforeEach(() => {
        document.body.innerHTML = "";
        grid = {};
        mockCell = {
            x: 1,
            y: 1,
            rawValue: "2",
            displayValue: "2",
            isSelected: false,
            isEditing: false,
        };
        handleEnterPress = vi.fn();
        handleMouseDown = vi.fn();
        handleMouseUp = vi.fn();
    });

    test("通常セルはボタンとして表示される", () => {
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        expect(button).toBeTruthy();
        expect(button?.textContent).toBe("2");
        expect(button?.className).toContain("bg-white");

        unmount(component);
    });

    test("選択されたセルは背景色が変わる", () => {
        mockCell.isSelected = true;
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        expect(button?.className).toContain("bg-gray-200");

        unmount(component);
    });

    test("列ヘッダーセルは列名を表示する", () => {
        mockCell.x = 1;
        mockCell.y = 0;
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const div = document.querySelector("div");
        expect(div?.textContent?.trim()).toBe("A");

        unmount(component);
    });

    test("行ヘッダーセルは行番号を表示する", () => {
        mockCell.x = 0;
        mockCell.y = 3;
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const div = document.querySelector("div");
        expect(div?.textContent?.trim()).toBe("3");

        unmount(component);
    });

    test("ダブルクリックで編集モードに入る", async () => {
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        expect(button).toBeTruthy();

        if (button) {
            fireEvent.dblClick(button)
        }
        flushSync()

        const input = document.querySelector("input");
        expect(input).toBeTruthy()
        expect(input?.value).toBe("2");

        unmount(component);
    });

    test("編集モードでEnterキーを押すと処理が実行される", () => {
        mockCell.isEditing = true;
        mockCell.rawValue = "5";
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const input = document.querySelector("input");
        input?.dispatchEvent(
            new KeyboardEvent("keydown", { key: "Enter", bubbles: true })
        );
        flushSync();

        expect(handleEnterPress).toHaveBeenCalled();
        expect(mockCell.displayValue).toBe("5");

        unmount(component);
    });

    test("編集モードでEscapeキーを押すと編集モードが終了する", () => {
        mockCell.isEditing = true;
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const input = document.querySelector("input");
        input?.dispatchEvent(
            new KeyboardEvent("keydown", { key: "Escape", bubbles: true })
        );
        flushSync();

        expect(mockCell.isEditing).toBe(false);
        expect(mockCell.isSelected).toBe(false);

        unmount(component);
    });

    test("通常モードで文字キーを押すと編集モードに入り、その文字が入力される", () => {
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        button?.dispatchEvent(
            new KeyboardEvent("keydown", { key: "a", bubbles: true })
        );
        flushSync();

        expect(mockCell.isEditing).toBe(true);
        expect(mockCell.rawValue).toBe("a");

        unmount(component);
    });

    test("数式（=で始まる）が評価される", () => {
        mockCell.isEditing = true;
        mockCell.rawValue = "=2+3";
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const input = document.querySelector("input");
        input?.dispatchEvent(
            new KeyboardEvent("keydown", { key: "Enter", bubbles: true })
        );
        flushSync();

        expect(mockCell.displayValue).toBe("5");

        unmount(component);
    });

    test("数式のエラーは#ERRORと表示される", () => {
        mockCell.isEditing = true;
        mockCell.rawValue = "=invalid formula";
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const input = document.querySelector("input");
        input?.dispatchEvent(
            new KeyboardEvent("keydown", { key: "Enter", bubbles: true })
        );
        flushSync();

        expect(mockCell.displayValue).toBe("#ERROR");

        unmount(component);
    });

    test("マウスダウンイベントが呼ばれる", () => {
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        button?.dispatchEvent(new MouseEvent("mousedown", { bubbles: true }));

        expect(handleMouseDown).toHaveBeenCalled();

        unmount(component);
    });

    test("マウスアップイベントが呼ばれる", () => {
        const component = mount(Cell, {
            target: document.body,
            props: {
                grid,
                cell: mockCell,
                onEnterPress: handleEnterPress,
                onMouseDown: handleMouseDown,
                onMouseUp: handleMouseUp,
            },
        });

        const button = document.querySelector("button");
        button?.dispatchEvent(new MouseEvent("mouseup", { bubbles: true }));

        expect(handleMouseUp).toHaveBeenCalled();

        unmount(component);
    });
});
