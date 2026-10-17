/**
 * Pure state logic for the user's synced data.
 *
 * Every change, whether it was made on this device (the API response) or on
 * another device (a sync event), goes through this reducer. Items are only
 * replaced by an equal or newer version, so a late or duplicate event can
 * never roll back newer data.
 */
import type { Deadline, Meta, Tag, Task } from "../types";
import type { SyncEvent } from "./eventStream.ts";

export interface DataState {
    tasks: Task[];
    deadlines: Deadline[];
    tags: Tag[];
    /** True until the first full load finished. */
    loading: boolean;
    /** Error from the last full load, if it failed. */
    error: string | null;
}

export type DataAction =
    | { type: "loaded"; tasks: Task[]; deadlines: Deadline[]; tags: Tag[] }
    | { type: "loadFailed"; error: string }
    | { type: "event"; event: SyncEvent };

export const initialDataState: DataState = {
    tasks: [],
    deadlines: [],
    tags: [],
    loading: true,
    error: null,
};

/** Inserts item, or replaces the existing item when item is not older. */
export const upsert = <T extends Meta>(items: T[], item: T): T[] => {
    const index = items.findIndex(existing => existing.id === item.id);
    if (index === -1) return [...items, item];
    if (items[index].version > item.version) return items;
    const next = items.slice();
    next[index] = item;
    return next;
};

/** Removes the item with the given ID. */
export const remove = <T extends Meta>(items: T[], id: string): T[] =>
    items.filter(item => item.id !== id);

/** Removes a linked ID from the given field of every item. */
const unlink = <T, K extends keyof T>(items: T[], field: K, id: string): T[] =>
    items.map(item => {
        const ids = item[field] as unknown as string[];
        return ids.includes(id) ? { ...item, [field]: ids.filter(x => x !== id) } : item;
    });

/** Applies one change event to the state. */
export const applyEvent = (state: DataState, event: SyncEvent): DataState => {
    switch (event.type) {
        case "task.created":
        case "task.updated":
            return { ...state, tasks: upsert(state.tasks, event.data) };
        case "task.deleted":
            return { ...state, tasks: remove(state.tasks, event.id) };
        case "deadline.created":
        case "deadline.updated":
            return { ...state, deadlines: upsert(state.deadlines, event.data) };
        case "deadline.deleted":
            // The server also unlinks it from tasks and sends task updates;
            // unlinking here too keeps the screen right in the meantime.
            return {
                ...state,
                deadlines: remove(state.deadlines, event.id),
                tasks: unlink(state.tasks, "deadline_ids", event.id),
            };
        case "tag.created":
        case "tag.updated":
            return { ...state, tags: upsert(state.tags, event.data) };
        case "tag.deleted":
            return {
                ...state,
                tags: remove(state.tags, event.id),
                tasks: unlink(state.tasks, "tag_ids", event.id),
                deadlines: unlink(state.deadlines, "tag_ids", event.id),
            };
        default:
            return state;
    }
};

export const dataReducer = (state: DataState, action: DataAction): DataState => {
    switch (action.type) {
        case "loaded":
            return { tasks: action.tasks, deadlines: action.deadlines, tags: action.tags, loading: false, error: null };
        case "loadFailed":
            return { ...state, loading: false, error: action.error };
        case "event":
            return applyEvent(state, action.event);
    }
};
