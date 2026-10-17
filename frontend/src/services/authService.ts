import { apiRequest } from "../api/client.ts";
import type { Session, User } from "../types";

interface Credentials {
    email: string;
    password: string;
}

export const register = (credentials: Credentials): Promise<Session> =>
    apiRequest<Session>("/auth/register", { method: "POST", body: credentials });

export const login = (credentials: Credentials): Promise<Session> =>
    apiRequest<Session>("/auth/login", { method: "POST", body: credentials });

export const getMe = async (): Promise<User> =>
    (await apiRequest<{ user: User }>("/auth/me")).user;
