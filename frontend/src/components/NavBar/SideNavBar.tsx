import React, {type ReactNode} from "react";
import '../../styles/NavBar/SideNavBar.css';

interface SideNavBarProps {
    width?: string;
    bgColor?: string;
    className?: string;
    children?: ReactNode;
}

const SideNavBar: React.FC<SideNavBarProps> = ({
    width = 'clamp(75px, 7%, 200px)',
    bgColor = '#232830',
    className = '',
    children
}) => {
    return (
        <nav
            className={`side-nav-bar ${className}`}
            style={{ width: width, backgroundColor: bgColor }}
        >
            {children}
        </nav>
    );
}

export default SideNavBar;
