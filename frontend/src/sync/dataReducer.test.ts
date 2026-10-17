import { describe, expect, it } from "vitest";
import { applyEvent, dataReducer, initialDataState, upsert, type DataState } from "./dataReducer.ts";
import { makeDeadline, makeTag, makeTask } from "../test/fixtures.tsx";

const stateWith = (overrides: Partial<DataState>): DataState => ({ ...initialDataState, loading: false, ...overrides });

describe("upsert", () => {
    it("adds new items and replaces with equal or newer versions only", () => {
        const v1 = makeTask({ title: "v1" });
        const v2 = { ...v1, title: "v2", version: 2 };
        expect(upsert([], v1)).toEqual([v1]);
        expect(upsert([v1], v2)).toEqual([v2]);
        // A late event for an older version must not roll back the newer data
        expect(upsert([v2], v1)).toEqual([v2]);
    });
});

describe("applyEvent", () => {
    it("creates, updates and deletes tasks", () => {
        const task = makeTask();
        let state = applyEvent(stateWith({}), { type: "task.created", id: task.id, data: task });
        expect(state.tasks).toEqual([task]);

        const done = { ...task, completed: true, version: 2 };
        state = applyEvent(state, { type: "task.updated", id: task.id, data: done });
        expect(state.tasks[0].completed).toBe(true);

        state = applyEvent(state, { type: "task.deleted", id: task.id });
        expect(state.tasks).toEqual([]);
    });

    it("ignores a duplicate created event", () => {
        const task = makeTask();
        const once = applyEvent(stateWith({}), { type: "task.created", id: task.id, data: task });
        const twice = applyEvent(once, { type: "task.created", id: task.id, data: task });
        expect(twice.tasks).toHaveLength(1);
    });

    it("unlinks a deleted deadline from tasks", () => {
        const deadline = makeDeadline();
        const task = makeTask({ deadline_ids: [deadline.id] });
        const state = applyEvent(stateWith({ tasks: [task], deadlines: [deadline] }), { type: "deadline.deleted", id: deadline.id });
        expect(state.deadlines).toEqual([]);
        expect(state.tasks[0].deadline_ids).toEqual([]);
    });

    it("unlinks a deleted tag from tasks and deadlines", () => {
        const tag = makeTag();
        const other = makeTag();
        const task = makeTask({ tag_ids: [tag.id, other.id] });
        const deadline = makeDeadline({ tag_ids: [tag.id] });
        const state = applyEvent(stateWith({ tasks: [task], deadlines: [deadline], tags: [tag, other] }), { type: "tag.deleted", id: tag.id });
        expect(state.tags).toEqual([other]);
        expect(state.tasks[0].tag_ids).toEqual([other.id]);
        expect(state.deadlines[0].tag_ids).toEqual([]);
    });

    it("upserts deadlines and tags", () => {
        const deadline = makeDeadline();
        const tag = makeTag();
        let state = applyEvent(stateWith({}), { type: "deadline.created", id: deadline.id, data: deadline });
        state = applyEvent(state, { type: "tag.updated", id: tag.id, data: tag });
        expect(state.deadlines).toEqual([deadline]);
        expect(state.tags).toEqual([tag]);
    });
});

describe("dataReducer", () => {
    it("replaces everything on load and records failures", () => {
        const task = makeTask();
        const loaded = dataReducer(initialDataState, { type: "loaded", tasks: [task], deadlines: [], tags: [] });
        expect(loaded).toEqual({ tasks: [task], deadlines: [], tags: [], loading: false, error: null });

        const failed = dataReducer(initialDataState, { type: "loadFailed", error: "boom" });
        expect(failed.loading).toBe(false);
        expect(failed.error).toBe("boom");
    });
});
