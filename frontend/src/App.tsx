import './styles/App.css'
import SideNavBar from "./components/SideNavBar.tsx";
import SideNavBarItem from "./components/SideNavBarItem.tsx";
import UserIcon from "./assets/icons/UserIcon.tsx";
import SearchIcon from "./assets/icons/SearchIcon.tsx";
import SettingsIcon from "./assets/icons/SettingsIcon.tsx";

function App() {
    return (
        <SideNavBar links={[]}>
            <SideNavBarItem icon={<UserIcon/>} />
            <SideNavBarItem icon={<SearchIcon/>} />
            <SideNavBarItem icon={<SettingsIcon/>} />
        </SideNavBar>
    )
}

export default App
