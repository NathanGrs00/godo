import React from "react";
import type { Tag } from "../../types";
import "../../styles/Common/Chip.css";

/** A small coloured label showing a tag's name. */
const TagChip: React.FC<{ tag: Tag }> = ({ tag }) => (
    <span className="chip" style={{ borderColor: tag.color, color: tag.color }}>
        {tag.title}
    </span>
);

export default TagChip;
