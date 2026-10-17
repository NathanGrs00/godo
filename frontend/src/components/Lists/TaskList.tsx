import React, { useState } from "react";
import type { Task } from "../../types";
import TaskItem from "./ListItem/TaskItem.tsx";
import { useData } from "../../context/dataContext.ts";
import "../../styles/Lists/TaskList.css";

interface TaskListProps {
    tasks: Task[];
    onEdit: (task: Task) => void;
}

/** Shows tasks, open ones first, and handles completing and deleting them. */
const TaskList: React.FC<TaskListProps> = ({ tasks, onEdit }) => {
    const { tags, deadlines, updateTask, deleteTask } = useData();
    const [error, setError] = useState<string | null>(null);

    const run = async (action: () => Promise<unknown>) => {
        setError(null);
        try {
            await action();
        } catch (err) {
            setError(err instanceof Error ? err.message : "Something went wrong");
        }
    };

    const handleDelete = (task: Task) => {
        if (!window.confirm(`Delete "${task.title}"?`)) return;
        run(() => deleteTask(task.id));
    };

    if (tasks.length === 0) {
        return <p className="empty-message">No tasks yet. Add one below.</p>;
    }

    // Open tasks first; within each group keep the creation order
    const sorted = [...tasks].sort((a, b) => Number(a.completed) - Number(b.completed));

    return (
        <div className="task-list">
            {error && <p className="error-message" role="alert">{error}</p>}
            {sorted.map(task => (
                <TaskItem
                    key={task.id}
                    task={task}
                    tags={tags.filter(tag => task.tag_ids.includes(tag.id))}
                    deadlines={deadlines.filter(d => task.deadline_ids.includes(d.id))}
                    onToggle={completed => run(() => updateTask(task, { completed }))}
                    onEdit={() => onEdit(task)}
                    onDelete={() => handleDelete(task)}
                />
            ))}
        </div>
    );
};

export default TaskList;
