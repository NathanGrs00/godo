import React, { useState } from "react";
import { useAuth } from "../../context/authContext.ts";
import FieldError from "../Common/FieldError.tsx";
import { errorsFromSave, hasErrors, validateCredentials, type FormErrors } from "../../utils/validation.ts";
import "../../styles/Pages/LoginPage.css";

/**
 * Log in or create an account. Logging in with the same account on another
 * device shows the same tasks, kept in sync live.
 */
const LoginPage: React.FC = () => {
    const { login, register } = useAuth();
    const [mode, setMode] = useState<"login" | "register">("login");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [errors, setErrors] = useState<FormErrors>({});
    const [busy, setBusy] = useState(false);

    const isLogin = mode === "login";

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const found = validateCredentials({ email, password }, mode);
        setErrors(found);
        if (hasErrors(found)) return;

        setBusy(true);
        try {
            if (isLogin) await login(email, password);
            else await register(email, password);
        } catch (err) {
            setErrors(errorsFromSave(err, ["email", "password"]));
            setBusy(false);
        }
    };

    return (
        <main className="login-page">
            <form className="login-form" onSubmit={handleSubmit} noValidate aria-label={isLogin ? "Log in" : "Create account"}>
                <h1>GoDo</h1>
                <p className="login-subtitle">{isLogin ? "Log in to see your tasks on any device." : "Create an account to sync your tasks."}</p>

                <label htmlFor="email">Email</label>
                <input id="email" type="email" autoComplete="email" value={email} onChange={e => setEmail(e.target.value)} />
                <FieldError message={errors.email} />

                <label htmlFor="password">Password</label>
                <input
                    id="password"
                    type="password"
                    autoComplete={isLogin ? "current-password" : "new-password"}
                    value={password}
                    onChange={e => setPassword(e.target.value)}
                />
                <FieldError message={errors.password} />
                <FieldError message={errors.form} />

                <button type="submit" className="save-button" disabled={busy}>
                    {busy ? "Please wait…" : isLogin ? "Log in" : "Create account"}
                </button>
                <button
                    type="button"
                    className="link-button"
                    onClick={() => {
                        setMode(isLogin ? "register" : "login");
                        setErrors({});
                    }}
                >
                    {isLogin ? "No account yet? Create one" : "Already have an account? Log in"}
                </button>
            </form>
        </main>
    );
};

export default LoginPage;
