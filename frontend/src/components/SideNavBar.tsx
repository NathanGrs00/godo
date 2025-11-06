import React, {type ReactNode} from "react";
import '../styles/SideNavBar.css';

interface SideNavBarProps {
    links: { label: string; href: string }[];
    width?: string;
    bgColor?: string;
    className?: string;
    children?: ReactNode;
}

const SideNavBar: React.FC<SideNavBarProps> = ({
    links,
    width = '10%',
    bgColor = '#232830',
    className = '',
    children
}) => {
    return (
        <nav
            className={`side-nav-bar ${className}`}
            style={{ width: width, backgroundColor: bgColor }}
        >
            <ul>
                {links.map((link, index) => (
                    <li key={index}>
                        <a href={link.href}>
                            {link.label}
                        </a>
                    </li>
                ))}
            </ul>
            {children}
        </nav>
    );
}

export default SideNavBar;
