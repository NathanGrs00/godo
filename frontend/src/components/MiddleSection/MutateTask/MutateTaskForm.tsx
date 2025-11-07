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
            <input type="text" id="title" name="title" defaultValue={title} placeholder="Add a title here..." />
            <div className="form-label-row">
                <label htmlFor="description">Description:</label>
                <label htmlFor="priority">Priority:</label>
                <label htmlFor="completed">Completed:</label>
            </div>
            <div className="form-input-row">
                <textarea id="description" name="description" defaultValue={description} />
                <select id="priority" name="priority" defaultValue={priority}>
                    <option value="">Select priority</option>
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                </select>
                <input type="checkbox" id="completed" name="completed" defaultChecked={completed} />
            </div>
        </form>
    );
};

export default MutateTaskForm;