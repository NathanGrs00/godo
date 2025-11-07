import React, { useState } from "react";
import type { Task } from "../../hooks/useTasks.ts";
import "../../styles/Lists/TaskList.css";
import ListItem from "./ListItem/ListItem.tsx";

interface TaskListProps {
    tasks: Task[];
}

const TaskList: React.FC<TaskListProps> = ({ tasks }) => {
    const [taskList, setTaskList] = useState(tasks);

    const toggleTask = (id: string | undefined, newStatus: boolean) => {
        setTaskList((prev) =>
            prev.map((task) =>
                task.id === id ? { ...task, completed: newStatus } : task
            )
        );
    };

    if (taskList.length === 0) return <p>No tasks found.</p>;

    return (
        <div>
            {taskList.map((task) => (
                <ListItem
                    key={task.id}
                    kind="task"
                    title={task.title}
                    description={task.description}
                    completed={task.completed}
                    priority={task.priority}
                    onToggle={(newStatus) => toggleTask(task.id, newStatus)}
                />
            ))}
        </div>
    );
};

export default TaskList;