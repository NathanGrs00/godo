import React from 'react';
import '../../../styles/Lists/ListItem/ListItem.css';

export interface Deadline {
    id?: string;
    date: string;
}

export type TaskRef = {
    id?: string;
    title: string;
};

interface BaseListItem {
    title: string;
    description?: string;
    priority?: string;
    completed?: boolean;
    tags?: string[];
}

export type ListItemProps =
    | (BaseListItem & { kind: 'task'; deadlines?: Deadline[] })
    | (BaseListItem & { kind: 'deadline'; tasks?: TaskRef[] });

export const ListItem: React.FC<ListItemProps & { onToggle?: (completed: boolean) => void }> = (props) => {
    const { completed = false, onToggle, priority } = props;

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onToggle?.(e.target.checked);
    };

    const priorityClass = priority ? `priority-${priority.toLowerCase()}` : "";

    return (
        <div className={`list-item ${priorityClass}`}>
            <input
                type="checkbox"
                checked={completed}
                onChange={handleChange}
            />
            <h3>{props.title}</h3>

            {props.description && <p>{props.description}</p>}

            {props.tags && props.tags.length > 0 && (
                <div>
                    Tags:
                    <ul>
                        {props.tags.map((tag, index) => (
                            <li key={index}>{tag}</li>
                        ))}
                    </ul>
                </div>
            )}

            {props.kind === 'task' && props.deadlines && props.deadlines.length > 0 && (
                <div>
                    Deadlines:
                    <ul>
                        {props.deadlines.map((deadline) => (
                            <li key={deadline.id || deadline.date}>{deadline.date}</li>
                        ))}
                    </ul>
                </div>
            )}

            {props.kind === 'deadline' && props.tasks && props.tasks.length > 0 && (
                <div>
                    Associated Tasks:
                    <ul>
                        {props.tasks.map((task) => (
                            <li key={task.id || task.title}>{task.title}</li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default ListItem;