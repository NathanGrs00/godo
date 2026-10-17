import React, { useState } from "react";
import TaskIcon from "../../assets/icons/NavBar/TaskIcon.tsx";
import DeadlineIcon from "../../assets/icons/NavBar/DeadlineIcon.tsx";
import { useTasks } from "../../hooks/useTasks.ts";
import { useDeadlines } from "../../hooks/useDeadlines.ts";
import TaskList from "../Lists/TaskList.tsx";
import DeadlineList from "../Lists/DeadlineList.tsx";
import AddButton from "./AddButton.tsx";
import MutateTaskForm from "./MutateTask/MutateTaskForm.tsx";
import MutateDeadlineForm from "./MutateDeadline/MutateDeadlineForm.tsx";
import Modal from "../Common/Modal.tsx";
import PageLayout from "../Common/PageLayout.tsx";
import type { Deadline, Task } from "../../types";

interface MainWorkspaceProps {
    route: "tasks" | "deadlines";
}

const MainWorkspace: React.FC<MainWorkspaceProps> = ({ route }) =>
    route === "tasks" ? <TaskWorkspace /> : <DeadlineWorkspace />;

/**
 * Which form is open: closed (null), a new item ("new") or an existing item.
 */
type FormTarget<T> = null | "new" | T;

const TaskWorkspace: React.FC = () => {
    const { tasks, loading, error } = useTasks();
    const [target, setTarget] = useState<FormTarget<Task>>(null);
    const close = () => setTarget(null);

    return (
        <PageLayout title="Tasks" icon={<TaskIcon />}>
            {loading && <p className="empty-message">Loading tasks...</p>}
            {error && <p className="error-message">Error: {error}</p>}
            {!loading && !error && <TaskList tasks={tasks} onEdit={setTarget} />}

            <AddButton text="Add Task" onClick={() => setTarget("new")} />
            {target && (
                <Modal onClose={close}>
                    <MutateTaskForm
                        key={target === "new" ? "new" : target.id}
                        task={target === "new" ? undefined : target}
                        onClose={close}
                    />
                </Modal>
            )}
        </PageLayout>
    );
};

const DeadlineWorkspace: React.FC = () => {
    const { deadlines, loading, error } = useDeadlines();
    const [target, setTarget] = useState<FormTarget<Deadline>>(null);
    const close = () => setTarget(null);

    return (
        <PageLayout title="Deadlines" icon={<DeadlineIcon />}>
            {loading && <p className="empty-message">Loading deadlines...</p>}
            {error && <p className="error-message">Error: {error}</p>}
            {!loading && !error && <DeadlineList deadlines={deadlines} onEdit={setTarget} />}

            <AddButton text="Add Deadline" onClick={() => setTarget("new")} />
            {target && (
                <Modal onClose={close}>
                    <MutateDeadlineForm
                        key={target === "new" ? "new" : target.id}
                        deadline={target === "new" ? undefined : target}
                        onClose={close}
                    />
                </Modal>
            )}
        </PageLayout>
    );
};

export default MainWorkspace;
