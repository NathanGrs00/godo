import { useState, useEffect } from "react";
import type {Task} from "../types";
import * as taskService from "../services/taskService";

export const useTasks = () => {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchTasks = async () => {
            try {
                const data = await taskService.getTasks();
                setTasks(data);
            } catch (err: any) {
                setError(err.message ?? "An unknown error occurred");
            } finally {
                setLoading(false);
            }
        };

        fetchTasks();
    }, []);

    return { tasks, loading, error };
};
