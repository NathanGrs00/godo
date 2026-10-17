import React from "react";
import { useData } from "../../context/dataContext.ts";
import "../../styles/Common/SyncStatus.css";

const LABELS = {
    live: "Synced",
    connecting: "Connecting…",
    offline: "Offline – retrying",
} as const;

/** A dot and label showing whether live sync with other devices is working. */
const SyncStatus: React.FC = () => {
    const { status } = useData();
    return (
        <span className={`sync-status sync-${status}`} role="status" title="Live sync with your other devices">
            <span className="sync-dot" />
            {LABELS[status]}
        </span>
    );
};

export default SyncStatus;
