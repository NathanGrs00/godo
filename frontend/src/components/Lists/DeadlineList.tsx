import React, { useState } from "react";
import type { Deadline } from "../../hooks/useDeadlines.ts";
import "../../styles/Lists/DeadlineList.css";
import ListItem from "./ListItem/ListItem.tsx";

interface DeadlineListProps {
    deadlines: Deadline[];
}

const DeadlineList: React.FC<DeadlineListProps> = ({ deadlines }) => {
    const [deadlineList, setDeadlineList] = useState(deadlines);

    const toggleDeadline = (id: string | undefined, newStatus: boolean) => {
        setDeadlineList((prev) =>
            prev.map((deadline) =>
                deadline.id === id ? { ...deadline, completed: newStatus } : deadline
            )
        );
    };

    if (deadlineList.length === 0) return <p>No deadlines found.</p>;

    return (
        <div>
            {deadlineList.map((deadline) => (
                <ListItem
                    key={deadline.id}
                    kind="deadline"
                    title={deadline.title}
                    description={deadline.description}
                    completed={deadline.completed}
                    priority={deadline.priority}
                    onToggle={(newStatus) => toggleDeadline(deadline.id, newStatus)}
                />
            ))}
        </div>
    );
};

export default DeadlineList;