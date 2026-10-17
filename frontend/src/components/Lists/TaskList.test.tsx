import { describe, expect, it, vi } from "vitest";
import { screen, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TaskList from "./TaskList.tsx";
import { fakeData, makeDeadline, makeTag, makeTask, renderWithData } from "../../test/fixtures.tsx";

describe("TaskList", () => {
    it("shows an empty message", () => {
        renderWithData(<TaskList tasks={[]} onEdit={vi.fn()} />);
        expect(screen.getByText(/No tasks yet/)).toBeInTheDocument();
    });

    it("lists open tasks first with their tags and deadlines", () => {
        const tag = makeTag({ title: "home" });
        const deadline = makeDeadline({ title: "Weekend" });
        const done = makeTask({ title: "Done one", completed: true });
        const open = makeTask({ title: "Open one", tag_ids: [tag.id], deadline_ids: [deadline.id] });
        renderWithData(<TaskList tasks={[done, open]} onEdit={vi.fn()} />, fakeData({ tags: [tag], deadlines: [deadline] }));

        const titles = screen.getAllByRole("heading").map(h => h.textContent);
        expect(titles).toEqual(["Open one", "Done one"]);
        const openRow = screen.getByRole("button", { name: 'Edit "Open one"' });
        expect(within(openRow).getByText("home")).toBeInTheDocument();
        expect(within(openRow).getByText(/Weekend/)).toBeInTheDocument();
    });

    it("completes, edits and deletes tasks", async () => {
        const task = makeTask({ title: "Buy milk" });
        const onEdit = vi.fn();
        const data = fakeData({ tasks: [task] });
        vi.spyOn(window, "confirm").mockReturnValue(true);
        renderWithData(<TaskList tasks={[task]} onEdit={onEdit} />, data);

        await userEvent.click(screen.getByRole("checkbox"));
        expect(data.updateTask).toHaveBeenCalledWith(task, { completed: true });

        await userEvent.click(screen.getByRole("button", { name: 'Edit "Buy milk"' }));
        expect(onEdit).toHaveBeenCalledWith(task);

        await userEvent.click(screen.getByRole("button", { name: 'Delete "Buy milk"' }));
        expect(data.deleteTask).toHaveBeenCalledWith(task.id);
    });

    it("does not delete when the user cancels, and shows errors", async () => {
        const task = makeTask({ title: "Keep me" });
        const data = fakeData({ updateTask: vi.fn().mockRejectedValue(new Error("Cannot reach the server")) });
        vi.spyOn(window, "confirm").mockReturnValue(false);
        renderWithData(<TaskList tasks={[task]} onEdit={vi.fn()} />, data);

        await userEvent.click(screen.getByRole("button", { name: 'Delete "Keep me"' }));
        expect(data.deleteTask).not.toHaveBeenCalled();

        await userEvent.click(screen.getByRole("checkbox"));
        expect(await screen.findByRole("alert")).toHaveTextContent("Cannot reach the server");
    });
});
