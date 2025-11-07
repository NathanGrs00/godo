import React from "react";
import "../../../styles/Lists/ListItem/ListItem.css";
import type { Task } from "../../../hooks/useTasks";

interface TaskItemProps extends Task {
    onToggle?: (completed: boolean) => void;
}

const TaskItem: React.FC<TaskItemProps> = ({
                                               title,
                                               description,
                                               priority,
                                               completed = false,
                                               tags,
                                               onToggle,
                                           }) => {
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) =>
        onToggle?.(e.target.checked);

    const priorityClass = priority ? `priority-${priority.toLowerCase()}` : "";

    return (
        <div className={`list-item ${priorityClass}`}>
            <input type="checkbox" checked={completed} onChange={handleChange} />
            <h3>{title}</h3>
            {description && <p>{description}</p>}
            {tags && tags.length > 0 && (
                <div>
                    Tags:
                    <ul>{tags.map((tag, i) => <li key={i}>{tag}</li>)}</ul>
                </div>
            )}
        </div>
    );
};

export default TaskItem;
