import React from "react";
import "../../styles/TitleHeader/TitleHeader.css";

interface TitleHeaderProps {
    icon: React.ReactNode;
    title: string;
}

export const TitleHeader: React.FC<TitleHeaderProps> = ({ icon, title }) => {
    return (
        <div className="title-header">
            <div className="title-header-icon">{icon}</div>
            <h1 className="title-header-text">{title}</h1>
        </div>
    )
}

export default TitleHeader;