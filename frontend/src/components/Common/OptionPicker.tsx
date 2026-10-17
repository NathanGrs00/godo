import React from "react";
import "../../styles/Common/Chip.css";

export interface PickerOption {
    id: string;
    label: string;
    /** Optional accent colour, e.g. a tag's colour. */
    color?: string;
}

interface OptionPickerProps {
    label: string;
    options: PickerOption[];
    value: string[];
    onChange: (value: string[]) => void;
    /** Shown when there are no options. */
    emptyText: string;
}

/** Lets the user pick any number of options by toggling chips. */
const OptionPicker: React.FC<OptionPickerProps> = ({ label, options, value, onChange, emptyText }) => {
    const toggle = (id: string) =>
        onChange(value.includes(id) ? value.filter(v => v !== id) : [...value, id]);

    return (
        <div className="option-picker" role="group" aria-label={label}>
            {options.length === 0 && <span className="option-picker-empty">{emptyText}</span>}
            {options.map(option => {
                const selected = value.includes(option.id);
                return (
                    <button
                        key={option.id}
                        type="button"
                        className={`chip chip-toggle ${selected ? "selected" : ""}`}
                        aria-pressed={selected}
                        style={option.color ? { borderColor: option.color } : undefined}
                        onClick={() => toggle(option.id)}
                    >
                        {option.label}
                    </button>
                );
            })}
        </div>
    );
};

export default OptionPicker;
