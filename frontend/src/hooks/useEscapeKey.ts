import { useEffect } from "react";

/** Calls onEscape when the Escape key is pressed. */
export const useEscapeKey = (onEscape: () => void) => {
    useEffect(() => {
        const handleKey = (event: KeyboardEvent) => {
            if (event.key === "Escape") onEscape();
        };
        window.addEventListener("keydown", handleKey);
        return () => window.removeEventListener("keydown", handleKey);
    }, [onEscape]);
};
