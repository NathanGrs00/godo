import React, { useState } from "react";
import TagSelect from "./TagSelect.tsx";
import DeadlineSelect from "./DeadlineSelect.tsx";
import FieldError from "../../Common/FieldError.tsx";
import { useData } from "../../../context/dataContext.ts";
import { PRIORITIES, type Task, type TaskInput } from "../../../types";
import { sameIds } from "../../../utils/inputs.ts";
import { errorsFromSave, hasErrors, validateTask, type FormErrors } from "../../../utils/validation.ts";
import "../../../styles/MiddleSection/MutateForm.css";

interface MutateTaskFormProps {
    /** The task to edit. Leave out to create a new task. */
    task?: Task;
    /** Pre-selected deadline for a new task. */
    deadlineId?: string;
    onClose: () => void;
}

type FormState = Omit<TaskInput, "priority" | "version"> & { priority: string };

const initialState = (task?: Task, deadlineId?: string): FormState => ({
    title: task?.title ?? "",
    description: task?.description ?? "",
    priority: task?.priority ?? "",
    completed: task?.completed ?? false,
    tag_ids: task?.tag_ids ?? [],
    deadline_ids: task?.deadline_ids ?? (deadlineId ? [deadlineId] : []),
});

/** Creates a task, or edits one when `task` is given. */
const MutateTaskForm: React.FC<MutateTaskFormProps> = ({ task, deadlineId, onClose }) => {
    const { tasks, createTask, updateTask } = useData();
    const [initial] = useState(() => initialState(task, deadlineId));
    const [form, setForm] = useState<FormState>(initial);
    const [errors, setErrors] = useState<FormErrors>({});
    const [saving, setSaving] = useState(false);
    // Set after a conflict: the user has been warned, so the next save may overwrite
    const [overwrite, setOverwrite] = useState(false);

    const set = <K extends keyof FormState>(key: K, value: FormState[K]) =>
        setForm(prev => ({ ...prev, [key]: value }));

    const changed =
        form.title !== initial.title ||
        form.description !== initial.description ||
        form.priority !== initial.priority ||
        !sameIds(form.tag_ids, initial.tag_ids) ||
        !sameIds(form.deadline_ids, initial.deadline_ids);
    // A new task can always be submitted (to show what's missing); an edit only when something changed
    const canSave = !saving && (!task || changed);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const found = validateTask(form);
        setErrors(found);
        if (hasErrors(found)) return;

        const input = { ...form, title: form.title.trim(), description: form.description.trim() } as TaskInput;
        setSaving(true);
        try {
            if (task) {
                // Normally save against the version this form opened with, so an
                // edit made meanwhile on another device is detected (409). After
                // the user saw that warning, save against the newest copy instead.
                const base = overwrite ? tasks.find(item => item.id === task.id) ?? task : task;
                await updateTask(base, input);
            }
            else await createTask(input);
            onClose();
        } catch (err) {
            if ((err as { status?: number }).status === 409) setOverwrite(true);
            setErrors(errorsFromSave(err, ["title", "description", "priority"]));
            setSaving(false);
        }
    };

    return (
        <form className="mutate-form" onSubmit={handleSubmit} noValidate aria-label={task ? "Edit task" : "New task"}>
            <input
                type="text"
                className="mutate-form-title"
                aria-label="Title"
                value={form.title}
                onChange={e => set("title", e.target.value)}
                placeholder="Add a title here..."
                autoFocus
            />
            <FieldError message={errors.title} />

            <div className="form-row">
                <label htmlFor="task-description">Description:</label>
                <textarea
                    id="task-description"
                    value={form.description}
                    onChange={e => set("description", e.target.value)}
                    rows={2}
                />
            </div>
            <FieldError message={errors.description} />

            <div className="form-row">
                <label htmlFor="task-priority">Priority:</label>
                <select id="task-priority" value={form.priority} onChange={e => set("priority", e.target.value)}>
                    <option value="">Select priority</option>
                    {PRIORITIES.map(p => (
                        <option key={p} value={p}>{p[0].toUpperCase() + p.slice(1)}</option>
                    ))}
                </select>
            </div>
            <FieldError message={errors.priority} />

            <div className="form-row">
                <span className="form-label">Tags:</span>
                <TagSelect value={form.tag_ids} onChange={ids => set("tag_ids", ids)} />
            </div>

            <div className="form-row">
                <span className="form-label">Deadlines:</span>
                <DeadlineSelect value={form.deadline_ids} onChange={ids => set("deadline_ids", ids)} />
            </div>

            <FieldError message={errors.form} />

            <div className="form-row form-actions">
                <button type="button" className="cancel-button" onClick={onClose}>Cancel</button>
                <button type="submit" className="save-button" disabled={!canSave}>
                    {saving ? "Saving…" : "Save"}
                </button>
            </div>
        </form>
    );
};

export default MutateTaskForm;
