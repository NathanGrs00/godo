import React from "react";
import type {Task} from "../../hooks/useTasks.ts";
import "../../styles/Lists/TaskList.css";

interface TaskListProps {
    tasks: Task[];
}

const TaskList: React.FC<TaskListProps> = ({ tasks }) => {
    if (tasks.length === 0) return <p>No tasks found.</p>;

    return (
        <ul>
            {tasks.map((task) => (
                <li key={task.id} className="task-item">
                    <h3>{task.title}</h3>
                    <p>{task.description}</p>
                    <p>Priority: {task.priority}</p>
                    <p>Status: {task.completed ? "Done" : "Pending"}</p>
                </li>
            ))}
        </ul>
    );
};

export default TaskList;
