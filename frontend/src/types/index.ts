/** Fields every stored item has. Mirrors models.Meta on the backend. */
export interface Meta {
    id: string;
    /** Increases on every change; sent back on update to detect edits made on another device. */
    version: number;
    created_at: string;
    updated_at: string;
}

export type Priority = "low" | "medium" | "high";

export const PRIORITIES: Priority[] = ["low", "medium", "high"];

export interface Task extends Meta {
    title: string;
    description: string;
    priority: Priority;
    completed: boolean;
    deadline_ids: string[];
    tag_ids: string[];
}

export interface Deadline extends Meta {
    title: string;
    description: string;
    /** ISO 8601 date-time in UTC. */
    date: string;
    /** True when the date is in the past (computed by the server). */
    passed: boolean;
    tag_ids: string[];
}

export interface Tag extends Meta {
    title: string;
    /** Hex colour such as "#8CFBDE". */
    color: string;
}

export interface User {
    id: string;
    email: string;
    created_at: string;
}

/** What the client sends to create or update a task. */
export interface TaskInput {
    title: string;
    description: string;
    priority: Priority;
    completed: boolean;
    deadline_ids: string[];
    tag_ids: string[];
    /** Required for updates: the version being edited. */
    version?: number;
}

/** What the client sends to create or update a deadline. */
export interface DeadlineInput {
    title: string;
    description: string;
    date: string;
    tag_ids: string[];
    version?: number;
}

/** What the client sends to create or update a tag. */
export interface TagInput {
    title: string;
    color: string;
    version?: number;
}

/** Returned by register and login. */
export interface Session {
    token: string;
    user: User;
}
