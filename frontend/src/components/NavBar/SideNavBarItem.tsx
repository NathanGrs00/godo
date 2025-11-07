import React, {type ReactNode} from "react";
import '../../styles/NavBar/SideNavBarItem.css';

interface SideNavBarItemProps {
    icon: ReactNode;
    onClick?: () => void;
}

const SideNavBarItem: React.FC<SideNavBarItemProps> = ({ icon, onClick }) => {
    return (
        <div className="side-nav-bar-item" onClick={onClick}>
            {icon}
        </div>
    );
}

export default SideNavBarItem;