import React from "react";
import { useTags } from "../../../hooks/useTags";

interface TagSelectProps {
    value: string[];
    onChange: (tags: string[]) => void;
}

const TagSelect: React.FC<TagSelectProps> = ({ value, onChange }) => {
    const { tags, loading, error } = useTags();

    if (loading) return <p>Loading tags...</p>;
    if (error) return <p>{error}</p>;

    return (
        <select
            multiple
            value={value}
            onChange={(e) =>
                onChange([...e.target.selectedOptions].map(o => o.value))
            }
        >
            {tags.map(tag => (
                <option key={tag._id} value={tag._id}>
                    {tag.title}
                </option>
            ))}
        </select>
    );
};

export default TagSelect;
