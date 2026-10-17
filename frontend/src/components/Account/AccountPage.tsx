import React from "react";
import UserIcon from "../../assets/icons/NavBar/UserIcon.tsx";
import PageLayout from "../Common/PageLayout.tsx";
import { useAuth } from "../../context/authContext.ts";
import { useData } from "../../context/dataContext.ts";
import "../../styles/Pages/Pages.css";

/** Shows who is logged in and the sync state, and lets the user log out. */
const AccountPage: React.FC = () => {
    const { user, logout } = useAuth();
    const { tasks, deadlines, tags, reload } = useData();

    return (
        <PageLayout title="Account" icon={<UserIcon />}>
            <section className="page-section">
                <p>Logged in as <strong>{user?.email}</strong></p>
                <p className="hint">
                    Log in with this account on another device to see the same tasks.
                    Changes show up on all your open devices straight away.
                </p>
                <p>{tasks.length} tasks · {deadlines.length} deadlines · {tags.length} tags</p>
                <div className="page-actions">
                    <button type="button" className="save-button" onClick={() => reload()}>Sync now</button>
                    <button type="button" className="cancel-button" onClick={logout}>Log out</button>
                </div>
            </section>
        </PageLayout>
    );
};

export default AccountPage;
