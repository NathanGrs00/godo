import React from "react";
import "../../../styles/Lists/ListItem/ListItem.css";
import type { Deadline } from "../../../hooks/useDeadlines";

interface DeadlineItemProps extends Deadline {
    onToggle?: (completed: boolean) => void; // optional if you want checkbox
}

const DeadlineItem: React.FC<DeadlineItemProps> = ({ title, description, date, passed, tasks }) => {
    return (
        <div className="list-item">
            <h3>{title}</h3>
            {description && <p>{description}</p>}
            <div className="deadline-info">
                <span>{new Date(date).toLocaleString()}</span>
                {passed && <span className="deadline-passed"> (Passed)</span>}
            </div>
            {tasks && tasks.length > 0 && (
                <div>
                    Associated Tasks:
                    <ul>{tasks.map((task) => <li key={task.id || task.title}>{task.title}</li>)}</ul>
                </div>
            )}
        </div>
    );
};

export default DeadlineItem;
