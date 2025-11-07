import React from "react";
import "../../../styles/MiddleSection/MutateTask/MutateTaskForm.css"

interface MutateTaskFormProps {
    title?: string;
    description?: string;
    priority?: string;
    completed?: boolean;
    //TODO: Add tag and deadline fields
}

const MutateTaskForm: React.FC<MutateTaskFormProps> = ({
    title = '',
    description = '',
    priority = '',
    completed = false,
} ) => {
    return (
        <form className="mutate-task-form">
            <div>
                <label htmlFor="title">Title:</label>
                <input type="text" id="title" name="title" defaultValue={title} />
            </div>
            <div>
                <label htmlFor="description">Description:</label>
                <textarea id="description" name="description" defaultValue={description} />
            </div>
            <div>
                <label htmlFor="priority">Priority:</label>
                <select id="priority" name="priority" defaultValue={priority}>
                    <option value="">Select priority</option>
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                </select>
            </div>
            <div>
                <label htmlFor="completed">Completed:</label>
                <input type="checkbox" id="completed" name="completed" defaultChecked={completed} />
            </div>
        </form>
    );
};

export default MutateTaskForm;