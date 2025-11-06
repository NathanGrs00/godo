import React, {type ReactNode} from "react";
import '../../styles/NavBar/SideNavBarItem.css';

interface SideNavBarItemProps {
    icon: ReactNode
}

const SideNavBarItem: React.FC<SideNavBarItemProps> = ({
    icon,
}) => {
    return (
        <div className="side-nav-bar-item">
            {icon}
        </div>
    );
}

export default SideNavBarItem;