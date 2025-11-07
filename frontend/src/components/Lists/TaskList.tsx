import React, { useState } from "react";
import type { Task } from "../../hooks/useTasks";
import TaskItem from "./ListItem/TaskItem.tsx";

interface TaskListProps {
    tasks: Task[];
}

const TaskList: React.FC<TaskListProps> = ({ tasks }) => {
    const [taskList, setTaskList] = useState(tasks);

    const toggleTask = (id: string, newStatus: boolean) => {
        setTaskList(prev =>
            prev.map(t => (t.id === id ? { ...t, completed: newStatus } : t))
        );
    };

    if (taskList.length === 0) return <p>No tasks found.</p>;

    return (
        <div>
            {taskList.map(task => (
                <TaskItem
                    key={task.id}
                    {...task}
                    onToggle={status => toggleTask(task.id, status)}
                />
            ))}
        </div>
    );
};

export default TaskList;
