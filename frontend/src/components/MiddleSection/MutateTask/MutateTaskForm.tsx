import React from "react";
import TagSelect from "./TagSelect.tsx";
import * as taskService from "../../../services/taskService.ts";
import "../../../styles/MiddleSection/MutateTask/MutateTaskForm.css"

interface MutateTaskFormProps {
    title?: string;
    description?: string;
    priority?: string;
    //TODO: Add tag and deadline fields
    onClose?: () => void;
}

const MutateTaskForm: React.FC<MutateTaskFormProps> = ({title = "",
                                                           description = "",
                                                           priority = "",
                                                           onClose,
                                                       }) => {
    const [selectedTags, setSelectedTags] = React.useState<string[]>([]);
    const [formTitle, setFormTitle] = React.useState(title);
    const [formDescription, setFormDescription] = React.useState(description);
    const [formPriority, setFormPriority] = React.useState(priority);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            await taskService.createTask({
                title: formTitle,
                description: formDescription,
                priority: formPriority,
                tags: selectedTags, // if you later add tags in Task type
            });

            // Optionally close the form or refresh task list
            onClose?.();
        } catch (err: any) {
            console.error("Failed to create task:", err.message);
        }
    };

    return (
        <form className="mutate-task-form" onSubmit={handleSubmit}>
            <input
                type="text"
                id="title"
                name="title"
                value={formTitle}
                onChange={(e) => setFormTitle(e.target.value)}
                placeholder="Add a title here..."
            />
            <div className="form-row">
                <label htmlFor="description">Description:</label>
                <input
                    type="text"
                    id="description"
                    name="description"
                    value={formDescription}
                    onChange={(e) => setFormDescription(e.target.value)}
                />
            </div>
            <div className="form-row">
                <label htmlFor="priority">Priority:</label>
                <select
                    id="priority"
                    name="priority"
                    value={formPriority}
                    onChange={(e) => setFormPriority(e.target.value)}
                >
                    <option value="">Select priority</option>
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                </select>
            </div>
            <div className="form-row">
                <TagSelect value={selectedTags} onChange={setSelectedTags} />
            </div>
            <div className="form-row">
                <button type="button" onClick={onClose}>
                    Cancel
                </button>
                <button type="submit">Save</button>
            </div>
        </form>
    );
};

export default MutateTaskForm;