import { useState, useEffect } from "react";
import type {Deadline} from "../types";
import * as deadlineService from "../services/deadlineService";

export const useDeadlines = () => {
    const [deadlines, setDeadlines] = useState<Deadline[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchDeadlines = async () => {
            try {
                const data = await deadlineService.getDummyDeadlines();
                setDeadlines(data);
            } catch (err: any) {
                setError(err.message ?? "An unknown error occurred");
            } finally {
                setLoading(false);
            }
        };

        fetchDeadlines();
    }, []);

    return { deadlines, loading, error };
};
