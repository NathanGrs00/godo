/**
 * Small wrapper around fetch for talking to the GoDo API.
 *
 * It adds the base URL and the login token, sends and parses JSON, and turns
 * error responses into ApiError so callers can react to the status code
 * (for example 409 when an item was changed on another device).
 */

export const API_URL: string = (import.meta.env.VITE_API_URL ?? "http://localhost:8080").replace(/\/$/, "");

/** An error response from the API. */
export class ApiError extends Error {
    readonly status: number;
    /** The invalid field, for validation errors (status 400). */
    readonly field?: string;

    constructor(status: number, message: string, field?: string) {
        super(message);
        this.name = "ApiError";
        this.status = status;
        this.field = field;
    }
}

let authToken: string | null = null;
let onUnauthorized: (() => void) | null = null;

/** Sets the token sent as "Authorization: Bearer <token>" on every request. */
export const setAuthToken = (token: string | null) => {
    authToken = token;
};

/** Returns the current login token. */
export const getAuthToken = () => authToken;

/** Registers a callback for 401 responses, used to log the user out when the token expires. */
export const setUnauthorizedHandler = (handler: (() => void) | null) => {
    onUnauthorized = handler;
};

/** Headers for an authenticated request. */
export const authHeaders = (): Record<string, string> =>
    authToken ? { Authorization: `Bearer ${authToken}` } : {};

interface RequestOptions {
    method?: "GET" | "POST" | "PUT" | "DELETE";
    body?: unknown;
}

/** Sends a request and returns the parsed JSON body (undefined for 204 No Content). */
export const apiRequest = async <T>(path: string, { method = "GET", body }: RequestOptions = {}): Promise<T> => {
    let res: Response;
    try {
        res = await fetch(`${API_URL}${path}`, {
            method,
            headers: {
                ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
                ...authHeaders(),
            },
            body: body !== undefined ? JSON.stringify(body) : undefined,
        });
    } catch {
        throw new ApiError(0, "Cannot reach the server. Check your connection.");
    }

    if (res.status === 401 && authToken) {
        onUnauthorized?.();
    }

    if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new ApiError(res.status, data.error ?? `Request failed (${res.status})`, data.field);
    }

    if (res.status === 204) {
        return undefined as T;
    }
    return res.json() as Promise<T>;
};
