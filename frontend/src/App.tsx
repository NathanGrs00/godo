import './styles/App.css'
import { Routes, Route, useNavigate } from "react-router-dom";
import SideNavBar from "./components/NavBar/SideNavBar.tsx";
import SideNavBarItem from "./components/NavBar/SideNavBarItem.tsx";
import UserIcon from "./assets/icons/NavBar/UserIcon.tsx";
import SearchIcon from "./assets/icons/NavBar/SearchIcon.tsx";
import SettingsIcon from "./assets/icons/NavBar/SettingsIcon.tsx";
import DeadlineIcon from "./assets/icons/NavBar/DeadlineIcon.tsx";
import TaskIcon from "./assets/icons/NavBar/TaskIcon.tsx";
import MainWorkspace from "./components/MiddleSection/MainWorkspace.tsx";

function App() {
    const navigate = useNavigate();

    return (
        <>
            <SideNavBar>
                <SideNavBarItem icon={<UserIcon/>} onClick={() => navigate("/users")} />
                <SideNavBarItem icon={<SearchIcon/>} onClick={() => navigate("/search")} />
                <SideNavBarItem icon={<DeadlineIcon/>} onClick={() => navigate("/deadlines")} />
                <SideNavBarItem icon={<TaskIcon/>} onClick={() => navigate("/")} />
                <SideNavBarItem icon={<SettingsIcon/>} onClick={() => navigate("/settings")} />
            </SideNavBar>

            <Routes>
                <Route path="/" element={<MainWorkspace />} />
                <Route path="/users" element={<div>Users Page</div>} />
                <Route path="/search" element={<div>Search Page</div>} />
                <Route path="/deadlines" element={<div>Deadlines Page</div>} />
                <Route path="/settings" element={<div>Settings Page</div>} />
            </Routes>
        </>
    );
}

export default App
