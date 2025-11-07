import {useState, useEffect} from "react";

export interface Deadline {
    id: string;
    title: string;
    description: string;
    priority: string;
    completed: boolean;
}

export const useDeadlines = () => {
    const [deadlines, setDeadlines] = useState<Deadline[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchDeadlines = async () => {
            try {
                const response = await fetch("http://localhost:8080/deadlines/dummies");
                if (!response.ok) throw new Error("Failed to fetch tasks");
                const data = await response.json();
                setDeadlines(data.deadlines);
            } catch (err: any) {
                setError(err.message || "An unknown error occurred");
            } finally {
                setLoading(false);
            }
        };

        fetchDeadlines();
    }, []);

    return { deadlines, loading, error };
};