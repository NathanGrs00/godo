import { apiRequest } from "../api/client.ts";
import type { Task, TaskInput } from "../types";

export const getTasks = async (): Promise<Task[]> =>
    (await apiRequest<{ tasks: Task[] }>("/tasks")).tasks;

export const createTask = async (input: TaskInput): Promise<Task> =>
    (await apiRequest<{ task: Task }>("/tasks", { method: "POST", body: input })).task;

export const updateTask = async (id: string, input: TaskInput): Promise<Task> =>
    (await apiRequest<{ task: Task }>(`/tasks/${id}`, { method: "PUT", body: input })).task;

export const deleteTask = (id: string): Promise<void> =>
    apiRequest<void>(`/tasks/${id}`, { method: "DELETE" });
