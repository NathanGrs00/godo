import React from "react";
import { useTags } from "../../../hooks/useTags.ts";
import OptionPicker from "../../Common/OptionPicker.tsx";

interface TagSelectProps {
    value: string[];
    onChange: (tagIds: string[]) => void;
}

/** Pick tags for a task or deadline. Tags are created on the Settings page. */
const TagSelect: React.FC<TagSelectProps> = ({ value, onChange }) => {
    const { tags } = useTags();
    return (
        <OptionPicker
            label="Tags"
            options={tags.map(tag => ({ id: tag.id, label: tag.title, color: tag.color }))}
            value={value}
            onChange={onChange}
            emptyText="No tags yet – create them in Settings"
        />
    );
};

export default TagSelect;
