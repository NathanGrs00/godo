import { describe, expect, it, vi } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import MutateTaskForm from "./MutateTaskForm.tsx";
import { fakeData, makeDeadline, makeTag, makeTask, renderWithData } from "../../../test/fixtures.tsx";
import { ApiError } from "../../../api/client.ts";
import { CONFLICT_MESSAGE } from "../../../utils/validation.ts";

describe("MutateTaskForm", () => {
    it("shows validation errors and does not save", async () => {
        const { data } = renderWithData(<MutateTaskForm onClose={vi.fn()} />);
        await userEvent.click(screen.getByRole("button", { name: "Save" }));

        expect(screen.getByText("Title is required")).toBeInTheDocument();
        expect(screen.getByText("Priority must be selected")).toBeInTheDocument();
        expect(data.createTask).not.toHaveBeenCalled();
    });

    it("creates a task with the chosen tags and deadlines", async () => {
        const tag = makeTag({ title: "work" });
        const deadline = makeDeadline({ title: "Launch" });
        const onClose = vi.fn();
        const { data } = renderWithData(<MutateTaskForm onClose={onClose} />, fakeData({ tags: [tag], deadlines: [deadline] }));

        await userEvent.type(screen.getByLabelText("Title"), "  Ship it ");
        await userEvent.selectOptions(screen.getByLabelText("Priority:"), "high");
        await userEvent.click(screen.getByRole("button", { name: "work" }));
        await userEvent.click(screen.getByRole("button", { name: /Launch/ }));
        await userEvent.click(screen.getByRole("button", { name: "Save" }));

        expect(data.createTask).toHaveBeenCalledWith({
            title: "Ship it",
            description: "",
            priority: "high",
            completed: false,
            tag_ids: [tag.id],
            deadline_ids: [deadline.id],
        });
        expect(onClose).toHaveBeenCalled();
    });

    it("edits a task and only enables Save after a change", async () => {
        const task = makeTask({ title: "Old", priority: "low" });
        const { data } = renderWithData(<MutateTaskForm task={task} onClose={vi.fn()} />);

        const save = screen.getByRole("button", { name: "Save" });
        expect(save).toBeDisabled();

        const title = screen.getByLabelText("Title");
        await userEvent.clear(title);
        await userEvent.type(title, "New");
        expect(save).toBeEnabled();
        await userEvent.click(save);

        expect(data.updateTask).toHaveBeenCalledWith(task, expect.objectContaining({ title: "New", priority: "low" }));
    });

    it("shows a conflict from another device and stays open", async () => {
        const task = makeTask({ title: "Old", priority: "low" });
        const onClose = vi.fn();
        const data = fakeData({
            updateTask: vi.fn().mockRejectedValue(new ApiError(409, "this item was changed on another device; reload and try again")),
        });
        renderWithData(<MutateTaskForm task={task} onClose={onClose} />, data);

        await userEvent.type(screen.getByLabelText("Title"), "!");
        await userEvent.click(screen.getByRole("button", { name: "Save" }));

        expect(await screen.findByText(CONFLICT_MESSAGE)).toBeInTheDocument();
        expect(onClose).not.toHaveBeenCalled();
    });

    it("detects edits made on another device, then lets the user overwrite them", async () => {
        const stale = makeTask({ title: "Old", priority: "low" });
        const latest = { ...stale, title: "Changed on phone", version: 2 };
        const updateTask = vi.fn()
            .mockRejectedValueOnce(new ApiError(409, "conflict"))
            .mockResolvedValueOnce({ ...latest, title: "Old!", version: 3 });
        renderWithData(<MutateTaskForm task={stale} onClose={vi.fn()} />, fakeData({ tasks: [latest], updateTask }));

        await userEvent.type(screen.getByLabelText("Title"), "!");
        await userEvent.click(screen.getByRole("button", { name: "Save" }));
        // First save uses the version the form was opened with
        expect(updateTask).toHaveBeenLastCalledWith(stale, expect.objectContaining({ title: "Old!" }));
        await screen.findByText(CONFLICT_MESSAGE);

        await userEvent.click(screen.getByRole("button", { name: "Save" }));
        // After the warning it saves against the newest copy
        expect(updateTask).toHaveBeenLastCalledWith(latest, expect.objectContaining({ title: "Old!" }));
    });
});
