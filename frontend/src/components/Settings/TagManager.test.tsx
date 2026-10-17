import { describe, expect, it, vi } from "vitest";
import { fireEvent, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import TagManager from "./TagManager.tsx";
import { fakeData, makeTag, renderWithData } from "../../test/fixtures.tsx";

describe("TagManager", () => {
    it("creates a tag", async () => {
        const { data } = renderWithData(<TagManager />);
        await userEvent.type(screen.getByLabelText("Tag name"), "school");
        fireEvent.change(screen.getByLabelText("Tag colour"), { target: { value: "#db2763" } });
        await userEvent.click(screen.getByRole("button", { name: "Add tag" }));
        expect(data.createTag).toHaveBeenCalledWith({ title: "school", color: "#DB2763" });
    });

    it("requires a name", async () => {
        const { data } = renderWithData(<TagManager />);
        await userEvent.click(screen.getByRole("button", { name: "Add tag" }));
        expect(screen.getByText("Title is required")).toBeInTheDocument();
        expect(data.createTag).not.toHaveBeenCalled();
    });

    it("renames and deletes tags", async () => {
        const tag = makeTag({ title: "old" });
        const data = fakeData({ tags: [tag] });
        vi.spyOn(window, "confirm").mockReturnValue(true);
        renderWithData(<TagManager />, data);

        const name = screen.getByLabelText("Name of old");
        await userEvent.clear(name);
        await userEvent.type(name, "new{Enter}");
        expect(data.updateTag).toHaveBeenCalledWith(tag, { title: "new", color: tag.color });

        await userEvent.click(screen.getByRole("button", { name: "Delete tag old" }));
        expect(data.deleteTag).toHaveBeenCalledWith(tag.id);
    });
});
