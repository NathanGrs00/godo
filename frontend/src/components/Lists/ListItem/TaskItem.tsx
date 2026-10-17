import React from "react";
import "../../../styles/Lists/ListItem/ListItem.css";
import type { Deadline, Tag, Task } from "../../../types";
import DeleteIcon from "../../../assets/icons/DeleteIcon.tsx";
import TagChip from "../../Common/TagChip.tsx";
import { relativeTime } from "../../../utils/date.ts";

interface TaskItemProps {
    task: Task;
    /** The task's tags, already looked up. */
    tags: Tag[];
    /** The task's deadlines, already looked up. */
    deadlines: Deadline[];
    onToggle: (completed: boolean) => void;
    onEdit: () => void;
    onDelete: () => void;
}

/** One row in the task list. Click the text to edit; tick the box to complete. */
const TaskItem: React.FC<TaskItemProps> = ({ task, tags, deadlines, onToggle, onEdit, onDelete }) => {
    const priorityClass = `priority-${task.priority}`;

    return (
        <div className={`list-item ${priorityClass} ${task.completed ? "completed" : ""}`}>
            <input
                type="checkbox"
                checked={task.completed}
                onChange={e => onToggle(e.target.checked)}
                aria-label={`Mark "${task.title}" as ${task.completed ? "not done" : "done"}`}
            />
            <button type="button" className="list-item-body" onClick={onEdit} aria-label={`Edit "${task.title}"`}>
                <h3>{task.title}</h3>
                {task.description && <p>{task.description}</p>}
                {(tags.length > 0 || deadlines.length > 0) && (
                    <div className="list-item-meta">
                        {tags.map(tag => <TagChip key={tag.id} tag={tag} />)}
                        {deadlines.map(d => (
                            <span key={d.id} className={`deadline-ref ${d.passed ? "deadline-passed" : ""}`}>
                                {d.title} · {relativeTime(d.date)}
                            </span>
                        ))}
                    </div>
                )}
            </button>
            <button type="button" className="icon-button" onClick={onDelete} aria-label={`Delete "${task.title}"`}>
                <DeleteIcon />
            </button>
        </div>
    );
};

export default TaskItem;
