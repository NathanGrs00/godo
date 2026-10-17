import { describe, expect, it, vi } from "vitest";
import { screen } from "@testing-library/react";
import DeadlineList from "./DeadlineList.tsx";
import { fakeData, makeDeadline, makeTask, renderWithData } from "../../test/fixtures.tsx";

describe("DeadlineList", () => {
    it("shows upcoming deadlines soonest first, then passed ones, with task progress", () => {
        const later = makeDeadline({ title: "Later", date: "2031-01-01T00:00:00Z" });
        const sooner = makeDeadline({ title: "Sooner", date: "2030-01-01T00:00:00Z" });
        const passed = makeDeadline({ title: "Gone", date: "2020-01-01T00:00:00Z", passed: true });
        const tasks = [
            makeTask({ title: "a", completed: true, deadline_ids: [sooner.id] }),
            makeTask({ title: "b", deadline_ids: [sooner.id] }),
        ];
        renderWithData(<DeadlineList deadlines={[passed, later, sooner]} onEdit={vi.fn()} />, fakeData({ tasks }));

        expect(screen.getAllByRole("heading").map(h => h.textContent)).toEqual(["Sooner", "Later", "Gone"]);
        expect(screen.getByText("1/2 tasks done")).toBeInTheDocument();
        expect(screen.getByText("Passed")).toBeInTheDocument();
    });
});
