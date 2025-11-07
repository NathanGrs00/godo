import { useState, useEffect } from "react";
import type {Tag} from "../types";
import * as tagService from "../services/tagService";

export const useTags = () => {
    const [tags, setTags] = useState<Tag[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchTags = async () => {
            setLoading(true);
            setError(null);
            try {
                const data = await tagService.getAllTags();
                setTags(data);
            } catch (err) {
                setError("Failed to load tags");
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchTags();
    }, []);

    return { tags, loading, error };
};
