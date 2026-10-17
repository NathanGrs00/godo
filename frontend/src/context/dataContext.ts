import { createContext, useContext } from "react";
import type { DataState } from "../sync/dataReducer.ts";
import type { ConnectionStatus } from "../sync/eventStream.ts";
import type { Deadline, DeadlineInput, Tag, TagInput, Task, TaskInput } from "../types";

export interface DataContextValue extends DataState {
    /** State of the live connection to the server. */
    status: ConnectionStatus;
    /** Reloads everything from the server. */
    reload: () => Promise<void>;

    createTask: (input: TaskInput) => Promise<Task>;
    /** Saves changes to a task. The task's version is sent along automatically. */
    updateTask: (task: Task, changes: Partial<TaskInput>) => Promise<Task>;
    deleteTask: (id: string) => Promise<void>;

    createDeadline: (input: DeadlineInput) => Promise<Deadline>;
    updateDeadline: (deadline: Deadline, changes: Partial<DeadlineInput>) => Promise<Deadline>;
    deleteDeadline: (id: string) => Promise<void>;

    createTag: (input: TagInput) => Promise<Tag>;
    updateTag: (tag: Tag, changes: Partial<TagInput>) => Promise<Tag>;
    deleteTag: (id: string) => Promise<void>;
}

export const DataContext = createContext<DataContextValue | null>(null);

/** Returns the user's synced data and actions. Must be used inside <DataProvider>. */
export const useData = (): DataContextValue => {
    const value = useContext(DataContext);
    if (!value) throw new Error("useData must be used inside <DataProvider>");
    return value;
};
