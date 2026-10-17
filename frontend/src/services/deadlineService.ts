import { apiRequest } from "../api/client.ts";
import type { Deadline, DeadlineInput } from "../types";

export const getDeadlines = async (): Promise<Deadline[]> =>
    (await apiRequest<{ deadlines: Deadline[] }>("/deadlines")).deadlines;

export const createDeadline = async (input: DeadlineInput): Promise<Deadline> =>
    (await apiRequest<{ deadline: Deadline }>("/deadlines", { method: "POST", body: input })).deadline;

export const updateDeadline = async (id: string, input: DeadlineInput): Promise<Deadline> =>
    (await apiRequest<{ deadline: Deadline }>(`/deadlines/${id}`, { method: "PUT", body: input })).deadline;

export const deleteDeadline = (id: string): Promise<void> =>
    apiRequest<void>(`/deadlines/${id}`, { method: "DELETE" });
