import React, { useCallback, useEffect, useMemo, useReducer, useState } from "react";
import { ApiError } from "../api/client.ts";
import * as deadlineService from "../services/deadlineService.ts";
import * as tagService from "../services/tagService.ts";
import * as taskService from "../services/taskService.ts";
import { dataReducer, initialDataState } from "../sync/dataReducer.ts";
import { connectEventStream, type ConnectionStatus, type SyncEvent } from "../sync/eventStream.ts";
import { toDeadlineInput, toTaskInput } from "../utils/inputs.ts";
import { DataContext, type DataContextValue } from "./dataContext.ts";
import { useAuth } from "./authContext.ts";

/**
 * Loads the user's tasks, deadlines and tags and keeps them in sync with the
 * server and with the user's other devices.
 *
 * - Changes made here are applied as soon as the server confirms them.
 * - Changes made on other devices arrive through the live event stream.
 * - After a reconnect everything is reloaded, so nothing missed offline is lost.
 * - If an edit is rejected because another device changed the item first
 *   (409), the data is reloaded and the error is passed on to the form.
 */
const DataProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { logout } = useAuth();
    const [state, dispatch] = useReducer(dataReducer, initialDataState);
    const [status, setStatus] = useState<ConnectionStatus>("connecting");

    const apply = useCallback((event: SyncEvent) => dispatch({ type: "event", event }), []);

    const reload = useCallback(async () => {
        try {
            const [tasks, deadlines, tags] = await Promise.all([
                taskService.getTasks(),
                deadlineService.getDeadlines(),
                tagService.getTags(),
            ]);
            dispatch({ type: "loaded", tasks, deadlines, tags });
        } catch (err) {
            dispatch({ type: "loadFailed", error: err instanceof Error ? err.message : "Failed to load data" });
        }
    }, []);

    // Subscribe first and load once the stream is ready, so no change can
    // slip in between the load and the subscription. If the stream cannot
    // connect, still load once so the app works without live updates.
    useEffect(() => {
        let loadedOnce = false;
        const disconnect = connectEventStream({
            onEvent: apply,
            onReady: () => {
                loadedOnce = true;
                reload();
            },
            onStatus: next => {
                setStatus(next);
                if (next === "offline" && !loadedOnce) {
                    loadedOnce = true;
                    reload();
                }
            },
            onUnauthorized: logout,
        });
        return disconnect;
    }, [apply, reload, logout]);

    /** Runs a mutation; on a version conflict reloads before rethrowing. */
    const guard = useCallback(async <T,>(mutation: () => Promise<T>): Promise<T> => {
        try {
            return await mutation();
        } catch (err) {
            if (err instanceof ApiError && err.status === 409) await reload();
            throw err;
        }
    }, [reload]);

    const value = useMemo<DataContextValue>(() => ({
        ...state,
        status,
        reload,

        createTask: input => guard(async () => {
            const task = await taskService.createTask(input);
            apply({ type: "task.created", id: task.id, data: task });
            return task;
        }),
        updateTask: (task, changes) => guard(async () => {
            const updated = await taskService.updateTask(task.id, { ...toTaskInput(task), ...changes, version: task.version });
            apply({ type: "task.updated", id: updated.id, data: updated });
            return updated;
        }),
        deleteTask: id => guard(async () => {
            await taskService.deleteTask(id);
            apply({ type: "task.deleted", id });
        }),

        createDeadline: input => guard(async () => {
            const deadline = await deadlineService.createDeadline(input);
            apply({ type: "deadline.created", id: deadline.id, data: deadline });
            return deadline;
        }),
        updateDeadline: (deadline, changes) => guard(async () => {
            const updated = await deadlineService.updateDeadline(deadline.id, { ...toDeadlineInput(deadline), ...changes, version: deadline.version });
            apply({ type: "deadline.updated", id: updated.id, data: updated });
            return updated;
        }),
        deleteDeadline: id => guard(async () => {
            await deadlineService.deleteDeadline(id);
            apply({ type: "deadline.deleted", id });
        }),

        createTag: input => guard(async () => {
            const tag = await tagService.createTag(input);
            apply({ type: "tag.created", id: tag.id, data: tag });
            return tag;
        }),
        updateTag: (tag, changes) => guard(async () => {
            const updated = await tagService.updateTag(tag.id, { title: tag.title, color: tag.color, ...changes, version: tag.version });
            apply({ type: "tag.updated", id: updated.id, data: updated });
            return updated;
        }),
        deleteTag: id => guard(async () => {
            await tagService.deleteTag(id);
            apply({ type: "tag.deleted", id });
        }),
    }), [state, status, reload, guard, apply]);

    return <DataContext.Provider value={value}>{children}</DataContext.Provider>;
};

export default DataProvider;
