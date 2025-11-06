import './styles/App.css'
import SideNavBar from "./components/NavBar/SideNavBar.tsx";
import SideNavBarItem from "./components/NavBar/SideNavBarItem.tsx";
import UserIcon from "./assets/icons/NavBar/UserIcon.tsx";
import SearchIcon from "./assets/icons/NavBar/SearchIcon.tsx";
import SettingsIcon from "./assets/icons/NavBar/SettingsIcon.tsx";
import DeadlineIcon from "./assets/icons/NavBar/DeadlineIcon.tsx";
import TaskIcon from "./assets/icons/NavBar/TaskIcon.tsx";
import MainWorkspace from "./components/MainWorkspace.tsx";

function App() {
    return (
        <>
            <SideNavBar>
                <SideNavBarItem icon={<UserIcon/>} />
                <SideNavBarItem icon={<SearchIcon/>} />
                <SideNavBarItem icon={<DeadlineIcon/>} />
                <SideNavBarItem icon={<TaskIcon/>} />
                <SideNavBarItem icon={<SettingsIcon/>} />
            </SideNavBar>
            <MainWorkspace/>
        </>
    )
}

export default App
