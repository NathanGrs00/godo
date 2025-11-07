import React from 'react';
import AddIcon from "../../assets/icons/AddIcon.tsx";
import "../../styles/MiddleSection/AddButton.css";

interface AddButtonProps {
    text: string;
    onClick: () => void;
}

const AddButton: React.FC<AddButtonProps> = ({ text, onClick }) => {
    return (
        <button className="add-button"
            onClick={onClick}
        >
            <AddIcon className={"add-icon"}/>
            {text}
        </button>
    );
}

export default AddButton;