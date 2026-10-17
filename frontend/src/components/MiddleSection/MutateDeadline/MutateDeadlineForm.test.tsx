import { describe, expect, it, vi } from "vitest";
import { fireEvent, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import MutateDeadlineForm from "./MutateDeadlineForm.tsx";
import { renderWithData } from "../../../test/fixtures.tsx";

describe("MutateDeadlineForm", () => {
    it("requires a date", async () => {
        const { data } = renderWithData(<MutateDeadlineForm onClose={vi.fn()} />);
        await userEvent.type(screen.getByLabelText("Title"), "Exam");
        await userEvent.click(screen.getByRole("button", { name: "Save" }));
        expect(screen.getByText("Date is required")).toBeInTheDocument();
        expect(data.createDeadline).not.toHaveBeenCalled();
    });

    it("creates a deadline with the date converted to UTC", async () => {
        const onClose = vi.fn();
        const { data } = renderWithData(<MutateDeadlineForm onClose={onClose} />);
        await userEvent.type(screen.getByLabelText("Title"), "Exam");
        fireEvent.change(screen.getByLabelText("Due:"), { target: { value: "2030-06-01T09:00" } });
        await userEvent.click(screen.getByRole("button", { name: "Save" }));

        expect(data.createDeadline).toHaveBeenCalledWith({
            title: "Exam",
            description: "",
            date: new Date("2030-06-01T09:00").toISOString(),
            tag_ids: [],
        });
        expect(onClose).toHaveBeenCalled();
    });
});
