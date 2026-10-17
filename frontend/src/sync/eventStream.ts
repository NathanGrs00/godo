/**
 * Live connection to the server's /events stream (Server-Sent Events).
 *
 * The browser's EventSource cannot send an Authorization header, so this
 * reads the stream with fetch instead. When the connection drops it
 * reconnects with exponential backoff; every (re)connection starts with a
 * "ready" message, after which the app reloads its data to catch up on
 * anything it missed while offline.
 */
import { API_URL, authHeaders } from "../api/client.ts";
import type { Deadline, Tag, Task } from "../types";

/** A change made on any of the user's devices. */
export type SyncEvent =
    | { type: "task.created" | "task.updated"; id: string; data: Task }
    | { type: "task.deleted"; id: string }
    | { type: "deadline.created" | "deadline.updated"; id: string; data: Deadline }
    | { type: "deadline.deleted"; id: string }
    | { type: "tag.created" | "tag.updated"; id: string; data: Tag }
    | { type: "tag.deleted"; id: string };

/** One raw message from the stream. */
export interface SSEMessage {
    /** The "event:" field; "message" when absent. */
    event: string;
    data: string;
}

/**
 * Incremental Server-Sent Events parser. Feed it text chunks as they arrive
 * (chunks may split lines anywhere); it returns the complete messages.
 */
export class SSEParser {
    private buffer = "";
    private event = "";
    private data: string[] = [];

    push(chunk: string): SSEMessage[] {
        this.buffer += chunk;
        const messages: SSEMessage[] = [];
        const lines = this.buffer.split(/\r\n|\r|\n/);
        // The last piece may be an incomplete line; keep it for the next chunk
        this.buffer = lines.pop() ?? "";

        for (const line of lines) {
            if (line === "") {
                // A blank line ends a message
                if (this.data.length > 0) {
                    messages.push({ event: this.event || "message", data: this.data.join("\n") });
                }
                this.event = "";
                this.data = [];
            } else if (line.startsWith(":")) {
                // Comment, used by the server as a keep-alive
            } else {
                const colon = line.indexOf(":");
                const field = colon === -1 ? line : line.slice(0, colon);
                let value = colon === -1 ? "" : line.slice(colon + 1);
                if (value.startsWith(" ")) value = value.slice(1);
                if (field === "event") this.event = value;
                if (field === "data") this.data.push(value);
            }
        }
        return messages;
    }
}

export type ConnectionStatus = "connecting" | "live" | "offline";

export interface StreamHandlers {
    /** Called for every change event. */
    onEvent: (event: SyncEvent) => void;
    /** Called each time the stream (re)connects. */
    onReady: () => void;
    /** Called whenever the connection status changes. */
    onStatus: (status: ConnectionStatus) => void;
    /** Called when the server rejects the token. Reconnecting stops. */
    onUnauthorized: () => void;
}

/** Delay before reconnect attempt n (0-based): 1s, 2s, 4s ... up to 30s. */
export const backoffDelay = (attempt: number): number => Math.min(30_000, 1000 * 2 ** attempt);

/**
 * Opens the event stream and keeps it open until the returned function is
 * called. Uses the current login token from the API client.
 */
export const connectEventStream = (handlers: StreamHandlers): (() => void) => {
    let stopped = false;
    let controller: AbortController | null = null;
    let retryTimer: ReturnType<typeof setTimeout> | undefined;
    let attempt = 0;

    const scheduleReconnect = () => {
        if (stopped) return;
        handlers.onStatus("offline");
        retryTimer = setTimeout(run, backoffDelay(attempt++));
    };

    const run = async () => {
        if (stopped) return;
        handlers.onStatus("connecting");
        controller = new AbortController();
        try {
            const res = await fetch(`${API_URL}/events`, {
                headers: { Accept: "text/event-stream", ...authHeaders() },
                signal: controller.signal,
            });
            if (res.status === 401) {
                stopped = true;
                handlers.onUnauthorized();
                return;
            }
            if (!res.ok || !res.body) throw new Error(`stream failed (${res.status})`);

            const reader = res.body.pipeThrough(new TextDecoderStream()).getReader();
            const parser = new SSEParser();
            for (;;) {
                const { value, done } = await reader.read();
                if (done) break;
                for (const message of parser.push(value)) {
                    if (message.event === "ready") {
                        attempt = 0;
                        handlers.onStatus("live");
                        handlers.onReady();
                    } else {
                        handlers.onEvent(JSON.parse(message.data) as SyncEvent);
                    }
                }
            }
        } catch (err) {
            if (stopped) return;
            console.warn("Sync stream interrupted:", err);
        }
        scheduleReconnect();
    };

    // Reconnect straight away when the device comes back online
    const handleOnline = () => {
        if (stopped) return;
        clearTimeout(retryTimer);
        controller?.abort();
        attempt = 0;
        run();
    };
    window.addEventListener("online", handleOnline);

    run();

    return () => {
        stopped = true;
        clearTimeout(retryTimer);
        controller?.abort();
        window.removeEventListener("online", handleOnline);
    };
};
