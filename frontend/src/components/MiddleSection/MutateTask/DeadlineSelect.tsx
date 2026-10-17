import React from "react";
import { useDeadlines } from "../../../hooks/useDeadlines.ts";
import OptionPicker from "../../Common/OptionPicker.tsx";
import { byDate } from "../../../utils/date.ts";

interface DeadlineSelectProps {
    value: string[];
    onChange: (deadlineIds: string[]) => void;
}

/** Pick which deadlines a task belongs to. */
const DeadlineSelect: React.FC<DeadlineSelectProps> = ({ value, onChange }) => {
    const { deadlines } = useDeadlines();
    return (
        <OptionPicker
            label="Deadlines"
            options={[...deadlines].sort(byDate).map(d => ({
                id: d.id,
                label: `${d.title} (${new Date(d.date).toLocaleDateString()})`,
            }))}
            value={value}
            onChange={onChange}
            emptyText="No deadlines yet"
        />
    );
};

export default DeadlineSelect;
