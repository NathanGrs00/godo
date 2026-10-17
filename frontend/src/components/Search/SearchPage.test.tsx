import { describe, expect, it } from "vitest";
import { screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import SearchPage from "./SearchPage.tsx";
import { fakeData, makeDeadline, makeTask, renderWithData } from "../../test/fixtures.tsx";

describe("SearchPage", () => {
    it("finds tasks and deadlines", async () => {
        renderWithData(<SearchPage />, fakeData({
            tasks: [makeTask({ title: "Read book" }), makeTask({ title: "Cook" })],
            deadlines: [makeDeadline({ title: "Book club" })],
        }));
        await userEvent.type(screen.getByLabelText("Search"), "book");
        expect(screen.getByRole("button", { name: /Read book/ })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: /Book club/ })).toBeInTheDocument();
        expect(screen.queryByText("Cook")).not.toBeInTheDocument();

        await userEvent.clear(screen.getByLabelText("Search"));
        await userEvent.type(screen.getByLabelText("Search"), "zzz");
        expect(screen.getByText(/Nothing found/)).toBeInTheDocument();
    });
});
