import { createContext, useContext } from "react";
import type { User } from "../types";

export interface AuthContextValue {
    /** The logged-in user, or null when logged out. */
    user: User | null;
    /** True while a stored login is being checked on start-up. */
    restoring: boolean;
    login: (email: string, password: string) => Promise<void>;
    register: (email: string, password: string) => Promise<void>;
    logout: () => void;
}

export const AuthContext = createContext<AuthContextValue | null>(null);

/** Returns the auth state. Must be used inside <AuthProvider>. */
export const useAuth = (): AuthContextValue => {
    const value = useContext(AuthContext);
    if (!value) throw new Error("useAuth must be used inside <AuthProvider>");
    return value;
};
