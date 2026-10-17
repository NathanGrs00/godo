import type { Deadline, Tag, Task } from "../types";

/** Lower-cases and trims text for case-insensitive matching. */
const normalize = (text: string) => text.toLowerCase().trim();

/** Whether any of the texts contains the query. */
const matches = (query: string, ...texts: string[]) =>
    texts.some(text => normalize(text).includes(query));

/** Names of the tags with the given IDs. */
const tagNames = (tags: Tag[], ids: string[]) =>
    tags.filter(tag => ids.includes(tag.id)).map(tag => tag.title);

export interface SearchResults {
    tasks: Task[];
    deadlines: Deadline[];
}

/**
 * Finds tasks and deadlines whose title, description or tag names contain
 * the query. An empty query finds nothing.
 */
export const search = (query: string, data: { tasks: Task[]; deadlines: Deadline[]; tags: Tag[] }): SearchResults => {
    const q = normalize(query);
    if (!q) return { tasks: [], deadlines: [] };
    return {
        tasks: data.tasks.filter(t => matches(q, t.title, t.description, ...tagNames(data.tags, t.tag_ids))),
        deadlines: data.deadlines.filter(d => matches(q, d.title, d.description, ...tagNames(data.tags, d.tag_ids))),
    };
};
