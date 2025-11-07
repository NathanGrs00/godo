import {useState, useEffect} from "react";

export interface Deadline {
    id: string;
    title: string;
    description: string;
    date: string;
    passed: boolean;
    tasks: {
        id?: string;
        title: string;
    }[];
}

export const useDeadlines = () => {
    const [deadlines, setDeadlines] = useState<Deadline[]>([]); // ✅ empty array
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchDeadlines = async () => {
            try {
                const response = await fetch("http://localhost:8080/deadlines/dummies");
                if (!response.ok) throw new Error("Failed to fetch deadlines");
                const data = await response.json();

                const deadlinesWithDate = (data.deadlines || []).map((d: any) => ({
                    id: d.ID || d.id || "",
                    title: d.Title || d.title,
                    description: d.Description || d.description,
                    date: new Date(d.Date || d.date).toISOString(),
                    passed: d.Passed ?? false,
                    tasks: d.tasks || [], // ensure tasks is always an array
                }));

                setDeadlines(deadlinesWithDate);
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
