import React from "react";
import TitleHeader from "../MiddleSection/TitleHeader.tsx";
import SyncStatus from "./SyncStatus.tsx";
import "../../styles/MiddleSection/MainWorkspace.css";

interface PageLayoutProps {
    title: string;
    icon: React.ReactNode;
    children: React.ReactNode;
}

/** The shared frame of every page: title, divider, sync status and content. */
const PageLayout: React.FC<PageLayoutProps> = ({ title, icon, children }) => (
    <main className="main-workspace">
        <div className="page-header">
            <TitleHeader title={title} icon={icon} />
            <SyncStatus />
        </div>
        <hr className="title-header-line" />
        {children}
    </main>
);

export default PageLayout;
