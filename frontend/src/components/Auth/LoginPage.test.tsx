import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import LoginPage from "./LoginPage.tsx";
import { AuthContext, type AuthContextValue } from "../../context/authContext.ts";
import { ApiError } from "../../api/client.ts";

const renderLogin = (overrides: Partial<AuthContextValue> = {}) => {
    const auth: AuthContextValue = {
        user: null,
        restoring: false,
        login: vi.fn(async () => {}),
        register: vi.fn(async () => {}),
        logout: vi.fn(),
        ...overrides,
    };
    render(<AuthContext.Provider value={auth}><LoginPage /></AuthContext.Provider>);
    return auth;
};

describe("LoginPage", () => {
    it("logs in", async () => {
        const auth = renderLogin();
        await userEvent.type(screen.getByLabelText("Email"), "a@example.com");
        await userEvent.type(screen.getByLabelText("Password"), "secret");
        await userEvent.click(screen.getByRole("button", { name: "Log in" }));
        expect(auth.login).toHaveBeenCalledWith("a@example.com", "secret");
    });

    it("switches to registration and enforces the password length", async () => {
        const auth = renderLogin();
        await userEvent.click(screen.getByRole("button", { name: /Create one/ }));
        await userEvent.type(screen.getByLabelText("Email"), "a@example.com");
        await userEvent.type(screen.getByLabelText("Password"), "short");
        await userEvent.click(screen.getByRole("button", { name: "Create account" }));
        expect(screen.getByText(/at least 8 characters/)).toBeInTheDocument();
        expect(auth.register).not.toHaveBeenCalled();
    });

    it("shows the server's error", async () => {
        renderLogin({ login: vi.fn().mockRejectedValue(new ApiError(401, "invalid email or password")) });
        await userEvent.type(screen.getByLabelText("Email"), "a@example.com");
        await userEvent.type(screen.getByLabelText("Password"), "wrong");
        await userEvent.click(screen.getByRole("button", { name: "Log in" }));
        expect(await screen.findByText("invalid email or password")).toBeInTheDocument();
    });
});
