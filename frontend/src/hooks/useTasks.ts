import {useState, useEffect} from "react";

export interface Task {
    id: string;
    title: string;
    description: string;
    priority: string;
    completed: boolean;
}

export function useTasks = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const response = await fetch("http://localhost:8080/tasks/dummies");
                if (!response.ok) throw new Error("Failed to fetch tasks");
                const data = await response.json();
                setTasks(data.tasks);
            } catch (err: any) {
                setError(err.message || "An unknown error occurred");
            } finally {
                setLoading(false);
            }
        };

        fetchTasks();
    }, []);

    return { tasks, loading, error };
};