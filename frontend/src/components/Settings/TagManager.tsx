import React, { useState } from "react";
import { useData } from "../../context/dataContext.ts";
import type { Tag } from "../../types";
import FieldError from "../Common/FieldError.tsx";
import DeleteIcon from "../../assets/icons/DeleteIcon.tsx";
import { errorsFromSave, hasErrors, validateTag, type FormErrors } from "../../utils/validation.ts";
import "../../styles/Pages/Pages.css";

const DEFAULT_COLOR = "#8CFBDE";

/** Create, rename, recolour and delete tags. */
const TagManager: React.FC = () => {
    const { tags, createTag, deleteTag } = useData();
    const [title, setTitle] = useState("");
    const [color, setColor] = useState(DEFAULT_COLOR);
    const [errors, setErrors] = useState<FormErrors>({});

    const handleCreate = async (e: React.FormEvent) => {
        e.preventDefault();
        const found = validateTag({ title, color });
        setErrors(found);
        if (hasErrors(found)) return;
        try {
            await createTag({ title: title.trim(), color });
            setTitle("");
        } catch (err) {
            setErrors(errorsFromSave(err, ["title", "color"]));
        }
    };

    const handleDelete = async (tag: Tag) => {
        if (!window.confirm(`Delete tag "${tag.title}"? It is removed from all tasks and deadlines.`)) return;
        try {
            await deleteTag(tag.id);
        } catch (err) {
            setErrors(errorsFromSave(err, []));
        }
    };

    return (
        <section className="page-section">
            <h2>Tags</h2>
            <form className="tag-form" onSubmit={handleCreate} noValidate aria-label="New tag">
                <input
                    type="text"
                    aria-label="Tag name"
                    placeholder="New tag name"
                    value={title}
                    onChange={e => setTitle(e.target.value)}
                />
                <input type="color" aria-label="Tag colour" value={color} onChange={e => setColor(e.target.value.toUpperCase())} />
                <button type="submit" className="save-button">Add tag</button>
            </form>
            <FieldError message={errors.title ?? errors.color ?? errors.form} />

            {tags.length === 0 && <p className="empty-message">No tags yet.</p>}
            <ul className="tag-list">
                {tags.map(tag => <TagRow key={`${tag.id}-${tag.version}`} tag={tag} onDelete={() => handleDelete(tag)} />)}
            </ul>
        </section>
    );
};

/** One tag that can be renamed or recoloured in place. */
const TagRow: React.FC<{ tag: Tag; onDelete: () => void }> = ({ tag, onDelete }) => {
    const { updateTag } = useData();
    const [title, setTitle] = useState(tag.title);
    const [error, setError] = useState<string | undefined>();

    const save = async (changes: { title?: string; color?: string }) => {
        const next = { title: changes.title ?? tag.title, color: changes.color ?? tag.color };
        if (next.title === tag.title && next.color === tag.color) return;
        const found = validateTag(next);
        if (hasErrors(found)) {
            setError(found.title ?? found.color);
            return;
        }
        setError(undefined);
        try {
            await updateTag(tag, next);
        } catch (err) {
            setError(err instanceof Error ? err.message : "Failed to save");
            setTitle(tag.title);
        }
    };

    return (
        <li className="tag-row">
            <input
                type="color"
                aria-label={`Colour of ${tag.title}`}
                value={tag.color}
                onChange={e => save({ color: e.target.value.toUpperCase() })}
            />
            <input
                type="text"
                aria-label={`Name of ${tag.title}`}
                value={title}
                onChange={e => setTitle(e.target.value)}
                onBlur={() => save({ title: title.trim() })}
                onKeyDown={e => e.key === "Enter" && (e.target as HTMLInputElement).blur()}
                style={{ color: tag.color }}
            />
            <button type="button" className="icon-button" onClick={onDelete} aria-label={`Delete tag ${tag.title}`}>
                <DeleteIcon />
            </button>
            <FieldError message={error} />
        </li>
    );
};

export default TagManager;
