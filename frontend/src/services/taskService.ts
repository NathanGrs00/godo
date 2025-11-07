import type {Task} from "../types";

export const getDummyTasks = async (): Promise<Task[]> => {
    const res = await fetch("http://localhost:8080/tasks/dummies");

    if (!res.ok) {
        throw new Error("Failed to fetch tasks");
    }

    const data = await res.json();
    return data.tasks;
};
