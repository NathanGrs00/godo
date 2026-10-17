import { useData } from "../context/dataContext.ts";

/** The user's tags, kept in sync across devices. */
export const useTags = () => {
    const { tags, loading, error } = useData();
    return { tags, loading, error };
};
