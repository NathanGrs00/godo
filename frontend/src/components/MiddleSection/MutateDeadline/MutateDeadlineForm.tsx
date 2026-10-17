import React, { useState } from "react";
import TagSelect from "../MutateTask/TagSelect.tsx";
import FieldError from "../../Common/FieldError.tsx";
import { useData } from "../../../context/dataContext.ts";
import type { Deadline } from "../../../types";
import { fromLocalInputValue, toLocalInputValue } from "../../../utils/date.ts";
import { sameIds } from "../../../utils/inputs.ts";
import { errorsFromSave, hasErrors, validateDeadline, type FormErrors } from "../../../utils/validation.ts";
import "../../../styles/MiddleSection/MutateForm.css";

interface MutateDeadlineFormProps {
    /** The deadline to edit. Leave out to create a new one. */
    deadline?: Deadline;
    onClose: () => void;
}

interface FormState {
    title: string;
    description: string;
    /** Local date-time as used by <input type="datetime-local">. */
    date: string;
    tag_ids: string[];
}

const initialState = (deadline?: Deadline): FormState => ({
    title: deadline?.title ?? "",
    description: deadline?.description ?? "",
    date: deadline ? toLocalInputValue(deadline.date) : "",
    tag_ids: deadline?.tag_ids ?? [],
});

/** Creates a deadline, or edits one when `deadline` is given. */
const MutateDeadlineForm: React.FC<MutateDeadlineFormProps> = ({ deadline, onClose }) => {
    const { deadlines, createDeadline, updateDeadline } = useData();
    const [initial] = useState(() => initialState(deadline));
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
        form.date !== initial.date ||
        !sameIds(form.tag_ids, initial.tag_ids);
    const canSave = !saving && (!deadline || changed);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const found = validateDeadline(form);
        setErrors(found);
        if (hasErrors(found)) return;

        const input = {
            title: form.title.trim(),
            description: form.description.trim(),
            date: fromLocalInputValue(form.date),
            tag_ids: form.tag_ids,
        };
        setSaving(true);
        try {
            if (deadline) {
                // Normally save against the version this form opened with, so an
                // edit made meanwhile on another device is detected (409). After
                // the user saw that warning, save against the newest copy instead.
                const base = overwrite ? deadlines.find(item => item.id === deadline.id) ?? deadline : deadline;
                await updateDeadline(base, input);
            }
            else await createDeadline(input);
            onClose();
        } catch (err) {
            if ((err as { status?: number }).status === 409) setOverwrite(true);
            setErrors(errorsFromSave(err, ["title", "description", "date"]));
            setSaving(false);
        }
    };

    return (
        <form className="mutate-form" onSubmit={handleSubmit} noValidate aria-label={deadline ? "Edit deadline" : "New deadline"}>
            <input
                type="text"
                className="mutate-form-title"
                aria-label="Title"
                value={form.title}
                onChange={e => set("title", e.target.value)}
                placeholder="Name this deadline..."
                autoFocus
            />
            <FieldError message={errors.title} />

            <div className="form-row">
                <label htmlFor="deadline-date">Due:</label>
                <input
                    id="deadline-date"
                    type="datetime-local"
                    value={form.date}
                    onChange={e => set("date", e.target.value)}
                />
            </div>
            <FieldError message={errors.date} />

            <div className="form-row">
                <label htmlFor="deadline-description">Description:</label>
                <textarea
                    id="deadline-description"
                    value={form.description}
                    onChange={e => set("description", e.target.value)}
                    rows={2}
                />
            </div>
            <FieldError message={errors.description} />

            <div className="form-row">
                <span className="form-label">Tags:</span>
                <TagSelect value={form.tag_ids} onChange={ids => set("tag_ids", ids)} />
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

export default MutateDeadlineForm;
