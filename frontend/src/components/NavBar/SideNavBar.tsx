import React, {type ReactNode} from "react";
import '../../styles/NavBar/SideNavBar.css';

interface SideNavBarProps {
    children?: ReactNode;
}

const SideNavBar: React.FC<SideNavBarProps> = ({
    children
}) => {
    return (
        <nav className={`side-nav-bar`}>
            {children}
        </nav>
    );
}

export default SideNavBar;
