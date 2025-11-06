import React from "react";
import TitleHeader from "./TitleHeader.tsx";
import TaskIcon from "../../assets/icons/NavBar/TaskIcon.tsx";
import "../../styles/MainWorkspace.css";
import {useTasks} from "../../hooks/useTasks.ts";
import TaskList from "../Lists/TaskList.tsx";

interface MainWorkspaceProps {
}

const MainWorkspace: React.FC<MainWorkspaceProps> = () => {
    const { tasks, loading, error } = useTasks();

    return (
        <main className="main-workspace">
            <TitleHeader title={"Tasks"} icon={<TaskIcon/>} />
            <hr className="title-header-line"/>

            {loading && <p>Loading tasks...</p>}
            {error && <p className="error-message">Error: {error}</p>}
            {!loading && !error && <TaskList tasks={tasks} />}
        </main>
    );
}

export default MainWorkspace;