import type {Tag} from "../types";

export const getAllTags = async (): Promise<Tag[]> => {
    const res = await fetch("http://localhost:8080/tags/dummies");

    if (!res.ok) {
        throw new Error("Failed to fetch tags");
    }

    const data = await res.json();
    return data.tags;
};
