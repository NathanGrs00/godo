/** Test helpers: factories for data and a render function with a fake data context. */
import React from "react";
import { render } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { vi } from "vitest";
import { DataContext, type DataContextValue } from "../context/dataContext.ts";
import type { Deadline, Tag, Task } from "../types";

let counter = 0;
const nextId = () => `id${++counter}`.padEnd(24, "0");

const meta = () => ({
    id: nextId(),
    version: 1,
    created_at: "2026-01-01T00:00:00Z",
    updated_at: "2026-01-01T00:00:00Z",
});

export const makeTask = (overrides: Partial<Task> = {}): Task => ({
    ...meta(),
    title: "A task",
    description: "",
    priority: "medium",
    completed: false,
    deadline_ids: [],
    tag_ids: [],
    ...overrides,
});

export const makeDeadline = (overrides: Partial<Deadline> = {}): Deadline => ({
    ...meta(),
    title: "A deadline",
    description: "",
    date: "2030-01-01T09:00:00.000Z",
    passed: false,
    tag_ids: [],
    ...overrides,
});

export const makeTag = (overrides: Partial<Tag> = {}): Tag => ({
    ...meta(),
    title: "work",
    color: "#8CFBDE",
    ...overrides,
});

/** A DataContext value with empty data and mocked actions. */
export const fakeData = (overrides: Partial<DataContextValue> = {}): DataContextValue => ({
    tasks: [],
    deadlines: [],
    tags: [],
    loading: false,
    error: null,
    status: "live",
    reload: vi.fn(async () => {}),
    createTask: vi.fn(async input => makeTask(input)),
    updateTask: vi.fn(async (task, changes) => ({ ...task, ...changes, version: task.version + 1 })),
    deleteTask: vi.fn(async () => {}),
    createDeadline: vi.fn(async input => makeDeadline(input)),
    updateDeadline: vi.fn(async (d, changes) => ({ ...d, ...changes, version: d.version + 1 })),
    deleteDeadline: vi.fn(async () => {}),
    createTag: vi.fn(async input => makeTag(input)),
    updateTag: vi.fn(async (tag, changes) => ({ ...tag, ...changes, version: tag.version + 1 })),
    deleteTag: vi.fn(async () => {}),
    ...overrides,
});

/** Renders ui inside a router and the given fake data context. */
export const renderWithData = (ui: React.ReactElement, data: DataContextValue = fakeData()) => ({
    data,
    ...render(
        <MemoryRouter>
            <DataContext.Provider value={data}>{ui}</DataContext.Provider>
        </MemoryRouter>,
    ),
});
