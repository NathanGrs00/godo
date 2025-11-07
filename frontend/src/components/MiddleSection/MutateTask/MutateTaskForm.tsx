import React from "react";
import TagSelect from "./TagSelect.tsx";
import "../../../styles/MiddleSection/MutateTask/MutateTaskForm.css"

interface MutateTaskFormProps {
    title?: string;
    description?: string;
    priority?: string;
    //TODO: Add tag and deadline fields
    onClose?: () => void;
}

const MutateTaskForm: React.FC<MutateTaskFormProps> = ({
    title = '',
    description = '',
    priority = '',
    onClose
} ) => {

    const [selectedTags, setSelectedTags] = React.useState<string[]>([]);

    return (
        <form className="mutate-task-form" onSubmit={(e) => e.preventDefault()}>
            <input type="text" id="title" name="title" defaultValue={title} placeholder="Add a title here..." />
            <div className="form-row">
                <label htmlFor="description">Description:</label>
                <input type="text" id="description" name="description" defaultValue={description} />
            </div>
            <div className="form-row">
                <label htmlFor="priority">Priority:</label>
                <select id="priority" name="priority" defaultValue={priority}>
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
                <button type="button" onClick={onClose}>Cancel</button>
                <button>Save</button>
            </div>
        </form>
    );
};

export default MutateTaskForm;