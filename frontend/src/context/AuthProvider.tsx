import React, { useCallback, useEffect, useMemo, useState } from "react";
import { setAuthToken, setUnauthorizedHandler } from "../api/client.ts";
import * as authService from "../services/authService.ts";
import type { Session, User } from "../types";
import { AuthContext, type AuthContextValue } from "./authContext.ts";

/** Where the login token is kept so the user stays logged in after a reload. */
export const TOKEN_STORAGE_KEY = "godo.token";

const readToken = (): string | null => {
    try {
        return localStorage.getItem(TOKEN_STORAGE_KEY);
    } catch {
        return null;
    }
};

const writeToken = (token: string | null) => {
    try {
        if (token) localStorage.setItem(TOKEN_STORAGE_KEY, token);
        else localStorage.removeItem(TOKEN_STORAGE_KEY);
    } catch {
        // Storage can be unavailable (private mode); the session then lasts until reload
    }
};

/**
 * Keeps track of who is logged in. Each device logs in separately and gets
 * its own token for the same account, which is what lets the data sync.
 */
const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [restoring, setRestoring] = useState(() => readToken() !== null);

    const logout = useCallback(() => {
        setAuthToken(null);
        writeToken(null);
        setUser(null);
    }, []);

    const startSession = useCallback((session: Session) => {
        setAuthToken(session.token);
        writeToken(session.token);
        setUser(session.user);
    }, []);

    // Log out automatically when the server says the token expired
    useEffect(() => {
        setUnauthorizedHandler(logout);
        return () => setUnauthorizedHandler(null);
    }, [logout]);

    // Restore a login saved on this device
    useEffect(() => {
        const token = readToken();
        if (!token) return;
        setAuthToken(token);
        authService.getMe()
            .then(setUser)
            .catch(() => logout())
            .finally(() => setRestoring(false));
    }, [logout]);

    const value = useMemo<AuthContextValue>(() => ({
        user,
        restoring,
        logout,
        login: async (email, password) => startSession(await authService.login({ email, password })),
        register: async (email, password) => startSession(await authService.register({ email, password })),
    }), [user, restoring, logout, startSession]);

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export default AuthProvider;
