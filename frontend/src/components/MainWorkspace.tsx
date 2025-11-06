import React from "react";
import TitleHeader from "./TitleHeader.tsx";
import TaskIcon from "../assets/icons/NavBar/TaskIcon.tsx";
import "../styles/MainWorkspace.css";

interface MainWorkspaceProps {
}

const MainWorkspace: React.FC<MainWorkspaceProps> = () => {
    return (
        <main className="main-workspace">
            <TitleHeader title={"Tasks"} icon={<TaskIcon/>} />
            <hr className="title-header-line"/>
            /* List of tasks will go here */
        </main>
    );
}

export default MainWorkspace;