import { apiRequest } from "../api/client.ts";
import type { Tag, TagInput } from "../types";

export const getTags = async (): Promise<Tag[]> =>
    (await apiRequest<{ tags: Tag[] }>("/tags")).tags;

export const createTag = async (input: TagInput): Promise<Tag> =>
    (await apiRequest<{ tag: Tag }>("/tags", { method: "POST", body: input })).tag;

export const updateTag = async (id: string, input: TagInput): Promise<Tag> =>
    (await apiRequest<{ tag: Tag }>(`/tags/${id}`, { method: "PUT", body: input })).tag;

export const deleteTag = (id: string): Promise<void> =>
    apiRequest<void>(`/tags/${id}`, { method: "DELETE" });
