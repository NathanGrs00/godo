import React, { useState } from "react";
import type { Deadline } from "../../types";
import DeadlineItem from "./ListItem/DeadlineItem.tsx";
import { useData } from "../../context/dataContext.ts";
import { byDate } from "../../utils/date.ts";
import "../../styles/Lists/DeadlineList.css";

interface DeadlineListProps {
    deadlines: Deadline[];
    onEdit: (deadline: Deadline) => void;
}

/** Shows deadlines soonest first (passed ones last) with their linked tasks. */
const DeadlineList: React.FC<DeadlineListProps> = ({ deadlines, onEdit }) => {
    const { tasks, tags, deleteDeadline } = useData();
    const [error, setError] = useState<string | null>(null);

    const handleDelete = async (deadline: Deadline) => {
        if (!window.confirm(`Delete "${deadline.title}"? Linked tasks are kept.`)) return;
        setError(null);
        try {
            await deleteDeadline(deadline.id);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Something went wrong");
        }
    };

    if (deadlines.length === 0) return <p className="empty-message">No deadlines yet. Add one below.</p>;

    const upcoming = deadlines.filter(d => !d.passed).sort(byDate);
    const passed = deadlines.filter(d => d.passed).sort(byDate).reverse();

    return (
        <div className="deadline-list">
            {error && <p className="error-message" role="alert">{error}</p>}
            {[...upcoming, ...passed].map(deadline => (
                <DeadlineItem
                    key={deadline.id}
                    deadline={deadline}
                    tasks={tasks.filter(task => task.deadline_ids.includes(deadline.id))}
                    tags={tags.filter(tag => deadline.tag_ids.includes(tag.id))}
                    onEdit={() => onEdit(deadline)}
                    onDelete={() => handleDelete(deadline)}
                />
            ))}
        </div>
    );
};

export default DeadlineList;
