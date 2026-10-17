/** Helpers for showing and editing deadline dates. */

const pad = (n: number) => String(n).padStart(2, "0");

/**
 * Formats an ISO date for an <input type="datetime-local">, in the device's
 * own time zone (e.g. "2026-03-01T09:00").
 */
export const toLocalInputValue = (iso: string): string => {
    const d = new Date(iso);
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
};

/** Converts a datetime-local value (local time) to an ISO string in UTC. */
export const fromLocalInputValue = (value: string): string => new Date(value).toISOString();

/** Describes how far away a date is, e.g. "in 3 days", "2 hours ago", "today". */
export const relativeTime = (iso: string, now: Date = new Date()): string => {
    const diffMs = new Date(iso).getTime() - now.getTime();
    const minutes = Math.round(diffMs / 60_000);
    const hours = Math.round(diffMs / 3_600_000);
    const days = Math.round(diffMs / 86_400_000);

    const phrase = (n: number, unit: string) => {
        const abs = Math.abs(n);
        const label = `${abs} ${unit}${abs === 1 ? "" : "s"}`;
        return n >= 0 ? `in ${label}` : `${label} ago`;
    };

    if (Math.abs(minutes) < 1) return "now";
    if (Math.abs(minutes) < 60) return phrase(minutes, "minute");
    if (Math.abs(hours) < 24) return phrase(hours, "hour");
    return phrase(days, "day");
};

/** Sorts deadlines (or anything with a date) soonest first. */
export const byDate = <T extends { date: string }>(a: T, b: T): number =>
    new Date(a.date).getTime() - new Date(b.date).getTime();
