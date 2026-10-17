import { useData } from "../context/dataContext.ts";

/** The user's tasks, kept in sync across devices. */
export const useTasks = () => {
    const { tasks, loading, error } = useData();
    return { tasks, loading, error };
};
