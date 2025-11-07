import React from "react";
import type { Deadline } from "../../hooks/useDeadlines";
import DeadlineItem from "./ListItem/DeadlineItem.tsx";

interface DeadlineListProps {
    deadlines: Deadline[];
}

const DeadlineList: React.FC<DeadlineListProps> = ({ deadlines }) => {
    if (deadlines.length === 0) return <p>No deadlines found.</p>;

    return (
        <div>
            {deadlines.map(deadline => (
                <DeadlineItem key={deadline.id} {...deadline} />
            ))}
        </div>
    );
};

export default DeadlineList;
