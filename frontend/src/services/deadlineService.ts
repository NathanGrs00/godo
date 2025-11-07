import type {Deadline} from "../types";

export const getDummyDeadlines = async (): Promise<Deadline[]> => {
    const res = await fetch("http://localhost:8080/deadlines/dummies");
    if (!res.ok) {
        throw new Error("Failed to fetch deadlines");
    }

    const data = await res.json();

    return (data.deadlines || []).map((d: any) => ({
        id: d.ID || d.id || "",
        title: d.Title || d.title,
        description: d.Description || d.description,
        date: new Date(d.Date || d.date).toISOString(),
        passed: d.Passed ?? false,
        tasks: d.tasks || [],
    }));
};
