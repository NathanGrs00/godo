import React from "react";
import TitleHeader from "./TitleHeader.tsx";
import TaskIcon from "../../assets/icons/NavBar/TaskIcon.tsx";
import "../../styles/MiddleSection/MainWorkspace.css";
import {useTasks} from "../../hooks/useTasks.ts";
import TaskList from "../Lists/TaskList.tsx";
import DeadlineIcon from "../../assets/icons/NavBar/DeadlineIcon.tsx";
import DeadlineList from "../Lists/DeadlineList.tsx";
import {useDeadlines} from "../../hooks/useDeadlines.ts";
import AddButton from "./AddButton.tsx";
import MutateTaskForm from "./MutateTask/MutateTaskForm.tsx";

interface MainWorkspaceProps {
    route?: string;
}

const MainWorkspace: React.FC<MainWorkspaceProps> = ({route}) => {
    if (route === "tasks") {
        return <TaskWorkspace />;
    } else if (route === "deadlines") {
        return <DeadlineWorkspace />;
    } else {
        return <div className="main-workspace"><p>Please select a valid route.</p></div>;
    }
}

// TODO: onClick should open a form to add a new task
const TaskWorkspace: React.FC = () => {
    const { tasks, loading, error } = useTasks();
    const [showForm, setShowForm] = React.useState(false);

    React.useEffect(() => {
        const handleEscape = (event: KeyboardEvent) => {
            if (event.key === 'Escape') setShowForm(false);
        }
        window.addEventListener('keydown', handleEscape);
        return () => window.removeEventListener('keydown', handleEscape);
    })

    return (
        <main className="main-workspace">
            <TitleHeader title={"Tasks"} icon={<TaskIcon/>} />
            <hr className="title-header-line"/>

            {loading && <p>Loading tasks...</p>}
            {error && <p className="error-message">Error: {error}</p>}
            {!loading && !error && <TaskList tasks={tasks} />}

            <AddButton text={"Add Task"} onClick={() => setShowForm(true)} />
            {showForm && (
                <>
                    <div className="mutate-task-backdrop" onClick={() => setShowForm(false)}/>
                    <MutateTaskForm onClose={() => setShowForm(false)}/>
                </>
            )}
        </main>
    );
}

const DeadlineWorkspace: React.FC = () => {
    const { deadlines, loading, error } = useDeadlines();

    return (
        <main className="main-workspace">
            <TitleHeader title={"Deadlines"} icon={<DeadlineIcon/>} />
            <hr className="title-header-line"/>

            {loading && <p>Loading deadlines...</p>}
            {error && <p className="error-message">Error: {error}</p>}
            {!loading && !error && <DeadlineList deadlines={deadlines} />}
            <AddButton text={"Add Deadline"} />
        </main>
    );
}
export default MainWorkspace;