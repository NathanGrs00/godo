import React, {type ReactNode} from "react";
import '../../styles/NavBar/SideNavBarItem.css';

interface SideNavBarItemProps {
    icon: ReactNode;
    /** Accessible name and tooltip, e.g. "Tasks". */
    label: string;
    active?: boolean;
    onClick?: () => void;
}

const SideNavBarItem: React.FC<SideNavBarItemProps> = ({ icon, label, active = false, onClick }) => {
    return (
        <button
            type="button"
            className={`side-nav-bar-item ${active ? "active" : ""}`}
            onClick={onClick}
            aria-label={label}
            aria-current={active ? "page" : undefined}
            title={label}
        >
            {icon}
        </button>
    );
}

export default SideNavBarItem;
