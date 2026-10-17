import { afterEach, describe, expect, it, vi } from "vitest";
import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import DataProvider from "./DataProvider.tsx";
import { useData } from "./dataContext.ts";
import { AuthContext, type AuthContextValue } from "./authContext.ts";
import { makeTask } from "../test/fixtures.tsx";

const json = (body: unknown, status = 200) =>
    new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });

/**
 * A fake server: answers the list endpoints and keeps the /events stream
 * open so the test can push events, as if another device changed something.
 */
const fakeServer = (initialTasks = [makeTask({ title: "From server" })]) => {
    let push: (text: string) => void = () => {};
    let tasks = initialTasks;
    const fetchMock = vi.spyOn(globalThis, "fetch").mockImplementation(async (input, init) => {
        const url = String(input);
        const method = init?.method ?? "GET";
        if (url.endsWith("/events")) {
            const body = new ReadableStream<Uint8Array>({
                start(controller) {
                    push = text => controller.enqueue(new TextEncoder().encode(text));
                    push("event: ready\ndata: {}\n\n");
                },
            });
            return new Response(body, { status: 200 });
        }
        if (url.endsWith("/tasks") && method === "GET") return json({ tasks });
        if (url.endsWith("/deadlines")) return json({ deadlines: [] });
        if (url.endsWith("/tags")) return json({ tags: [] });
        if (url.includes("/tasks/") && method === "PUT") {
            return json({ error: "this item was changed on another device; reload and try again" }, 409);
        }
        throw new Error(`unexpected request ${method} ${url}`);
    });
    return {
        fetchMock,
        pushEvent: (event: unknown) => act(() => push(`data: ${JSON.stringify(event)}\n\n`)),
        setTasks: (next: typeof tasks) => { tasks = next; },
    };
};

const auth: AuthContextValue = { user: null, restoring: false, login: vi.fn(), register: vi.fn(), logout: vi.fn() };

/** Shows the task titles, connection status and a button that edits the first task. */
const Probe = () => {
    const { tasks, status, updateTask } = useData();
    return (
        <div>
            <span data-testid="status">{status}</span>
            <ul>{tasks.map(t => <li key={t.id}>{t.title} v{t.version}</li>)}</ul>
            <button onClick={() => updateTask(tasks[0], { completed: true }).catch(() => {})}>edit</button>
        </div>
    );
};

const renderProvider = () =>
    render(
        <AuthContext.Provider value={auth}>
            <DataProvider><Probe /></DataProvider>
        </AuthContext.Provider>,
    );

describe("DataProvider", () => {
    afterEach(() => vi.useRealTimers());

    it("loads data once the live stream is ready", async () => {
        fakeServer();
        renderProvider();
        expect(await screen.findByText("From server v1")).toBeInTheDocument();
        expect(screen.getByTestId("status")).toHaveTextContent("live");
    });

    it("applies changes made on another device", async () => {
        const server = fakeServer([]);
        renderProvider();
        await screen.findByText("live");

        const task = makeTask({ title: "Added on phone" });
        await server.pushEvent({ type: "task.created", id: task.id, data: task });
        expect(await screen.findByText("Added on phone v1")).toBeInTheDocument();

        await server.pushEvent({ type: "task.updated", id: task.id, data: { ...task, title: "Renamed on phone", version: 2 } });
        expect(await screen.findByText("Renamed on phone v2")).toBeInTheDocument();

        await server.pushEvent({ type: "task.deleted", id: task.id });
        expect(screen.queryByText(/on phone/)).not.toBeInTheDocument();
    });

    it("reloads when an edit conflicts with a newer version", async () => {
        const original = makeTask({ title: "Original" });
        const server = fakeServer([original]);
        renderProvider();
        await screen.findByText("Original v1");

        server.setTasks([{ ...original, title: "Changed elsewhere", version: 2 }]);
        await userEvent.click(screen.getByRole("button", { name: "edit" }));

        expect(await screen.findByText("Changed elsewhere v2")).toBeInTheDocument();
    });
});
