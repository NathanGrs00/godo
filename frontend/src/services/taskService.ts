import type {Task} from "../types";

export const getTasks = async (): Promise<Task[]> => {
    const res = await fetch("http://localhost:8080/tasks/" , {
        method: "GET",
        headers: { "Content-Type": "application/json" },
    });

    if (!res.ok) {
        throw new Error("Failed to fetch tasks");
    }

    const data = await res.json();
    return data.tasks;
};

export const createTask = async (task: Partial<Task>): Promise<Task> => {
    const res = await fetch("http://localhost:8080/tasks/", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(task),
    });

    if (!res.ok) {
        throw new Error("Failed to create task");
    }

    return res.json();
}

export const deleteTask = async (id: string): Promise<void> => {
    const res = await fetch(`http://localhost:8080/tasks/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
    });

    if (!res.ok) {
        throw new Error("Failed to delete task");
    }

    return;
}

