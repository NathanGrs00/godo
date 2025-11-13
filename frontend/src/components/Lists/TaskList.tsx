import React, { useState } from "react";
import type { Task } from "../../hooks/useTasks";
import TaskItem from "./ListItem/TaskItem.tsx";
import * as taskService from "../../services/taskService.ts";

interface TaskListProps {
    tasks: Task[];
}

const TaskList: React.FC<TaskListProps> = ({ tasks }) => {
    const [taskList, setTaskList] = useState(tasks || []);

    const toggleTask = (id: string, newStatus: boolean) => {
        setTaskList(prev =>
            prev.map(t => (t.id === id ? { ...t, completed: newStatus } : t))
        );
    };

    const deleteTask = async (_id: string) => {
        if (!window.confirm("Are you sure you want to delete this task?")) return;

        try {
            await taskService.deleteTask(_id);
            setTaskList(prev => prev.filter(t => t._id !== _id));
        } catch (err: any) {
            console.error("Failed to delete task:", err);
            alert("Failed to delete task: " + err.message);
        }
    }

    if (taskList === undefined || taskList.length === 0) {
        return <p>No tasks available.</p>;
    }

    return (
        <div>
            {taskList.map(task => (
                <TaskItem
                    key={task.id}
                    {...task}
                    onToggle={status => toggleTask(task.id, status)}
                    onDelete={() => deleteTask(task.id.toString())}
                />
            ))}
        </div>
    );
};

export default TaskList;
