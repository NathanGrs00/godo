import React from "react";
import SettingsIcon from "../../assets/icons/NavBar/SettingsIcon.tsx";
import PageLayout from "../Common/PageLayout.tsx";
import TagManager from "./TagManager.tsx";

/** Settings: currently tag management. */
const SettingsPage: React.FC = () => (
    <PageLayout title="Settings" icon={<SettingsIcon />}>
        <TagManager />
    </PageLayout>
);

export default SettingsPage;
