import React from "react";
import "../../../styles/Lists/ListItem/ListItem.css";
import type { Deadline, Tag, Task } from "../../../types";
import DeleteIcon from "../../../assets/icons/DeleteIcon.tsx";
import TagChip from "../../Common/TagChip.tsx";
import { relativeTime } from "../../../utils/date.ts";

interface DeadlineItemProps {
    deadline: Deadline;
    /** Tasks linked to this deadline. */
    tasks: Task[];
    tags: Tag[];
    onEdit: () => void;
    onDelete: () => void;
}

/** One row in the deadline list, with its date, progress and linked tasks. */
const DeadlineItem: React.FC<DeadlineItemProps> = ({ deadline, tasks, tags, onEdit, onDelete }) => {
    const done = tasks.filter(task => task.completed).length;

    return (
        <div className={`list-item deadline-item ${deadline.passed ? "deadline-item-passed" : ""}`}>
            <button type="button" className="list-item-body" onClick={onEdit} aria-label={`Edit "${deadline.title}"`}>
                <h3>{deadline.title}</h3>
                <div className="deadline-info">
                    <span>{new Date(deadline.date).toLocaleString()}</span>
                    <span className={deadline.passed ? "deadline-passed" : "deadline-upcoming"}>
                        {deadline.passed ? "Passed" : relativeTime(deadline.date)}
                    </span>
                    {tasks.length > 0 && <span>{done}/{tasks.length} tasks done</span>}
                </div>
                {deadline.description && <p>{deadline.description}</p>}
                {(tags.length > 0 || tasks.length > 0) && (
                    <div className="list-item-meta">
                        {tags.map(tag => <TagChip key={tag.id} tag={tag} />)}
                        {tasks.map(task => (
                            <span key={task.id} className={`task-ref ${task.completed ? "completed" : ""}`}>{task.title}</span>
                        ))}
                    </div>
                )}
            </button>
            <button type="button" className="icon-button" onClick={onDelete} aria-label={`Delete "${deadline.title}"`}>
                <DeleteIcon />
            </button>
        </div>
    );
};

export default DeadlineItem;
