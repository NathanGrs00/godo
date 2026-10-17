import './styles/App.css'
import { Routes, Route, useNavigate, useLocation, Navigate } from "react-router-dom";
import SideNavBar from "./components/NavBar/SideNavBar.tsx";
import SideNavBarItem from "./components/NavBar/SideNavBarItem.tsx";
import UserIcon from "./assets/icons/NavBar/UserIcon.tsx";
import SearchIcon from "./assets/icons/NavBar/SearchIcon.tsx";
import SettingsIcon from "./assets/icons/NavBar/SettingsIcon.tsx";
import DeadlineIcon from "./assets/icons/NavBar/DeadlineIcon.tsx";
import TaskIcon from "./assets/icons/NavBar/TaskIcon.tsx";
import MainWorkspace from "./components/MiddleSection/MainWorkspace.tsx";
import LoginPage from "./components/Auth/LoginPage.tsx";
import AccountPage from "./components/Account/AccountPage.tsx";
import SearchPage from "./components/Search/SearchPage.tsx";
import SettingsPage from "./components/Settings/SettingsPage.tsx";
import DataProvider from "./context/DataProvider.tsx";
import { useAuth } from "./context/authContext.ts";

/** Navigation entries, top to bottom. The last one is pinned to the bottom. */
const NAV_ITEMS = [
    { path: "/account", label: "Account", icon: <UserIcon/> },
    { path: "/search", label: "Search", icon: <SearchIcon/> },
    { path: "/deadlines", label: "Deadlines", icon: <DeadlineIcon/> },
    { path: "/", label: "Tasks", icon: <TaskIcon/> },
    { path: "/settings", label: "Settings", icon: <SettingsIcon/> },
];

function App() {
    const { user, restoring } = useAuth();

    if (restoring) return <p className="empty-message">Loading...</p>;
    if (!user) return <LoginPage />;

    return (
        <DataProvider>
            <AppShell />
        </DataProvider>
    );
}

/** The logged-in app: side navigation and the current page. */
function AppShell() {
    const navigate = useNavigate();
    const { pathname } = useLocation();

    return (
        <>
            <SideNavBar>
                {NAV_ITEMS.map(item => (
                    <SideNavBarItem
                        key={item.path}
                        icon={item.icon}
                        label={item.label}
                        active={pathname === item.path}
                        onClick={() => navigate(item.path)}
                    />
                ))}
            </SideNavBar>

            <Routes>
                <Route path="/" element={<MainWorkspace route="tasks"/>} />
                <Route path="/deadlines" element={<MainWorkspace route="deadlines"/>} />
                <Route path="/search" element={<SearchPage/>} />
                <Route path="/account" element={<AccountPage/>} />
                <Route path="/settings" element={<SettingsPage/>} />
                <Route path="*" element={<Navigate to="/" replace/>} />
            </Routes>
        </>
    );
}

export default App
