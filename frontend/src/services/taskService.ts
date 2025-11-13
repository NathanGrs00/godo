import type {Task} from "../types";

export const getDummyTasks = async (): Promise<Task[]> => {
    const res = await fetch("http://localhost:8080/tasks/dummies");

    if (!res.ok) {
        throw new Error("Failed to fetch tasks");
    }

    const data = await res.json();
    return data.tasks;
};

export const createTask = async (task: Partial<Task>): Promise<Task> => {
    const res = await fetch("http://localhost:8080/tasks", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(task),
    });

    if (!res.ok) {
        throw new Error("Failed to create task");
    }

    return res.json();
}

