/**
 * Form validation that mirrors the server's rules, so users see mistakes
 * before saving. The server still validates everything.
 */
import { PRIORITIES, type Priority } from "../types";

/** Field name to error message. Empty when the form is valid. */
export type FormErrors = Partial<Record<string, string>>;

export const LIMITS = {
    taskTitle: 200,
    deadlineTitle: 200,
    description: 2000,
    tagTitle: 50,
} as const;

const checkTitle = (errors: FormErrors, title: string, max: number) => {
    if (!title.trim()) errors.title = "Title is required";
    else if (title.trim().length > max) errors.title = `Title must be at most ${max} characters`;
};

const checkDescription = (errors: FormErrors, description: string) => {
    if (description.trim().length > LIMITS.description) {
        errors.description = `Description must be at most ${LIMITS.description} characters`;
    }
};

export const validateTask = (form: { title: string; description: string; priority: string }): FormErrors => {
    const errors: FormErrors = {};
    checkTitle(errors, form.title, LIMITS.taskTitle);
    checkDescription(errors, form.description);
    if (!PRIORITIES.includes(form.priority as Priority)) errors.priority = "Priority must be selected";
    return errors;
};

export const validateDeadline = (form: { title: string; description: string; date: string }): FormErrors => {
    const errors: FormErrors = {};
    checkTitle(errors, form.title, LIMITS.deadlineTitle);
    checkDescription(errors, form.description);
    if (!form.date || Number.isNaN(new Date(form.date).getTime())) errors.date = "Date is required";
    return errors;
};

export const validateTag = (form: { title: string; color: string }): FormErrors => {
    const errors: FormErrors = {};
    checkTitle(errors, form.title, LIMITS.tagTitle);
    if (!/^#[0-9a-fA-F]{6}$/.test(form.color)) errors.color = "Colour must look like #8CFBDE";
    return errors;
};

/**
 * Checks login or registration input. The minimum password length only
 * applies when creating an account.
 */
export const validateCredentials = (form: { email: string; password: string }, mode: "login" | "register"): FormErrors => {
    const errors: FormErrors = {};
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) errors.email = "Enter a valid email address";
    if (mode === "login" && !form.password) errors.password = "Password is required";
    if (mode === "register" && form.password.length < 8) errors.password = "Password must be at least 8 characters";
    return errors;
};

export const hasErrors = (errors: FormErrors): boolean => Object.keys(errors).length > 0;

/** Shown when another device changed the item while this form was open. */
export const CONFLICT_MESSAGE =
    "This was changed on another device in the meantime. Save again to keep your version, or Cancel to keep theirs.";

/**
 * Turns an error thrown while saving into form errors. Validation errors
 * from the server are shown under their field; anything else goes in "form".
 */
export const errorsFromSave = (err: unknown, fields: string[]): FormErrors => {
    if ((err as { status?: number })?.status === 409) {
        return { form: CONFLICT_MESSAGE };
    }
    if (err instanceof Error) {
        const field = (err as { field?: string }).field;
        if (field && fields.includes(field)) return { [field]: err.message };
        return { form: err.message };
    }
    return { form: "Something went wrong" };
};
