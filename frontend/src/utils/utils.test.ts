import { describe, expect, it } from "vitest";
import { byDate, fromLocalInputValue, relativeTime, toLocalInputValue } from "./date.ts";
import { CONFLICT_MESSAGE, errorsFromSave, hasErrors, validateCredentials, validateDeadline, validateTag, validateTask } from "./validation.ts";
import { search } from "./search.ts";
import { sameIds, toTaskInput } from "./inputs.ts";
import { ApiError } from "../api/client.ts";
import { makeDeadline, makeTag, makeTask } from "../test/fixtures.tsx";

describe("date helpers", () => {
    it("round-trips datetime-local values", () => {
        const iso = fromLocalInputValue("2026-03-01T09:30");
        expect(toLocalInputValue(iso)).toBe("2026-03-01T09:30");
    });

    it("describes relative times", () => {
        const now = new Date("2026-01-01T12:00:00Z");
        expect(relativeTime("2026-01-01T12:00:20Z", now)).toBe("now");
        expect(relativeTime("2026-01-01T12:30:00Z", now)).toBe("in 30 minutes");
        expect(relativeTime("2026-01-01T11:00:00Z", now)).toBe("1 hour ago");
        expect(relativeTime("2026-01-04T12:00:00Z", now)).toBe("in 3 days");
    });

    it("sorts soonest first", () => {
        const later = { date: "2026-02-01T00:00:00Z" };
        const sooner = { date: "2026-01-01T00:00:00Z" };
        expect([later, sooner].sort(byDate)).toEqual([sooner, later]);
    });
});

describe("validation", () => {
    it("validates tasks", () => {
        expect(validateTask({ title: " ", description: "", priority: "" })).toEqual({
            title: "Title is required",
            priority: "Priority must be selected",
        });
        expect(validateTask({ title: "x".repeat(201), description: "", priority: "low" }).title).toMatch(/at most 200/);
        expect(hasErrors(validateTask({ title: "ok", description: "", priority: "high" }))).toBe(false);
    });

    it("validates deadlines", () => {
        expect(validateDeadline({ title: "x", description: "", date: "" })).toEqual({ date: "Date is required" });
        expect(hasErrors(validateDeadline({ title: "x", description: "", date: "2026-01-01T10:00" }))).toBe(false);
    });

    it("validates tags", () => {
        expect(validateTag({ title: "x", color: "red" }).color).toBeDefined();
        expect(hasErrors(validateTag({ title: "x", color: "#AbCdEf" }))).toBe(false);
    });

    it("only enforces password length when registering", () => {
        expect(validateCredentials({ email: "a@b.co", password: "short" }, "login")).toEqual({});
        expect(validateCredentials({ email: "a@b.co", password: "short" }, "register").password).toBeDefined();
        expect(validateCredentials({ email: "nope", password: "" }, "login")).toEqual({
            email: "Enter a valid email address",
            password: "Password is required",
        });
    });

    it("maps server errors to fields", () => {
        expect(errorsFromSave(new ApiError(400, "title is required", "title"), ["title"])).toEqual({ title: "title is required" });
        expect(errorsFromSave(new ApiError(409, "changed elsewhere"), ["title"])).toEqual({ form: CONFLICT_MESSAGE });
        expect(errorsFromSave(new ApiError(500, "boom"), ["title"])).toEqual({ form: "boom" });
        expect(errorsFromSave("weird", [])).toEqual({ form: "Something went wrong" });
    });
});

describe("search", () => {
    const tag = makeTag({ title: "Urgent" });
    const data = {
        tags: [tag],
        tasks: [makeTask({ title: "Buy milk" }), makeTask({ title: "Call mum", tag_ids: [tag.id] })],
        deadlines: [makeDeadline({ title: "Exam", description: "Bring milk money" })],
    };

    it("matches titles, descriptions and tag names case-insensitively", () => {
        expect(search("MILK", data).tasks.map(t => t.title)).toEqual(["Buy milk"]);
        expect(search("milk", data).deadlines.map(d => d.title)).toEqual(["Exam"]);
        expect(search("urgent", data).tasks.map(t => t.title)).toEqual(["Call mum"]);
    });

    it("finds nothing for an empty query", () => {
        expect(search("  ", data)).toEqual({ tasks: [], deadlines: [] });
    });
});

describe("inputs", () => {
    it("compares id lists ignoring order", () => {
        expect(sameIds(["a", "b"], ["b", "a"])).toBe(true);
        expect(sameIds(["a"], ["a", "b"])).toBe(false);
    });

    it("keeps only editable task fields", () => {
        const input = toTaskInput(makeTask({ title: "x" }));
        expect(Object.keys(input).sort()).toEqual(["completed", "deadline_ids", "description", "priority", "tag_ids", "title"]);
    });
});
