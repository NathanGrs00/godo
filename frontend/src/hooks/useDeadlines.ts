import { useData } from "../context/dataContext.ts";

/** The user's deadlines, kept in sync across devices. */
export const useDeadlines = () => {
    const { deadlines, loading, error } = useData();
    return { deadlines, loading, error };
};
